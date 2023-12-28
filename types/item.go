/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	// ID contains the unique ID of an item used to identify it in the backend
	ID *primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// Name contains the displayed name of the item in the frontend
	Name string `bson:"name" json:"name"`
	// Price contains the price of the item which can be negative to allow
	// discounts
	Price float64 `bson:"price" json:"price"`
	// Icon contains the identifier of the icon that is displayed next to the
	// item's name to allow an easier identification of the item in the forntend
	Icon string `bson:"icon" json:"icon"`
	// IssueTicket is a boolean specifying if a ticket is issued on the
	// purchase of this item
	IssueTicket bool `bson:"issueTicket" json:"issueTicket"`
	// TicketCount is an unsigned integer specifying how many tickets are to
	// be issued on the purchase of the item. If it is not set it will not be
	TicketCount uint `bson:"ticketCount" json:"ticketCount"`
}
