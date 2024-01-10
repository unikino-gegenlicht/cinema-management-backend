/*
 * Copyright (c) 2024.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "github.com/jackc/pgx/v5/pgtype"

// Reservation reflects a single reservation present in the system
type Reservation struct {
	// ID contains the reservation id
	ID pgtype.UUID `json:"id" db:"id"`

	// Movie contains the uuid of the movie the reservation is for
	Movie pgtype.UUID `json:"movie" db:"movie"`

	// ScreeningTime contains the uuid of the screening time for the movie
	ScreeningTime pgtype.UUID `json:"screeningTime" db:"screening_time"`

	// At contains the timestamp at which the reservation has been made
	At pgtype.Timestamptz `json:"at" db:"at"`

	// FirstName contains the first name of the person that made the reservation
	FirstName pgtype.Text `json:"firstName" db:"first_name"`

	// LastName contains the last name of the person that made the reservation
	LastName pgtype.Text `json:"lastName" db:"last_names"`

	// EmailAddress contains the email address of the person that made the
	// reservation
	EmailAddress pgtype.Text `json:"emailAddress" db:"email_address"`

	// Tickets is an integer indicating the number of tickets that have been
	// reserved
	Tickets uint `json:"tickets" db:"tickets"`
}
