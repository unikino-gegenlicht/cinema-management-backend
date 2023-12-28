/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	// ID contains the unique id of the ticket used to identify it in
	// the backend and in other non-human-readable media types
	ID *primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	// Holder contains the name of the person that ticket has been issued to
	Holder string `bson:"holder,omitempty" json:"holder,omitempty"`

	// IssuedAt contains the date and time the ticket has been issued
	IssuedAt time.Time `bson:"issuedAt" json:"issuedAt"`
}
