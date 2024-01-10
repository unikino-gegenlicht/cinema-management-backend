/*
 * Copyright (c) 2024.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "github.com/jackc/pgx/v5/pgtype"

// Register reflects a single instance of a database-stored register.
// A register is used to differentiate between sales and other transactions
type Register struct {
	// ID contains the identifier of the register in backend applications
	ID pgtype.UUID `json:"id" db:"id"`

	// Name contains the name that identifies the register in the frontend
	Name pgtype.Text `json:"name" db:"name"`

	// Description contains a more in depth description of the register that may
	// be displayed in the frontend
	Description pgtype.Text `json:"description" db:"description"`
}
