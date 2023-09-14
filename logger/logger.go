package logger

import "lineblocs.com/api/model"

/*
Interface of Logger Store.
Implementation of Logger Store is located /store/logger
*/
type LoggerStoreInterface interface {
	StartLogRoutine(*model.Workspace, *model.LogRoutine) (*string, error)
}
