package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/sirupsen/logrus"
	"lineblocs.com/api/helpers"
	"lineblocs.com/api/utils"
)

/*
Implementation of Carrier Store
*/

type CarrierStore struct {
	db *sql.DB
}

func NewCarrierStore(db *sql.DB) *CarrierStore {
	return &CarrierStore{
		db: db,
	}
}

/*
Input: callid, status
Todo : Update sip_status of calls with matching sip_call_id
Output: If success return nil else return err
*/
func (crs *CarrierStore) CreateSIPReport(callid string, status string) error {
	stmt, err := crs.db.Prepare("UPDATE `calls` SET sip_status = ? WHERE sip_call_id = ?")
	if err != nil {
		utils.Log(logrus.ErrorLevel, "CreateSIPReport 2 Could not execute query..")
		return err
	}
	defer stmt.Close()

	statusAsInt, err := strconv.Atoi(status)
	if err != nil {
		utils.Log(logrus.ErrorLevel, "CreateSIPReport 3 error in convert...")
		return err
	}
	_, err = stmt.Exec(statusAsInt, callid)
	return err
}

/*
Input: originCode, destCode, userid
Todo : Create and Start Router Flow
Output: First value: Flow model, Second Value: error
If success return (Flow model, nil) else (nil, err)
*/
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

/*
Input: Flow model, data map
Todo : Start Processing Flow
Output: First value: RoutablePSTNProvider model, Second Value: error
If success return (RoutablePSTNProvider model, nil) else (nil, err)
*/
func (crs *CarrierStore) StartProcessingFlow(flow *helpers.Flow, data map[string]string) ([]*helpers.RoutablePSTNProvider, error) {
	providers, err := helpers.StartProcessingFlow(flow, flow.Cells[0], data, crs.db)
	return providers, err
}
