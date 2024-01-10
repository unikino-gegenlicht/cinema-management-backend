/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/pelletier/go-toml/v2"
	"github.com/qustavo/dotsql"
	"github.com/rs/zerolog/log"

	"github.com/unikino-gegenlicht/cinema-management-backend/database"
	backendErrors "github.com/unikino-gegenlicht/cinema-management-backend/errors"
)
import (
	"github.com/rs/zerolog"
)

// this function initializes the backend before the main function is executed
func init() {
	// start the initialization by configuring the logger and reading the
	// desired logging level from the environment
	rawLoggingLevel, loggingLevelSet := os.LookupEnv("LOG_LEVEL")
	// now try to parse the logging level
	loggingLevel, err := zerolog.ParseLevel(rawLoggingLevel)
	if err != nil {
		log.Warn().Msg("invalid logging level supplied via the environment. defaulting to warnings only")
		loggingLevel = zerolog.WarnLevel
	}
	if !loggingLevelSet {
		// since no logging level has explicitly been set, default to warnings
		// only to generate less logging output in production environments
		loggingLevel = zerolog.WarnLevel
	}
	// now set up the logger
	zerolog.SetGlobalLevel(loggingLevel)

	// since the logger is configured, start reading the configuration file
	configurationFile, err := os.OpenFile("./config.toml", os.O_RDONLY, 0660)
	// now check if the errors indicated that there is no configuration file
	// found
	if errors.Is(err, fs.ErrNotExist) {
		// since there is no configuration file to be found, exit the backend
		// with a fatal errors
		log.Fatal().Msg("no configuration file found. please check the documentation")
	}
	// since the configuration file exists, read the configuration
	err = toml.NewDecoder(configurationFile).Decode(&configuration)
	if err != nil {
		log.Error().Err(err).Msg("unable to read configuration file")
		os.Exit(1)
	}

	// now check if the open id connect configuration should utilize the
	// automatic discovery process
	if configuration.OpenIdConnect.UseDiscovery {
		err = configuration.OpenIdConnect.Discover()
		if err != nil {
			log.Fatal().Err(err).Msg("unable to discover open id configuration")
		}
	}

	// now check the endpoints in the configuration
	if configuration.OpenIdConnect.UserinfoEndpointUri == nil {
		log.Fatal().Msg("empty userinfo endpoint in configuration")
	}
	if configuration.OpenIdConnect.JWKSEndpointUri == nil {
		log.Fatal().Msg("empty jwks endpoint in configuration")
	}

	// now connect to the postgres database
	database.Postgres, err = pgx.Connect(context.Background(), configuration.Database.ToDSN())
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to database")
	}

	// now load the tables and type queries
	tableQueries, err := dotsql.LoadFromFile("./database/tables-and-types.sql")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load table and type queries")
	}

	var queryNames []string

	for queryName, _ := range tableQueries.QueryMap() {
		queryNames = append(queryNames, queryName)
	}
	slices.Sort(queryNames)

	for _, queryName := range queryNames {
		query, err := tableQueries.Raw(queryName)
		if err != nil {
			log.Fatal().Err(err).Str("query", queryName).Msg("unable to parse query")
		}
		_, err = database.Postgres.Exec(context.Background(), query)
		if err != nil {
			log.Fatal().Err(err).Str("query", queryName).Msg("unable to execute query")
		}
	}

	// now register the composite types
	compositeTypes := []string{
		"screening_time",
		"screening_time[]",
	}
	for _, compositeType := range compositeTypes {
		dataType, err := database.Postgres.LoadType(context.Background(), compositeType)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to load composite types")
		}
		database.Postgres.TypeMap().RegisterType(dataType)
	}

	// now load the error file from the disk
	errorFile, err := os.Open("./errors.json")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open error file")
	}
	var errs []backendErrors.APIError
	err = json.NewDecoder(errorFile).Decode(&errs)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load errors")
	}
	apiErrors = make(map[string]backendErrors.APIError)
	for _, apiError := range errs {
		apiErrors[apiError.Code] = apiError
	}

	// now the init process is done
	log.Info().Msg("startup validation done")

}
