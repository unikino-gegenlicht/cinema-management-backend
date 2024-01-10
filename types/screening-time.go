/*
 * Copyright (c) 2024.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "github.com/jackc/pgx/v5/pgtype"

// ScreeningTime represents the composite type screening_time in the database
// which is used to keep track of every screening time for a movie
type ScreeningTime struct {
	// ID contains the id of the screening time which helps to identify the
	// screening time a ticket is for
	ID pgtype.UUID `json:"id"`

	// At contains the timestamp showing at which time the movie starts
	At pgtype.Timestamptz `json:"at"`

	// AvailableSeats contains the number of generally available seats for
	// this screening
	AvailableSeats int `json:"availableSeats"`

	// AllowReservations indicates if reservations are allowed for this
	// screening time
	AllowReservations bool `json:"allowReservations"`
}
