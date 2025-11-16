package config

import (
	"time"

	"github.com/fatkulnurk/foundation/mailer"
	"github.com/fatkulnurk/foundation/queue"
)

type Config struct {
	App           *App
	Database      *Database
	DeliveryHttp  *DeliveryHttp
	DeliveryQueue *DeliveryQueue
	Redis         *Redis
	Queue         *queue.Config
	Schedule      *Schedule
	SMTP          *mailer.SMTPConfig
}

// App only this struct can deliver to module
type App struct {
	Name        string
	Environment string
	Version     string
}

type DeliveryHttp struct {
	Prefork       bool
	CaseSensitive bool
	StrictRouting bool
	BodyLimit     int
	ServerHeader  string
}

type DeliveryQueue struct {
	Concurrency int
}

type Database struct {
	User            string
	Password        string
	Host            string
	Port            int
	Database        string
	Params          string // opsional: tambahan param seperti charset=utf8mb4
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type Redis struct {
	Addr            string
	Password        string
	DB              int
	PoolSize        int
	MinIdleConns    int
	ConnMaxLifetime time.Duration
	PoolTimeout     time.Duration
	ConnMaxIdleTime time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	DialTimeout     time.Duration
}

type Schedule struct {
	Timezone string
}
