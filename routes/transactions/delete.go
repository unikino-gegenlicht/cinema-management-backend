/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package transactionRoutes

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func deleteTransaction(w http.ResponseWriter, r *http.Request) {
	// get the transaction id which shall be deleted
	transactionId := chi.URLParam(r, "transactionId")
	// now convert that value into a document id for mongo
	documentId, err := primitive.ObjectIDFromHex(transactionId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// now get the collection for the transactions
	collection, err := middleware.ExtractCollection(r)
	if err != nil {
		log.Error().Err(err).Msg("no collection available")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now find the document related to the document id and delete it
	result := collection.FindOneAndDelete(r.Context(), bson.D{{"_id", documentId}})
	if result.Err() != nil && !errors.Is(result.Err(), mongo.ErrNoDocuments) {
		log.Error().Err(err).Msg("unable to find and delete transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// since there either is no such object or it now has been deleted, return
	// a 204 No Content as response
	w.WriteHeader(http.StatusNoContent)

	// TODO: Implement logging of transaction deletion
}
