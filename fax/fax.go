package fax

import "lineblocs.com/api/model"

/*
Interface of Fax Store.
Implementation of Fax Store is located /store/fax
*/
type Store interface {
	GetFaxCount(int) (*int, error)
	CreateFax(*model.Fax, string, int64, string, string) (int64, error)
}
