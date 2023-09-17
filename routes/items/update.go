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

func updateItem(w http.ResponseWriter, r *http.Request) {
	// extract the itemid from the url
	rawItemId := chi.URLParam(r, "itemId")
	// now convert the raw id into an object it
	itemId, err := primitive.ObjectIDFromHex(rawItemId)
	if err != nil {
		log.Warn().Msg("item id not object id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// build a filter with the id
	filter := bson.D{{"_id", itemId}}

	// now get the collection for the items
	collection, err := middleware.ExtractCollection(r)
	if err != nil {
		log.Error().Err(err).Msg("no collection available")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now query the collection to check if an item with this id exists
	err = collection.FindOne(r.Context(), filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("unable to query collection")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now read the request body into an item
	var item types.Item
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Error().Err(err).Msg("unable to read request body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// now replace the representation stored in the mongodb
	_, err = collection.ReplaceOne(r.Context(), filter, item)
	if err != nil {
		log.Error().Err(err).Msg("update failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// now use the get one handler to return the updated item
	oneItem(w, r)
	return
}
