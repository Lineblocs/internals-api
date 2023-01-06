package store

import (
	"database/sql"
	"fmt"
	"time"

	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

type RecordingStore struct {
	db *sql.DB
}

func NewRecordingStore(db *sql.DB) *RecordingStore {
	return &RecordingStore{
		db: db,
	}
}

func (rs *RecordingStore) CreateRecording(workspace *model.Workspace, recording *model.Recording) (int64, error) {
	now := time.Now()

	// perform a db.Query insert
	stmt, err := rs.db.Prepare("INSERT INTO recordings (`user_id`, `call_id`, `workspace_id`, `status`, `name`, `uri`, `tag`, `api_id`, `plan_snapshot`, `storage_id`, `storage_server_ip`, `trim`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(
		recording.UserId,
		recording.CallId,
		recording.WorkspaceId,
		"started",
		"",
		"",
		"",
		recording.APIId,
		workspace.Plan,
		recording.StorageId,
		recording.StorageServerIp,
		recording.Trim,
		now,
		now)
	if err != nil {
		return -1, err
	}
	recId, err := res.LastInsertId()
	if err != nil {
		return recId, err
	}
	if recording.Tags != nil {
		for _, v := range *recording.Tags {
			fmt.Printf("adding tag to recording %s\r\n", v)
			stmt, err := rs.db.Prepare("INSERT INTO recording_tags (`recording_id`, `tag`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?)")
			if err != nil {
				return recId, err
			}

			defer stmt.Close()
			res, err = stmt.Exec(recId, v, now, now)
			if err != nil {
				return recId, err
			}
		}
	}

	defer stmt.Close()
	return recId, nil
}

func (rs *RecordingStore) GetRecordingFromDB(id int) (*model.Recording, error) {
	var apiId string
	var ready int
	var size int
	var text string
	row := rs.db.QueryRow("SELECT api_id, transcription_ready, transcription_text, size FROM recordings WHERE id=?", id)

	err := row.Scan(&apiId, &ready, &text, &size)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if ready == 1 {
		return &model.Recording{APIId: apiId, Id: id, TranscriptionReady: true, TranscriptionText: text, Size: size}, nil
	}
	return &model.Recording{APIId: apiId, Id: id, Size: size}, nil
}

func (rs *RecordingStore) GetRecordingSpace(id int) (int, error) {
	var bytes int
	row := rs.db.QueryRow(`SELECT SUM(size) FROM recordings WHERE workspace_id=?`, id)

	err := row.Scan(&bytes)
	if err == sql.ErrNoRows {
		return 0, err
	}
	if err != nil { //another error
		return 0, err
	}
	return bytes, nil
}

func (rs *RecordingStore) UpdateRecording(apiId string, status string, size int64, recordingId int) error {
	now := time.Now()
	uri := utils.CreateS3URL("recordings", apiId)
	stmt, err := rs.db.Prepare("UPDATE `recordings` SET `status` = ?, `uri` = ?, `size` = ?, `updated_at` = ? WHERE `id` = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(status, uri, size, now, recordingId)
	if err != nil {
		return err
	}
	return nil
}

func (rs *RecordingStore) UpdateRecordingTranscription(update *model.RecordingTranscription) error {
	stmt, err := rs.db.Prepare("UPDATE recordings SET `transcription_ready` = ?, `transcription_text` = ? WHERE `id` = ?")
	_, err = stmt.Exec("1", update.Text, update.RecordingId)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return nil
}
