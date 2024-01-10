/*
 * Copyright (c) 2024.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "github.com/jackc/pgx/v5/pgtype"

// Event reflects a single event that happened in the application. These events
// are tracked automatically.
type Event struct {
	// ID is the event id that is automatically set by the database
	ID int `json:"id" db:"id"`

	// User contains the uuid of the user responsible for the event
	User pgtype.UUID `json:"user" db:"user_id"`

	// At contains the timestamp showing at which time the event happened
	At pgtype.Timestamp `json:"at" db:"at"`

	// Event contains the actual event that happened
	Event pgtype.Text `json:"event" db:"event"`
}
