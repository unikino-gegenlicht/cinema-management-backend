/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package registerRoutes

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/unikino-gegenlicht/cinema-management-backend/database"
	"github.com/unikino-gegenlicht/cinema-management-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func _Update(w http.ResponseWriter, r *http.Request) {
	// extract the register id from the url
	stringRegisterId := chi.URLParam(r, "registerId")
	// convert the string into a object id
	registerId, err := primitive.ObjectIDFromHex(stringRegisterId)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	collection := database.Database.Collection("registers")
	// since getting a single register is not a protected action. just query the
	// collection
	err = collection.FindOne(r.Context(), bson.D{{"_id", registerId}}).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		// since the register does not exist, return a 404 error
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		return
	}
	// now parse the request body
	var newRegisterRepresentation types.Register
	err = json.NewDecoder(r.Body).Decode(&newRegisterRepresentation)

	// now update the document in the mongodb
	_, err = collection.ReplaceOne(r.Context(), bson.D{{"_id", registerId}}, newRegisterRepresentation)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	// now just get the single register
	_GetOne(w, r)
	return
}
