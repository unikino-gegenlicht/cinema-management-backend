/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package transactionRoutes

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
)

// startTransaction creates an empty entry in the transaction log to generate
// a transaction id which can be used by clients to keep track and display the
// transactions
func startTransaction(w http.ResponseWriter, r *http.Request) {
	// now get the collection for the transactions
	collection, err := middleware.ExtractCollection(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now insert an empty transaction
	t, err := collection.InsertOne(r.Context(), bson.D{})
	// now return the inserted id to the user
	id := t.InsertedID.(primitive.ObjectID)
	w.Write([]byte(id.Hex()))
}
