package debit

import (
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

type Store interface {
	CreateDebit(*utils.CallRate, *model.Debit) error
	CreateAPIUsageDebit(*model.Workspace, *model.DebitAPI) error
}
