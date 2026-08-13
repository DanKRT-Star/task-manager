package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init() {
	env := os.Getenv("APP_ENV")

	if env == "production" {
		// Production: JSON thuần, dễ cho log aggregator (Datadog, Loki...) parse
		Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		// Development: có màu, dễ đọc trên terminal
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		Log = zerolog.New(output).With().Timestamp().Logger()
	}
}