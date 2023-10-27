/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package transactionRoutes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.AppendMongoCollection("transactions"))
	r.Post("/", newTransaction)
	r.Post("/start", startTransaction)
	r.Get("/", getAllTransactions)
	r.Get("/by-register/{registerId}", getTransactionsForRegister)
	r.Get("/{transactionId}", getTransaction)
	r.Delete("/{transactionId}", deleteTransaction)
	return r
}
