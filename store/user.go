package store

import (
	"database/sql"
	"fmt"

	"github.com/ttacon/libphonenumber"
	"lineblocs.com/api/model"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) DoVerifyCaller(workspace *model.Workspace, number string) (bool, error) {
	num, err := libphonenumber.Parse(number, "US")
	if err != nil {
		return false, err
	}
	formattedNum := libphonenumber.Format(num, libphonenumber.E164)
	fmt.Printf("looking up number %s\r\n", formattedNum)
	fmt.Printf("domain isr %s\r\n", workspace.Name)
	var id string
	row := us.db.QueryRow("SELECT id FROM `did_numbers` WHERE `number` = ? AND `workspace_id` = ?", formattedNum, workspace.Id)
	err = row.Scan(&id)
	if err != sql.ErrNoRows {
		return true, nil
	}
	return false, nil
}

func (us *UserStore) GetWorkspaceParams(workspaceId int) (*[]model.WorkspaceParam, error) {
	// Execute the query
	results, err := us.db.Query("SELECT `key`, `value` FROM workspace_params WHERE `workspace_id` = ?", workspaceId)
	defer results.Close()
	params := []model.WorkspaceParam{}
	if err == sql.ErrNoRows {
		// no records setup were setup, just return empty
		return &params, nil
	}
	if err != nil {
		return nil, err
	}

	for results.Next() {
		param := model.WorkspaceParam{}
		// for each row, scan the result into our tag composite object
		err = results.Scan(&param.Key, &param.Value)
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return &params, nil
}

func (us *UserStore) GetUserByDID(did string) (string, error) {
	result := us.db.QueryRow(`SELECT
	workspaces.name
	FROM did_numbers
	INNER JOIN workspaces ON workspaces.id = did_numbers.workspace_id
	WHERE did_numbers.api_number = ?`, did)
	var domain string
	err := result.Scan(&domain)
	return domain, err
}

func (us *UserStore) GetUserByTrunkSourceIp(sourceIp string) (string, error) {
	// todo get ipv6
	sourceIpv6 := sourceIp
	result := us.db.QueryRow(`SELECT
		workspaces.name
		FROM workspaces
		INNER JOIN sip_trunks ON sip_trunks.workspace_id = workspaces.id
		INNER JOIN sip_trunks_origination_endpoints ON sip_trunks_origination_endpoints.trunk_id = sip_trunks.id
		WHERE sip_trunks_origination_endpoints.ipv4 = ?  OR sip_trunks_origination_endpoints.ipv6 = ?`, sourceIp, sourceIpv6)
	var domain string
	err := result.Scan(&domain)
	return domain, err
}

func (us *UserStore) GetWorkspaceMacros(workspace string) ([]model.MacroFunction, error) {
	// Execute the query
	results, err := us.db.Query("SELECT title, code, compiled_code FROM macro_functions WHERE `workspace_id` = ?", workspace)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	values := []model.MacroFunction{}

	for results.Next() {
		value := model.MacroFunction{}
		err = results.Scan(&value.Title, &value.Code, &value.CompiledCode)
		if err != nil {
			return nil, err
		}

		// for each row, scan the result into our tag composite object
		values = append(values, value)
	}
	return values, err
}

func (us *UserStore) GetDIDNumberData(number string) (*model.WorkspaceDIDInfo, sql.NullString, error) {
	var info model.WorkspaceDIDInfo
	var flowJson sql.NullString
	fmt.Printf("Looking up number: %s", number)
	// Execute the query
	row := us.db.QueryRow(`SELECT 
		flows.id AS flow_id,
		flows.workspace_id, 
		flows.flow_json, 
		did_numbers.number, 
		workspaces.name, 
		workspaces.name AS workspace_name, 
        workspaces.plan,
        workspaces.byo_enabled,
        workspaces.creator_id,
        workspaces.id AS workspace_id,
        workspaces.api_token,
		workspaces.api_secret 
		FROM workspaces
		INNER JOIN did_numbers ON did_numbers.workspace_id = workspaces.id	
		INNER JOIN flows ON flows.id = did_numbers.flow_id	
		INNER JOIN users ON users.id = workspaces.creator_id
		WHERE did_numbers.api_number = ?	
		`, number)
	err := row.Scan(
		&info.FlowId,
		&info.WorkspaceId,
		&flowJson,
		&info.Number,
		&info.Name,
		&info.WorkspaceName,
		&info.Plan,
		&info.BYOEnabled,
		&info.CreatorId,
		&info.WorkspaceId,
		&info.APIToken,
		&info.APISecret)
	return &info, flowJson, err
}

func (us *UserStore) GetBYODIDNumberData(number string) (*model.WorkspaceDIDInfo, sql.NullString, error) {
	var info model.WorkspaceDIDInfo
	var flowJson sql.NullString
	// Execute the query
	row := us.db.QueryRow(`SELECT 
		flows.id AS flow_id,
		flows.workspace_id, 
		flows.flow_json, 
		byo_did_numbers.number, 
		workspaces.name, 
		workspaces.name AS workspace_name, 
        workspaces.plan,
        workspaces.byo_enabled,
        workspaces.creator_id,
        workspaces.id AS workspace_id,
        workspaces.api_token,
		workspaces.api_secret FROM workspaces
		INNER JOIN byo_did_numbers ON byo_did_numbers.workspace_id = workspaces.id	
		INNER JOIN flows ON flows.id = byo_did_numbers.flow_id	
		INNER JOIN users ON users.id = workspaces.creator_id
		WHERE byo_did_numbers.number = ?	
		`, number)
	err := row.Scan(
		&info.FlowId,
		&info.WorkspaceId,
		&flowJson,
		&info.Number,

		&info.Name,
		&info.WorkspaceName,
		&info.Plan,

		&info.BYOEnabled,
		&info.CreatorId,
		&info.WorkspaceId,
		&info.APIToken,

		&info.APISecret)
	return &info, flowJson, err
}
