package util

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func LoadDotEnv(path string) error {
	if path == "" {
		path = ".env"
	}

	err := godotenv.Load(path)
	if err == nil {
		slog.Debug("loaded dotenv file", "path", path)
		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		slog.Debug("dotenv file not found, skipping", "path", path)
		return nil
	}
	return fmt.Errorf("load %s failed: %w", path, err)
}
