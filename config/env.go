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
	if err != nil {
		return nil, err
	}

	return &Config{
		AppPort:        os.Getenv("APP_PORT"),
		GinMode:        os.Getenv("GIN_MODE"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpireHours: getIntEnv("JWT_EXPIRE_HOURS"),
	}, nil
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
