/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"context"
	"errors"
	"github.com/pelletier/go-toml/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/unikino-gegenlicht/cinema-management-backend/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/fs"
	"os"
	"strings"
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
	// now check if the error indicated that there is no configuration file
	// found
	if errors.Is(err, fs.ErrNotExist) {
		// since there is no configuration file to be found, exit the backend
		// with a fatal error
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

	// now check that the mongo db uri is not empty
	if strings.TrimSpace(configuration.MongoDbUri) == "" {
		log.Fatal().Msg("empty mongodb-uri in configuration")
	}

	// since the configuration has been validated, connect to the mongodb
	mongoOptions := options.Client().ApplyURI(configuration.MongoDbUri)
	mongoOptions.SetAppName("cinema-management-backend")

	// now connect to the database
	mongoClient, err := mongo.Connect(context.TODO(), mongoOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to database")
	}

	// now check the connection by a ping
	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("database did not answer to ping")
	}

	// now get the database for the client
	database.Database = mongoClient.Database("cinema-management")

	// now connect to the redis datbase
	redisOptions, err := redis.ParseURL(configuration.RedisUri)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse redis uri")
	}

	// now create the client
	database.RedisClient = redis.NewClient(redisOptions)

	// now the init process is done
	log.Info().Msg("startup validation done")

}
