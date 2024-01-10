/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/unikino-gegenlicht/cinema-management-backend/errors"
	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	configurationTypes "github.com/unikino-gegenlicht/cinema-management-backend/types/configuration"
)
import chiMiddleware "github.com/go-chi/chi/v5/middleware"

var configuration configurationTypes.Configuration
var apiErrors map[string]errors.APIError

func main() {
	// to manage the different backend routes, create a new router
	mainRouter := chi.NewRouter()

	// now enable some middleware to help identify request and their origins
	mainRouter.Use(chiMiddleware.RealIP)
	mainRouter.Use(chiMiddleware.RequestID)

	// now create a router which is used for public-facing routes (e.g., reservation management)
	publicRouter := chi.NewRouter()
	publicRouter.HandleFunc("/*", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotImplemented)
	})

	// now create a router which is used for private-facing routes (e.g., ticket issuing, sales)
	privateRouter := chi.NewRouter()
	privateRouter.Use(middleware.OpenIDConnectJWTAuthentication(configuration.OpenIdConnect, apiErrors))
	privateRouter.HandleFunc("/*", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotImplemented)
	})

	// now mount the public and private routers
	mainRouter.Mount("/public", publicRouter)
	mainRouter.Mount("/", privateRouter)

	// now create a http server
	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: mainRouter,
	}

	// now spin up the server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal().Err(err).Msg("An errors occurred while starting the backend http server")
		}
	}()

	// now set up some signal handling to let the backend run indefinitely
	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, os.Interrupt)

	log.Info().Msg("backend started successfully")
	// now block the execution of the shutdown code
	<-cancelSignal

	log.Info().Msg("shutting down backend")
}
