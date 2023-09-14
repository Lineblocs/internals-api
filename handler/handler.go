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
)

/*
Handler Setting for all services.
You can add new service here.
*/

type Handler struct {
	adminStore     admin.AdminStoreInterface
	callStore      call.CallStoreInterface
	carrierStore   carrier.CarrierStoreInterface
	debitStore     debit.DebitStoreInterface
	faxStore       fax.FaxStoreInterface
	loggerStore    logger.LoggerStoreInterface
	recordingStore recording.RecordingStoreInterface
	userStore      user.UserStoreInterface
}

func NewHandler(as admin.AdminStoreInterface, cs call.CallStoreInterface, crs carrier.CarrierStoreInterface, ds debit.DebitStoreInterface, fs fax.FaxStoreInterface, ls logger.LoggerStoreInterface, rs recording.RecordingStoreInterface, us user.UserStoreInterface) *Handler {
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
