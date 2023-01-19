package store

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

/*
Implementation of Call Store
*/

type CallStore struct {
	db *sql.DB
}

func NewCallStore(db *sql.DB) *CallStore {
	return &CallStore{
		db: db,
	}
}

/*
Input: Call model
Todo : Create new call and store to db
Output: First Value: callId, Second Value:error
If success return (callid, nil) else return (nil, err)
*/
func (cs *CallStore) CreateCall(call *model.Call) (string, error) {
	now := time.Now()
	call.StartedAt = now.Format("MM/DD/YYYY")
	call.CreatedAt = now.Format("MM/DD/YYYY")
	call.UpdatedAt = now.Format("MM/DD/YYYY")
	workspace, err := cs.GetWorkspaceFromDB(call.WorkspaceId)

	if err != nil {
		return "-1", err
	}

	stmt, err := cs.db.Prepare("INSERT INTO calls ( `from`, `to`, `channel_id`, `status`, `direction`, `duration`, `sip_call_id`, `user_id`, `workspace_id`, `started_at`, `created_at`, `updated_at`, `api_id`, `plan_snapshot`, `notes`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, '' )")
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

/*
Input: CallUpdate model
Todo : Update existing call with matching id
Output: If success return nil else return err
*/
func (cs *CallStore) UpdateCall(update *model.CallUpdate) error {
	// Perform a db.Query insert
	stmt, err := cs.db.Prepare("UPDATE calls SET `status` = ?, `ended_at` = ?, `updated_at` = ? WHERE `api_id` = ?")
	if err != nil {
		utils.Log(logrus.InfoLevel, "UpdateCall 2 Could not execute query..")
		utils.Log(logrus.InfoLevel, err.Error())
		return err
	}
	defer stmt.Close()
	endedAt := time.Now()
	updatedAt := time.Now()
	_, err = stmt.Exec(update.Status, endedAt, updatedAt, update.CallId)
	if err != nil {
		return err
	}
	return nil
}

/*
Input: Call model
Todo : Check first call and send email
*/
func (cs *CallStore) CheckIsMakingOutboundCallFirstTime(call model.Call) {
	var id string
	row := cs.db.QueryRow("SELECT id FROM `calls` WHERE `workspace_id` = ? AND `from` LIKE '?%s' AND `direction = 'outbound'", call.WorkspaceId, call.From, call.Direction)
	err := row.Scan(&id)
	if err != sql.ErrNoRows {
		// All ok
		return
	}
	//Send notification
	user, err := cs.GetUserFromDB(call.UserId)
	if err != nil {
		panic(err)
	}
	body := `A call was made to ` + call.To + ` for the first time on your account.`
	sendEmail(user, "First call to destination country", body)
}

/*
Input: id
Todo : Fetch a call with call_id
Output: First Value: Call model,Second Value: error
If success return Call model else return err
*/
func (cs *CallStore) GetCallFromDB(id int) (*model.Call, error) {
	row := cs.db.QueryRow("SELECT `from`, `to`, `channel_id`, `status`, `direction`, `duration`, `user_id`, `workspace_id`, `started_at`, `created_at`, `updated_at`, `api_id`, `plan_snapshot`) FROM calls WHERE id = ?", id)
	call := model.Call{}
	err := row.Scan(
		&call.From,
		&call.To,
		&call.ChannelId,
		&call.Status,
		&call.Direction,
		&call.Duration,
		&call.UserId,
		&call.WorkspaceId,
		&call.StartedAt,
		&call.CreatedAt,
		&call.UpdatedAt,
		&call.APIId,
		&call.PlanSnapshot)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &call, nil
}

/*
Input: callid, apiid
Todo : Set sip_call_id field with matching id
Output: If success return nil else return err
*/
func (cs *CallStore) SetSIPCallID(callid string, apiid string) error {
	stmt, err := cs.db.Prepare("UPDATE calls SET sip_call_id = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(callid, apiid)

	if err != nil {
		return err
	}
	return nil
}

/*
Input: ip, apiid
Todo : Update provider_id of call table with matching ip address
Output: If success return nil else return err
*/
func (cs *CallStore) SetProviderByIP(ip string, apiid string) error {
	results, err := cs.db.Query(`SELECT sip_providers_hosts.provider_id FROM sip_providers_hosts WHERE sip_providers_hosts.ip_address = ?`, ip)
	if err != nil {
		return err
	}
	defer results.Close()
	for results.Next() {
		var providerId int
		results.Scan(&providerId)
		stmt, err := cs.db.Prepare("UPDATE calls SET provider_id = ? WHERE id = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(providerId, apiid)

		if err != nil {
			return err
		}
	}
	return nil
}

/*
Input: Conference model
Todo : Create new conference and store to db
Output: First Value: ConferenceId, Second Value: error
If success return (conferenceId, nil) else return (-1, err)
*/
func (cs *CallStore) CreateConference(conference *model.Conference) (string, error) {
	var id int
	var name string
	row := cs.db.QueryRow("SELECT id, name FROM conferences WHERE workspace_id=? AND name=?", conference.WorkspaceId, conference.Name)
	err := row.Scan(&id, &name)
	if err == sql.ErrNoRows {
		conference.APIId = utils.CreateAPIID("conf")
		// perform a db.Query insert
		now := time.Now()
		stmt, err := cs.db.Prepare("INSERT INTO conferences (`name`, `workspace_id`, `api_id`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ? )")
		if err != nil {
			return "-1", err
		}
		defer stmt.Close()
		res, err := stmt.Exec(conference.Name, conference.WorkspaceId, conference.APIId, now, now)

		if err != nil {
			return "-1", err
		}
		conferenceId, err := res.LastInsertId()
		if err != nil {
			return "-1", err
		}
		return strconv.FormatInt(conferenceId, 10), nil
	}
	return strconv.Itoa(id), nil
}

func sendEmail(user *model.User, subject string, body string) {
}

/*
Input: id
Todo : Create new conference and store to db
Output: First Value: ConferenceId, Second Value: error
If success return (conferenceId, nil) else return (-1, err)
*/
func (cs *CallStore) GetWorkspaceFromDB(id int) (*model.Workspace, error) {
	var workspaceId int
	var name string
	var creatorId int
	var outboundMacroId sql.NullInt64
	var plan string
	row := cs.db.QueryRow(`SELECT id, name, creator_id, outbound_macro_id, plan FROM workspaces WHERE id=?`, id)

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

/*
Input: domain
Todo : Get Workspace with matching domain
Output: First Value: Workspace model, Second Value: error
If success return (Workspace model, nil) else return (nil, err)
*/
func (cs *CallStore) GetWorkspaceByDomain(domain string) (*model.Workspace, error) {
	var workspaceId int
	var name string
	var byo bool
	var ipWhitelist bool
	var creatorId int
	s := strings.Split(domain, ".")
	workspaceName := s[0]
	row := cs.db.QueryRow("SELECT id, creator_id, name, byo_enabled, ip_whitelist_disabled FROM workspaces WHERE name=?", workspaceName)

	err := row.Scan(&workspaceId, &creatorId, &name, &byo, &ipWhitelist)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &model.Workspace{Id: workspaceId, CreatorId: creatorId, Name: name, BYOEnabled: byo, IPWhitelistDisabled: ipWhitelist}, nil
}

/*
Input: id
Todo : Get User with matching id
Output: First Value: User model, Second Value: error
If success return (User model, nil) else return (nil, err)
*/
func (cs *CallStore) GetUserFromDB(id int) (*model.User, error) {
	var userId int
	var username string
	var fname string
	var lname string
	var email string
	utils.Log(logrus.InfoLevel, fmt.Sprintf("looking up user %d\r\n", id))
	row := cs.db.QueryRow(`SELECT id, username, first_name, last_name, email FROM users WHERE id=?`, id)

	err := row.Scan(&userId, &username, &fname, &lname, &email)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil { // Another error
		return nil, err
	}

	return &model.User{Id: userId, Username: username, FirstName: fname, LastName: lname, Email: email}, nil
}
