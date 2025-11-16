package app

import (
	"sync"

	"github.com/fatkulnurk/foundation/shared"
)

var (
	app  *App
	once sync.Once
)

type App struct {
	cfg *Config
}

func init() {
	New()
}

func New() *App {
	once.Do(func() {
		cfg := LoadConfig()
		app = &App{
			cfg: cfg,
		}
	})

	return app
}

func Name() string {
	return app.cfg.name
}

func Version() string {
	return app.cfg.version
}

func Env() string {
	return app.cfg.env
}

func IsDevelopment() bool {
	return app.cfg.env == shared.EnvironmentDevelopment
}

func IsTesting() bool {
	return app.cfg.env == shared.EnvironmentTest
}

func IsStaging() bool {
	return app.cfg.env == shared.EnvironmentStaging
}

func IsProduction() bool {
	return app.cfg.env == shared.EnvironmentProduction
}
