/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
)
import chiMiddleware "github.com/go-chi/chi/v5/middleware"

func main() {
	// to manage the different backend routes, create a new router
	mainRouter := chi.NewRouter()

	// now enable some middleware to help identify request and their origins
	mainRouter.Use(chiMiddleware.RealIP)
	mainRouter.Use(chiMiddleware.RequestID)

	// TODO: Implement middleware for authentication using a OpenID connect
	//  userinfo endpoint

	// now mount the different sub-routers that are maintained in this file
	mainRouter.Mount("/registers", registerRouter())
	mainRouter.Mount("/transactions", transactionRouter())

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

func registerRouter() http.Handler {
	// create a new router
	r := chi.NewRouter()
	// todo: implement routes related to the registers and add them here
	return r
}

func transactionRouter() http.Handler {
	// create a new router
	r := chi.NewRouter()
	// todo: implement routes related to transactions and add them here
	return r
}
