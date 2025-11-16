package cache

import "github.com/fatkulnurk/foundation/support"

type Config struct {
	Prefix string
}

func LoadConfig() *Config {
	return &Config{
		Prefix: support.GetEnv("CACHE_PREFIX", ""), // example: foundation:
	}
}
