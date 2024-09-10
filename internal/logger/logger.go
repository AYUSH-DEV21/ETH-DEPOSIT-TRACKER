package logger

import (
	"log/slog"
	"os"
)

func Init() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	slog.SetDefault(logger)
}
