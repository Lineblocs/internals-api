package user

import (
	"database/sql"

	helpers "github.com/Lineblocs/go-helpers"
	"lineblocs.com/api/model"
)

/*
Interface of User Store.
Implementation of User Store is located /store/user
*/
type Store interface {
	ValidateAccess(string, string) bool
	DoVerifyCaller(*model.Workspace, string) (bool, error)
	GetWorkspaceParams(int) (*[]model.WorkspaceParam, error)
	GetUserByDID(did string) (string, error)
	GetUserByTrunkSourceIp(string) (string, error)
	GetWorkspaceMacros(string) ([]model.MacroFunction, error)
	GetDIDNumberData(string) (*model.WorkspaceDIDInfo, sql.NullString, error)
	GetBYODIDNumberData(string) (*model.WorkspaceDIDInfo, sql.NullString, error)
	GetBYOPSTNProvider(string, string, int) (*model.PSTNInfo, error)
	GetBestPSTNProvider(string, string) (*model.PSTNInfo, error)
	IPWhitelistLookup(string, *model.Workspace) (bool, error)
	HostedSIPTrunkLookup(string, *model.Workspace) (bool, error)
	GetDIDAcceptOption(string) ([]byte, error)
	GetUserRoutedServer2(bool, *model.Workspace, string) (*helpers.MediaServer, error)
	GetCallerIdToUse(*model.Workspace, string) (string, error)
	GetExtensionFlowInfo(string, string) (*model.ExtensionFlowInfo, error)
	GetFlowInfo(string, string) (*model.ExtensionFlowInfo, error)
	GetCodeFlowInfo(string, string) (*model.CodeFlowInfo, error)
	IncomingDIDValidation(string) (*model.DidNumberInfo, error)
	CheckPSTNIPWhitelist(string, string) (bool, error)
	FinishValidation(string, string) (bool, error)
	IncomingBYODIDValidation(string) (*model.DidNumberInfo, error)
	CheckBYOPSTNIPWhitelist(string, string) (bool, error)
	IncomingTrunkValidation(string) ([]byte, error)
	LookupSIPTrunkByDID(string) ([]byte, error)
	IncomingMediaServerValidation(string) (bool, error)
	StoreRegistration(string, int, *model.Workspace) error
	GetSettings() (*model.Settings, error)
	ProcessSIPTrunkCall(string) ([]byte, error)
	ProcessDialplan(string) ([]byte, error)
	CaptureSIPMessage(string, string) ([]byte, error)
}
