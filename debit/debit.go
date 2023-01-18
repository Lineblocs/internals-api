package debit

import (
	"lineblocs.com/api/model"
)

/*
Interface of Debit Store.
Implementation of Debit Store is located /store/debit
*/
type Store interface {
	CreateDebit(*model.CallRate, *model.Debit) error
	CreateAPIUsageDebit(*model.Workspace, *model.DebitAPI) error
}
