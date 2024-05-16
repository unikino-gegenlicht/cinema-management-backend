package db

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"

	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
)

func Migrate(l *slog.Logger) error {
	conString, err := buildConnectionString(l)
	if err != nil {
		return fmt.Errorf("db: migrate: error while building connection string: %w", err)
	}
	s, err := iofs.New(Migrations, "migrations")
	if err != nil {
		return fmt.Errorf("db: migrate: unable to load embeds: %w", err)
	}
	conString = strings.ReplaceAll(conString, "postgres://", "pgx5://")
	m, err := migrate.NewWithSourceInstance("iofs", s, conString)
	if err != nil {
		return fmt.Errorf("db: migrate: error creating migrations: %w", err)
	}
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return err
		}
		return fmt.Errorf("db: migrate: error during migrations: %w", err)
	}
	return nil
}
