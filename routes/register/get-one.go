/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package registerRoutes

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/unikino-gegenlicht/cinema-management-backend/database"
	"github.com/unikino-gegenlicht/cinema-management-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func oneRegister(w http.ResponseWriter, r *http.Request) {
	// extract the register id from the url
	stringRegisterId := chi.URLParam(r, "registerId")
	// convert the string into a object id
	registerId, err := primitive.ObjectIDFromHex(stringRegisterId)
	if err != nil {
	}
	collection := database.Database.Collection("registers")
	// since getting a single register is not a protected action. just query the
	// collection
	result := collection.FindOne(r.Context(), bson.D{{"_id", registerId}})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		// since the register does not exist, return a 404 error
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// now parse the result into a register
	var register types.Register
	err = result.Decode(&register)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse response from database")
		w.WriteHeader(500)
		return
	}
	// now send back the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(register)
	if err != nil {
		log.Error().Err(err).Msg("unable ro encode response")
	}
}
