/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import (
	"encoding/json"
	"errors"
)

type Action uint16

const (
	CreateRegister Action = iota
	UpdateRegister
	CreateItem
	UpdateItem
	DeleteItem
	CreateTransaction
	EditTransaction
	DeleteTransaction
)

const (
	createRegisterString    = "createRegister"
	updateRegisterString    = "updateRegister"
	createItemString        = "createItem"
	updateItemString        = "updateItem"
	deleteItemString        = "deleteItem"
	createTransactionString = "createTransaction"
	editTransactionString   = "editTransaction"
	deleteTransactionString = "deleteTransaction"
)

func (a Action) String() string {
	switch a {
	case CreateRegister:
		return createRegisterString
	case UpdateRegister:
		return updateRegisterString
	case CreateItem:
		return createItemString
	case UpdateItem:
		return updateItemString
	case DeleteItem:
		return deleteItemString
	case CreateTransaction:
		return createTransactionString
	case EditTransaction:
		return editTransactionString
	case DeleteTransaction:
		return deleteTransactionString
	default:
		return ""
	}
}

func (a Action) MarshalJSON() ([]byte, error) {
	s := a.String()
	if s == "" {
		return nil, errors.New("unknown action")
	}
	return json.Marshal(s)
}
