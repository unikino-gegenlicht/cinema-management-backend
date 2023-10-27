/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package statistics

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.AppendMongoCollection("transactions"))
	r.Get("/sales", getTotalSales)
	r.Get("/sales/{registerId}", getRegisterSales)
	return r
}
