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
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func getAllTransactions(w http.ResponseWriter, r *http.Request) {
	// get the collection from the request context
	collection, err := middleware.ExtractCollection(r)
	if err != nil {
		log.Error().Err(err).Msg("no collection available")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now get all records from the collection
	records, err := collection.Find(r.Context(), bson.D{})
	if err != nil {
		log.Error().Err(err).Msg("unable to pull all records from collection")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now parse all records
	var transactions []types.Transaction
	err = records.All(r.Context(), &transactions)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse transactions in database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now output the transactions
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		log.Error().Err(err).Msg("unable to output transactions")
	}
}
