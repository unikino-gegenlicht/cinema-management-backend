/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package errors

import "net/http"

func SendInternalError(err error, w http.ResponseWriter) {
	e := APIError{
		storedError: storedError{
			Code:           "INTERNAL_ERROR",
			Title:          "Internal Error",
			Description:    err.Error(),
			HttpStatusCode: 500,
		},
		HttpStatusMessage: http.StatusText(500),
	}
	e.SendError(w)
}
