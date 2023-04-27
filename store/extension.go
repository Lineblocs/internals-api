package store

import (
	"database/sql"
	"strconv"

	"lineblocs.com/api/model"
)

/*
Implementation of Extension Store
*/

type ExtensionStore struct {
	db *sql.DB
}

func NewExtensionStore(db *sql.DB) *ExtensionStore {
	return &ExtensionStore{
		db: db,
	}
}

/*
Input: CallRate model, Extension Model
Todo : Create new user_extension and store to db
Output: If success return nil else return err
*/
func (ds *ExtensionStore) GetExtensionByUsername(workspaceId int, username string) (*model.Extension, error) {
	row := ds.db.QueryRow("SELECT id FROM extensions WHERE username = ? AND workspace_id = ?", username, strconv.Itoa(workspaceId))
	extension := model.Extension{}
	err := row.Scan(
		&extension.Id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &extension, nil
}