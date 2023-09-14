package store

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"lineblocs.com/api/model"
)

func TestCreateCall_Failure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	callStore := NewCallStore(db)

	mock.ExpectPrepare("INSERT INTO calls").ExpectExec().
		WillReturnError(sql.ErrNoRows)

	call := &model.Call{
		From:        "from_val",
		To:          "to_val",
		ChannelId:   "channel_val",
		Status:      "status_val",
		Direction:   "direction_val",
		SIPCallId:   "sip_call_id_val",
		UserId:      1,
		WorkspaceId: 2,
		APIId:       "api_id_val",
	}

	callID, err := callStore.CreateCall(call)

	assert.Error(t, err)
	assert.Equal(t, callID, "-1")
}
