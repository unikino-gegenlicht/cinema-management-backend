/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package registerRoutes

import "github.com/unikino-gegenlicht/cinema-management-backend/database"

// collection contains the link to the MongoDB database for the operations
// that are implemented in this package
var collection = database.Database.Collection("registers")
