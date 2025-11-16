package app

import (
	"github.com/fatkulnurk/foundation/shared"
)

type App struct {
	cfg *Config
}

func New() *App {
	cfg := LoadConfig()
	return &App{
		cfg: cfg,
	}
}

func (a *App) Name() string {
	return a.cfg.name
}

func (a *App) Version() string {
	return a.cfg.version
}

func (a *App) Env() string {
	return a.cfg.env
}

func (a *App) IsDevelopment() bool {
	return a.cfg.env == shared.EnvironmentDevelopment
}

func (a *App) IsTesting() bool {
	return a.cfg.env == shared.EnvironmentTest
}

func (a *App) IsStaging() bool {
	return a.cfg.env == shared.EnvironmentStaging
}

func (a *App) IsProduction() bool {
	return a.cfg.env == shared.EnvironmentProduction
}
