/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LogEntry struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Time      time.Time          `bson:"at" json:"at"`
	User      string             `bson:"username" bson:"username"`
	FullName  string             `bson:"full-name" json:"fullName"`
	SubjectID string             `bson:"subject-id" json:"subject-id"`
	Action    Action             `bson:"action" json:"action"`
	Entity    json.RawMessage    `bson:"additionalInformation" json:"additionalInformation"`
}
