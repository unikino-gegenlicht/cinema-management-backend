/*
 * Copyright (c) 2024.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "github.com/jackc/pgx/v5/pgtype"

// Item represents a button that is available in the register view in the
// frontend available for sale.
type Item struct {
	// ID contains the uuid used to identify this item in sales and in requests
	ID pgtype.UUID `json:"id" db:"id"`

	// Name sets the name of the item that is displayed in the frontend
	Name pgtype.Text `json:"name" db:"name"`

	// Price contains the price of this Item
	Price pgtype.Numeric `json:"price" db:"price"`

	// Icon contains the name of the icon that is displayed in the frontend.
	// The icons are taken from the Iconify (https://icon-sets.iconify.design/)
	// framework.
	Icon pgtype.Text `json:"icon" db:"icon"`

	// IssueTicket indicates if a ticket is automatically issued when selling
	// this item
	IssueTicket bool `json:"issueTicket" db:"issue_ticket"`
}
