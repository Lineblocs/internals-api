package main

import (
	"fmt"
	"errors"
	"net/http"
	"database/sql"
	"encoding/json"
	lineblocs "github.com/Lineblocs/go-helpers"
	"lineblocs.com/api/helpers"
)
var db *sql.DB


func createRoutingFlow(originCode, destCode, userId *string) (*helpers.Flow, error) {
	var info helpers.FlowInfo
	var flowJson helpers.FlowVars

	// find flow by user id
	// if no flow available, use country flow
	row := db.QueryRow(`SELECT router_flows.id AS flow_id,
router_flows.flow_json
FROM workspaces_users
INNER JOIN router_flows ON router_flows.id = workspaces.flow_id
INNER JOIN workspaces ON workspaces.id = workspaces_users.workspace_id
INNER JOIN workspaces_routing_flows ON workspaces_routing_flows.workspace_id = workspaces.id
WHERE workspaces_users.user_id= ?
AND workspaces_routing_flows.dest_code= ?
`, *userId, *destCode)
	err := row.Scan(&info.FlowId,&info.FlowJSON)

	if err != sql.ErrNoRows { //lookup country flow
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal( []byte(info.FlowJSON), &flowJson )
		if err != nil {
			return nil, err
		}

		return helpers.NewFlow( info.FlowId, &flowJson	 ), nil
	}

	// lookup by country
	row = db.QueryRow(`SELECT router_flows.id AS flow_id,
router_flows.flow_json
FROM sip_countries
INNER JOIN router_flows ON router_flows.id = sip_countries.flow_id
WHERE sip_countries.country_code= ?`, *destCode)
	err = row.Scan(&info.FlowId,&info.FlowJSON)

	if err != sql.ErrNoRows { //lookup country flow
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal( []byte(info.FlowJSON), &flowJson )
		if err != nil {
			return nil, err
		}

		return helpers.NewFlow( info.FlowId, &flowJson	 ), nil
	}

	return nil, errors.New("no routing flow found...")
}


func main() {
	var w http.ResponseWriter
	var err error
	db, err = lineblocs.CreateDBConn()



	if err != nil {
		panic(err)
	}

	callto := "+17808503688"
	callfrom := "+17808503688"
	userId := "2"

	destCode, err := helpers.ParseCountryCode(callto)

	if err != nil {
		panic(err)
	}
	fmt.Println( "code is: " + destCode )

	originCode, err := helpers.ParseCountryCode(callfrom)

	if err != nil {
		panic(err)
	}
	fmt.Println( "code is: " + originCode )

	flow,err :=createRoutingFlow( &originCode, &destCode, &userId )
	if err != nil {
		panic(err)
	}

	data := make(map[string]string)
	data["origin_code"] = originCode
	data["dest_code"] = destCode
	data["from"] = callfrom
	data["to"] = callto

	providers, err := helpers.StartProcessingFlow( flow, flow.Cells[ 0 ], data, db )

	if err != nil {
		panic(err)
	}
	if len( providers ) == 0 {
		fmt.Println("No providers available..")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	provider:=providers[0]
	if len( provider.Hosts ) == 0 {
		fmt.Println("No IPs to route to..")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	host:=provider.Hosts[0]
	w.Write([]byte(host.IPAddr))
}