package store

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	lineblocs "github.com/Lineblocs/go-helpers"
	"github.com/ttacon/libphonenumber"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
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

func (us *UserStore) GetBYOPSTNProvider(from, to string, workspaceId int) (*model.PSTNInfo, error) {
	fmt.Println("Checking BYO..")
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
			fmt.Printf("skipping 1 PSTN IP result as private IP is empty..\r\n")
			continue
		}
		valid, err := utils.CheckRouteMatches(from, to, prefix, prepend, match)
		if err != nil {
			fmt.Printf("error occured when trying to match from: %s, to: %s, prefix: %s, prepend: %s, match: %s", from, to, prefix, prepend, match)
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
		fmt.Println("Checking non BYO..")
		var id int
		var dialPrefix string
		var name string
		var rateRefId int
		var rate float64
		err = results.Scan(&id, &dialPrefix, &name, &rateRefId, &rate)
		if err != nil {
			return nil, err
		}
		fmt.Println("Checking rate from provider: " + name)
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
			fmt.Printf("checking rate dial prefix %s\r\n", rateDialPrefix)
			full := rateDialPrefix + ".*"
			valid, err := regexp.MatchString(full, to)
			if err != nil {
				return nil, err
			}
			if valid {
				fmt.Println("found matching route...")
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

		// lookup hosts
		fmt.Printf("Looking up hosts..\r\n")
		// do LCR based on dial prefixes
		results1, err := us.db.Query(`SELECT sip_providers_hosts.id, sip_providers_hosts.ip_address, sip_providers_hosts.name, sip_providers_hosts.priority_prefixes
	FROM sip_providers_hosts
	WHERE sip_providers_hosts.provider_id = ?
	`, *lowestProviderId)
		if err != nil {
			return nil, err
		}
		defer results1.Close()
		// TODO check which host is best for routing
		// add area code checking
		var info *model.PSTNInfo
		var bestProviderId *int
		var bestIpAddr *string
		for results1.Next() {
			var id int
			var ipAddr string
			var name string
			var prefixPriorities string
			results1.Scan(&id, &ipAddr, &name, &prefixPriorities)
			fmt.Printf("Checking SIP host %s, IP: %s\r\n", name, ipAddr)
			prefixArr := strings.Split(prefixPriorities, ",")
			info = &model.PSTNInfo{IPAddr: ipAddr, DID: number}
			if bestProviderId == nil {
				bestProviderId = &id
				bestIpAddr = &ipAddr
			}
			// take priority
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
	return nil, errors.New("no available routes for LCR...")
}

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
	fmt.Println("err check is ", err)
	if err == nil {
		return []byte(action), nil
	}
	return nil, err
}

func (us *UserStore) GetUserRoutedServer2(rtcOptimized bool, workspace *model.Workspace, routerip string) (*lineblocs.MediaServer, error) {
	servers, err := createMediaServersForRouter(routerip)

	if err != nil {
		return nil, err
	}
	var result *lineblocs.MediaServer
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

func createMediaServersForRouter(routerip string) ([]*lineblocs.MediaServer, error) {
	var servers []*lineblocs.MediaServer
	db, err := lineblocs.CreateDBConn()
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
		fmt.Println("query error occurred..")
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		value := lineblocs.MediaServer{}
		err := results.Scan(&value.Id, &value.IpAddress, &value.PrivateIpAddress, &value.RtcOptimized, &value.LiveCallCount, &value.LiveCPUPCTUsed, &value.Status)
		if err != nil {
			return nil, err
		}
		servers = append(servers, &value)
	}
	return servers, nil
}

func (us *UserStore) GetCallerIdToUse(workspace *model.Workspace, extension string) (string, error) {
	var callerId string
	fmt.Printf("Looking up caller ID in domain %s, ID %d, extension %s\r\n", workspace.Name, workspace.Id, extension)
	row := us.db.QueryRow("SELECT caller_id FROM extensions WHERE workspace_id=? AND username = ?", workspace.Id, extension)
	err := row.Scan(&callerId)

	return callerId, err
}

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
	return &info, err
}

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

func (us *UserStore) IncomingDIDValidation(did string) (*model.DidNumberInfo, error) {
	var info model.DidNumberInfo
	// Execute the query
	row := us.db.QueryRow(`SELECT did_numbers.number, did_numbers.api_number, did_numbers.workspace_id, did_numbers.trunk_id FROM did_numbers WHERE did_numbers.api_number = ?`, did)
	err := row.Scan(&info.DidNumber,
		&info.DidApiNumber,
		&info.DidWorkspaceId,
		&info.TrunkId)
	return &info, err
}

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
			fmt.Printf("error matching CIDR source %s, full %s\r\n", sourceIp, ipWithCidr)
			continue
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

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

func (us *UserStore) IncomingBYODIDValidation(did string) (*model.DidNumberInfo, error) {
	var info model.DidNumberInfo

	row := us.db.QueryRow(`SELECT byo_did_numbers.number, byo_did_numbers.workspace_id FROM byo_did_numbers WHERE byo_did_numbers.number = ?`, did)
	err := row.Scan(&info.DidNumber,
		&info.DidWorkspaceId)
	return &info, err
}

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
			fmt.Printf("error matching CIDR source %s, full %s\r\n", sourceIp, ipWithCidr)
			continue
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

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
		fmt.Printf("trying to route to SIP server..\r\n")
		var routingSIPURI string
		var recoverySIPURI string
		err := results.Scan(
			&recoverySIPURI,
			&routingSIPURI)
		if err != nil {
			return nil, err
		}
		fmt.Printf("SIP routing URI = %s SIP recovery URI = %s\r\n", routingSIPURI, recoverySIPURI)
		// TODO do some health checks here to see if SIP server is actually up..
		ips, err := utils.LookupSIPAddresses(routingSIPURI)
		if err != nil {
			fmt.Printf("failed to lookup SIP server %s\r\n", routingSIPURI)
			continue
		}
		for _, ip := range *ips {
			ipAddr := ip.String()
			fmt.Printf("found IP = %s\r\n", ipAddr)
			fmt.Printf("comparing with source IP = %s\r\n", trunkip)
			if ipAddr == trunkip {
				return []byte(ipAddr), nil
			}
		}
	}
	return nil, nil
}

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

		fmt.Printf("SIP routing URI = %s SIP recovery URI = %s\r\n", sipRoutingUri, sipRecoveryUri)
		// TODO do some health checks here to see if SIP server is actually up..
		isOnline, err := utils.CheckSIPServerHealth(sipRoutingUri)
		if err != nil {
			fmt.Printf("failed to verify health of SIP server %s\r\n", sipRoutingUri)
			continue
		}
		if isOnline {
			return []byte(sipRoutingUri), nil
		}
		fmt.Printf("routing server %s is offline, checking next server...\r\n", sipRoutingUri)
	}

	// no SIP servers were online try to route to recovery URI
	isOnline, err := utils.CheckSIPServerHealth(sipRecoveryUri)
	fmt.Printf("no SIP servers were online. routing to recovery URI\r\n")
	if isOnline {
		return []byte(sipRecoveryUri), nil
	}

	return nil, nil
}

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
		fmt.Printf("checking IP = %s", ipAddr)
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

func (us *UserStore) GetSettings() (*model.Settings, error) {
	results, err := us.db.Query("SELECT `aws_access_key_id`, `aws_secret_access_key`, `aws_region`, `google_service_account_json`, `stripe_pub_key`, `stripe_private_key`, `stripe_test_pub_key`, `stripe_test_private_key`, `stripe_mode`, `smtp_host`, `smtp_port`, `smtp_user`, `smtp_password`, `smtp_tls` FROM api_credentials")
	defer results.Close()
	if err != nil {
		return nil, err
	}

	settings := model.Settings{}
	for results.Next() {

		err := results.Scan(&settings.AwsAccessKeyId,
			&settings.AwsSecretAccessKey,
			&settings.AwsRegion,
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
		fmt.Printf("trying to route to user trunk server..\r\n")
		var trunkSIPURI string
		results.Scan(&trunkSIPURI)
		fmt.Printf("found SIP trunk server %s\r\n", trunkSIPURI)
		return []byte(trunkSIPURI), nil
	}
	return nil, nil
}