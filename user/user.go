package user

import (
	"database/sql"

	"lineblocs.com/api/model"
)

type Store interface {
	DoVerifyCaller(*model.Workspace, string) (bool, error)
	GetWorkspaceParams(int) (*[]model.WorkspaceParam, error)
	GetUserByDID(did string) (string, error)
	GetUserByTrunkSourceIp(string) (string, error)
	GetWorkspaceMacros(string) ([]model.MacroFunction, error)
	GetDIDNumberData(string) (*model.WorkspaceDIDInfo, sql.NullString, error)
	GetBYODIDNumberData(string) (*model.WorkspaceDIDInfo, sql.NullString, error)
}
