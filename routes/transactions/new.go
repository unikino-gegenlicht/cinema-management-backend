/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package transactionRoutes

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/unikino-gegenlicht/cinema-management-backend/database"
	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	"github.com/unikino-gegenlicht/cinema-management-backend/types"
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

	// to allow tracking the numver of issued tickets now get the items from the
	// database
	itemCollection := database.Database.Collection("items")
	results, err := itemCollection.Find(r.Context(), bson.D{})
	if err != nil {
		log.Error().Err(err).Msg("unable to retrieve items")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var items []types.Item
	err = results.All(r.Context(), &items)
	if err != nil {
		log.Error().Err(err).Msg("unable to retrieve items")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ticketCount uint
	for _, itemId := range transaction.Items {
		for _, item := range items {
			if item.ID == itemId {
				ticketCount += item.TicketCount
			}
		}
	}

	log.Debug().Uint("tickets", ticketCount).Msg("issuing new tickets")
	// todo: issue tickets

	// now insert the transaction
	_, err = collection.ReplaceOne(r.Context(), bson.D{{"_id", transaction.ID}}, transaction)
	if err != nil {
		log.Error().Err(err).Msg("unable to insert transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// since everything worked, send back that the transaction has been created
	w.WriteHeader(http.StatusCreated)
	return
}
