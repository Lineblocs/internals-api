package store

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"database/sql"
	helpers "github.com/Lineblocs/go-helpers"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/ttacon/libphonenumber"
	"golang.org/x/crypto/bcrypt"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
	"lineblocs.com/api/database"
)

/*
Implementation of User Store
*/

type UserStore struct {
	db *database.MySQLConn
	rdb *redis.Client
}

func NewUserStore(db *database.MySQLConn, rdb *redis.Client) *UserStore {
	return &UserStore{
		db:  db,
		rdb: rdb,
	}
}

func (us *UserStore) ValidateAccess(service_name string, api_key string) bool {
	utils.Log(logrus.DebugLevel, "Start Validate Access...")
	results, err := us.db.Query("SELECT service_name, token FROM microservice_api_keys WHERE service_name='" + service_name + "'")
	defer results.Close()
	if err != nil {
		utils.Log(logrus.DebugLevel, "ValidateAccess query error occurred...")
		return false
	}

	var serviceName string
	var token []byte
	for results.Next() {
		err := results.Scan(&serviceName, &token)
		if err != nil {
			return false
		}
		if err != nil {
			utils.Log(logrus.DebugLevel, "bcrypt generate error occurred...")
			return false
		}
		err = bcrypt.CompareHashAndPassword(token, []byte(api_key))
		if err == nil {
			return true
		}
	}
	utils.Log(logrus.DebugLevel, "token match failed...")
	return false
}

/*
Input: Workspace model, number
Todo : Check number is valid with workspace and number?
Output: First Value: valid boolean Second Value: error
If success return (valid, nil) else (false, err)
*/
func (us *UserStore) DoVerifyCaller(workspace *model.Workspace, number string) (bool, error) {
	if !utils.GetSetting().ValidateCallerId {
		return true, nil
	}

	num, err := libphonenumber.Parse(number, "US")
	if err != nil {
		return false, err
	}
	formattedNum := libphonenumber.Format(num, libphonenumber.E164)
	utils.Log(logrus.InfoLevel, fmt.Sprintf("looking up number %s\r\n", formattedNum))
	utils.Log(logrus.InfoLevel, fmt.Sprintf("domain isr %s\r\n", workspace.Name))
	var id string
	row := us.db.QueryRow("SELECT id FROM `did_numbers` WHERE `number` = ? AND `workspace_id` = ?", formattedNum, workspace.Id)
	err = row.Scan(&id)
	if err != sql.ErrNoRows {
		return true, nil
	}
	return false, nil
}

/*
Input: workspace_id
Todo : Get WorkspaceParam with matching workspace_id
Output: First Value: WorkspacePram model Second Value: error
If success return (valid, nil) else (false, err)
*/
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

/*
Input: did
Todo : Get Workspace name from did_numbers with matching did
Output: First Value: domain Second Value: error
return (domain, err)
*/
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

/*
Input: sourceIp
Todo : Get Workspace name from workspaces with matching sourceIp
Output: First Value: domain Second Value: error
return (domain, err)
*/
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

/*
Input: workspace_id
Todo : Get macro_function with matching workspace_id
Output: First Value: MacroFunction model Second Value: error
return (MacroFunction model, err)
*/
func (us *UserStore) GetWorkspaceMacros(workspaceId string) ([]model.MacroFunction, error) {
	// Execute the query
	results, err := us.db.Query("SELECT title, code, compiled_code FROM macro_functions WHERE `workspace_id` = ?", workspaceId)
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

/*
Input: number
Todo : Get WorkspaceDIDinfo with matching number
Output: First Value: WorkspaceDIDInfo model, Second Value: flowJson(sql.NullString), Third Value: error
return (WorkspaceDIDInfo model, flowJson, err)
*/
func (us *UserStore) GetDIDNumberData(number string) (*model.WorkspaceDIDInfo, sql.NullString, error) {
	var info model.WorkspaceDIDInfo
	var flowJson sql.NullString
	utils.Log(logrus.InfoLevel, fmt.Sprintf("Looking up number: %s", number))
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

/*
Input: number
Todo : Get WorkspaceDIDinfo with matching byo_did_number
Output: First Value: WorkspaceDIDInfo model, Second Value: flowJson(sql.NullString), Third Value: error
return (WorkspaceDIDInfo model, flowJson, err)
*/
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

/*
Input: from, to, workspace_id
Todo : Get PSTNInfo with matching from, to, workspace_id
Output: First Value: PSTNInfo model, Second Value: error
If success return (PSTNInfo model, nil) else return (nil, err)
*/
func (us *UserStore) GetBYOPSTNProvider(from, to string, workspaceId int) (*model.PSTNInfo, error) {
	utils.Log(logrus.InfoLevel, "Checking BYO..")
	results, err := us.db.Query(`SELECT byo_carriers.name, byo_carriers.ip_address, byo_carriers_routes.prefix, byo_carriers_routes.prepend, byo_carriers_routes.match
	FROM byo_carriers_routes
	INNER JOIN byo_carriers  ON byo_carriers.id = byo_carriers_routes.carrier_id
	INNER JOIN workspaces ON workspaces.id = byo_carriers.workspace_id
	INNER JOIN users ON users.id = workspaces.creator_id
	WHERE byo_carriers.workspace_id = ?`, workspaceId)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	for results.Next() {
		var name string
		var ip sql.NullString
		var prefix string
		var prepend string
		var match string
		err = results.Scan(&name, &ip, &prefix, &prepend, &match)
		if err != nil {
			return nil, err
		}
		if !ip.Valid {
			utils.Log(logrus.InfoLevel, "skipping 1 PSTN IP result as private IP is empty..\r\n")
			continue
		}
		valid, err := utils.CheckRouteMatches(from, to, prefix, prepend, match)
		if err != nil {
			utils.Log(logrus.InfoLevel, fmt.Sprintf("error occured when trying to match from: %s, to: %s, prefix: %s, prepend: %s, match: %s", from, to, prefix, prepend, match))
			continue
		}
		if valid {
			var number string
			number = prepend + to
			info := &model.PSTNInfo{IPAddr: ip.String, DID: number}
			return info, err
		}
	}
	return nil, err
}

/*
Input: from, to
Todo : Get PSTNInfo with matching from, to
Output: First Value: PSTNInfo model, Second Value: error
If success return (PSTNInfo model, nil) else return (nil, err)
*/
func (us *UserStore) GetBestPSTNProvider(from, to string) (*model.PSTNInfo, error) {
	// do LCR based on dial prefixes

	results, err := us.db.Query(`SELECT sip_providers.id, sip_providers.dial_prefix, sip_providers.name, sip_providers_rates.rate_ref_id, sip_providers_rates.rate
		FROM sip_providers
		INNER JOIN sip_providers_rates ON sip_providers_rates.provider_id = sip_providers.id
		WHERE (sip_providers.type_of_provider = 'outbound'
		OR sip_providers.type_of_provider = 'both')
		AND sip_providers.active = 1`)
	if err != nil {
		return nil, err
	}

	var routableProviders []*model.RoutableProvider
	var lowestRate *float64 = nil
	var lowestProviderId *int
	var lowestDialPrefix *string
	var longestMatch *int

	routableProviders = make([]*model.RoutableProvider, 0)

	defer results.Close()
	for results.Next() {
		utils.Log(logrus.InfoLevel, "Checking non BYO..")
		var id int
		var dialPrefix string
		var name string
		var rateRefId int
		var rate float64
		err = results.Scan(&id, &dialPrefix, &name, &rateRefId, &rate)
		if err != nil {
			return nil, err
		}
		utils.Log(logrus.InfoLevel, "Checking rate from provider: "+name)
		results1, err := us.db.Query(`SELECT dial_prefix
	FROM call_rates_dial_prefixes
	WHERE call_rates_dial_prefixes.call_rate_id = ?
	`, rateRefId)
		if err != nil {
			return nil, err
		}
		defer results1.Close()
		// TODO check which host is best for routing

		var rateDialPrefix string
		for results1.Next() {
			results1.Scan(&rateDialPrefix)
			utils.Log(logrus.InfoLevel, fmt.Sprintf("checking rate dial prefix %s\r\n", rateDialPrefix))
			full := rateDialPrefix + ".*"
			valid, err := regexp.MatchString(full, to)
			if err != nil {
				return nil, err
			}
			if valid {
				utils.Log(logrus.InfoLevel, "found matching route...")
				fullLen := len(full)

				if longestMatch == nil || fullLen >= *longestMatch {
					provider := model.RoutableProvider{
						Provider:   id,
						Rate:       rate,
						DialPrefix: dialPrefix}
					routableProviders = append(routableProviders, &provider)
				}
				if (longestMatch == nil || fullLen >= *longestMatch) && (lowestRate == nil || rate < *lowestRate) {
					lowestProviderId = &id
					lowestRate = &rate
					lowestDialPrefix = &dialPrefix
					longestMatch = &fullLen
				}
			}
		}
	}
	if lowestProviderId != nil {
		var number string
		number = *lowestDialPrefix + to

		// Lookup hosts
		utils.Log(logrus.InfoLevel, "Looking up hosts..\r\n")
		// Do LCR based on dial prefixes
		results1, err := us.db.Query(`SELECT sip_providers_hosts.id, sip_providers_hosts.ip_address, sip_providers_hosts.name, sip_providers_hosts.priority_prefixes
	FROM sip_providers_hosts
	WHERE sip_providers_hosts.provider_id = ?
	`, *lowestProviderId)
		if err != nil {
			return nil, err
		}
		defer results1.Close()
		// TODO check which host is best for routing
		// Add area code checking
		var info *model.PSTNInfo
		var bestProviderId *int
		var bestIpAddr *string
		for results1.Next() {
			var id int
			var ipAddr string
			var name string
			var prefixPriorities string
			results1.Scan(&id, &ipAddr, &name, &prefixPriorities)
			utils.Log(logrus.InfoLevel, fmt.Sprintf("Checking SIP host %s, IP: %s\r\n", name, ipAddr))
			prefixArr := strings.Split(prefixPriorities, ",")
			info = &model.PSTNInfo{IPAddr: ipAddr, DID: number}
			if bestProviderId == nil {
				bestProviderId = &id
				bestIpAddr = &ipAddr
			}
			// Take priority
			if len(prefixArr) != 0 {
				for _, prefix := range prefixArr {
					valid, err := regexp.MatchString(prefix, to)
					if err != nil {
						return nil, err
					}
					if valid {
						bestProviderId = &id
						bestIpAddr = &ipAddr
					}
				}
			}
		}
		info = &model.PSTNInfo{IPAddr: *bestIpAddr, DID: number}
		return info, nil
	}
	return nil, errors.New("No available routes for LCR...")
}

/*
Input: source, Workspace model
Todo : Check source ip from ip_whitelist with matching source and workspace_id
Output: First Value: match boolean, Second Value: error
If success return (true, nil) else return (false, err)
*/
func (us *UserStore) IPWhitelistLookup(source string, workspace *model.Workspace) (bool, error) {
	results, err := us.db.Query("SELECT ip, `range` FROM ip_whitelist WHERE `workspace_id` = ?", workspace.Id)
	if err != nil {
		return false, err
	}
	defer results.Close()
	if workspace.IPWhitelistDisabled {
		return true, nil
	}

	for results.Next() {
		var ip string
		var ipRange string
		err = results.Scan(&ip, &ipRange)
		if err != nil {
			return false, err
		}
		ipWithCidr := ip + ipRange
		match, err := utils.CheckCIDRMatch(source, ipWithCidr)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

/*
Input: source, Workspace model
Todo : Check source ip from ip_whitelist with matching source and workspace_id
Output: First Value: match boolean, Second Value: error
If success return (true, nil) else return (false, err)
*/
func (us *UserStore) HostedSIPTrunkLookup(source string, workspace *model.Workspace) (bool, error) {
	return false, nil
}

/*
Input: did
Todo : Get did_action with matching did number
Output: First Value: did_action, Second Value: error
If success return (did_action, nil) else return (nil, err)
*/
func (us *UserStore) GetDIDAcceptOption(did string) ([]byte, error) {
	row := us.db.QueryRow(`SELECT did_action FROM did_numbers WHERE did_numbers.api_number = ?`, did)
	var action string
	err := row.Scan(&action)
	if err == nil {
		return []byte(action), nil
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	row = us.db.QueryRow(`SELECT did_action FROM byo_did_numbers WHERE byo_did_numbers.number = ?`, did)
	err = row.Scan(&action)
	utils.Log(logrus.InfoLevel, "err check is "+err.Error())
	if err == nil {
		return []byte(action), nil
	}
	return nil, err
}

/*
Input: rtcOptimized, Workspace model, routerip
Todo : Get UserAssignedIP
Output: First Value: MediaServer model, Second Value: error
If success return (MediaServer model, nil) else return (nil, err)
*/
func (us *UserStore) GetUserRoutedServer2(rtcOptimized bool, workspace *model.Workspace, routerip string) (*helpers.MediaServer, error) {
	servers, err := createMediaServersForRouter(routerip)

	if err != nil {
		return nil, err
	}
	var result *helpers.MediaServer
	for _, server := range servers {
		// class of server
		// type of call
		// use all metrics

		//
		//if result == nil || result != nil && server.LiveCallCount < result.LiveCallCount && rtcOptimized == server.RtcOptimized {
		if result == nil || result != nil && server.LiveCPUPCTUsed < result.LiveCPUPCTUsed && rtcOptimized == server.RtcOptimized {
			result = server
		}
	}
	return result, nil

}

/*
Input: routerip
Todo : Get MediaServer list with matching routerip
Output: First Value: MediaServer model Slice, Second Value: error
If success return (MediaServer model slice, nil) else return (nil, err)
*/
func createMediaServersForRouter(routerip string) ([]*helpers.MediaServer, error) {
	var servers []*helpers.MediaServer
	db, err := helpers.CreateDBConn()
	if err != nil {
		return nil, err
	}

	results, err := db.Query(`SELECT 
		media_servers.id,
		media_servers.ip_address,
		media_servers.private_ip_address,
		media_servers.webrtc_optimized,
		media_servers.live_call_count,
		media_servers.live_cpu_pct_used,
		media_servers.live_status 
		FROM sip_routers 
		INNER JOIN sip_routers_media_servers ON sip_routers_media_servers.router_id = sip_routers.id
		INNER JOIN media_servers ON media_servers.id =  sip_routers_media_servers.server_id
		WHERE sip_routers.ip_address = ?`, routerip)
	if err != nil {
		utils.Log(logrus.InfoLevel, "query error occurred..")
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		value := helpers.MediaServer{}
		err := results.Scan(&value.Id, &value.IpAddress, &value.PrivateIpAddress, &value.RtcOptimized, &value.LiveCallCount, &value.LiveCPUPCTUsed, &value.Status)
		if err != nil {
			return nil, err
		}
		servers = append(servers, &value)
	}
	return servers, nil
}

/*
Input: Workspace model, extension
Todo : Get CallerId with mathcing workspace and extension
Output: First Value: callerId, Second Value: error
return (callerId, err)
*/
func (us *UserStore) GetCallerIdToUse(workspace *model.Workspace, extension string) (string, error) {
	var callerId string
	utils.Log(logrus.InfoLevel, fmt.Sprintf("Looking up caller ID in domain %s, ID %d, extension %s\r\n", workspace.Name, workspace.Id, extension))
	row := us.db.QueryRow("SELECT caller_id FROM extensions WHERE workspace_id=? AND username = ?", workspace.Id, extension)
	err := row.Scan(&callerId)

	return callerId, err
}

/*
Input: extension, workspace_id
Todo : Get ExtensionFlowInfo with matching workspace and extension
Output: First Value: ExtensionFlowInfo model, Second Value: error
If success return (ExtensionFlowInfo model, nil) else return (nil, err)
*/
func (us *UserStore) GetExtensionFlowInfo(workspaceId string, extension string) (*model.ExtensionFlowInfo, error) {
	var info model.ExtensionFlowInfo
	var trialStartedTime time.Time
	row := us.db.QueryRow(`SELECT flows.workspace_id,
		flows.id AS flow_id,
		flows.flow_json,
		extensions.username,
		workspaces.name,
		workspaces.name AS workspace_name,
		workspaces.plan,
		workspaces.creator_id,
		workspaces.id AS workspace_id,
		workspaces.api_token,
		workspaces.api_secret,
		users.free_trial_started
		FROM workspaces
		INNER JOIN extensions ON extensions.workspace_id = workspaces.id
		INNER JOIN flows ON flows.id = extensions.flow_id
		INNER JOIN users ON users.id = workspaces.creator_id
		WHERE extensions.username = ?
		AND extensions.workspace_id = ?`, extension, workspaceId)
	err := row.Scan(&info.FlowId, &info.WorkspaceId, &info.FlowJSON, &info.Username, &info.Name, &info.WorkspaceName, &info.Plan,
		&info.CreatorId, &info.Id, &info.APIToken, &info.APISecret, &trialStartedTime)
	if err == sql.ErrNoRows {
		return nil, err
	}
	info.FreeTrialStatus = utils.CheckFreeTrialStatus(info.Plan, trialStartedTime)
	return &info, nil
}

/*
Input: workspace_id, flow_id
Todo : Get ExtensionFlowInfo with matching flow_id and workspace_id
Output: First Value: ExtensionFlowInfo model, Second Value: error
If success return (ExtensionFlowInfo model, nil) else return (nil, err)
*/
func (us *UserStore) GetFlowInfo(workspaceId string, flowId string) (*model.ExtensionFlowInfo, error) {
	var info model.ExtensionFlowInfo
	var trialStartedTime time.Time
	row := us.db.QueryRow(`SELECT flows.workspace_id,
		flows.id AS flow_id,
		flows.flow_json,
		extensions.username,
		workspaces.name,
		workspaces.name AS workspace_name,
		workspaces.plan,
		workspaces.creator_id,
		workspaces.id AS workspace_id,
		workspaces.api_token,
		workspaces.api_secret,
		users.free_trial_started
		FROM workspaces
		INNER JOIN extensions ON extensions.workspace_id = workspaces.id
		INNER JOIN flows ON flows.id = extensions.flow_id
		INNER JOIN users ON users.id = workspaces.creator_id
		WHERE flows.public_id = ?
		AND extensions.workspace_id = ?`, flowId, workspaceId)
	err := row.Scan(&info.FlowId, &info.WorkspaceId, &info.FlowJSON, &info.Name, &info.WorkspaceName, &info.Plan,
		&info.CreatorId, &info.Id, &info.APIToken, &info.APISecret, &trialStartedTime)
	if err == sql.ErrNoRows {
		return nil, err
	}
	info.FreeTrialStatus = utils.CheckFreeTrialStatus(info.Plan, trialStartedTime)
	return &info, err
}

/*
Input: workspace_id, code
Todo : Get CodeFlowInfo with matching code and workspace_id
Output: First Value: CodeFlowInfo model, Second Value: error
If success return (CodeFlowInfo model, nil) else return (nil, err)
*/
func (us *UserStore) GetCodeFlowInfo(workspaceId string, code string) (*model.CodeFlowInfo, error) {
	var info model.CodeFlowInfo
	var trialStartedTime time.Time
	row := us.db.QueryRow(`SELECT 
		flows.workspace_id, 
		flows.flow_json, 
		extension_codes.code, 
		workspaces.name, 
		workspaces.name AS workspace_name, 
        users.plan,
        workspaces.creator_id,
		workspaces.id AS workspace_id,
		workspaces.api_token,
		workspaces.api_secret,
		users.free_trial_started
		FROM workspaces
		INNER JOIN extension_codes ON extension_codes.workspace_id = workspaces.id
		INNER JOIN flows ON flows.id = extension_codes.flow_id
		INNER JOIN users ON users.id = workspaces.creator_id
		WHERE extension_codes.code = ? 
		AND extension_codes.workspace_id = ? 
		`, code, workspaceId)
	err := row.Scan(&info.WorkspaceId,
		&info.FlowJSON,
		&info.Code,
		&info.Name,
		&info.WorkspaceName,
		&info.Plan,
		&info.CreatorId,
		&info.Id,
		&info.APIToken,
		&info.APISecret,
		&trialStartedTime)
	if err == sql.ErrNoRows {
		info.FoundCode = false
		return &info, nil
	}

	if err != nil {
		return nil, err
	}
	info.FreeTrialStatus = utils.CheckFreeTrialStatus(info.Plan, trialStartedTime)
	info.FoundCode = true
	return &info, nil
}

/*
Input: did
Todo : Get DidNumberInfo with matching did number
Output: First Value: DidNumberInfo model, Second Value: error
return (DidNumberInfo, err)
*/
func (us *UserStore) IncomingDIDValidation(did string) (*model.DidNumberInfo, error) {
	var info model.DidNumberInfo
	// Execute the query

	row := us.db.QueryRow(`SELECT did_numbers.number, did_numbers.api_number, did_numbers.workspace_id, COALESCE(0, did_numbers.trunk_id) FROM did_numbers WHERE did_numbers.api_number = ?`, did)
	err := row.Scan(&info.DidNumber,
		&info.DidApiNumber,
		&info.DidWorkspaceId,
		&info.TrunkId)
	return &info, err
}

/*
Input: did, sourceIp
Todo : Check sourceIp is matched with sip_providers_whitelist_ips
Output: First Value: match boolean, Second Value: error
If successfully matched return (true, nil), not matched return (false, nil), error return (nil, err)
*/
func (us *UserStore) CheckPSTNIPWhitelist(did string, sourceIp string) (bool, error) {
	results, err := us.db.Query(`SELECT 
	sip_providers_whitelist_ips.ip_address, 
	sip_providers_whitelist_ips.range
	FROM sip_providers_whitelist_ips
	INNER JOIN sip_providers ON sip_providers.id = sip_providers_whitelist_ips.provider_id`)
	if err != nil {
		return false, err
	}
	defer results.Close()
	for results.Next() {
		var ipAddr string
		var ipAddrRange string
		err = results.Scan(&ipAddr, &ipAddrRange)
		if err != nil {
			return false, err

		}
		ipWithCidr := ipAddr + ipAddrRange
		match, err := utils.CheckCIDRMatch(sourceIp, ipWithCidr)
		if err != nil {
			utils.Log(logrus.InfoLevel, fmt.Sprintf("error matching CIDR source %s, full %s\r\n", sourceIp, ipWithCidr))
			continue
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

/*
Input: number, didWorkspaceId
Todo : Validate number by checking number is in blocked_numbers
Output: First Value: match boolean, Second Value: error
If successfully matched return (true, nil), not matched return (false, nil), error return (nil, err)
*/
func (us *UserStore) FinishValidation(number string, didWorkspaceId string) (bool, error) {
	num, err := libphonenumber.Parse(number, "US")
	if err != nil {
		return false, err
	}
	formattedNum := libphonenumber.Format(num, libphonenumber.E164)
	row := us.db.QueryRow("SELECT id FROM `blocked_numbers` WHERE `workspace_id` = ? AND `number` = ?", didWorkspaceId, formattedNum)
	var id string
	err = row.Scan(&id)
	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, nil
}

/*
Input: did
Todo : Get DidNumberInfo with matching did number
Output: First Value: DidNumberInfo model, Second Value: error
return (DidNumberInfo model, err)
*/
func (us *UserStore) IncomingBYODIDValidation(did string) (*model.DidNumberInfo, error) {
	var info model.DidNumberInfo

	row := us.db.QueryRow(`SELECT byo_did_numbers.number, byo_did_numbers.workspace_id FROM byo_did_numbers WHERE byo_did_numbers.number = ?`, did)
	err := row.Scan(&info.DidNumber,
		&info.DidWorkspaceId)
	return &info, err
}

/*
Input: did, sourceIp
Todo : Check sourceIp is matched with byo_did_numbers
Output: First Value: match boolean, Second Value: error
If successfully matched return (true, nil), not matched return (false, nil), error return (nil, err)
*/
func (us *UserStore) CheckBYOPSTNIPWhitelist(did string, sourceIp string) (bool, error) {
	results, err := us.db.Query(`SELECT 
	byo_carriers_ips.ip,
	byo_carriers_ips.range
	FROM byo_carriers_ips
	INNER JOIN byo_carriers ON byo_carriers.id = byo_carriers_ips.carrier_id
	INNER JOIN byo_did_numbers ON byo_did_numbers.workspace_id = byo_carriers.workspace_id
	WHERE byo_did_numbers.number = ?
	`, did)
	if err != nil {
		return false, err
	}
	defer results.Close()
	for results.Next() {
		var ipAddr string
		var ipAddrRange string
		err = results.Scan(&ipAddr, &ipAddrRange)
		if err != nil {
			return false, err
		}
		ipWithCidr := ipAddr + ipAddrRange
		match, err := utils.CheckCIDRMatch(sourceIp, ipWithCidr)
		if err != nil {
			utils.Log(logrus.InfoLevel, fmt.Sprintf("error matching CIDR source %s, full %s\r\n", sourceIp, ipWithCidr))
			continue
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

/*
Input: trunkip
Todo : Looking up SIP Server and find matched one with trunkip
Output: First Value: SIP IP address, Second Value: error
If success return (ip, nil) else return (nil, err)
*/
func (us *UserStore) IncomingTrunkValidation(trunkip string) ([]byte, error) {
	results, err := us.db.Query(`SELECT 
	sip_trunks_origination_settings.recovery_sip_uri,
	sip_trunks_origination_endpoints.sip_uri
	FROM sip_trunks_origination_endpoints
	INNER JOIN sip_trunks ON sip_trunks.id = sip_trunks_origination_endpoints.trunk_id
	INNER JOIN sip_trunks_origination_settings ON sip_trunks_origination_settings.trunk_id = sip_trunks.id`)

	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		utils.Log(logrus.InfoLevel, "trying to route to SIP server..\r\n")
		var routingSIPURI string
		var recoverySIPURI string
		err := results.Scan(
			&recoverySIPURI,
			&routingSIPURI)
		if err != nil {
			return nil, err
		}
		utils.Log(logrus.InfoLevel, fmt.Sprintf("SIP routing URI = %s SIP recovery URI = %s\r\n", routingSIPURI, recoverySIPURI))
		// TODO do some health checks here to see if SIP server is actually up..
		ips, err := utils.LookupSIPAddresses(routingSIPURI)
		if err != nil {
			utils.Log(logrus.InfoLevel, fmt.Sprintf("failed to lookup SIP server %s\r\n", routingSIPURI))
			continue
		}
		for _, ip := range *ips {
			ipAddr := ip.String()
			utils.Log(logrus.InfoLevel, fmt.Sprintf("found IP = %s\r\n", ipAddr))
			utils.Log(logrus.InfoLevel, fmt.Sprintf("comparing with source IP = %s\r\n", trunkip))
			if ipAddr == trunkip {
				return []byte(ipAddr), nil
			}
		}
	}
	return nil, nil
}

/*
Input: trunkip
Todo : Looking up SIP Server and find matched one with did number
Output: First Value: SIP IP address, Second Value: error
If success return (ip, nil) else return (nil, err)
*/
func (us *UserStore) LookupSIPTrunkByDID(did string) ([]byte, error) {
	results, err := us.db.Query(`SELECT
		sip_trunks_origination_endpoints.sip_uri,
		sip_trunks_origination_settings.recovery_sip_uri
		FROM sip_trunks_origination_endpoints
		INNER JOIN did_numbers ON did_numbers.trunk_id = sip_trunks_origination_endpoints.trunk_id
		INNER JOIN sip_trunks_origination_settings  ON sip_trunks_origination_settings.trunk_id = sip_trunks_origination_endpoints.trunk_id
		WHERE did_numbers.api_number = ?`, did)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var sipRoutingUri string
	var sipRecoveryUri string
	for results.Next() {
		err := results.Scan(&sipRoutingUri, &sipRecoveryUri)
		if err != nil {
			return nil, err
		}

		utils.Log(logrus.InfoLevel, fmt.Sprintf("SIP routing URI = %s SIP recovery URI = %s\r\n", sipRoutingUri, sipRecoveryUri))
		// TODO do some health checks here to see if SIP server is actually up..
		isOnline, err := utils.CheckSIPServerHealth(sipRoutingUri)
		if err != nil {
			utils.Log(logrus.InfoLevel, fmt.Sprintf("failed to verify health of SIP server %s\r\n", sipRoutingUri))
			continue
		}
		if isOnline {
			return []byte(sipRoutingUri), nil
		}
		utils.Log(logrus.InfoLevel, fmt.Sprintf("routing server %s is offline, checking next server...\r\n", sipRoutingUri))
	}

	// no SIP servers were online try to route to recovery URI
	isOnline, err := utils.CheckSIPServerHealth(sipRecoveryUri)
	utils.Log(logrus.InfoLevel, "no SIP servers were online. routing to recovery URI\r\n")
	if isOnline {
		return []byte(sipRecoveryUri), nil
	}

	return nil, nil
}

/*
Input: source
Todo : Looking up MediaServer and find matched one with source
Output: First Value: match boolean, Second Value: error
If successfully matched return (true, nil), not matched return (false, nil), errors return (nil, err)
*/
func (us *UserStore) IncomingMediaServerValidation(source string) (bool, error) {
	results, err := us.db.Query(`SELECT media_servers.ip_address, media_servers.ip_address_range FROM media_servers`)
	// Execute the query
	if err != nil {
		return false, err
	}
	defer results.Close()
	for results.Next() {
		var ipAddr string
		var ipRange string
		err = results.Scan(&ipAddr, &ipRange)
		if err != nil {
			return false, err
		}
		full := ipAddr + ipRange
		utils.Log(logrus.InfoLevel, fmt.Sprintf("checking IP = %s", ipAddr))
		match, err := utils.CheckCIDRMatch(source, full)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

/*
Input: user, expires, Workspace model
Todo : Update extensions with domain, user, workspace
Output: If success return nil else return err
*/
func (us *UserStore) StoreRegistration(user string, expires int, workspace *model.Workspace) error {
	now := time.Now()
	stmt, err := us.db.Prepare("UPDATE extensions SET `last_registered` = ?, `register_expires`  = ? WHERE `username` = ? AND `workspace_id` = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(now, expires, user, workspace.Id)
	if err != nil {
		return err
	}
	return nil
}

/*
Input: _
Todo : Get Settings
Output: return (Settings model, err)
*/
func (us *UserStore) GetSettings() (*model.Settings, error) {
	results, err := us.db.Query("SELECT `aws_access_key_id`, `aws_secret_access_key`, `aws_region`, `s3_bucket`, `google_service_account_json`, `stripe_pub_key`, `stripe_private_key`, `stripe_test_pub_key`, `stripe_test_private_key`, `stripe_mode`, `smtp_host`, `smtp_port`, `smtp_user`, `smtp_password`, `smtp_tls` FROM api_credentials")
	defer results.Close()
	if err != nil {
		return nil, err
	}

	settings := model.Settings{}
	for results.Next() {

		err := results.Scan(&settings.AwsAccessKeyId,
			&settings.AwsSecretAccessKey,
			&settings.AwsRegion,
			&settings.S3Bucket,
			&settings.GoogleServiceAccountJson,
			&settings.StripePubKey,
			&settings.StripePrivateKey,
			&settings.StripeTestPubKey,
			&settings.StripeTestPrivateKey,
			&settings.StripeMode,
			&settings.SmtpHost,
			&settings.SmtpPort,
			&settings.SmtpUser,
			&settings.SmtpPassword,
			&settings.SmtpTls)
		if err != nil {
			return nil, err
		}
		return &settings, nil
	}
	return nil, nil
}

/*
Input: did
Todo : Get SIP URI with matching did number
Output: First Value: SIP Uri, Second Value: error
If success return (SIP Uri, nil) else return (nil, err)
*/
func (us *UserStore) ProcessSIPTrunkCall(did string) ([]byte, error) {
	// get trunk from
	results, err := us.db.Query(`SELECT 
		sip_trunks_origination_endpoints.sip_uri
		FROM did_numbers
		INNER JOIN sip_trunks ON sip_trunks.id = did_numbers.trunk_id
		INNER JOIN sip_trunks_origination_endpoints ON sip_trunks_origination_endpoints.trunk_id = sip_trunks.id
		WHERE did_numbers.api_number = ?`, did)

	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		utils.Log(logrus.InfoLevel, fmt.Sprintf("Trying to route to user trunk server..\r\n"))
		var trunkSIPURI string
		results.Scan(&trunkSIPURI)
		utils.Log(logrus.InfoLevel, fmt.Sprintf("Found SIP trunk server %s\r\n", trunkSIPURI))
		return []byte(trunkSIPURI), nil
	}
	return nil, nil
}

/*
Input: did
Todo : Get SIP URI with matching did number
Output: First Value: SIP Uri, Second Value: error
If success return (SIP Uri, nil) else return (nil, err)
*/
func (us *UserStore) ProcessDialplan(requestUser string) ([]byte, error) {
	return []byte("pstn"), nil
}

/*
Input: sip_msg
Todo : Get SIP URI with matching did number
Output: First Value: SIP Uri, Second Value: error
If success return (SIP Uri, nil) else return (nil, err)
*/
func (us *UserStore) CaptureSIPMessage(domain string, sipMsg string) ([]byte, error) {
	// get any configured webhooks for workspace
	results, err := us.db.Query(`SELECT 
		workspaces_sip_webhooks.http_uri
		FROM workspaces_sip_webhooks
		INNER JOIN workspaces ON workspaces.id = workspaces_sip_webhooks.workspace_id
		WHERE workspaces.name = ?`, domain)

	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		utils.Log(logrus.InfoLevel, fmt.Sprintf("Trying to send HTTP request to webserver..\r\n"))
		var httpURL string
		results.Scan(&httpURL)
		utils.Log(logrus.InfoLevel, fmt.Sprintf("HTTP URL =  %s\r\n", httpURL))
		data := url.Values{
			"sip_msg": {sipMsg},
		}

		_, err := http.PostForm(httpURL, data)

		if err != nil {
			utils.Log(logrus.InfoLevel, fmt.Sprintf("error occured while sending SIP message to webhook. %s..\r\n", err.Error()))
		}

	}
	return nil, nil
}

/*
Input: invite_ip
Todo : Get SIP URI with matching did number
Output: First Value: SIP Uri, Second Value: error
If success return (SIP Uri, nil) else return (nil, err)
*/
func (us *UserStore) LogCallInviteEvent(inviteIp string) error {
	// increment CPS counter for this SIP trunk
	key := fmt.Sprintf("%s_cps", inviteIp)
	err := us.rdb.Incr(key).Err()
	if err != nil {
		return err
	}
	return nil
}

/*
Input: invite_ip
Todo : Get SIP URI with matching did number
Output: First Value: SIP Uri, Second Value: error
If success return (SIP Uri, nil) else return (nil, err)
*/
func (us *UserStore) LogCallByeEvent(inviteIp string) error {
	// decrement CPS counter for this SIP trunk
	key := fmt.Sprintf("%s_cps", inviteIp)
	err := us.rdb.Decr(key).Err()
	if err != nil {
		return err
	}
	return nil
}
