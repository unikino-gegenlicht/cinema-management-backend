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

	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	itemRoutes "github.com/unikino-gegenlicht/cinema-management-backend/routes/items"
	registerRoutes "github.com/unikino-gegenlicht/cinema-management-backend/routes/register"
	"github.com/unikino-gegenlicht/cinema-management-backend/routes/statistics"
	"github.com/unikino-gegenlicht/cinema-management-backend/routes/transactions"
	configurationTypes "github.com/unikino-gegenlicht/cinema-management-backend/types/configuration"
)
import chiMiddleware "github.com/go-chi/chi/v5/middleware"

var configuration configurationTypes.Configuration

func main() {
	// to manage the different backend routes, create a new router
	mainRouter := chi.NewRouter()

	// now enable some middleware to help identify request and their origins
	mainRouter.Use(chiMiddleware.RealIP)
	mainRouter.Use(chiMiddleware.RequestID)

	mainRouter.Use(middleware.OpenIDConnectJWTAuthentication(configuration.OpenIdConnect))
	// now mount the different sub-routers that are maintained in this file
	mainRouter.Mount("/registers", registerRoutes.Router())
	mainRouter.Mount("/items", itemRoutes.Router())
	mainRouter.Mount("/transactions", transactionRoutes.Router())
	mainRouter.Mount("/statistics", statistics.Router())

	// now create a http server
	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: mainRouter,
	}

	// now spin up the server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal().Err(err).Msg("An error occurred while starting the backend http server")
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
