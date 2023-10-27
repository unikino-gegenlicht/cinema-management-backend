/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID    *primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string              `bson:"name" json:"name"`
	Price float64             `bson:"price" json:"price"`
	Icon  string              `bson:"icon" json:"icon"`
	Flags *[]string           `bson:"flags,omitempty" json:"flags,omitempty"`
}
