/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package itemRoutes

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	"github.com/unikino-gegenlicht/cinema-management-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func allItems(w http.ResponseWriter, r *http.Request) {
	// get the collection
	collection, err := middleware.ExtractCollection(r)
	if err != nil {
		log.Error().Err(err).Msg("no mongo collection set by middleware")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get all entries that are in the collection
	results, err := collection.Find(r.Context(), bson.D{})
	if err != nil {
		log.Error().Err(err).Msg("unable to execute database query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now parse the results
	var items []types.Item
	err = results.All(r.Context(), &items)

	// now check that there are items that can be returned
	if len(items) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// send the available items
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		log.Error().Msg("unable to encode response")
	}

}
