/*
 * Copyright (c) 2024.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "github.com/jackc/pgx/v5/pgtype"

// Ticket represents a single ticket that has been issued for a movie
type Ticket struct {
	// ID contains the internally used id for this ticket and is hidden from
	// any JSON output
	ID int64 `json:"-" db:"id"`

	// ExternalID contains the externally used id for this ticket and is sent
	// as `id` in JSON outputs
	ExternalID pgtype.UUID `json:"id" db:"external_id"`

	// IssuedAt contains the timestamp at which the ticket has been issued
	// through the system
	IssuedAt pgtype.Timestamptz `json:"iat" db:"issued_at"`

	// Movie contains the uuid of the movie that the ticket is for
	Movie pgtype.UUID `json:"movie" db:"movie"`

	// ScreeningTime contains the uuid of the screening time the ticket has
	// been issued for
	ScreeningTime pgtype.UUID `json:"screeningTime" db:"screening_time"`
}
