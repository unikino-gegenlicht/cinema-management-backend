/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package itemRoutes

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

func oneItem(w http.ResponseWriter, r *http.Request) {
	// extract the itemid from the url
	rawItemId := chi.URLParam(r, "itemId")
	// now convert the raw id into an object it
	itemId, err := primitive.ObjectIDFromHex(rawItemId)
	if err != nil {
		log.Warn().Msg("item id not object id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// now get the collection for the items
	collection, err := middleware.ExtractCollection(r)
	if err != nil {
		log.Error().Err(err).Msg("no collection available")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now query the collection
	result := collection.FindOne(r.Context(), bson.D{{"_id", itemId}})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("unable to query collection")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now parse the result into an item
	var item types.Item
	err = result.Decode(&item)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse query result into struct")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now return the object
	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		log.Error().Err(err).Msg("unable to return response")
	}
}
