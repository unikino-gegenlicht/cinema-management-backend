/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package registerRoutes

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/unikino-gegenlicht/cinema-management-backend/database"
	"github.com/unikino-gegenlicht/cinema-management-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func newRegister(w http.ResponseWriter, r *http.Request) {
	// parse the body of the request which should contain a register definition
	var register types.Register
	err := json.NewDecoder(r.Body).Decode(&register)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// now insert the register into the mongodb
	collection := database.Database.Collection("registers")
	res, err := collection.InsertOne(r.Context(), register)
	if err != nil {
		log.Error().Err(err).Msg("unable to create register")
		w.WriteHeader(500)
		return
	}

	// now get the object back from the mongodb and output it
	result := collection.FindOne(r.Context(), bson.D{{"_id", res.InsertedID}})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		// since the register does not exist, return a 404 error
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// now parse the result into a register
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
