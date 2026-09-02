package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %v", err)
	}

	// Get the directory of the executable
	ex, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting executable path: %v", err)
	}
	exPath := filepath.Dir(ex)

	// List of possible locations for .env file
	envLocations := []string{
		filepath.Join(cwd, ".env"),
		filepath.Join(exPath, ".env"),
		"/var/www/schematics26-backend/.env", // Add the expected location when run as a service
	}

	// Try to load .env from each location
	for _, loc := range envLocations {
		err := godotenv.Load(loc)
		if err == nil {
			fmt.Printf("Loaded .env from: %s\n", loc)
			return nil
		}
	}

	return fmt.Errorf("no .env file found in any of the expected locations")
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
