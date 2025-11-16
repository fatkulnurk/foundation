package app

import (
	"sync"

	"github.com/fatkulnurk/foundation/shared"
	"github.com/fatkulnurk/foundation/support"
)

type Config struct {
	env     string
	name    string
	version string
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		cfg = &Config{
			env:     support.GetEnv("APP_ENV", shared.EnvironmentDevelopment),
			name:    support.GetEnv("APP_NAME", "Foundation"),
			version: support.GetEnv("APP_VERSION", "1.0.0"),
		}
	})

	return cfg
}
