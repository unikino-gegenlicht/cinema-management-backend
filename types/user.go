/*
 * Copyright (c) 2024.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "github.com/jackc/pgx/v5/pgtype"

// User reflects a single user stored in the database. The users are stored in
// the database to keep track of all users that have accessed the application
// without needing to manually look up the users from the web logs
type User struct {
	// ID contains the internally used id of the user
	ID pgtype.UUID `json:"id" db:"id"`

	// ExternalID contains the identifier used by the external user management
	// system that handles authenticating and authorizing users to access the
	// application
	ExternalID string `json:"externalID" db:"external_id"`

	// Name contains the name of the user
	Name string `json:"name" db:"name"`

	// Active displays if the user is marked as active in the application
	// allowing to block the access to the application even though the login
	// is enabled
	Active bool `json:"active" db:"active"`
}
