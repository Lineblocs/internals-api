package main

import (
	"database/sql"
	"encoding/json"
		lineblocs "bitbucket.org/infinitet3ch/lineblocs-go-helpers"
		"lineblocs.com/api/helpers"
)
var db *sql.DB

func createRoutingFlow(callfrom, callto, workspaceid *string) (*helpers.Flow, error) {
	var info helpers.FlowInfo
	var flowJson helpers.FlowVars

	// find flow by user id
	// if no flow available, use country flow
	row := db.QueryRow(`SELECT router_flows.id AS flow_id,
router_flows.flow_json
FROM router_flows
WHERE router_flows.workspace_id= ?`, *workspaceid)
	err := row.Scan(&info.FlowId,&info.FlowJSON)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal( []byte(info.FlowJSON), &flowJson )
	if err != nil {
		return nil, err
	}

	return helpers.NewFlow( info.FlowId, &flowJson	 ), nil
}
func main() {
	var err error
	db, err = lineblocs.CreateDBConn()
	if err != nil {
		panic(err)
	}

	callto := ""
	callfrom := ""
	userId := ""

	flow,err :=createRoutingFlow( &callfrom, &callto, &userId )
	if err != nil {
		panic(err)
	}
	flowctx := helpers.FlowContext{ DbConn: db }
	helpers.ProcessFlow( &flowctx, flow, flow.Cells[ 0 ] )
}