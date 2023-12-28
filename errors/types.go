/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package errors

import (
	"encoding/json"
	"net/http"
)

type storedError struct {
	Code           string `json:"code"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	HttpStatusCode int    `json:"httpCode"`
}

type APIError struct {
	storedError
	HttpStatusMessage string `json:"httpStatusMessage"`
}

func (a APIError) SendError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(a.HttpStatusCode)
	json.NewEncoder(w).Encode(a)
}

func (a *APIError) UnmarshalJSON(src []byte) error {
	var apiError storedError
	err := json.Unmarshal(src, &apiError)
	if err != nil {
		return err
	}
	*a = APIError{
		storedError:       apiError,
		HttpStatusMessage: http.StatusText(apiError.HttpStatusCode),
	}
	return nil
}
