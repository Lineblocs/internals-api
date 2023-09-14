package call

import "lineblocs.com/api/model"

/*
Interface of Call Store.
Implementation of Call Store is located /store/call
*/
type CallStoreInterface interface {
	CreateCall(*model.Call) (string, error)
	UpdateCall(*model.CallUpdate) error
	GetCallFromDB(int) (*model.Call, error)
	SetSIPCallID(string, string) error
	SetProviderByIP(string, string) error
	CreateConference(*model.Conference) (string, error)
	CheckIsMakingOutboundCallFirstTime(call model.Call)
	GetWorkspaceFromDB(int) (*model.Workspace, error)
	GetWorkspaceByDomain(string) (*model.Workspace, error)
	GetUserFromDB(id int) (*model.User, error)
}
