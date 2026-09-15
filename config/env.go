package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	GinMode        string
	JWTSecret      string
	JWTExpireHours int
}

func LoadEnv() (*Config, error) {
	err := godotenv.Load()

	// Environment variables may be supplied by the process (for example, in a
	// container or deployment platform), so a missing .env file is not fatal.
	// Always return a usable config to avoid a nil dereference in callers.
	cfg := &Config{
		AppPort:        os.Getenv("APP_PORT"),
		GinMode:        os.Getenv("GIN_MODE"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpireHours: getIntEnv("JWT_EXPIRE_HOURS"),
	}

	return cfg, err
}

func getIntEnv(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0
	}
	return value
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
