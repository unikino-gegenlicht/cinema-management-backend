/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package registerRoutes

import (
	"github.com/go-chi/chi/v5"
	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	"net/http"
)

func Router() http.Handler {
	// create a new router
	r := chi.NewRouter()
	r.Use(middleware.AppendMongoCollection("registers"))
	r.Get("/", allRegisters)
	r.Post("/", newRegister)
	r.Get("/{registerId}", oneRegister)
	r.Patch("/{registerId}", updateRegister)
	return r
}
