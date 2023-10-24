/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package transactionRoutes

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	"github.com/unikino-gegenlicht/cinema-management-backend/types"
	"net/http"
)

func newTransaction(w http.ResponseWriter, r *http.Request) {
	// parse the request body which should contain the transaction data
	var transaction types.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// now get the collection for the transactions
	collection, err := middleware.ExtractCollection(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now insert the transaction
	_, err = collection.InsertOne(r.Context(), transaction)
	if err != nil {
		log.Error().Err(err).Msg("unable to insert transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// since everything worked, send back that the transaction has been created
	w.WriteHeader(http.StatusCreated)
	return
}
