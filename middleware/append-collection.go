/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package middleware

import (
	"context"
	"errors"
	"github.com/unikino-gegenlicht/cinema-management-backend/database"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// AppendMongoCollection attaches the collection needed for the request to the
// context of the request
func AppendMongoCollection(collectionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handlerFunction := func(w http.ResponseWriter, r *http.Request) {
			// get the collection with the requested name from the database
			// connection
			collection := database.Database.Collection(collectionName)
			// now append it to the context of the request
			ctx := context.WithValue(r.Context(), "collection", collection)
			// now let the request continue
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(handlerFunction)
	}
}

func ExtractCollection(r *http.Request) (collection *mongo.Collection, err error) {
	// get the context
	ctx := r.Context()
	// now check if a collection has been set
	collection, collectionSet := ctx.Value("collection").(*mongo.Collection)
	if !collectionSet {
		return nil, errors.New("no collection set by middleware")
	}
	return collection, nil
}
