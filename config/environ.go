package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	envFilePath := filepath.Join(cwd, ".env")

	err = godotenv.Load(envFilePath)
	if err != nil {
		return err
	}
	return nil
}

func ComposeDBConnectionString() string {
	return os.Getenv("MONGO_CONNECTION_STRING")
}
