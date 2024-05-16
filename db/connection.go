package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

const (
	envUser     = "PG_USER"
	envPassword = "PG_PASS"
	envHost     = "PG_HOST"
	envPort     = "PG_PORT"
	envDatabase = "PG_DB"
)

var errNoUserSet = errors.New("db: no user set in environment")
var errNoPasswordSet = errors.New("db: no password set in environment")
var errNoHostSet = errors.New("db: no host set in environment")

func buildConnectionString(l *slog.Logger) (string, error) {
	l.Debug("getting username from environment", "envKey", envUser)
	username, set := os.LookupEnv(envUser)
	if !set {
		return "", errNoUserSet
	}

	l.Debug("getting password from environment", "envKey", envPassword)
	password, set := os.LookupEnv(envPassword)
	if !set {
		return "", errNoPasswordSet
	}

	l.Debug("getting host from environment", "envKey", envHost)
	host, set := os.LookupEnv(envHost)
	if !set {
		return "", errNoHostSet
	}

	l.Debug("getting port from environment or using default", "envKey", envPort, "default", 5432)
	var port uint16
	rawPort, set := os.LookupEnv(envPort)
	if !set {
		l.Debug("no port set, using default")
		port = 5432
	} else {
		port64, err := strconv.ParseUint(rawPort, 10, 16)
		if err != nil {
			return "", fmt.Errorf("db: unable to use port: %w", err)
		}
		port = uint16(port64)
	}

	l.Debug("getting database name from environment", "envKey", envDatabase, "default", "cmp")
	database, set := os.LookupEnv(envDatabase)
	if !set {
		l.Debug("no database name set, using default")
		database = "cmp"
	}
	address := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, host, port, database)
	return address, nil
}

func Connect(l *slog.Logger) error {
	l.Debug("connecting to postgresql database")

	address, err := buildConnectionString(l)
	config, err := pgxpool.ParseConfig(address)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	Pool, err = pgxpool.NewWithConfig(context.TODO(), config)
	if err != nil {
		return fmt.Errorf("db: connection error: %w", err)
	}
	return nil
}
