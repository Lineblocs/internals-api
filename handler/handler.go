package handler

import (
	"lineblocs.com/api/admin"
	"lineblocs.com/api/call"
	"lineblocs.com/api/carrier"
	"lineblocs.com/api/debit"
	"lineblocs.com/api/fax"
	"lineblocs.com/api/logger"
	"lineblocs.com/api/recording"
	"lineblocs.com/api/user"
	"lineblocs.com/api/extension"
)

/*
Handler Setting for all services.
You can add new service here.
*/

type Handler struct {
	adminStore     admin.Store
	callStore      call.Store
	carrierStore   carrier.Store
	debitStore     debit.Store
	faxStore       fax.Store
	loggerStore    logger.Store
	recordingStore recording.Store
	userStore      user.Store
	extensionStore      extension.Store
}

func NewHandler(as admin.Store, cs call.Store, crs carrier.Store, ds debit.Store, fs fax.Store, ls logger.Store, rs recording.Store, us user.Store, es extension.Store) *Handler {
	return &Handler{
		adminStore:     as,
		callStore:      cs,
		carrierStore:   crs,
		debitStore:     ds,
		faxStore:       fs,
		loggerStore:    ls,
		recordingStore: rs,
		userStore:      us,
	}
}
