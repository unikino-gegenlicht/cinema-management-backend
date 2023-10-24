/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package transactionRoutes

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	"github.com/unikino-gegenlicht/cinema-management-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func getTransaction(w http.ResponseWriter, r *http.Request) {
	// get the transaction id from the url
	transactionId := chi.URLParam(r, "transactionId")
	// now try to parse the transaction id into an object id
	documentId, err := primitive.ObjectIDFromHex(transactionId)
	if err != nil {
		// since the object ID could not be parsed, it is invalid and therefore
		// respond with a bad request sttus
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// now get the collection containing the transactions
	collection, err := middleware.ExtractCollection(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// now try to get the document with this id
	document := collection.FindOne(r.Context(), bson.D{{"_id", documentId}})
	if errors.Is(document.Err(), mongo.ErrNoDocuments) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// now parse the result and send it back
	var transaction types.Transaction
	err = document.Decode(&transaction)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse document from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(transaction)
	if err != nil {
		log.Error().Err(err).Msg("unable to encode response into json")
	}
}
