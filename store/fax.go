package store

import (
	"database/sql"
	"time"

	"lineblocs.com/api/model"
)

type FaxStore struct {
	db *sql.DB
}

func NewFaxStore(db *sql.DB) *FaxStore {
	return &FaxStore{
		db: db,
	}
}

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
