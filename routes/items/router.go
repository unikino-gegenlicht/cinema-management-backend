/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package itemRoutes

import (
	"github.com/go-chi/chi/v5"
	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	"net/http"
)

func Router() http.Handler {
	// create a new router
	r := chi.NewRouter()
	r.Use(middleware.AppendMongoCollection("items"))
	r.Get("/", allItems)
	r.Post("/", newItem)
	r.Get("/{itemId}", oneItem)
	r.Patch("/{itemId}", updateItem)
	r.Delete("/{itemId}", deleteItem)
	return r
}
