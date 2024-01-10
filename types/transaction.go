/*
 * Copyright (c) 2024.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "github.com/jackc/pgx/v5/pgtype"

// Transaction reflects a single stored transaction in the database
type Transaction struct {
	// ID contains the transaction id that has been assigned by the database
	// to this instance
	ID pgtype.UUID `json:"id" db:"id"`

	// At contains the timestamp that indicates at which point in time the
	// transaction happened
	At pgtype.Timestamptz `json:"at" db:"at"`

	// Amount contains the amount of money the transaction is about
	Amount pgtype.Numeric `json:"amount" db:"amount"`

	// By contains the uuid of the user responsible for the transaction
	By pgtype.UUID `json:"by" db:"by"`

	// Title contains the title of the transaction
	Title pgtype.Text `json:"title" db:"tile"`

	// Description contains a more-indepth description of the transaction
	Description pgtype.Text `json:"description" db:"description"`

	// Register contains the uuid of the register the transaction is associated
	// with
	Register pgtype.UUID `json:"register" db:"register"`

	// SaleID contains the id of a sale if the transaction has been the source
	// of a transaction
	SaleID pgtype.UUID `json:"saleID" db:"sale_id"`
}
