package store

import (
	"database/sql"
	"reflect"
	"strconv"
	"time"

	"lineblocs.com/api/model"
)

var db *sql.DB

func CreateCall(call *model.Call) (string, error) {
	now := time.Now()
	call.StartedAt = now.Format("MM/DD/YYYY")
	call.CreatedAt = now.Format("MM/DD/YYYY")
	call.UpdatedAt = now.Format("MM/DD/YYYY")
	workspace, err := getWorkspaceFromDB(call.WorkspaceId)

	if err != nil {
		return "-1", err
	}

	stmt, err := db.Prepare("INSERT INTO calls ( `from`, `to`, `channel_id`, `status`, `direction`, `duration`, `sip_call_id`, `user_id`, `workspace_id`, `started_at`, `created_at`, `updated_at`, `api_id`, `plan_snapshot`, `notes`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, '' )")
	if err != nil {
		return "-1", err
	}
	defer stmt.Close()

	res, err := stmt.Exec(call.From, call.To, call.ChannelId, call.Status, call.Direction, "8", call.SIPCallId, call.UserId, call.WorkspaceId, now, now, now, call.APIId, workspace.Plan)

	if err != nil {
		return "-1", err
	}

	callId, err := res.LastInsertId()
	if err != nil {
		return "-1", err
	}
	return strconv.FormatInt(callId, 10), err
}

func CheckIsMakingOutboundCallFirstTime(call model.Call) {
	var id string
	row := db.QueryRow("SELECT id FROM `calls` WHERE `workspace_id` = ? AND `from` LIKE '?%s' AND `direction = 'outbound'", call.WorkspaceId, call.From, call.Direction)
	err := row.Scan(&id)
	if err != sql.ErrNoRows {
		// all ok
		return
	}
	//send notification
	user, err := getUserFromDB(call.UserId)
	if err != nil {
		panic(err)
	}
	body := `A call was made to ` + call.To + ` for the first time on your account.`
	sendEmail(user, "First call to destination country", body)
}

func sendEmail(user *model.User, subject string, body string) {
}

func getWorkspaceFromDB(id int) (*model.Workspace, error) {
	var workspaceId int
	var name string
	var creatorId int
	var outboundMacroId sql.NullInt64
	var plan string
	row := db.QueryRow(`SELECT id, name, creator_id, outbound_macro_id, plan FROM workspaces WHERE id=?`, id)

	err := row.Scan(&workspaceId, &name, &creatorId, &outboundMacroId, &plan)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil { //another error
		return nil, err
	}
	if reflect.TypeOf(outboundMacroId) == nil {
		return &model.Workspace{Id: workspaceId, Name: name, CreatorId: creatorId, Plan: plan}, nil
	}
	return &model.Workspace{Id: workspaceId, Name: name, CreatorId: creatorId, OutboundMacroId: int(outboundMacroId.Int64), Plan: plan}, nil
}
