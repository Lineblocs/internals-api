package logger

import "lineblocs.com/api/model"

type Store interface {
	StartLogRoutine(*model.Workspace, *model.LogRoutine) (*string, error)
}
