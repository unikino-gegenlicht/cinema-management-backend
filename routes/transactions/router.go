/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package transactionRoutes

import (
	"github.com/go-chi/chi/v5"
	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	"net/http"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.AppendMongoCollection("transactions"))
	r.Post("/", newTransaction)
	r.Get("/", getAllTransactions)
	r.Get("/{transactionId}", getTransaction)
	r.Delete("/{transactionId}", deleteTransaction)
	return r
}
