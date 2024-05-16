package logging

import (
	"log/slog"

	"github.com/pterm/pterm"
)

var Logger *slog.Logger

func Configure() {
	handler := pterm.NewSlogHandler(&loggerConfig)
	Logger = slog.New(handler)
}
