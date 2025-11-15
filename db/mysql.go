package db

import (
	"database/sql"
	"fmt"

	"github.com/fatkulnurk/gostarter/pkg/config"
	_ "github.com/go-sql-driver/mysql"
)

// NewMySQL membuat koneksi MySQL baru berdasarkan konfigurasi
func NewMySQL(cfg *config.Database) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Params,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql connection: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Cek koneksi
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	return db, nil
}
