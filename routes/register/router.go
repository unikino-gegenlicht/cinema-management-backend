/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package registerRoutes

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router() http.Handler {
	// create a new router
	r := chi.NewRouter()
	r.Get("/all", _GetAll)
	r.Get("/{registerId}", _GetOne)
	r.Post("/new", _NewRegister)
	r.Patch("/{registerId}", _Update)
	return r
}
