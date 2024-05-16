package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/unikino-gegenlicht/cinema-management/backend/db"
	"github.com/unikino-gegenlicht/cinema-management/backend/logging"

	"github.com/pterm/pterm"
)

var logger *slog.Logger

func init() {
	printStartupMessage()
	logging.Configure()
	logger = logging.Logger
	logger.Info("starting backend application")
	err := db.Connect(logger)
	if err != nil {
		logger.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}
	logger.Info("connected to the database")
	logger.Info("running database migrations")
	err = db.Migrate(logger)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("no migrations required")
		} else {
			logger.Error("unable to run database migrations", "error", err)
			os.Exit(1)
		}
	}
}

func printStartupMessage() {
	p := pterm.DefaultCenter
	p.Println("Cinema Management Plaform — Backend")
	msg := fmt.Sprintf("Version: %s\n", Commit)
	msg += fmt.Sprintf("Commit Date: %s\n", CommitTime)
	msg += "\n© Unikino GEGENLICHT (gegenlicht.net) — Licensed under the MIT License"
	p.WithCenterEachLineSeparately().Println(msg)
}
