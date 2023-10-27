/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package statistics

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/unikino-gegenlicht/cinema-management-backend/middleware"
	"github.com/unikino-gegenlicht/cinema-management-backend/types"
)

func getSales(w http.ResponseWriter, r *http.Request) {
	// check if the "from" and "until" parameters have been set
	rId, registerIdSet := r.URL.Query()["register"]
	fromTimeStamp, fromFilterSet := r.URL.Query()["from"]
	untilTimeStamp, untilFilterSet := r.URL.Query()["until"]

	// create a filter document
	idFilter := bson.D{}
	fromFilter := bson.D{}
	untilFilter := bson.D{}
	if registerIdSet {
		registerID, err := primitive.ObjectIDFromHex(rId[0])
		if err != nil {
			w.WriteHeader(421)
		}
		idFilter = bson.D{
			{"$or", []bson.D{
				bson.D{{"register", registerID}},
				bson.D{{"register", rId}},
			}},
		}
	}

	if fromFilterSet {
		from, err := time.Parse(time.RFC3339, fromTimeStamp[0])
		if err != nil {
			w.WriteHeader(400)
		}
		fromFilter = bson.D{{"at", bson.D{{"$gte", from}}}}

	}
	if untilFilterSet {
		from, err := time.Parse(time.RFC3339, untilTimeStamp[0])
		if err != nil {
			w.WriteHeader(400)
		}
		untilFilter = bson.D{{"at", bson.D{{"$lt", from}}}}
	}
	// now create the complete filter
	filter := bson.D{{"$and", []bson.D{
		idFilter,
		fromFilter,
		untilFilter,
	}}}

	// now extract the collection from the request context
	collection, err := middleware.ExtractCollection(r)
	if err != nil {
		log.Error().Err(err).Msg("no collection available")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now get the sale records that match the filter
	records, err := collection.Find(r.Context(), filter)
	// now parse all records
	var transactions []types.Transaction
	err = records.All(r.Context(), &transactions)

	// now go over all transactions that have been found and count the items
	// by their ids
	itemSalesCount := map[primitive.ObjectID]int{}
	customArticleSalesCount := map[string]int{}
	for _, transaction := range transactions {
		if transaction.Items == nil {
			continue
		}

		// now iterate over the items listed in the transaction
		for _, itemId := range *transaction.Items {
			itemSalesCount[itemId] += 1
		}

		// now iterate over the custom items if set
		if transaction.CustomItems == nil {
			continue
		}

		for _, item := range *transaction.CustomItems {
			statsKey := strings.TrimSpace(item.Name)
			customArticleSalesCount[statsKey] += 1
		}
	}

	statistics := types.StatisticsResponse{KnownArticles: itemSalesCount, CustomArticles: customArticleSalesCount}
	json.NewEncoder(w).Encode(statistics)
}
