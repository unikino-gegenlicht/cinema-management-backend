/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type StatisticsResponse struct {
	KnownArticles  map[primitive.ObjectID]int `json:"knownArticles"`
	CustomArticles map[string]int             `json:"customArticles"`
}
