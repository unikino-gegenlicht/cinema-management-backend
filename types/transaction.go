/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID          primitive.ObjectID    `bson:"_id,omitempty" json:"id,omitempty"`
	Time        time.Time             `bson:"at" json:"at"`
	Amount      float64               `bson:"amount" json:"amount"`
	Title       string                `bson:"title" json:"title"`
	Description string                `bson:"description" json:"description"`
	Items       *[]primitive.ObjectID `bson:"items,omitempty" json:"items,omitempty"`
	Register    primitive.ObjectID    `bson:"register" json:"register"`
	CustomItems *[]Item               `bson:"custom-items,omitempty" json:"customItems,omitempty"`
}
