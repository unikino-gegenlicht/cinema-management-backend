/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package itemRoutes

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	"github.com/unikino-gegenlicht/cinema-management-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func newItem(w http.ResponseWriter, r *http.Request) {
	// parse the request body
	var item types.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Warn().Err(err).Msg("unable to decode request body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// now get the collection for the items
	collection, err := middleware.ExtractCollection(r)
	if err != nil {
		log.Error().Err(err).Msg("no mongo collection found")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now insert the item into the database
	insertResult, err := collection.InsertOne(r.Context(), item)
	if err != nil {
		log.Error().Err(err).Msg("unable to insert new register item")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now get the inserted item from the database
	res := collection.FindOne(r.Context(), bson.D{{"_id", insertResult.InsertedID}})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		log.Error().Msg("something went wrong during the insertion")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res.Err() != nil {
		log.Error().Err(res.Err()).Msg("unable to get inserted document")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now decode the response
	err = res.Decode(&item)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse mongodb document into item")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now return the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		log.Error().Err(err).Msg("unable to return response")
	}
}
