package main

import (
	"context"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"errors"

	//"errors"
	"database/sql"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"reflect"
	"regexp"

	lineblocs "github.com/Lineblocs/go-helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "github.com/go-sql-driver/mysql"
	guuid "github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mailgun/mailgun-go/v4"
	limiter "github.com/nadirhamid/golang-etcd-rate-limiter"
	libphonenumber "github.com/ttacon/libphonenumber"
	"go.etcd.io/etcd/clientv3"
	"lineblocs.com/api/helpers"
)

type Call struct {
	From         string `json:"from"`
	To           string `json:"to"`
	Status       string `json:"status"`
	Direction    string `json:"direction"`
	Duration     string `json:"duration"`
	UserId       int    `json:"user_id"`
	WorkspaceId  int    `json:"workspace_id"`
	APIId        string `json:"api_id"`
	SourceIp     string `json:"source_ip"`
	ChannelId    string `json:"channel_id"`
	SIPCallId    string `json:"call_id"`
	StartedAt    string `json:"started_at"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	PlanSnapshot string `json:"plan_snapshot"`
}
type CallUpdateReq struct {
	CallId   int    `json:"call_id"`
	Status   string `json:"status"`
	SourceIp string `json:"source_ip"`
}
type RecordingTranscriptionReq struct {
	RecordingId int    `json:"recording_id"`
	Ready       bool   `json:"ready"`
	Text        string `json:"text"`
}
type Conference struct {
	Name        string `json:"name"`
	WorkspaceId int    `json:"workspace_id"`
	APIId       string `json:"api_id"`
}
type DebitCreateReq struct {
	UserId      int `json:"user_id"`
	WorkspaceId int `json:"workspace_id"`
	ModuleId    int `json:"module_id"`

	Source  string  `json:"source"`
	Number  string  `json:"number"`
	Type    string  `json:"type"`
	Seconds float64 `json:"seconds"`
}

type CallRate struct {
	CallRate float64
}

type DebitAPIParams struct {
	Length          int     `json:"length"`
	RecordingLength float64 `json:"recording_length"`
}
type DebitAPICreateReq struct {
	UserId      int            `json:"user_id"`
	WorkspaceId int            `json:"workspace_id"`
	Type        string         `json:"type"`
	Source      string         `json:"source"`
	Params      DebitAPIParams `json:"params"`
}

type LogCreateReq struct {
	UserId      int     `json:"user_id"`
	WorkspaceId int     `json:"workspace_id"`
	Title       string  `json:"title"`
	Report      string  `json:"report"`
	FlowId      int     `json:"flow_id"`
	Level       *string `json:"report"`
	From        *string `json:"from"`
	To          *string `json:"to"`
}
type LogSimpleCreateReq struct {
	Type  string  `json:"type"`
	Level *string `json:"level"`
}
type Fax struct {
	UserId      int    `json:"user_id"`
	WorkspaceId int    `json:"workspace_id"`
	CallId      int    `json:"call_id"`
	Uri         string `json:"uri"`
	APIId       string `json:"api_id"`
}

type Recording struct {
	Id                 int       `json:"id"`
	UserId             int       `json:"user_id"`
	CallId             *int      `json:"call_id"`
	Size               int       `json:"size"`
	WorkspaceId        int       `json:"workspace_id"`
	APIId              string    `json:"api_id"`
	Tags               *[]string `json:"tags"`
	Trim bool      `json:"trim"`
	TranscriptionReady bool      `json:"transcription_ready"`
	TranscriptionText  string    `json:"transcription_text"`
	StorageId          string    `json:"storage_id"`
	StorageServerIp    string    `json:"storage_server_ip"`
}

type VerifyNumber struct {
	Valid bool `json:"valid"`
}

type LogRoutine struct {
	UserId      int
	WorkspaceId int
	Title       string
	Report      string
	FlowId      int
	Level       string
	From        string
	To          string
}
type User struct {
	Id        int
	Username  string
	FirstName string
	LastName  string
	Email     string
}

type Workspace struct {
	Id                  int    `json:"id"`
	CreatorId           int    `json:"creator_id"`
	Name                string `json:"name"`
	BYOEnabled          bool   `json:"byo_enabled"`
	IPWhitelistDisabled bool   `json:"ip_whitelist_disabled"`
	OutboundMacroId     int    `json:"outbound_macro_id"`
	Region              string `json:"region"`
	Plan                string `json:"plan"`
}

type WorkspaceParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type WorkspaceCreatorFullInfo struct {
	Id              int               `json:"id"`
	Workspace       *Workspace        `json:"workspace"`
	WorkspaceName   string            `json:"workspace_name"`
	WorkspaceDomain string            `json:"workspace_domain"`
	WorkspaceId     int               `json:"workspace_id"`
	WorkspaceParams *[]WorkspaceParam `json:"workspace_params"`
	OutboundMacroId int               `json:"outbound_macro_id"`
}
type WorkspaceDIDInfo struct {
	WorkspaceId         int               `json:"workspace_id"`
	Number              string            `json:"number"`
	FlowId              int               `json:"flow_id"`
	FlowJSON            string            `json:"flow_json"`
	WorkspaceName       string            `json:"workspace_name"`
	Name                string            `json:"name"`
	Plan                string            `json:"plan"`
	BYOEnabled          bool              `json:"byo_enabled"`
	IPWhitelistDisabled bool              `json:"ip_whitelist_disabled"`
	OutboundMacroId     int               `json:"outbound_macro_id"`
	CreatorId           int               `json:"creator_id"`
	APIToken            string            `json:"api_token"`
	APISecret           string            `json:"api_secret"`
	WorkspaceParams     *[]WorkspaceParam `json:"workspace_params"`
}
type PSTNInfo struct {
	IPAddr string `json:"ip_addr"`
	DID    string `json:"did"`
}
type CallerIDInfo struct {
	CallerID string `json:"caller_id"`
}

type RoutableProvider struct {
	Rate float64 `json:"rate"`
	DialPrefix string `json:"dial_prefix"`
	Provider int `json:"provider"`
	IPAddress int `json:"ip_address"`
}

type SIPTrunkInfo struct {
	Domain        string            `json:"domain"`
	WorkspaceId     int               `json:"workspace_id"`
	WorkspaceName   string            `json:"workspace_name"`
	CreatorId       int               `json:"creator_id"`
}


type ExtensionFlowInfo struct {
	FlowId          int               `json:"flow_id"`
	CallerID        string            `json:"caller_id"`
	WorkspaceId     int               `json:"workspace_id"`
	FlowJSON        string            `json:"flow_json"`
	Username        string            `json:"username"`
	Name            string            `json:"name"`
	WorkspaceName   string            `json:"workspace_name"`
	Plan            string            `json:"plan"`
	CreatorId       int               `json:"creator_id"`
	Id              int               `json:"id"`
	APIToken        string            `json:"api_token"`
	APISecret       string            `json:"api_secret"`
	WorkspaceParams *[]WorkspaceParam `json:"workspace_params"`
	FreeTrialStatus string            `json:"workspace_params"`
}

type CodeFlowInfo struct {
	WorkspaceId     int    `json:"workspace_id"`
	Code            string `json:"code"`
	FlowJSON        string `json:"flow_json"`
	Name            string `json:"name"`
	WorkspaceName   string `json:"workspace_name"`
	Plan            string `json:"plan"`
	CreatorId       int    `json:"creator_id"`
	Id              int    `json:"id"`
	APIToken        string `json:"api_token"`
	APISecret       string `json:"api_secret"`
	FreeTrialStatus string `json:"workspace_params"`
	FoundCode       bool   `json:"found_code"`
}

type MacroFunction struct {
	Title        string `json:"title"`
	Code         string `json:"code"`
	CompiledCode string `json:"compiled_code"`
}

type EmailInfo struct {
	Message string `json:"message"`
}

type Settings struct {
	AwsAccessKeyId           string `json:"aws_access_key_id"`
	AwsSecretAccessKey       string `json:"aws_secret_access_key"`
	AwsRegion                string `json:"aws_region"`
	StripePubKey             string `json:"stripe_pub_key"`
	StripePrivateKey         string `json:"stripe_private_key"`
	StripeTestPubKey         string `json:"stripe_test_pub_key"`
	StripeTestPrivateKey     string `json:"stripe_test_private_key"`
	StripeMode               string `json:"stripe_mode"`
	SmtpHost                 string `json:"smtp_host"`
	SmtpPort                 string `json:"smtp_port"`
	SmtpUser                 string `json:"smtp_user"`
	SmtpPassword             string `json:"smtp_password"`
	SmtpTls                  string `json:"smtp_tls"`
	GoogleServiceAccountJson string `json:"google_service_account_json"`
}
type GlobalSettings struct {
	ValidateCallerId bool
}

type ServerData struct {
	mu      sync.RWMutex
	servers []*lineblocs.MediaServer
}

var db *sql.DB
var settings *GlobalSettings
var data *ServerData

func createETCDClient() (*clientv3.Client, error) {
	dialTimeout := time.Duration(time.Second * 10)
	etcdAddr := os.Getenv("ETCD_ENDPOINT")
	user := os.Getenv("ETCD_USERNAME")
	pass := os.Getenv("ETCD_PASSWORD")
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{etcdAddr},
		Username:    user,
		Password:    pass,
	})
	return cli, err
}

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
func createAPIID(prefix string) string {
	id := guuid.New()
	return prefix + "-" + id.String()
}
func lookupBestCallRate(number string, typeRate string) *CallRate {
	return &CallRate{CallRate: 9.99}
}
func handleInternalErr(msg string, err error, w http.ResponseWriter) {
	fmt.Printf(msg)
	fmt.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
}
func calculateTTSCosts(length int) float64 {
	var result float64 = float64(length) * .000005
	return result
}
func calculateSTTCosts(recordingLength float64) float64 {
	// Google cloud bills .006 per 15 seconds
	billable := recordingLength / 15
	var result float64 = 0.006 * billable
	return result
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
func getUserFromDB(id int) (*User, error) {
	var userId int
	var username string
	var fname string
	var lname string
	var email string
	fmt.Printf("looking up user %d\r\n", id)
	row := db.QueryRow(`SELECT id, username, first_name, last_name, email FROM users WHERE id=?`, id)

	err := row.Scan(&userId, &username, &fname, &lname, &email)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil { //another error
		return nil, err
	}

	return &User{Id: userId, Username: username, FirstName: fname, LastName: lname, Email: email}, nil
}
func getWorkspaceFromDB(id int) (*Workspace, error) {
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
		return &Workspace{Id: workspaceId, Name: name, CreatorId: creatorId, Plan: plan}, nil
	}
	return &Workspace{Id: workspaceId, Name: name, CreatorId: creatorId, OutboundMacroId: int(outboundMacroId.Int64), Plan: plan}, nil
}

func getRecordingSpace(id int) (int, error) {
	var bytes int
	row := db.QueryRow(`SELECT SUM(size) FROM recordings WHERE workspace_id=?`, id)

	err := row.Scan(&bytes)
	if err == sql.ErrNoRows {
		return 0, err
	}
	if err != nil { //another error
		return 0, err
	}
	return bytes, nil
}
func getFaxCount(id int) (*int, error) {
	var count int
	row := db.QueryRow(`SELECT COUNT(*) FROM faxes WHERE workspace_id=?`, id)

	err := row.Scan(&count)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil { //another error
		return nil, err
	}
	return &count, nil
}
func getWorkspaceByDomain(domain string) (*Workspace, error) {
	var workspaceId int
	var name string
	var byo bool
	var ipWhitelist bool
	var creatorId int
	s := strings.Split(domain, ".")
	workspaceName := s[0]
	row := db.QueryRow("SELECT id, creator_id, name, byo_enabled, ip_whitelist_disabled FROM workspaces WHERE name=?", workspaceName)

	err := row.Scan(&workspaceId, &creatorId, &name, &byo, &ipWhitelist)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &Workspace{Id: workspaceId, CreatorId: creatorId, Name: name, BYOEnabled: byo, IPWhitelistDisabled: ipWhitelist}, nil
}

func getWorkspaceParams(workspaceId int) (*[]WorkspaceParam, error) {
	// Execute the query
	results, err := db.Query("SELECT `key`, `value` FROM workspace_params WHERE `workspace_id` = ?", workspaceId)
	defer results.Close()
	params := []WorkspaceParam{}
	if err == sql.ErrNoRows {
		// no records setup were setup, just return empty
		return &params, nil
	}
	if err != nil {
		return nil, err
	}

	for results.Next() {
		param := WorkspaceParam{}
		// for each row, scan the result into our tag composite object
		err = results.Scan(&param.Key, &param.Value)
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return &params, nil
}

func getUserByDomain(domain string) (*WorkspaceCreatorFullInfo, error) {
	workspace, err := getWorkspaceByDomain(domain)
	if err != nil {
		return nil, err
	}

	// Execute the query
	params, err := getWorkspaceParams(workspace.Id)
	if err != nil {
		return nil, err
	}
	full := &WorkspaceCreatorFullInfo{
		Id:              workspace.CreatorId,
		Workspace:       workspace,
		WorkspaceParams: params,
		WorkspaceName:   workspace.Name,
		WorkspaceDomain: fmt.Sprintf("%s.lineblocs.com", workspace.Name),
		WorkspaceId:     workspace.Id,
		OutboundMacroId: workspace.OutboundMacroId}

	return full, nil
}

func getUserByDID(did string) (*WorkspaceCreatorFullInfo, error) {
	result:= db.QueryRow(`SELECT
		workspaces.name
		FROM did_numbers
		INNER JOIN workspaces ON workspaces.id = did_numbers.workspace_id
		WHERE did_numbers.api_number = ?`, did)
	var domain string
	err := result.Scan( &domain )
	if err != nil {
		return nil, err
	}

	workspace, err := getWorkspaceByDomain(domain)
	if err != nil {
		return nil, err
	}

	// Execute the query
	params, err := getWorkspaceParams(workspace.Id)
	if err != nil {
		return nil, err
	}
	full := &WorkspaceCreatorFullInfo{
		Id:              workspace.CreatorId,
		Workspace:       workspace,
		WorkspaceParams: params,
		WorkspaceName:   workspace.Name,
		WorkspaceDomain: fmt.Sprintf("%s.lineblocs.com", workspace.Name),
		WorkspaceId:     workspace.Id,
		OutboundMacroId: workspace.OutboundMacroId}

	return full, nil
}

func getUserByTrunkSourceIp(sourceIp string) (*WorkspaceCreatorFullInfo, error) {

	// todo get ipv6
	sourceIpv6 := sourceIp
	result:= db.QueryRow(`SELECT
		workspaces.name
		FROM workspaces
		INNER JOIN sip_trunks ON sip_trunks.workspace_id = workspaces.id
		INNER JOIN sip_trunks_origination_endpoints ON sip_trunks_origination_endpoints.trunk_id = sip_trunks.id
		WHERE sip_trunks_origination_endpoints.ipv4 = ?  OR sip_trunks_origination_endpoints.ipv6 = ?`, sourceIp, sourceIpv6)
	var domain string
	err := result.Scan( &domain )
	if err != nil {
		return nil, err
	}
	workspace, err := getWorkspaceByDomain(domain)
	if err != nil {
		return nil, err
	}

	// Execute the query
	params, err := getWorkspaceParams(workspace.Id)
	if err != nil {
		return nil, err
	}
	full := &WorkspaceCreatorFullInfo{
		Id:              workspace.CreatorId,
		Workspace:       workspace,
		WorkspaceParams: params,
		WorkspaceName:   workspace.Name,
		WorkspaceDomain: fmt.Sprintf("%s.lineblocs.com", workspace.Name),
		WorkspaceId:     workspace.Id,
		OutboundMacroId: workspace.OutboundMacroId}

	return full, nil
}


func getRecordingFromDB(id int) (*Recording, error) {
	var apiId string
	var ready int
	var size int
	var text string
	row := db.QueryRow("SELECT api_id, transcription_ready, transcription_text, size FROM recordings WHERE id=?", id)

	err := row.Scan(&apiId, &ready, &text, &size)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if ready == 1 {
		return &Recording{APIId: apiId, Id: id, TranscriptionReady: true, TranscriptionText: text, Size: size}, nil
	}
	return &Recording{APIId: apiId, Id: id, Size: size}, nil
}

func getCallFromDB(id int) (*Call, error) {
	row := db.QueryRow("SELECT `from`, `to`, `channel_id`, `status`, `direction`, `duration`, `user_id`, `workspace_id`, `started_at`, `created_at`, `updated_at`, `api_id`, `plan_snapshot`) FROM calls WHERE id = ?", id)
	call := Call{}
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

//todo move to microservice
func getPlanRecordingLimit(workspace *Workspace) (int, error) {
	if workspace.Plan == "pay-as-you-go" {
		return 1024, nil
	} else if workspace.Plan == "starter" {
		return 1024 * 2, nil
	} else if workspace.Plan == "pro" {
		return 1024 * 32, nil
	} else if workspace.Plan == "starter" {
		return 1024 * 128, nil
	}
	return 0, nil
}

//todo move to microservice
func getPlanFaxLimit(workspace *Workspace) (*int, error) {
	var res *int
	if workspace.Plan == "pay-as-you-go" {
		*res = 100
	} else if workspace.Plan == "starter" {
		*res = 100
	} else if workspace.Plan == "pro" {
		res = nil
	} else if workspace.Plan == "starter" {
		res = nil
	}
	return res, nil
}
func sendLogRoutineEmail(log *LogRoutine, user *lineblocs.User, workspace *Workspace) error {
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
	m := mg.NewMessage(
		"Lineblocs <monitor@lineblocs.com>",
		"Debug Monitor",
		"Debug Monitor",
		user.Email)
	m.AddCC("contact@lineblocs.com")
	//m.AddBCC("bar@example.com")

	body := `<html>
<head></head>
<body>
	<h1>Lineblocs Monitor Report</h1>
	<h5>` + log.Title + `</h5>
	<p>` + log.Report + `</p>
</body>
</html>`

	m.SetHtml(body)
	//m.AddAttachment("files/test.jpg")
	//m.AddAttachment("files/test.txt")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := mg.Send(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func startLogRoutine(log *LogRoutine) (*string, error) {
	var user *lineblocs.User
	var workspace *Workspace

	user, err := lineblocs.GetUserFromDB(log.UserId)
	if err != nil {
		fmt.Printf("could not get user..")
		return nil, err
	}

	workspace, err = getWorkspaceFromDB(log.WorkspaceId)
	if err != nil {
		fmt.Printf("could not get workspace..")
		return nil, err
	}
	now := time.Now()
	apiId := createAPIID("log")
	stmt, err := db.Prepare("INSERT INTO debugger_logs (`from`, `to`, `title`, `report`, `workspace_id`, `level`, `api_id`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? )")

	if err != nil {
		fmt.Printf("could not prepare query..")
		return nil, err
	}

	defer stmt.Close()
	res, err := stmt.Exec(log.From, log.To, log.Title, log.Report, workspace.Id, log.Level, apiId, now, now)
	if err != nil {
		fmt.Printf("could not execute query..")
		return nil, err
	}

	logId, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("could not get insert id..")
		return nil, err
	}
	logIdStr := strconv.FormatInt(logId, 10)

	go sendLogRoutineEmail(log, user, workspace)

	return &logIdStr, err
}
func incrementServerCallCount(call *Call) error {
	data.mu.Lock()
	defer data.mu.Unlock()
	for _, server := range data.servers {
		if server.IpAddress == call.SourceIp {
			server.LiveCallCount = server.LiveCallCount + 1
		}
	}
	return nil
}
func decrementServerCallCount(call *CallUpdateReq) error {
	data.mu.Lock()
	defer data.mu.Unlock()
	for _, server := range data.servers {
		if server.IpAddress == call.SourceIp {
			server.LiveCallCount = server.LiveCallCount - 1
		}
	}
	return nil
}

func checkRouteMatches(from string, to string, prefix string, prepend string, match string) (bool, error) {
	full := prefix + match
	valid, err := regexp.MatchString(full, to)
	if err != nil {
		return false, err
	}
	return valid, err
}
func shouldUseHostNext(name string, ipPrivate string) (bool, error) {
	return true, nil
}
func checkCIDRMatch(sourceIp string, ipWithCidr string) (bool, error) {
	// remove port if needed

	ipSlice1 := strings.Split(sourceIp, ":")
	ip1 := ipSlice1[0]

	// check if there is port
	ipSlice2 := strings.Split(ipWithCidr, ":")
	var ip2 string
	if len(ipSlice2) > 1 {
		cidr := strings.Split(ipSlice2[1], "/")
		ip2 = ipSlice2[0] + "/" + cidr[1]
	} else {
		ip2 = ipSlice2[0]
	}

	_, net1, err := net.ParseCIDR(ip1 + "/32")
	if err != nil {
		return false, err
	}
	_, net2, err := net.ParseCIDR(ip2)
	if err != nil {
		return false, err
	}

	return net2.Contains(net1.IP), nil
}
func checkPSTNIPWhitelist(did string, sourceIp string) (bool, error) {
	results, err := db.Query(`SELECT 
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
		match, err := checkCIDRMatch(sourceIp, ipWithCidr)
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
func checkBYOPSTNIPWhitelist(did string, sourceIp string) (bool, error) {
	results, err := db.Query(`SELECT 
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
		match, err := checkCIDRMatch(sourceIp, ipWithCidr)
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

func finishValidation(number string, didWorkspaceId string) (bool, error) {
	num, err := libphonenumber.Parse(number, "US")
	if err != nil {
		return false, err
	}
	formattedNum := libphonenumber.Format(num, libphonenumber.E164)
	row := db.QueryRow("SELECT id FROM `blocked_numbers` WHERE `workspace_id` = ? AND `number` = ?", didWorkspaceId, formattedNum)
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
func checkFreeTrialStatus(plan string, started time.Time) string {
	if plan == "trial" {
		now := time.Now()
		//make configurable
		expireDays := 10
		expireHours := expireDays * 24
		started.Add(time.Hour * time.Duration(expireHours))
		if started.After(now) {
			return "expired"
		}
		return "pending-expiry"
	}
	return "not-applicable"
}
func checkIsMakingOutboundCallFirstTime(call Call) {
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
func sendEmail(user *User, subject string, body string) {
}

func createLBResult(results *sql.Rows) (*lineblocs.MediaServer, error) {
	for results.Next() {
		value := lineblocs.MediaServer{}
		err := results.Scan(&value.IpAddress, &value.PrivateIpAddress)
		if err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil

}

// get the server with the least amount of
// calls
func getUserRoutedServer(rtcOptimized bool, workspace *Workspace) (*lineblocs.MediaServer, error) {
	data.mu.Lock()
	defer data.mu.Unlock()
	var result *lineblocs.MediaServer
	for _, server := range data.servers {
		//if result == nil || result != nil && server.LiveCallCount < result.LiveCallCount && rtcOptimized == server.RtcOptimized {
		if result == nil || result != nil && server.LiveCPUPCTUsed < result.LiveCPUPCTUsed && rtcOptimized == server.RtcOptimized {
			result = server
		}
	}
	return result, nil

}

// get the server with the least amount of
// CPU
func getUserRoutedServer2(rtcOptimized bool, workspace *Workspace, routerip string) (*lineblocs.MediaServer, error) {
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

// get the server closest to the DID profile
func getDIDRoutedServer(rtcOptimized bool) (*lineblocs.MediaServer, error) {
	data.mu.Lock()
	defer data.mu.Lock()
	var result *lineblocs.MediaServer
	for _, server := range data.servers {
		if result == nil || result != nil && server.LiveCallCount < result.LiveCallCount && rtcOptimized == server.RtcOptimized {
			result = server
		}
	}
	return result, nil

}

// get the server closest to the DID profile
func getDIDRoutedServer2(rtcOptimized bool) (*lineblocs.MediaServer, error) {
	servers, err := lineblocs.CreateMediaServers()

	if err != nil {
		return nil, err
	}

	var result *lineblocs.MediaServer
	for _, server := range servers {
		//if result == nil || result != nil && server.LiveCallCount < result.LiveCallCount && rtcOptimized == server.RtcOptimized {
		if result == nil || result != nil && server.LiveCPUPCTUsed < result.LiveCPUPCTUsed && rtcOptimized == server.RtcOptimized {
			result = server
		}
	}
	return result, nil
}

func getBestPSTNProvider(from, to *string) (*PSTNInfo, error) {
	// do LCR based on dial prefixes

	results, err := db.Query(`SELECT sip_providers.id, sip_providers.dial_prefix, sip_providers.name, sip_providers_rates.rate_ref_id, sip_providers_rates.rate
		FROM sip_providers
		INNER JOIN sip_providers_rates ON sip_providers_rates.provider_id = sip_providers.id
		WHERE (sip_providers.type_of_provider = 'outbound'
    OR sip_providers.type_of_provider = 'both')
		AND sip_providers.active = 1`)
	if err != nil {
		return nil, err
	}

	var routableProviders  []*RoutableProvider
	var lowestRate *float64 = nil
	var lowestProviderId *int
	var lowestDialPrefix *string
	var longestMatch *int

	routableProviders = make([]*RoutableProvider,0)

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
		results1, err := db.Query(`SELECT dial_prefix
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
			valid, err := regexp.MatchString(full, *to)
			if err != nil {
				return nil, err
			}
			if valid {
				fmt.Println("found matching route...")
				fullLen := len(full)

 				if (longestMatch == nil || fullLen >= *longestMatch) {
					provider := RoutableProvider{
						Provider: id,
						Rate: rate,
						DialPrefix: dialPrefix }
					routableProviders = append( routableProviders, &provider )
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
		number = *lowestDialPrefix + *to

		// lookup hosts
		fmt.Printf("Looking up hosts..\r\n")
		// do LCR based on dial prefixes
		results1, err := db.Query(`SELECT sip_providers_hosts.id, sip_providers_hosts.ip_address, sip_providers_hosts.name, sip_providers_hosts.priority_prefixes
			FROM sip_providers_hosts
			WHERE sip_providers_hosts.provider_id = ?
			`, *lowestProviderId)
		if err != nil {
			return nil, err
		}
		defer results1.Close()
		// TODO check which host is best for routing
		// add area code checking
		var info *PSTNInfo
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
			info = &PSTNInfo{IPAddr: ipAddr, DID: number}
			if bestProviderId == nil {
				bestProviderId = &id
				bestIpAddr = &ipAddr
			}
			// take priority
			if len(prefixArr) != 0 {
				for _, prefix := range prefixArr {
					valid, err := regexp.MatchString(prefix, *to)
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
		info = &PSTNInfo{IPAddr: *bestIpAddr, DID: number}
		return info, nil
	}
	return nil, errors.New("no available routes for LCR...")
}


func doVerifyCaller(workspaceId int, number string) (bool, error) {
	var workspace *Workspace

	if !settings.ValidateCallerId {
		return true, nil
	}

	workspace, err := getWorkspaceFromDB(workspaceId)
	if err != nil {
		return false, err
	}

	num, err := libphonenumber.Parse(number, "US")
	if err != nil {
		return false, err
	}
	formattedNum := libphonenumber.Format(num, libphonenumber.E164)
	fmt.Printf("looking up number %s\r\n", formattedNum)
	fmt.Printf("domain isr %s\r\n", workspace.Name)
	var id string
	row := db.QueryRow("SELECT id FROM `did_numbers` WHERE `number` = ? AND `workspace_id` = ?", formattedNum, workspace.Id)
	err = row.Scan(&id)
	if err != sql.ErrNoRows {
		return true, nil
	}
	return false, nil
}

func getQueryVariable(r *http.Request, key string) *string {
	vals := r.URL.Query()    // Returns a url.Values, which is a map[string][]string
	results, ok := vals[key] // Note type, not ID. ID wasn't specified anywhere.
	var value *string
	if ok {
		if len(results) >= 1 {
			value = &results[0] // The first `?type=model`
		}
	}
	return value
}
func uploadS3(folder string, name string, file multipart.File) error {
	bucket := "lineblocs"
	key := folder + "/" + name
	// The session the S3 Uploader will use
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("ca-central-1")})
	if err != nil {
		return fmt.Errorf("S3 session err: %s", err)
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(session)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	return nil
}

func createS3URL(folder string, id string) string {
	return "https://lineblocs.s3.ca-central-1.amazonaws.com/" + folder + "/" + id
}

func NoContent(w http.ResponseWriter, r *http.Request) {
	// Set up any headers you want here.
	w.WriteHeader(http.StatusNoContent) // send the headers with a 204 response code.
}
func toCents(dollars float64) int {
	result := dollars * 100
	return int(result)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("checking health...\r\n")
	w.Header().Set("Content-Type", "text/plain")
	// execute test query...
	results, err := db.Query("SELECT k8s_pod_id FROM sip_routers")
	if err != nil {
		handleInternalErr("healthz error", err, w)
		return
	}
	defer results.Close()
	fmt.Fprintf(w, "OK\n")
}

func CreateCall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var call Call
	now := time.Now()

	err := json.NewDecoder(r.Body).Decode(&call)
	if err != nil {
		handleInternalErr("CreateCall Could not decode JSON", err, w)
		return
	}
	workspace, err := getWorkspaceFromDB(call.WorkspaceId)
	if err != nil {
		handleInternalErr("CreateCall Could not decode JSON", err, w)
		return
	}

	call.APIId = createAPIID("call")

	if call.Direction == "outbound" {
		//check if this is the first time we are making a call to this destination
		go checkIsMakingOutboundCallFirstTime(call)
	}

	// perform a db.Query insert
	stmt, err := db.Prepare("INSERT INTO calls ( `from`, `to`, `channel_id`, `status`, `direction`, `duration`, `sip_call_id`, `user_id`, `workspace_id`, `started_at`, `created_at`, `updated_at`, `api_id`, `plan_snapshot`, `notes`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, '' )")
	if err != nil {
		handleInternalErr("CreateCall Could not execute query..", err, w)
		return
	}
	defer stmt.Close()
	fmt.Printf("CreateCall args from=%s, to=%s, channel_id=%s, status=%s, direction=%s, user_id=%s, workspace_id=%s, started=%s, plan=%s",
		call.From, call.To, call.ChannelId, call.Status, call.Direction, call.UserId, call.WorkspaceId, now, call.APIId, workspace.Plan)

	res, err := stmt.Exec(call.From, call.To, call.ChannelId, call.Status, call.Direction, "8", call.SIPCallId, call.UserId, call.WorkspaceId, now, now, now, call.APIId, workspace.Plan)

	if err != nil {
		handleInternalErr("CreateCall Could not execute query", err, w)
		return
	}

	callId, err := res.LastInsertId()
	if err != nil {
		handleInternalErr("CreateCall Could not execute query..", err, w)
		return
	}
	//increment call count on server

	/*
			err = incrementServerCallCount(&call)
		if err != nil {
				handleInternalErr("CreateCall Could not increment call count..", err, w);
			return
		}
	*/

	w.Header().Set("X-Call-ID", strconv.FormatInt(callId, 10))
	json.NewEncoder(w).Encode(&call)
}

func UpdateCall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var update CallUpdateReq
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		handleInternalErr("UpdateCall 1 Could not decode JSON", err, w)

		return
	}

	if update.Status == "ended" {
		// perform a db.Query insert
		stmt, err := db.Prepare("UPDATE calls SET `status` = ?, `ended_at` = ?, `updated_at` = ? WHERE `api_id` = ?")
		if err != nil {
			fmt.Printf("UpdateCall 2 Could not execute query..")
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		endedAt := time.Now()
		updatedAt := time.Now()
		_, err = stmt.Exec(update.Status, endedAt, updatedAt, update.CallId)
		if err != nil {
			handleInternalErr("UpdateCall 3 Could not execute query", err, w)
			return
		}
		//increment call count on server
		/*
			err = decrementServerCallCount(&update)
			if err != nil {
				handleInternalErr("UpdateCall Could not decrement call count..", err, w);
				return
			}
		*/
	}

	w.WriteHeader(http.StatusNoContent)
}

func FetchCall(w http.ResponseWriter, r *http.Request) {
	id := getQueryVariable(r, "id")
	id_int, err := strconv.Atoi(*id)
	if err != nil {
		handleInternalErr("FetchCall error occured", err, w)
		return
	}
	call, err := getCallFromDB(id_int)
	if err != nil {
		handleInternalErr("FetchCall error occured", err, w)
		return
	}
	json.NewEncoder(w).Encode(&call)

}

func SetSIPCallID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	callid := r.FormValue("callid")
	apiid := r.FormValue("apiid")
	stmt, err := db.Prepare("UPDATE calls SET sip_call_id = ? WHERE id = ?")
	if err != nil {
		handleInternalErr("SetSIPCallID 1 could not execute query..", err, w)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(callid, apiid)

	if err != nil {
		handleInternalErr("SetSIPCallID 2 could not execute query..", err, w)
		return
	}
}

func SetProviderByIP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ip := r.FormValue("ip")
	apiid := r.FormValue("apiid")

	results, err := db.Query(`SELECT sip_providers_hosts.provider_id FROM sip_providers_hosts WHERE sip_providers_hosts.ip_address = ?`, ip)
	if err != nil {
		handleInternalErr("SetProviderByID error", err, w)
		return
	}
	defer results.Close()
	for results.Next() {
		var providerId int
		results.Scan(&providerId)
		stmt, err := db.Prepare("UPDATE calls SET provider_id = ? WHERE id = ?")
		if err != nil {
			handleInternalErr("SetProviderByID 1 could not execute query..", err, w)
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(providerId, apiid)

		if err != nil {
			handleInternalErr("SetProviderByID 2 could not execute query..", err, w)
			return
		}
	}

}

func CreateConference(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var conference Conference
	err := json.NewDecoder(r.Body).Decode(&conference)
	if err != nil {
		handleInternalErr("CreateConference 1 Could not decode JSON", err, w)
		return
	}
	var id int
	var name string
	row := db.QueryRow("SELECT id, name FROM conferences WHERE workspace_id=? AND name=?", conference.WorkspaceId, conference.Name)
	err = row.Scan(&id, &name)
	if err == sql.ErrNoRows {
		conference.APIId = createAPIID("conf")
		// perform a db.Query insert
		now := time.Now()
		stmt, err := db.Prepare("INSERT INTO conferences (`name`, `workspace_id`, `api_id`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ? )")
		if err != nil {
			handleInternalErr("CreateConference 3 Could not execute query..", err, w)
			return
		}
		defer stmt.Close()
		res, err := stmt.Exec(conference.Name, conference.WorkspaceId, conference.APIId, now, now)

		if err != nil {
			handleInternalErr("CreateConference 4 Could not execute query", err, w)
			return
		}
		conferenceId, err := res.LastInsertId()
		if err != nil {
			handleInternalErr("CreateConference 5 Could not get ID..", err, w)
			return
		}
		w.Header().Set("X-Conference-ID", strconv.FormatInt(conferenceId, 10))
		json.NewEncoder(w).Encode(&conference)
	}

	w.Header().Set("X-Conference-ID", strconv.Itoa(id))
	json.NewEncoder(w).Encode(&conference)
}

func CreateDebit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var debitReq DebitCreateReq

	err := json.NewDecoder(r.Body).Decode(&debitReq)
	if err != nil {
		handleInternalErr("CreateDebit Could not decode JSON", err, w)
		return
	}
	workspace, err := getWorkspaceFromDB(debitReq.WorkspaceId)
	if err != nil {
		fmt.Printf("could not get workspace..")
		return
	}
	rate := lookupBestCallRate(debitReq.Number, debitReq.Type)
	if rate == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	minutes := math.Floor(debitReq.Seconds / 60)
	dollars := minutes * rate.CallRate
	cents := toCents(dollars)
	now := time.Now()
	stmt, err := db.Prepare("INSERT INTO users_debits (`user_id`, `cents`, `source`, `plan_snapshot`, `module_id`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ?, ? )")
	if err != nil {
		handleInternalErr("CreateDebit Could not execute query..", err, w)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(debitReq.UserId, cents, debitReq.Source, workspace.Plan, debitReq.ModuleId, now, now)
	if err != nil {
		handleInternalErr("CreateDebit Could not execute query..", err, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func CreateAPIUsageDebit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var debitReq DebitAPICreateReq

	err := json.NewDecoder(r.Body).Decode(&debitReq)
	if err != nil {
		handleInternalErr("CreateDebit Could not decode JSON", err, w)
		return
	}
	workspace, err := getWorkspaceFromDB(debitReq.WorkspaceId)
	if err != nil {
		fmt.Printf("could not get workspace..")
		return
	}

	if debitReq.Type == "TTS" {
		dollars := calculateTTSCosts(debitReq.Params.Length)
		cents := toCents(dollars)
		source := fmt.Sprintf("API usage - %s", debitReq.Type)
		now := time.Now()
		stmt, err := db.Prepare("INSERT INTO users_debits (`user_id, `cents`, `source`, `plan_snapshot`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ? )")
		if err != nil {
			handleInternalErr("CreateDebit Could not execute query..", err, w)
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(debitReq.UserId, cents, source, workspace.Plan, now, now)
		if err != nil {
			handleInternalErr("CreateAPIUsageDebit Could not execute query..", err, w)
			return
		}
	} else if debitReq.Type == "STT" {
		dollars := calculateSTTCosts(debitReq.Params.RecordingLength)
		cents := toCents(dollars)
		source := fmt.Sprintf("API usage - %s", debitReq.Type)
		now := time.Now()
		stmt, err := db.Prepare("INSERT INTO users_debits (`user_id, `cents`, `source`, `plan_snapshot`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ? )")
		if err != nil {
			handleInternalErr("CreateDebit Could not execute query..", err, w)
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(debitReq.UserId, cents, source, workspace.Plan, now, now)
		if err != nil {
			handleInternalErr("CreateAPIUsageDebit Could not execute query..", err, w)
			return
		}

	}
	w.WriteHeader(http.StatusNoContent)
}
func CreateLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var logReq LogCreateReq

	err := json.NewDecoder(r.Body).Decode(&logReq)
	if err != nil {
		handleInternalErr("CreateLog 1 Could not decode JSON", err, w)
		return
	}
	level := "info"
	if logReq.Level != nil {
		level = *logReq.Level
	}
	from := ""
	if logReq.From != nil {
		from = *logReq.From
	}

	to := ""
	if logReq.To != nil {
		to = *logReq.To
	}
	var log *LogRoutine = &LogRoutine{From: from,
		To:          to,
		Level:       level,
		Title:       logReq.Title,
		Report:      logReq.Report,
		FlowId:      logReq.FlowId,
		UserId:      logReq.UserId,
		WorkspaceId: logReq.WorkspaceId}
	_, err = startLogRoutine(log)
	if err != nil {
		handleInternalErr("CreateLog 2 log routine error", err, w)
		return
	}
}

func CreateLogSimple(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logType := r.FormValue("type")
	level := r.FormValue("level")
	domain := r.FormValue("domain")
	workspace, err := getWorkspaceByDomain(domain)
	if err != nil {
		handleInternalErr("CreateLog Could not decode JSON", err, w)
		return
	}

	if &level == nil {
		level = "infO"
	}

	var title string
	var report string
	switch logType {
	case "verify-callerid-cailed":
		title = "Caller ID Verify failed.."
		report = "Caller ID Verify failed.."
	}
	var log *LogRoutine = &LogRoutine{
		From:        "",
		To:          "",
		Level:       level,
		Title:       title,
		Report:      report,
		UserId:      workspace.CreatorId,
		WorkspaceId: workspace.Id}
	_, err = startLogRoutine(log)
	if err != nil {
		handleInternalErr("CreateLog log routine error", err, w)
		return
	}

}

func CreateFax(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var fax *Fax
	file, handler, err := r.FormFile("file")
	now := time.Now()
	if err != nil {
		handleInternalErr("CreateFax error occured", err, w)
		return

	}

	workspace, err := getWorkspaceFromDB(fax.WorkspaceId)
	if err != nil {
		fmt.Printf("could not get workspace..")
		return
	}

	userId := r.FormValue("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		handleInternalErr("CreateFax error occured user ID", err, w)
		return
	}

	workspaceId := r.FormValue("workspace_id")
	workspaceIdInt, err := strconv.Atoi(workspaceId)
	if err != nil {
		handleInternalErr("CreateFax error occured workspace ID", err, w)
		return
	}

	callId := r.FormValue("call_id")
	callIdInt, err := strconv.Atoi(callId)
	if err != nil {
		handleInternalErr("CreateFax error occured call ID", err, w)
		return

	}

	name := r.FormValue("name")
	defer file.Close()

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		handleInternalErr("CreateFax error occured", err, w)
		return
	}
	defer f.Close()
	apiId := createAPIID("fax")
	uri := createS3URL("faxes", apiId)
	count, err := getFaxCount(workspace.Id)
	if err != nil {
		handleInternalErr("CreateFax error occured", err, w)
		return
	}

	stmt, err := db.Prepare("INSERT INTO faxes (`uri`, `size`, `name`, `user_id`, `call_id`, `workspace_id`, `api_id`, `plan`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		handleInternalErr("CreateFax error occured", err, w)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(uri, handler.Size, name, userId, callId, workspaceId, apiId, workspace.Plan, now, now)
	if err != nil {
		handleInternalErr("CreateFax error occured", err, w)
		return
	}

	faxId, err := res.LastInsertId()
	if err != nil {
		handleInternalErr("CreateFax error occured", err, w)
		return
	}

	fax = &Fax{UserId: userIdInt, WorkspaceId: workspaceIdInt, CallId: callIdInt, Uri: uri}
	w.Header().Set("X-Fax-ID", strconv.FormatInt(faxId, 10))
	json.NewEncoder(w).Encode(fax)
	limit, err := getPlanFaxLimit(workspace)
	if err != nil {
		handleInternalErr("CreateFax error occured", err, w)
		return
	}
	newCount := (*count) + 1
	if newCount > *limit {
		fmt.Printf("not saving fax due to limit reached..")
		return
	}
	go uploadS3("faxes", apiId, file)
}

func CreateRecording(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var recording Recording

	now := time.Now()
	err := json.NewDecoder(r.Body).Decode(&recording)
	if err != nil {
		handleInternalErr("CreateCall Could not decode JSON", err, w)
		return
	}

	workspace, err := getWorkspaceFromDB(recording.WorkspaceId)
	if err != nil {
		fmt.Printf("could not get workspace..")
		return
	}

	recording.APIId = createAPIID("rec")

	// perform a db.Query insert
	stmt, err := db.Prepare("INSERT INTO recordings (`user_id`, `call_id`, `workspace_id`, `status`, `name`, `uri`, `tag`, `api_id`, `plan_snapshot`, `storage_id`, `storage_server_ip`, `trim`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		handleInternalErr("CreateRecording error.", err, w)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(
		recording.UserId,
		recording.CallId,
		recording.WorkspaceId,
		"started",
		"",
		"",
		"",
		recording.APIId,
		workspace.Plan,
		recording.StorageId,
		recording.StorageServerIp,
		recording.Trim,
		now,
		now)
	if err != nil {
		handleInternalErr("CreateRecording error.", err, w)
		return
	}
	recId, err := res.LastInsertId()
	if err != nil {
		handleInternalErr("CreateRecording error.", err, w)
		return
	}
	if recording.Tags != nil {
		for _, v := range *recording.Tags {
			fmt.Printf("adding tag to recording %s\r\n", v)
			stmt, err := db.Prepare("INSERT INTO recording_tags (`recording_id`, `tag`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?)")
			if err != nil {
				handleInternalErr("CreateRecording error.", err, w)
			}

			defer stmt.Close()
			res, err = stmt.Exec(recId, v, now, now)
			if err != nil {
				handleInternalErr("CreateRecording error.", err, w)
				return
			}
		}
	}

	defer stmt.Close()
	w.Header().Set("X-Recording-ID", strconv.FormatInt(recId, 10))
	json.NewEncoder(w).Encode(&recording)
}

func UpdateRecording(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//var recording Recording
	file, handler, err := r.FormFile("file")
	status := r.FormValue("status")
	now := time.Now()
	recordingId := r.FormValue("recording_id")
	recordingIdInt, err := strconv.Atoi(recordingId)
	record, err := getRecordingFromDB(recordingIdInt)
	if err != nil {
		fmt.Printf("could not get recording..")
		return
	}

	workspace, err := getWorkspaceFromDB(record.WorkspaceId)
	if err != nil {
		fmt.Printf("could not get workspace..")
		return
	}

	if err != nil {
		handleInternalErr("UpdateRecording error occured", err, w)
		return

	}
	defer file.Close()

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		handleInternalErr("UpdateRecording error occured", err, w)
		return
	}
	defer f.Close()

	size, err := getRecordingSpace(workspace.Id)
	if err != nil {
		fmt.Printf("could not get recording space..")
		return
	}
	apiId := createAPIID("rec")
	uri := createS3URL("recordings", apiId)
	stmt, err := db.Prepare("UPDATE `recordings` SET `status` = ?, `uri` = ?, `size` = ?, `updated_at` = ? WHERE `id` = ?")
	if err != nil {
		handleInternalErr("UpdateRecording error occured", err, w)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(status, uri, handler.Size, now, recordingIdInt)
	if err != nil {
		handleInternalErr("UpdateRecording error occured", err, w)
		return
	}
	//will not save
	limit, err := getPlanRecordingLimit(workspace)
	newSpace := size + int(handler.Size)
	if newSpace > limit {
		fmt.Printf("not saving recording due to space limit reached..")
		return
	}
	go uploadS3("recordings", apiId, file)
}

func UpdateRecordingTranscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var update RecordingTranscriptionReq
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		handleInternalErr("UpdateRecordingTranscription error", err, w)

		return
	}
	stmt, err := db.Prepare("UPDATE recordings SET `transcription_ready` = ?, `transcription_text` = ? WHERE `id` = ?")
	_, err = stmt.Exec("1", update.Text, update.RecordingId)
	if err != nil {
		handleInternalErr("UpdateCall Could not execute query", err, w)
		return
	}
	defer stmt.Close()
	w.WriteHeader(http.StatusNoContent)
}
func GetRecording(w http.ResponseWriter, r *http.Request) {
	id := getQueryVariable(r, "id")
	id_int, err := strconv.Atoi(*id)
	if err != nil {
		handleInternalErr("GetRecording error occured", err, w)
		return
	}
	record, err := getRecordingFromDB(id_int)
	if err != nil {
		handleInternalErr("GetRecording error occured", err, w)
		return
	}
	json.NewEncoder(w).Encode(&record)

}
func VerifyCaller(w http.ResponseWriter, r *http.Request) {
	workspaceId := getQueryVariable(r, "workspace_id")
	workspaceIdInt, err := strconv.Atoi(*workspaceId)
	if err != nil {
		handleInternalErr("VerifyCaller error occured", err, w)
		return
	}

	number := getQueryVariable(r, "number")

	valid, err := doVerifyCaller(workspaceIdInt, *number)
	if err != nil {
		handleInternalErr("VerifyCaller error occured", err, w)
		return
	}
	result := VerifyNumber{Valid: valid}
	json.NewEncoder(w).Encode(&result)
}
func VerifyCallerByDomain(w http.ResponseWriter, r *http.Request) {
	domain := getQueryVariable(r, "domain")
	number := getQueryVariable(r, "number")

	workspace, err := getWorkspaceByDomain(*domain)
	if err != nil {
		handleInternalErr("VerifyCallerByDomain error 1 occured", err, w)
		return
	}
	valid, err := doVerifyCaller(workspace.Id, *number)
	if err != nil {
		handleInternalErr("VerifyCaller error 2 occured", err, w)
		return
	}
	if !valid {
		handleInternalErr("VerifyCaller number not valid", err, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func GetUserAssignedIP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Get assigned IP called..\r\n")
	opt := getQueryVariable(r, "rtcOptimized")
	var err error
	var rtcOptimized bool
	// default
	rtcOptimized = false

	if opt != nil {
		rtcOptimized, err = strconv.ParseBool(*opt)
	}
	if err != nil {
		handleInternalErr("GetUserAssignedIP error", err, w)
		return
	}

	domain := getQueryVariable(r, "domain")
	routerip := getQueryVariable(r, "routerip")
	fmt.Printf("Finding server for domain " + *domain + "..\r\n")
	fmt.Printf("Router IP is " + *routerip + "..\r\n")
	//ru := getQueryVariable(r, "ru")
	workspace, err := getWorkspaceByDomain(*domain)
	if err != nil {
		handleInternalErr("GetUserAssignedIP error 1", err, w)
		return
	}

	if err != nil {
		handleInternalErr("GetUserAssignedIP error occured 2", err, w)
		return
	}

	server, err := getUserRoutedServer2(rtcOptimized, workspace, *routerip)

	if err != nil {
		handleInternalErr("GetUserAssignedIP error occured 3", err, w)
		return
	}
	if server == nil {
		handleInternalErr("GetUserAssignedIP could not get server", err, w)
		return
	}
	fmt.Printf("Found server " + server.PrivateIpAddress + "..\r\n")
	w.Write([]byte(server.PrivateIpAddress))
}

func GetUserByDomain(w http.ResponseWriter, r *http.Request) {
	domain := getQueryVariable(r, "domain")

	info, err := getUserByDomain(*domain)
	if err != nil {
		handleInternalErr("GetUserByDomain error occured", err, w)
		return
	}
	json.NewEncoder(w).Encode(&info)

}
func GetUserByDID(w http.ResponseWriter, r *http.Request) {
	did := getQueryVariable(r, "did")

	info, err := getUserByDID(*did)
	if err != nil {
		handleInternalErr("GetUserByDID error occured", err, w)
		return
	}
	json.NewEncoder(w).Encode(&info)
}


func GetUserByTrunkSourceIp(w http.ResponseWriter, r *http.Request) {
	sourceIp:= getQueryVariable(r, "source_ip")

	info, err := getUserByTrunkSourceIp(*sourceIp)
	if err != nil {
		handleInternalErr("GetUserByTrunkSourceIp error occured", err, w)
		return
	}
	json.NewEncoder(w).Encode(&info)

}

func GetWorkspaceMacros(w http.ResponseWriter, r *http.Request) {
	workspace := getQueryVariable(r, "workspace")
	// Execute the query
	results, err := db.Query("SELECT title, code, compiled_code FROM macro_functions WHERE `workspace_id` = ?", workspace)
	if err != nil {
		handleInternalErr("GetWorkspaceMacros error", err, w)
		return
	}
	defer results.Close()
	values := []MacroFunction{}

	for results.Next() {
		value := MacroFunction{}
		err = results.Scan(&value.Title, &value.Code, &value.CompiledCode)
		if err != nil {
			handleInternalErr("GetWorkspaceMacros error", err, w)
			return
		}

		// for each row, scan the result into our tag composite object
		values = append(values, value)
	}
	json.NewEncoder(w).Encode(&values)
}
func GetDIDNumberData(w http.ResponseWriter, r *http.Request) {
	number := getQueryVariable(r, "number")
	var info WorkspaceDIDInfo
	var flowJson sql.NullString
	fmt.Printf("Looking up number: %s", *number)
	// Execute the query
	row := db.QueryRow(`SELECT 
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
		`, *number)
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

	if err != nil && err != sql.ErrNoRows {
		handleInternalErr("GetDIDNumberData lookup error", err, w)
		return
	}

	checkOtherTable := false
	if err == sql.ErrNoRows {
		checkOtherTable = true
	}

	if !checkOtherTable {
		if flowJson.Valid {
			info.FlowJSON = flowJson.String
		}

		params, err := getWorkspaceParams(info.WorkspaceId)
		if err != nil {
			handleInternalErr("GetDIDNumberData 1 error", err, w)
			return
		}

		info.WorkspaceParams = params
		json.NewEncoder(w).Encode(&info)
		return
	}
	// Execute the query
	row = db.QueryRow(`SELECT 
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
		`, *number)
	err = row.Scan(
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
	if err != nil {
		handleInternalErr("GetDIDNumberData 3 error", err, w)
		return
	}

	if flowJson.Valid {
		info.FlowJSON = flowJson.String
	}

	params, err := getWorkspaceParams(info.WorkspaceId)
	info.WorkspaceParams = params
	json.NewEncoder(w).Encode(&info)
}

func buildProviderIPs() {

}
func GetPSTNProviderIP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("received PSTN request..\r\n")
	from := getQueryVariable(r, "from")
	to := getQueryVariable(r, "to")
	domain := getQueryVariable(r, "domain")
	//ru := getQueryVariable(r, "ru")
	workspace, err := getWorkspaceByDomain(*domain)
	if err != nil {
		handleInternalErr("GetPSTNProviderIP error", err, w)
		return
	}
	if workspace.BYOEnabled {
		fmt.Println("Checking BYO..")
		results, err := db.Query(`SELECT byo_carriers.name, byo_carriers.ip_address, byo_carriers_routes.prefix, byo_carriers_routes.prepend, byo_carriers_routes.match
		FROM byo_carriers_routes
		INNER JOIN byo_carriers  ON byo_carriers.id = byo_carriers_routes.carrier_id
		INNER JOIN workspaces ON workspaces.id = byo_carriers.workspace_id
		INNER JOIN users ON users.id = workspaces.creator_id
		WHERE byo_carriers.workspace_id = ? 

	`, workspace.Id)
		if err != nil {
			handleInternalErr("GetPSTNProviderIP error", err, w)
			return
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
				handleInternalErr("GetPSTNProviderIP error", err, w)
				return
			}
			if !ip.Valid {
				fmt.Printf("skipping 1 PSTN IP result as private IP is empty..\r\n")
				continue
			}
			valid, err := checkRouteMatches(*from, *to, prefix, prepend, match)
			if err != nil {
				fmt.Printf("error occured when trying to match from: %s, to: %s, prefix: %s, prepend: %s, match: %s", *from, *to, prefix, prepend, match)
				continue
			}
			if valid {
				var number string
				number = prepend + *to
				info := &PSTNInfo{IPAddr: ip.String, DID: number}
				json.NewEncoder(w).Encode(&info)
				return
			}
		}
	}

	info, err := getBestPSTNProvider(from, to)
	if err != nil {
		handleInternalErr("getPSTNProviderIp error 1 ", err, w)
		return
	}

	json.NewEncoder(w).Encode(&info)
	return

	/**
	TODO: this was the new code created for updating routing logic. Currently the default code does LCR, whereas this new code
	lets users create custom diagrams that basically let them decide how they want to route the call. If you want to enable code
	below, please ensure the default LCR logic is removed prior to that.
	**/
	/*
	
	NEW code
	destCode, err := helpers.ParseCountryCode(*to)

	if err != nil {
		fmt.Println("Couldnt parse country code for destination number... error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println( "code is: " + destCode )

	originCode, err := helpers.ParseCountryCode(*from)

	if err != nil {
		fmt.Println("Couldnt parse country code for origination number... error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var flow *helpers.Flow
	fmt.Println( "code is: " + originCode )

	creator := strconv.Itoa(workspace.CreatorId)
	flow,err =createRoutingFlow( from, to, &creator )
	if err != nil {
			handleInternalErr("ProcessRouterFlow error 1", err, w)
			return
	}

	data := make(map[string]string)
	data["origin_code"] = originCode
	data["dest_code"] = destCode
	data["from"] = *from
	data["to"] = *to

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

	if len( providers ) == 0  {
		fmt.Println("No providers available..")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bestProvider:=provider
	if len( bestProvider.Hosts)  == 0  {
		fmt.Println("No hosts available..")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bestIpAddr:=provider.Hosts[ 0 ].IPAddr
	var number string
	number = bestProvider.Hosts[0].Prefix + *to


	info := &PSTNInfo{IPAddr: bestIpAddr, DID: number}
	json.NewEncoder(w).Encode(info)
	*/
}

func GetPSTNProviderIPForTrunk(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("received PSTN request for trunk..\r\n")
	from := getQueryVariable(r, "from")
	to := getQueryVariable(r, "to")

	info, err := getBestPSTNProvider(from, to)
	if err != nil {
		handleInternalErr("getpstnprovideripfortrunk error 1", err, w)
		return
	}

	json.NewEncoder(w).Encode(&info)
	return
}


func IPWhitelistLookup(w http.ResponseWriter, r *http.Request) {
	source := getQueryVariable(r, "ip")
	domain := getQueryVariable(r, "domain")
	workspace, err := getWorkspaceByDomain(*domain)
	if err != nil {
		handleInternalErr("IPWhitelistLookup error occured", err, w)
		return
	}
	results, err := db.Query("SELECT ip, `range` FROM ip_whitelist WHERE `workspace_id` = ?", workspace.Id)
	if err != nil {
		handleInternalErr("IPWhitelistLookup error", err, w)
		return
	}
	defer results.Close()
	if workspace.IPWhitelistDisabled {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for results.Next() {
		var ip string
		var ipRange string
		err = results.Scan(&ip, &ipRange)
		if err != nil {
			handleInternalErr("IPWhitelistLookup error", err, w)
			return

		}
		ipWithCidr := ip + ipRange
		match, err := checkCIDRMatch(*source, ipWithCidr)
		if err != nil {
			handleInternalErr("IPWhitelistLookup error", err, w)
			return
		}
		if match {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
func GetDIDAcceptOption(w http.ResponseWriter, r *http.Request) {
	did := getQueryVariable(r, "did")

	row := db.QueryRow(`SELECT did_action FROM did_numbers WHERE did_numbers.api_number = ?`, did)
	var action string
	err := row.Scan(&action)
	if err == nil {
		w.Write([]byte(action))
		return
	}
	if err != nil && err != sql.ErrNoRows {
		handleInternalErr("GetDIDAcceptOption error 1 occured", err, w)
		return
	}

	row = db.QueryRow(`SELECT did_action FROM byo_did_numbers WHERE byo_did_numbers.number = ?`, did)
	err = row.Scan(&action)
	fmt.Println("err check is ", err)
	if err == nil {
		w.Write([]byte(action))
		return
	}
	handleInternalErr("GetDIDAcceptOption error 2 occured", err, w)
}
func GetDIDAssignedIP(w http.ResponseWriter, r *http.Request) {
	server, err := getDIDRoutedServer2(false)
	if err != nil {
		handleInternalErr("GetUserAssignedIP error occured", err, w)
		return
	}
	if server == nil {
		handleInternalErr("GetUserAssignedIP could not get server", err, w)
	}
	w.Write([]byte(server.PrivateIpAddress))
}
func GetTrunkAssignedIP(w http.ResponseWriter, r *http.Request) {
	server, err := getDIDRoutedServer2(false)
	if err != nil {
		handleInternalErr("GetUserAssignedIP error occured", err, w)
		return
	}
	if server == nil {
		handleInternalErr("GetUserAssignedIP could not get server", err, w)
	}
	w.Write([]byte(server.PrivateIpAddress))
}

func GetCallerIdToUse(w http.ResponseWriter, r *http.Request) {
	domain := getQueryVariable(r, "domain")
	extension := getQueryVariable(r, "extension")
	workspace, err := getWorkspaceByDomain(*domain)
	if err != nil {
		handleInternalErr("GetCallerIdToUse error 1 ", err, w)
		return
	}

	var callerId string
	fmt.Printf("Looking up caller ID in domain %s, ID %d, extension %s\r\n", workspace.Name, workspace.Id, *extension)
	row := db.QueryRow("SELECT caller_id FROM extensions WHERE workspace_id=? AND username = ?", workspace.Id, *extension)

	err = row.Scan(&callerId)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	info := &CallerIDInfo{CallerID: callerId}
	json.NewEncoder(w).Encode(&info)
}

func AddPSTNProviderTechPrefix(w http.ResponseWriter, r *http.Request) {

}
func ProcessRouterFlow(w http.ResponseWriter, r *http.Request) {
	var flow *helpers.Flow;
	callto := getQueryVariable(r, "callto")
	callfrom := getQueryVariable(r, "callfrom")
	userId := getQueryVariable(r, "userid")

	destCode, err := helpers.ParseCountryCode(*callto)

	if err != nil {
		panic(err)
	}
	fmt.Println( "code is: " + destCode )

	originCode, err := helpers.ParseCountryCode(*callfrom)

	if err != nil {
		panic(err)
	}
	fmt.Println( "code is: " + originCode )
	flow,err =createRoutingFlow( callfrom, callto, userId )
	if err != nil {
			handleInternalErr("ProcessRouterFlow error 1", err, w)
			return
	}

	data := make(map[string]string)
	data["origin_code"] = originCode
	data["dest_code"] = destCode
	data["from"] = *callfrom
	data["to"] = *callto

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

func ProcessSIPTrunkCall(w http.ResponseWriter, r *http.Request) {
	did := getQueryVariable(r, "did")

	// get trunk from 
	results, err := db.Query(`SELECT 
	sip_trunks_origination_endpoints.sip_uri
	FROM did_numbers
	INNER JOIN sip_trunks ON sip_trunks.id = did_numbers.trunk_id
	INNER JOIN sip_trunks_origination_endpoints ON sip_trunks_origination_endpoints.trunk_id = sip_trunks.id
	WHERE did_numbers.api_number = ?`, did)

	
	if err != nil {
		handleInternalErr("ProcessSIPTrunkCall error 1 valid", err, w)
		return
	}
	defer results.Close()


	for results.Next() {
		fmt.Printf("trying to route to user trunk server..\r\n");
		var trunkSIPURI string
		results.Scan( &trunkSIPURI )
		fmt.Printf("found SIP trunk server %s\r\n", trunkSIPURI);
		w.Write([]byte(trunkSIPURI));
		return;
	}
	fmt.Println("No trunks to route to..")
	w.WriteHeader(http.StatusInternalServerError)
}


func GetExtensionFlowInfo(w http.ResponseWriter, r *http.Request) {
	extension := getQueryVariable(r, "extension")
	workspaceId := getQueryVariable(r, "workspace")
	//number := getQueryVariable(r, "number")
	//workspace, err := getWorkspaceFromDB(*workspaceId)
	/*
		if err != nil {
			handleInternalErr("GetExtensionFlowInfo error", err, w)
			return
		}
	*/

	var info ExtensionFlowInfo
	var trialStartedTime time.Time
	row := db.QueryRow(`SELECT flows.workspace_id,
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
		w.WriteHeader(http.StatusNotFound)
		return
	}
	info.FreeTrialStatus = checkFreeTrialStatus(info.Plan, trialStartedTime)
	json.NewEncoder(w).Encode(&info)
}

func GetFlowInfo(w http.ResponseWriter, r *http.Request) {
	flowId := getQueryVariable(r, "flow_id")
	workspaceId := getQueryVariable(r, "workspace")
	//number := getQueryVariable(r, "number")
	//workspace, err := getWorkspaceFromDB(*workspaceId)
	/*
		if err != nil {
			handleInternalErr("GetExtensionFlowInfo error", err, w)
			return
		}
	*/

	var info ExtensionFlowInfo
	var trialStartedTime time.Time
	row := db.QueryRow(`SELECT flows.workspace_id,
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
		w.WriteHeader(http.StatusNotFound)
		return
	}
	info.FreeTrialStatus = checkFreeTrialStatus(info.Plan, trialStartedTime)
	json.NewEncoder(w).Encode(&info)
}
func GetDIDDomain(w http.ResponseWriter, r *http.Request) {

}
func GetCodeFlowInfo(w http.ResponseWriter, r *http.Request) {
	code := getQueryVariable(r, "code")
	workspaceId := getQueryVariable(r, "workspace")
	//number := getQueryVariable(r, "number")
	//workspace, err := getWorkspaceFromDB(*workspaceId)
	/*
		if err != nil {
			handleInternalErr("GetExtensionFlowInfo error", err, w)
			return
		}
	*/

	var info CodeFlowInfo
	var trialStartedTime time.Time
	row := db.QueryRow(`SELECT 
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
		json.NewEncoder(w).Encode(&info)
		return
	}

	if err != nil {
		handleInternalErr("GetCodeFlowInfo error", err, w)
		return

	}
	info.FreeTrialStatus = checkFreeTrialStatus(info.Plan, trialStartedTime)
	info.FoundCode = true
	json.NewEncoder(w).Encode(&info)

}
func IncomingDIDValidation(w http.ResponseWriter, r *http.Request) {
	// check for all types of call routing scenarios
	// 1. PSTN DID call
	// 2. Hosted SIP trunk call
	// 3. BYOC trunk call
	did := getQueryVariable(r, "did")
	number := getQueryVariable(r, "number")
	source := getQueryVariable(r, "source")
	// Execute the query
	row := db.QueryRow(`SELECT did_numbers.number, did_numbers.api_number, did_numbers.workspace_id, did_numbers.trunk_id FROM did_numbers WHERE did_numbers.api_number = ?`, did)
	var didNumber string
	var didApiNumber string
	var didWorkspaceId string
	var trunkId int
	err := row.Scan(&didNumber,
		&didApiNumber,
		&didWorkspaceId,
		&trunkId)
	if err == nil {

		// check if we're routing to user SIP turnk
		if trunkId != 0 {
			fmt.Printf("found trunk associated with DID number -- routing to user SIP trunk")
			w.Write([]byte("user_sip_trunk"))
			return
		}
		match, err := checkPSTNIPWhitelist(*did, *source)
		if err != nil {
			handleInternalErr("IncomingDIDValidation error 1", err, w)
			return
		}

		if !match {
			fmt.Println("IncomingDIDValidation no match found 1")
			w.WriteHeader(http.StatusInternalServerError) // send the headers with a 204 response code.
			return
		}
		fmt.Printf("Matched incoming DID..")
		valid, err := finishValidation(*number, didWorkspaceId)
		if err != nil {
			handleInternalErr("IncomingDIDValidation error 2 valid", err, w)
			return
		}
		if !valid {
			fmt.Printf("number not valid..")
			w.WriteHeader(http.StatusInternalServerError) // send the headers with a 204 response code.
			return
		}
		w.Write([]byte("network_managed"))
		return
	}

	fmt.Println("looking up BYO DIDs now...")
	//check BYO DIDs
	row = db.QueryRow(`SELECT byo_did_numbers.number, byo_did_numbers.workspace_id FROM byo_did_numbers WHERE byo_did_numbers.number = ?`, did)
	var byoDidNumber string
	var byoDidWorkspaceId string
	err = row.Scan(&byoDidNumber,
		&byoDidWorkspaceId)
	if err == nil {
		match, err := checkBYOPSTNIPWhitelist(*did, *source)
		if err != nil {
			handleInternalErr("IncomingDIDValidation error 3", err, w)
			return
		}
		if !match {
			fmt.Println("IncomingDIDValidation no match found 2")
			w.WriteHeader(http.StatusInternalServerError) // send the headers with a 204 response code.
			return
		}
		fmt.Printf("Matched incoming DID..")
		valid, err := finishValidation(*number, byoDidWorkspaceId)
		if err != nil {
			handleInternalErr("IncomingDIDValidation error 4 valid", err, w)
			return
		}
		if !valid {
			fmt.Printf("number not valid..")
			w.WriteHeader(http.StatusInternalServerError) // send the headers with a 204 response code.
			return
		}

		w.Write([]byte("byo_carrier"));
	}
	handleInternalErr("IncomingDIDValidation error 3", errors.New("no DID match found..."), w)
	//w.WriteHeader(http.StatusNoContent) // send the headers with a 204 response code.
}


func IncomingTrunkValidation(w http.ResponseWriter, r *http.Request) {
	//did := getQueryVariable(r, "did")
	//number := getQueryVariable(r, "number")
	//source := getQueryVariable(r, "source")
	fromdomain := getQueryVariable(r, "fromdomain")
	//destDomain := getQueryVariable(r, "destdomain")
	// Execute the query
	/*
	row:= db.QueryRow(`SELECT did_numbers.trunk_id FROM did_numbers WHERE did_numbers.api_number = ? AND did_numbers.trunk_id IS NOT NULL`, did)
	err := row.Scan(&trunkId)
	if err != nil {
		handleInternalErr("IncomingTrunkValidation error 1d", err, w)
		return
	}
	*/
	trunkip, err := lookupSIPAddress( *fromdomain );
	if err != nil {
		handleInternalErr("IncomingTrunkValidation error 4 valid", err, w)
		return
	}

	fmt.Printf("from domain %s trunk IP is %s..\r\n", *fromdomain, *trunkip);

	results, err := db.Query(`SELECT 
	sip_trunks_origination_settings.recovery_sip_uri,
	sip_trunks_origination_endpoints.sip_uri
	FROM sip_trunks_origination_endpoints
	INNER JOIN sip_trunks ON sip_trunks.id = sip_trunks_origination_endpoints.trunk_id
	INNER JOIN sip_trunks_origination_settings ON sip_trunks_origination_settings.trunk_id = sip_trunks.id`)

	
	if err != nil {
		handleInternalErr("IncomingTrunkValidation error 1 valid", err, w)
		return
	}
	defer results.Close()


	for results.Next() {
		fmt.Printf("trying to route to SIP server..\r\n");
		var routingSIPURI string
		var recoverySIPURI string
		err := results.Scan(
			&recoverySIPURI,
			&routingSIPURI)
		if err != nil {
			handleInternalErr("IncomingTrunkValidation error 1 valid", err, w)
			return
		}
		fmt.Printf("SIP routing URI = %s SIP recovery URI = %s\r\n", routingSIPURI, recoverySIPURI);
		// TODO do some health checks here to see if SIP server is actually up..
		ips, err := lookupSIPAddresses( routingSIPURI );
		if err != nil {
			fmt.Printf("failed to lookup SIP server %s\r\n", routingSIPURI)
			continue
		}
		for _,ip := range *ips {
			ipAddr := ip.String()
			fmt.Printf("found IP = %s\r\n", ipAddr)
			fmt.Printf("comparing with source IP = %s\r\n", *trunkip)
			if ipAddr == *trunkip {
				w.Write([]byte(ipAddr));
				return
			}
		}
	}

	fmt.Printf("checked all SIP trunks no matches were found... error.")
	w.WriteHeader(http.StatusInternalServerError) // send the headers with a 204 response code.
}


func lookupSIPAddresses( host string ) (*[]net.IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil,err
	}
	if len(ips) == 0 {
		return nil, errors.New("No IP match found..")
	}
	return &ips, nil
}

// get first match
func lookupSIPAddress( host string ) (*string, error) {
	ips, err := lookupSIPAddresses(host)
	if err != nil {
		return nil,err
	}
	ip := (*ips)[0].String()
	return &ip, nil
}
func checkSIPServerHealth( routingSIPURI string ) (bool, error) {
	return true, nil
}

func IncomingMediaServerValidation(w http.ResponseWriter, r *http.Request) {
	//number:= getQueryVariable(r, "number")
	source := getQueryVariable(r, "source")
	//did := getQueryVariable(r, "did")
	results, err := db.Query(`SELECT media_servers.ip_address, media_servers.ip_address_range FROM media_servers`)
	// Execute the query
	if err != nil {
		handleInternalErr("IncomingMediaServerValidation error 1", err, w)
		return
	}
	defer results.Close()
	for results.Next() {
		var ipAddr string
		var ipRange string
		err = results.Scan(&ipAddr, &ipRange)
		if err != nil {
			handleInternalErr("IncomingMediaServerValidation error 2", err, w)
			return
		}
		full := ipAddr + ipRange
		fmt.Printf("checking IP = %s", ipAddr)
		match, err := checkCIDRMatch(*source, full)
		if err != nil {
			handleInternalErr("IncomingMediaServerValidation error 3", err, w)
			return
		}
		if match {
			w.WriteHeader(http.StatusNoContent) // send the headers with a 204 response code.
			return
		}
	}
	fmt.Printf("No media server found..")
	w.WriteHeader(http.StatusInternalServerError) // send the headers with a 204 response code.

}
func SendAdminEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var emailInfo EmailInfo
	err := json.NewDecoder(r.Body).Decode(&emailInfo)
	if err != nil {
		handleInternalErr("SendAdminEmail Could not decode JSON", err, w)
		return
	}
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
	m := mg.NewMessage(
		"Lineblocs <monitor@lineblocs.com>",
		"Admin Error",
		"Admin Error",
		"contact@lineblocs.com")
	body := `<html>
<head></head>
<body>
	<h1>Lineblocs Admin Monitor</h1>
	<p>` + emailInfo.Message + `</p>
</body>
</html>`

	m.SetHtml(body)
	//m.AddAttachment("files/test.jpg")
	//m.AddAttachment("files/test.txt")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err = mg.Send(ctx, m)
	if err != nil {
		handleInternalErr("SendAdminEmail error", err, w)
		return
	}
	return
}

func LookupSIPTrunkByDID(w http.ResponseWriter, r *http.Request) {
	did := getQueryVariable(r, "did")
	results, err := db.Query(`SELECT
		sip_trunks_origination_endpoints.sip_uri,
		sip_trunks_origination_settings.recovery_sip_uri
		FROM sip_trunks_origination_endpoints
		INNER JOIN did_numbers ON did_numbers.trunk_id = sip_trunks_origination_endpoints.trunk_id
		INNER JOIN sip_trunks_origination_settings  ON sip_trunks_origination_settings.trunk_id = sip_trunks_origination_endpoints.trunk_id
		WHERE did_numbers.api_number = ?`, did)
	if err != nil {
		handleInternalErr("LookupSIPTrunkByDID error 1", err, w)
		return
	}
	defer results.Close()

	var sipRoutingUri string
	var sipRecoveryUri string
	for results.Next() {
		err := results.Scan(&sipRoutingUri, &sipRecoveryUri);
		if err != nil {
			handleInternalErr("LookupSIPTrunkByDID error 2", err, w)
			return
		}

		fmt.Printf("SIP routing URI = %s SIP recovery URI = %s\r\n", sipRoutingUri, sipRecoveryUri);
		// TODO do some health checks here to see if SIP server is actually up..
		isOnline, err := checkSIPServerHealth( sipRoutingUri );
		if err != nil {
			fmt.Printf("failed to verify health of SIP server %s\r\n", sipRoutingUri);
			continue
		}
		if isOnline {

			w.Write([]byte(sipRoutingUri))
			return
		}
		fmt.Printf("routing server %s is offline, checking next server...\r\n", sipRoutingUri);
	}

	// no SIP servers were online try to route to recovery URI
	isOnline, err := checkSIPServerHealth( sipRecoveryUri );
	fmt.Printf("no SIP servers were online. routing to recovery URI\r\n");
	if isOnline {
		w.Write([]byte(sipRecoveryUri));
		return
	}

	fmt.Printf("checked all SIP trunks and found that none were online... error.")
	w.WriteHeader(http.StatusInternalServerError) // send the headers with a 204 response code.
}

func StoreRegistration(w http.ResponseWriter, r *http.Request) {
	domain := r.FormValue("domain")
	//ip := r.FormValue("ip")
	user := r.FormValue("user")
	//contact := r.FormValue("contact")
	now := time.Now()
	workspace, err := getWorkspaceByDomain(domain)
	var expires int

	expires, err = strconv.Atoi(r.FormValue("expires"))

	if err != nil {
		fmt.Printf("could not get expiry.. not setting online\r\n")
		return
	}
	if err != nil {
		fmt.Printf("StoreRegistration error..")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("UPDATE extensions SET `last_registered` = ?, `register_expires`  = ? WHERE `username` = ? AND `workspace_id` = ?")
	if err != nil {
		fmt.Printf("StoreRegistration 2 Could not execute query..")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(now, expires, user, workspace.Id)
	if err != nil {
		handleInternalErr("StoreRegistration 3 Could not execute query", err, w)
		return
	}
}

func CreateSIPReport(w http.ResponseWriter, r *http.Request) {
	//requestto := r.FormValue("requestto")
	//user := r.FormValue("user")
	//domain := r.FormValue("domain")
	callid := r.FormValue("callid")
	status := r.FormValue("status")


	stmt, err := db.Prepare("UPDATE `calls` SET sip_status = ? WHERE sip_call_id = ?")
	if err != nil {
		fmt.Printf("CreateSIPReport 2 Could not execute query..")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	statusAsInt, err := strconv.Atoi(status)
	if err != nil {
		fmt.Printf("CreateSIPReport 3 error in convert..")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(statusAsInt, callid)
	if err != nil {
		handleInternalErr("CreateSIPReport 4 Could not execute query", err, w)
		return
	}
}

func GetBestRTPProxy(w http.ResponseWriter, r *http.Request) {
	//routerip := getQueryVariable(r, "routerip")
	results,err := db.Query(`SELECT rtpproxy_sock, set_id, cpu_pct, cpu_used FROM rtpproxy_sockets`)
	// Execute the query
	if err != nil {
		handleInternalErr("GetBestRTPProxy error 1", err, w)
		return
	}
	defer results.Close()
	var lowestCPU *float64 = nil
	var socketAddrResult string
	for results.Next() {
		var rtpSock string
		var setId int
		var cpuPct float64
		var cpuUsed float64
		err = results.Scan(&rtpSock, &setId, &cpuPct, &cpuUsed)
		if err != nil {
			handleInternalErr("GetBestRTPProxy error 2", err, w)
			return
		}
		if lowestCPU == nil || cpuPct < *lowestCPU{
			lowestCPU = &cpuPct
 			socketAddrResult =rtpSock
		}
	}
	w.Write([]byte(socketAddrResult))
}

func GetSettings(w http.ResponseWriter, r *http.Request) {
	results, err := db.Query("SELECT `aws_access_key_id`, `aws_secret_access_key`, `aws_region`, `google_service_account_json`, `stripe_pub_key`, `stripe_private_key`, `stripe_test_pub_key`, `stripe_test_private_key`, `stripe_mode`, `smtp_host`, `smtp_port`, `smtp_user`, `smtp_password`, `smtp_tls` FROM api_credentials")
	defer results.Close()
	if err == sql.ErrNoRows {
		// no records setup were setup, just return empty
		fmt.Printf("GetSettings no rows found..")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err != nil {
		fmt.Printf("GetSettings error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	settings := Settings{}
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
			fmt.Printf("GetSettings error: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(&settings)
	}
}

func startHTTPServer() {
	settings = &GlobalSettings{ValidateCallerId: false}
	r := mux.NewRouter()
	fmt.Println("starting HTTP server...")
	// Routes consist of a path and a handler function.
	r.HandleFunc("/healthz", healthz).Methods("GET")
	r.HandleFunc("/call/createCall", CreateCall).Methods("POST")
	r.HandleFunc("/call/updateCall", UpdateCall).Methods("POST")
	r.HandleFunc("/call/fetchCall", FetchCall).Methods("GET")
	r.HandleFunc("/call/setSIPCallID", SetSIPCallID).Methods("POST")
	r.HandleFunc("/call/setProviderByIP", SetProviderByIP).Methods("POST")

	r.HandleFunc("/conference/createConference", CreateConference).Methods("POST")

	//debits
	r.HandleFunc("/debit/createDebit", CreateDebit).Methods("POST")
	r.HandleFunc("/debit/createAPIUsageDebit", CreateAPIUsageDebit).Methods("POST")

	//logs
	r.HandleFunc("/debugger/createLog", CreateLog).Methods("POST")
	r.HandleFunc("/debugger/createLogSimple", CreateLogSimple).Methods("POST")

	//fax
	r.HandleFunc("/fax/createFax", CreateFax).Methods("POST")

	//recording
	r.HandleFunc("/recording/createRecording", CreateRecording).Methods("POST")
	r.HandleFunc("/recording/updateRecording", UpdateRecording).Methods("POST")
	r.HandleFunc("/recording/updateRecordingTranscription", UpdateRecordingTranscription).Methods("POST")
	r.HandleFunc("/recording/getRecording", GetRecording).Methods("GET")
	// carrier functions
	r.HandleFunc("/carrier/createSIPReport", CreateSIPReport).Methods("POST")

	r.HandleFunc("/carrier/processRouterFlow", ProcessRouterFlow).Methods("GET")



	// user functions
	r.HandleFunc("/user/verifyCaller", VerifyCaller).Methods("GET")
	r.HandleFunc("/user/verifyCallerByDomain", VerifyCallerByDomain).Methods("GET")
	r.HandleFunc("/user/getUserByDomain", GetUserByDomain).Methods("GET")
	r.HandleFunc("/user/getUserByDID", GetUserByDID).Methods("GET")
	r.HandleFunc("/user/getUserByTrunkSourceIp", GetUserByTrunkSourceIp).Methods("GET")
	r.HandleFunc("/user/getWorkspaceMacros", GetWorkspaceMacros).Methods("GET")
	r.HandleFunc("/user/getDIDNumberData", GetDIDNumberData).Methods("GET")
	r.HandleFunc("/user/getPSTNProviderIP", GetPSTNProviderIP).Methods("GET")
	r.HandleFunc("/user/getPSTNProviderIPForTrunk", GetPSTNProviderIPForTrunk).Methods("GET")
	r.HandleFunc("/user/ipWhitelistLookup", IPWhitelistLookup).Methods("GET")
	r.HandleFunc("/user/getDIDAcceptOption", GetDIDAcceptOption).Methods("GET")
	r.HandleFunc("/user/getDIDAssignedIP", GetDIDAssignedIP).Methods("GET")
	r.HandleFunc("/user/getUserAssignedIP", GetUserAssignedIP).Methods("GET")
	r.HandleFunc("/user/getTrunkAssignedIP", GetTrunkAssignedIP).Methods("GET")
	r.HandleFunc("/user/addPSTNProviderTechPrefix", AddPSTNProviderTechPrefix).Methods("GET")
	r.HandleFunc("/user/getCallerIdToUse", GetCallerIdToUse).Methods("GET")
	r.HandleFunc("/user/getExtensionFlowInfo", GetExtensionFlowInfo).Methods("GET")
	r.HandleFunc("/user/getFlowInfo", GetFlowInfo).Methods("GET")
	r.HandleFunc("/user/getDIDDomain", GetDIDDomain).Methods("GET")
	r.HandleFunc("/user/getCodeFlowInfo", GetCodeFlowInfo).Methods("GET")
	r.HandleFunc("/user/incomingDIDValidation", IncomingDIDValidation).Methods("GET")
	r.HandleFunc("/user/incomingTrunkValidation", IncomingTrunkValidation).Methods("GET")
	r.HandleFunc("/user/lookupSIPTrunkByDID", LookupSIPTrunkByDID).Methods("GET")
	r.HandleFunc("/user/incomingMediaServerValidation", IncomingMediaServerValidation).Methods("GET")
	r.HandleFunc("/user/storeRegistration", StoreRegistration).Methods("POST")

	r.HandleFunc("/user/getSettings", GetSettings).Methods("GET")

	r.HandleFunc("/user/processSIPTrunkCall", ProcessSIPTrunkCall).Methods("GET")


	// Send Admin email
	r.HandleFunc("/admin/sendAdminEmail", SendAdminEmail).Methods("POST")
	r.HandleFunc("/getBestRTPProxy", GetBestRTPProxy).Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	// Bind to a port and pass our router in

	routingHandler:=loggedRouter

	if os.Getenv("USE_LIMIT_MIDDLEWARE") == "on" {
		routingHandler= limitMiddleware(loggedRouter)
	}
	if os.Getenv("USE_TLS") == "on" {
		certPath := os.Getenv("TLS_CERT_PATH")
		keyPath := os.Getenv("TLS_KEY_PATH")


		fmt.Printf("Starting HTTP server with TLS. cert=%s,  key=%s\r\n", certPath, keyPath)
		log.Fatal( http.ListenAndServeTLS(":443", certPath,keyPath, routingHandler))
		//log.Fatal( http.ListenAndServeTLS(":443", certPath,keyPath, nil))
		fmt.Println("started server...")
		return
	}
	fmt.Printf("Starting HTTP server without TLS\r\n")

	httpPort := readEnv("HTTP_PORT", "80")
	fmt.Printf("HTTP port %s\r\n", httpPort)
	log.Fatal(http.ListenAndServe(":" + httpPort, routingHandler))
}

func readEnv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}

func checkIfCarrier(token string) bool {
	return true
}

func limitMiddleware(next http.Handler) http.Handler {
	cli, err := createETCDClient()
	usingEtcd := err == nil
	if !usingEtcd {
		fmt.Println("could not connect to ETCD - rate limits are disabled")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !usingEtcd {
			next.ServeHTTP(w, r)
			return
		}
		var addr string
		requestedAddr := getQueryVariable(r, "addr")
		if requestedAddr != nil {
			addr = *requestedAddr
		} else {
			addr = r.RemoteAddr
		}
		interval := time.Duration(time.Second * 60)
		flushinterval := time.Duration(time.Second * 30)
		carrier := r.Header.Get("X-Lineblocs-Carrier-Auth")

		isCarrier := false

		if carrier != "" {
			isCarrier = checkIfCarrier(carrier)
		}

		// limit for users

		var limit uint64 = 60
		if isCarrier {
			limit = 3600
		}
		rate := limiter.NewRateLimiter(cli, addr, limit, interval, flushinterval)
		rate.ProcessLimits()
		over, err := rate.IsOverLimit()
		if err != nil {
			handleInternalErr("rate_limit error", err, w)
			return
		}
		if over {
			fmt.Println("rate limit exhausted..")
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("starting API...")
	var err error
	servers, err := lineblocs.CreateMediaServers()

	data = &ServerData{
		mu:      sync.RWMutex{},
		servers: servers}

	if err != nil {
		panic(err)
	}
	db, err = lineblocs.CreateDBConn()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		startHTTPServer()
		wg.Done()
	}()
	wg.Wait()
}
