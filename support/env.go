package support

import (
	"os"
	"strconv"
	"time"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

func GetBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolVal
}

func GetDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	durationVal, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return durationVal
}
