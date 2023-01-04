package user

import "lineblocs.com/api/model"

type Store interface {
	getUserFromDB(id int) (*model.User, error)
}
