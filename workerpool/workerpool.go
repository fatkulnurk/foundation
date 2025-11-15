package workerpool

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sort"
	"sync"
	"time"
)

type Priority int

const (
	Low Priority = iota
	Normal
	High
)

// Job adalah tugas yang dijalankan oleh worker
type Job struct {
	Task       func(ctx context.Context) error
	Retry      int
	Timeout    time.Duration
	RetryDelay time.Duration
	Priority   Priority
	CreatedAt  time.Time // untuk antrian FIFO
}

type WorkerPool struct {
	mu          sync.Mutex
	workerCount int
	jobQueue    []Job
	cond        *sync.Cond
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	scaling     bool
}

func NewWorkerPool(initialWorkerCount int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	wp := &WorkerPool{
		workerCount: initialWorkerCount,
		jobQueue:    make([]Job, 0),
		ctx:         ctx,
		cancel:      cancel,
	}
	wp.cond = sync.NewCond(&wp.mu)

	// Mulai worker awal
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	return wp
}

func (wp *WorkerPool) Submit(job Job) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if wp.ctx.Err() != nil {
		return errors.New("worker pool has been stopped")
	}
	job.CreatedAt = time.Now()
	wp.jobQueue = append(wp.jobQueue, job)
	wp.cond.Signal() // bangunkan 1 worker
	return nil
}

func (wp *WorkerPool) getJob() (Job, bool) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for len(wp.jobQueue) == 0 {
		if wp.ctx.Err() != nil {
			return Job{}, false
		}
		wp.cond.Wait()
	}

	// Prioritaskan berdasarkan Priority & FIFO CreatedAt
	sort.SliceStable(wp.jobQueue, func(i, j int) bool {
		if wp.jobQueue[i].Priority == wp.jobQueue[j].Priority {
			return wp.jobQueue[i].CreatedAt.Before(wp.jobQueue[j].CreatedAt)
		}
		return wp.jobQueue[i].Priority > wp.jobQueue[j].Priority
	})

	job := wp.jobQueue[0]
	wp.jobQueue = wp.jobQueue[1:]
	return job, true
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for {
		select {
		case <-wp.ctx.Done():
			return
		default:
			job, ok := wp.getJob()
			if !ok {
				return
			}
			wp.executeWithRetry(id, job)
		}
	}
}

func (wp *WorkerPool) executeWithRetry(workerID int, job Job) {
	var err error
	for attempt := 1; attempt <= job.Retry+1; attempt++ {
		ctx, cancel := context.WithTimeout(wp.ctx, job.Timeout)

		func() {
			defer cancel()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[Worker %d] PANIC: %v\n%s", workerID, r, debug.Stack())
					err = fmt.Errorf("panic: %v", r)
				}
			}()

			err = job.Task(ctx)
		}()

		if err == nil {
			log.Printf("[Worker %d] ✅ Job success (Attempt %d)", workerID, attempt)
			return
		}

		log.Printf("[Worker %d] ❌ Job failed (Attempt %d): %v", workerID, attempt, err)

		if attempt <= job.Retry {
			time.Sleep(job.RetryDelay)
		}
	}
	log.Printf("[Worker %d] 🔥 Job permanently failed after %d attempts", workerID, job.Retry+1)
}

// ScaleTo menyesuaikan jumlah worker aktif
func (wp *WorkerPool) ScaleTo(newCount int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	diff := newCount - wp.workerCount
	if diff == 0 {
		return
	}

	wp.workerCount = newCount

	if diff > 0 {
		for i := 0; i < diff; i++ {
			wp.wg.Add(1)
			go wp.worker(wp.workerCount + i)
		}
		log.Printf("⬆️ Scaled up to %d workers", wp.workerCount)
	} else {
		log.Printf("⬇️ Scaled down to %d workers (pending stop after jobs)", wp.workerCount)
	}
}

func (wp *WorkerPool) Stop() {
	wp.cancel()
	wp.cond.Broadcast()
	wp.wg.Wait()
}
