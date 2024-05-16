//go:build !debug

package main

import (
	"os"
	"time"

	"github.com/pterm/pterm"
)

var loggerConfig pterm.Logger = pterm.Logger{
	TimeFormat: time.RFC3339,
	Level:      pterm.LogLevelInfo,
	Writer:     os.Stdout,
	ShowTime:   true,
	Formatter:  pterm.LogFormatterJSON,
	MaxWidth:   pterm.GetTerminalWidth(),
}
