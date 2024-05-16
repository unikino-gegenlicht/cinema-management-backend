//go:build debug

package logging

import (
	"os"
	"time"

	"github.com/pterm/pterm"
)

var loggerConfig pterm.Logger = pterm.Logger{
	TimeFormat: time.RFC3339,
	Level:      pterm.LogLevelDebug,
	Formatter:  pterm.LogFormatterColorful,
	ShowTime:   true,
	Writer:     os.Stdout,
	MaxWidth:   pterm.GetTerminalWidth(),
}
