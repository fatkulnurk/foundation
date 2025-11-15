package config

import (
	"time"

	"github.com/fatkulnurk/gostarter/pkg/support"
	"github.com/fatkulnurk/gostarter/shared/constant"
	"github.com/joho/godotenv"
)

func New(env string) *Config {
	// default environment is development
	envFile := ""
	if env == "" {
		env = constant.EnvironmentDevelopment
	}

	switch env {
	case constant.EnvironmentDevelopment:
		envFile = ".env"
	case constant.EnvironmentTest:
		envFile = ".env.test"
	}

	// load environment variables
	if err := godotenv.Load(envFile); err != nil {
		panic(err)
	}

	cfg := Config{
		App: &App{
			Environment: env,
			Name:        support.GetEnv("APP_NAME", "GoStarter"),
			Version:     support.GetEnv("APP_VERSION", "1.0.0"),
		},
		Database: &Database{
			User:            support.GetEnv("DB_USER", "root"),
			Password:        support.GetEnv("DB_PASSWORD", ""),
			Host:            support.GetEnv("DB_HOST", "localhost"),
			Port:            support.GetIntEnv("DB_PORT", 3306),
			Database:        support.GetEnv("DB_NAME", "gostarter"),
			Params:          support.GetEnv("DB_PARAMS", "charset=utf8mb4&parseTime=true"),
			MaxOpenConns:    support.GetIntEnv("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    support.GetIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: support.GetDurationEnv("DB_CONN_MAX_LIFETIME", time.Hour),
			ConnMaxIdleTime: support.GetDurationEnv("DB_CONN_MAX_IDLE_TIME", time.Minute*30),
		},
		DeliveryHttp: &DeliveryHttp{
			Prefork:       support.GetBoolEnv("HTTP_PREFORK", false),
			CaseSensitive: support.GetBoolEnv("HTTP_CASE_SENSITIVE", true),
			StrictRouting: support.GetBoolEnv("HTTP_STRICT_ROUTING", false),
			BodyLimit:     support.GetIntEnv("HTTP_BODY_LIMIT", 10*1024*1024),
			ServerHeader:  support.GetEnv("HTTP_SERVER_HEADER", "GoStarter"),
		},
		DeliveryQueue: &DeliveryQueue{
			Concurrency: support.GetIntEnv("QUEUE_CONCURRENCY", 10),
		},
		Redis: &Redis{
			Addr:            support.GetEnv("REDIS_ADDR", "localhost:6379"),
			Password:        support.GetEnv("REDIS_PASSWORD", ""),
			DB:              support.GetIntEnv("REDIS_DB", 0),
			PoolSize:        support.GetIntEnv("REDIS_POOL_SIZE", 10),
			MinIdleConns:    support.GetIntEnv("REDIS_MIN_IDLE_CONNS", 5),
			ConnMaxLifetime: support.GetDurationEnv("REDIS_CONN_MAX_LIFETIME", time.Hour),
			PoolTimeout:     support.GetDurationEnv("REDIS_POOL_TIMEOUT", time.Second*4),
			ConnMaxIdleTime: support.GetDurationEnv("REDIS_CONN_MAX_IDLE_TIME", time.Minute*30),
			ReadTimeout:     support.GetDurationEnv("REDIS_READ_TIMEOUT", time.Second*3),
			WriteTimeout:    support.GetDurationEnv("REDIS_WRITE_TIMEOUT", time.Second*3),
			DialTimeout:     support.GetDurationEnv("REDIS_DIAL_TIMEOUT", time.Second*5),
		},
		Queue: &Queue{
			Concurrency: support.GetIntEnv("QUEUE_WORKER_CONCURRENCY", 10),
		},
		Schedule: &Schedule{
			Timezone: support.GetEnv("SCHEDULE_TIMEZONE", "UTC"),
		},
		SMTP: &SMTP{
			Host:              support.GetEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:              support.GetIntEnv("SMTP_PORT", 587),
			Username:          support.GetEnv("SMTP_USERNAME", ""),
			Password:          support.GetEnv("SMTP_PASSWORD", ""),
			AuthType:          support.GetEnv("SMTP_AUTH_TYPE", "PLAIN"),
			WithTLSPortPolicy: support.GetIntEnv("SMTP_WITH_TLS_PORT_POLICY", 0),
		},
	}

	return &cfg
}
