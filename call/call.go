package call

import "lineblocs.com/api/model"

type Store interface {
	CreateCall(call *model.Call) (string, error)
	CheckIsMakingOutboundCallFirstTime(call model.Call)
}
