package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"lineblocs.com/api/helpers"
)

type CarrierStore struct {
	db *sql.DB
}

func NewCarrierStore(db *sql.DB) *CarrierStore {
	return &CarrierStore{
		db: db,
	}
}

func (crs *CarrierStore) CreateSIPReport(callid string, status string) error {
	stmt, err := crs.db.Prepare("UPDATE `calls` SET sip_status = ? WHERE sip_call_id = ?")
	if err != nil {
		fmt.Printf("CreateSIPReport 2 Could not execute query..")
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	statusAsInt, err := strconv.Atoi(status)
	if err != nil {
		fmt.Printf("CreateSIPReport 3 error in convert..")
		fmt.Println(err)
		return err
	}
	_, err = stmt.Exec(statusAsInt, callid)
	return err
}

func (crs *CarrierStore) CreateRoutingFlow(originCode, destCode, userId *string) (*helpers.Flow, error) {
	var info helpers.FlowInfo
	var flowJson helpers.FlowVars

	// find flow by user id
	// if no flow available, use country flow
	row := crs.db.QueryRow(`SELECT router_flows.id AS flow_id,
router_flows.flow_json
FROM workspaces_users
INNER JOIN router_flows ON router_flows.id = workspaces.flow_id
INNER JOIN workspaces ON workspaces.id = workspaces_users.workspace_id
INNER JOIN workspaces_routing_flows ON workspaces_routing_flows.workspace_id = workspaces.id
WHERE workspaces_users.user_id= ?
AND workspaces_routing_flows.dest_code= ?
`, *userId, *destCode)
	err := row.Scan(&info.FlowId, &info.FlowJSON)

	if err != sql.ErrNoRows { //lookup country flow
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(info.FlowJSON), &flowJson)
		if err != nil {
			return nil, err
		}

		return helpers.NewFlow(info.FlowId, &flowJson), nil
	}

	// lookup by country
	row = crs.db.QueryRow(`SELECT router_flows.id AS flow_id,
router_flows.flow_json
FROM sip_countries
INNER JOIN router_flows ON router_flows.id = sip_countries.flow_id
WHERE sip_countries.country_code= ?`, *destCode)
	err = row.Scan(&info.FlowId, &info.FlowJSON)

	if err != sql.ErrNoRows { //lookup country flow
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(info.FlowJSON), &flowJson)
		if err != nil {
			return nil, err
		}

		return helpers.NewFlow(info.FlowId, &flowJson), nil
	}

	return nil, errors.New("no routing flow found...")
}

func (crs *CarrierStore) StartProcessingFlow(flow *helpers.Flow, data map[string]string) ([]*helpers.RoutablePSTNProvider, error) {
	providers, err := helpers.StartProcessingFlow(flow, flow.Cells[0], data, crs.db)
	return providers, err
}
