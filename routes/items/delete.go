/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package itemRoutes

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func deleteItem(w http.ResponseWriter, r *http.Request) {
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

	// now delete all objects matching the item id
	_, err = collection.DeleteMany(r.Context(), filter)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete objects")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now respond that there is no content
	w.WriteHeader(http.StatusNoContent)
	return
}
