package config

import (
	"log/slog"
	"os"
)

func ValidateEnvVars(required []string) {
	exit := false
	for _, key := range required {
		_, exists := os.LookupEnv(key)
		if !exists {
			slog.Error("env var not set", "var", key)
			exit = true
		}
	}
	if exit {
		os.Exit(1)
	}
}
