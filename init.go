/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/unikino-gegenlicht/cinema-management-backend/database"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)
import (
	"github.com/rs/zerolog"

	"go.mongodb.org/mongo-driver/mongo/options"
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

	// since the logger is now set up, try to receive the mongodb url from the
	// environment
	mongoDbUri, mongoDbUriSet := os.LookupEnv("MONGODB_URI")
	if !mongoDbUriSet {
		log.Fatal().Msg("the mongo db uri has not been set. please check your environment")
	}

	// now parse the mongodb uri
	mongoOptions := options.Client().ApplyURI(mongoDbUri)
	mongoOptions.SetAppName("cinema-management-backend")

	// now connect to the mongodb server
	client, err := mongo.Connect(context.TODO(), mongoOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to the mongodb")
	}
	log.Info().Msg("connected to the mongodb")

	// now ping the database to verify the connection
	log.Info().Msg("validating connection")
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("mongodb server not reachable")
	}
	log.Info().Msg("connection validated")

	// now configure the connection to the cinema_management database
	database.Database = client.Database("cinema_management")

	// and that's it in initialization for the backends
}
