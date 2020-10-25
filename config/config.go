package config

import (
	"log"
	"os"
	"strconv"
)

// Config struct encapsulates the library's configuration.
// On initialization it will attempt to load all necessary configuration from environment variables.
type Config struct {
	SourceType    string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDatabase int
}

// New creates a new Config object based on the existing environment variables.
func New() (Config, error) {
	c := Config{}
	c.SourceType = loadEnvWithDefault("TBQ_SOURCE", "redis")
	c.RedisHost = loadEnvWithDefault("TBQ_REDIS_HOST", "localhost")
	c.RedisPort = loadEnvWithDefault("TBQ_REDIS_PORT", "6379")
	c.RedisPassword = loadEnvWithDefault("TBQ_REDIS_PWD", "")
	db, err := strconv.Atoi(loadEnvWithDefault("TBQ_REDIS_DB", "0"))
	if err != nil {
		return Config{}, err
	}
	c.RedisDatabase = db

	return c, nil
}

func loadEnvWithDefault(key, value string) string {

	v, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("env %s not found - falling back to default: %s", key, value)
		return value
	}

	return v
}
