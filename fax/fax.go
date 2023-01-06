package fax

import "lineblocs.com/api/model"

type Store interface {
	GetFaxCount(int) (*int, error)
	CreateFax(*model.Fax, string, int64, string, string) (int64, error)
}
