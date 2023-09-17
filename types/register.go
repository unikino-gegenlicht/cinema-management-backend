/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Register struct {
	ID          *primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string              `bson:"name" json:"name"`
	Description string              `bson:"description" json:"description"`
}
