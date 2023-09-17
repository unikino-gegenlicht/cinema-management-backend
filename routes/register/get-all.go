/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package registerRoutes

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/unikino-gegenlicht/cinema-management-backend/database"
	"github.com/unikino-gegenlicht/cinema-management-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

// allRegisters implements the call to the mongo db that outputs all
// registers currently available
func allRegisters(w http.ResponseWriter, r *http.Request) {
	collection := database.Database.Collection("registers")
	// since getting all registers is not a protected action. just query the
	// collection
	resultCursor, err := collection.Find(r.Context(), bson.D{})
	if err != nil {
		log.Error().Err(err).Msg("unable to execute database query")
	}
	// now parse the results into an array
	var registers []types.Register
	err = resultCursor.All(r.Context(), &registers)

	// now check the number of results to catch a 204 no content response
	if len(registers) == 0 {
		// since there are no entries in the array, respond with a no content
		// response
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// now send back the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(registers)
	if err != nil {
		log.Error().Err(err).Msg("unable ro encode response")
	}
}
