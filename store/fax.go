package store

import (
	"time"

	"database/sql"
	"lineblocs.com/api/model"
	"lineblocs.com/api/database"
)

/*
Implementation of Fax Store
*/

type FaxStore struct {
	db *database.MySQLConn
}

func NewFaxStore(db *database.MySQLConn) *FaxStore {
	return &FaxStore{
		db: db,
	}
}

/*
Input: Fax model, name, size, apiId, plan
Todo : Create fax and store to db,
Output: First Value: LastInsertId, Second Value: error
If success return (id, nil) else return (nil, err)
*/
func (fs *FaxStore) CreateFax(fax *model.Fax, name string, size int64, apiId string, plan string) (int64, error) {
	now := time.Now()

	stmt, err := fs.db.Prepare("INSERT INTO faxes (`uri`, `size`, `name`, `user_id`, `call_id`, `workspace_id`, `api_id`, `plan`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(fax.Uri, size, name, fax.UserId, fax.CallId, fax.WorkspaceId, apiId, plan, now, now)
	if err != nil {
		return -1, err
	}

	faxId, err := res.LastInsertId()
	return faxId, err
}

/*
Input: id
Todo : Get FaxCount with matching workspace id
Output: First Value: count, Second Value: error
If success return (count, nil) else return (nil, err)
*/
func (fs *FaxStore) GetFaxCount(id int) (*int, error) {
	var count int
	row := fs.db.QueryRow(`SELECT COUNT(*) FROM faxes WHERE workspace_id=?`, id)

	err := row.Scan(&count)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil { //another error
		return nil, err
	}
	return &count, nil
}
