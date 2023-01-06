package call

import "lineblocs.com/api/model"

type Store interface {
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
