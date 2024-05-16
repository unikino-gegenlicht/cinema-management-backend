package main

import (
	"fmt"
	"log/slog"

	"github.com/pterm/pterm"
)

var logger *slog.Logger

func init() {
	printStartupMessage()
	setupLogging()
	logger.Info("initializing backend")
	logger.Debug("connecting to database")
}

func printStartupMessage() {
	p := pterm.DefaultCenter
	p.Println("Cinema Management Plaform — Backend")
	msg := fmt.Sprintf("Version: %s\n", Commit)
	msg += fmt.Sprintf("Commit Date: %s\n", CommitTime)
	msg += "\n© Unikino GEGENLICHT (gegenlicht.net) — Licensed under the MIT License"
	p.WithCenterEachLineSeparately().Println(msg)
}

func setupLogging() {
	handler := pterm.NewSlogHandler(&loggerConfig)
	logger = slog.New(handler)
}
