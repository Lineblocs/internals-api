package extension

import "lineblocs.com/api/model"

/*
Interface of Call Store.
Implementation of Call Store is located /store/call
*/
type Store interface {
	GetExtensionByUsername(workspaceId int, username string) (*model.Extension, error)
}
