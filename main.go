package main

import (
    "net/http"
	"log"
	"os"
	"time"
	"strconv"
	"math"
	"net"
	"strings"
	"mime/multipart"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"reflect"
	"fmt"
	"database/sql"
    "encoding/json"
	_ "github.com/go-sql-driver/mysql"
	guuid "github.com/google/uuid"
	libphonenumber "github.com/ttacon/libphonenumber"
)

type Call struct {
  From string `json:"from"`
  To string `json:"to"`
  Status string `json:"status"`
  Direction string `json:"direction"`
  Duration string `json:"duration"`
  UserId int `json:"user_id"`
  WorkspaceId int  `json:"workspace_id"`
  APIId string `json:"api_id"`
}
type CallUpdateReq struct {
  CallId int `json:"call_id"`
  Status string `json:"status"`
}
type RecordingTranscriptionReq struct {
	Id string `json:"id"`
  Ready bool `json:"ready"`
  Text string `json:"text"`
}
type Conference struct {
  Name string `json:"name"`
  WorkspaceId int `json:"workspace_id"`
  APIId string `json:"api_id"`
}
type DebitCreateReq struct {
  UserId int `json:"user_id"`
  WorkspaceId int `json:"workspace_id"`

  Source string `json:"source"`
  Number string `json:"number"`
  Type string `json:"type"`
  Seconds float64 `json:"seconds"`
}


type CallRate struct {
	CallRate float64
}


type DebitAPIParams struct {
	Length int `json:"length"`	
	RecordingLength float64 `json:"recording_length"`	
}
type DebitAPICreateReq struct {
  UserId int `json:"user_id"`
  WorkspaceId int `json:"workspace_id"`
  Type string `json:"type"`
  Source string `json:"source"`
  Params DebitAPIParams `json:"params"`
}

type LogCreateReq struct {
  UserId int `json:"user_id"`
  WorkspaceId int `json:"workspace_id"`
  Title string `json:"title"`
  Report string `json:"report"`
  FlowId int `json:"flow_id"`
  Level *string `json:"report"`
  From *string `json:"from"`
  To *string `json:"to"`
}
type LogSimpleCreateReq struct {
  Type string `json:"type"`
  Level *string `json:"level"`
}
type Fax struct {
  UserId int `json:"user_id"`
  WorkspaceId int `json:"workspace_id"`
  CallId int `json:"call_id"`
  Uri string `json:"uri"`
  APIId string `json:"api_id"`
}

type Recording struct {
  UserId int `json:"user_id"`
  CallId int `json:"call_id"`
  WorkspaceId int `json:"workspace_id"`
  APIId string `json:"api_id"`
  Tags *[]string `json:"tags"`

}

type VerifyNumber struct {
	Valid bool `json:"valid"`
}




type LogRoutine struct {
  UserId int
  WorkspaceId int
  Title string
  Report string
  FlowId int
  Level string
  From string
  To string
}
type User struct {
  Id int
  Username string
  FirstName string
  LastName string
  Email string
}

type Workspace struct {
  Id int `json:"id"`
  CreatorId int `json:"creator_id"`
  Name string `json:"name"`
  BYOEnabled bool `json:"byo_enabled"`
  IPWhitelistDisabled bool `json:"ip_whitelist_disabled"`
  OutboundMacroId int `json:"outbound_macro_id"`
}

type WorkspaceParam struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
type WorkspaceFullInfo struct {
	Workspace *Workspace `json:"workspace"`
	WorkspaceName string `json:"workspace_name"`
	WorkspaceId int `json:"workspace_id"`
	WorkspaceParams *[]WorkspaceParam `json:"workspace_params"`
  	OutboundMacroId int `json:"outbound_macro_id"`
}
type WorkspaceDIDInfo struct {
  WorkspaceId int `json:"workspace_id"`
  Number string `json:"number"`
  FlowJSON string `json:"flow_json"`
  WorkspaceName string `json:"workspace_name"`
  Name string `json:"name"`
  Plan string `json:"plan"`
  BYOEnabled bool `json:"byo_enabled"`
  IPWhitelistDisabled bool `json:"ip_whitelist_disabled"`
  OutboundMacroId int `json:"outbound_macro_id"`
  CreatorId int `json:"creator_id"`
  APIToken string `json:"api_token"`
  APISecret string `json:"api_secret"`
  WorkspaceParams *[]WorkspaceParam `json:"workspace_params"`
}
type WorkspacePSTNInfo struct {
  IPAddr string `json:"ip_addr"`
  DID string `json:"did"`
}
type CallerIDInfo struct {
  CallerID string `json:"caller_id"`
}
type ExtensionFlowInfo struct {
  CallerID string `json:"caller_id"`
  WorkspaceId int `json:"workspace_id"`
  FlowJSON string `json:"flow_json"`
  Username string `json:"username"`
  Name string `json:"name"`
  WorkspaceName  string `json:"workspace_name"`
  Plan string `json:"plan"`
  CreatorId int `json:"creator_id"`
  Id int `json:"id"`
  APIToken string `json:"api_token"`
  APISecret string `json:"api_secret"`
  WorkspaceParams *[]WorkspaceParam `json:"workspace_params"`
  FreeTrialStatus string `json:"workspace_params"`
}

type CodeFlowInfo struct {
  WorkspaceId int `json:"workspace_id"`
  Code string `json:"code"`
  FlowJSON string `json:"flow_json"`
  Name string `json:"name"`
  WorkspaceName  string `json:"workspace_name"`
  Plan string `json:"plan"`
  CreatorId int `json:"creator_id"`
  Id int `json:"id"`
  APIToken string `json:"api_token"`
  APISecret string `json:"api_secret"`
  FreeTrialStatus string `json:"workspace_params"`
  FoundCode bool `json:"found_code"`
}


type MacroFunction struct {
	Title string `json:"title"`
	Code string `json:"code"`
	CompiledCode string `json:"compiled_code"`
}
type MediaServer struct {
	IpAddress string `json:"ip_address"`
	PrivateIpAddress string `json:"private_ip_address"`
}
var db* sql.DB;

func createAPIID(prefix string) string {
	id := guuid.New()
	return prefix + "-" + id.String()
}
func lookupBestCallRate(number string, typeRate string) *CallRate {
	return &CallRate{ CallRate: 9.99 };
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
func getUserFromDB(id int) (*User, error) {
	var userId int
	var username string
	var fname string
	var lname string
	var email string
	fmt.Printf("looking up user %d\r\n", id)
	row := db.QueryRow(`SELECT id, username, first_name, last_name, email FROM users WHERE id=?`, id)

	err := row.Scan(&userId, &username, &fname, &lname,  &email)
	if ( err == sql.ErrNoRows ) {  //create conference
		return nil, err
	}
	if ( err != nil ) {  //another error
		return nil, err
	}

	return &User{Id: userId, Username: username, FirstName: fname, LastName: lname, Email: email}, nil
}
func getWorkspaceFromDB(id int) (*Workspace, error) {
	var workspaceId int
	var name string
	var creatorId int
	var outboundMacroId sql.NullInt64
	row := db.QueryRow(`SELECT id, name, creator_id, outbound_macro_id FROM workspaces WHERE id=?`, id)

	err := row.Scan(&workspaceId, &name, &creatorId, &outboundMacroId)
	if ( err == sql.ErrNoRows ) {  //create conference
		return nil, err
	}
	if ( err != nil ) {  //another error
		return nil, err
	}
    if reflect.TypeOf(outboundMacroId) == nil {
		return &Workspace{Id: workspaceId, Name: name, CreatorId: creatorId}, nil
	}
	return &Workspace{Id: workspaceId, Name: name, CreatorId: creatorId, OutboundMacroId: int(outboundMacroId.Int64)}, nil
}
func getWorkspaceByDomain(domain string) (*Workspace, error) {
	var workspaceId int
	var name string
	var byo bool
	var ipWhitelist bool
	s := strings.Split(domain, ".")
	workspaceName := s[0]
	row := db.QueryRow("SELECT id, name, byo_enabled, ip_whitelist_disabled FROM workspaces WHERE name=?", workspaceName)

	err := row.Scan(&workspaceId, &name, &byo, &ipWhitelist)
	if ( err == sql.ErrNoRows ) {  //create conference
		return nil, err
	}
	return &Workspace{Id: workspaceId, Name: name, BYOEnabled: byo, IPWhitelistDisabled: ipWhitelist}, nil
}

func getWorkspaceParams(workspaceId int) (*[]WorkspaceParam, error) {
	// Execute the query
	results, err := db.Query("SELECT `key`, `value` FROM workspace_params WHERE `workspace_id` = ?", workspaceId)
    if err != nil {
		return nil, err;
	}
	params := []WorkspaceParam{};

    for results.Next() {
		param := WorkspaceParam{};
        // for each row, scan the result into our tag composite object
        err = results.Scan(&param.Key, &param.Value)
        if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return &params, nil;
}

func getUserByDomain(domain string) (*WorkspaceFullInfo, error) {
	workspace, err := getWorkspaceByDomain( domain )
	if err != nil {
		return nil, err
	}

	// Execute the query
	params, err  := getWorkspaceParams(workspace.Id)
    if err != nil {
		return nil, err;
	}
	full := &WorkspaceFullInfo{ Workspace: workspace, 
		WorkspaceParams: params,
		WorkspaceName: workspace.Name,
		WorkspaceId: workspace.Id,
		OutboundMacroId: workspace.OutboundMacroId	};

	return full, nil
}

func getRecordingFromDB(id int) *Recording {
	var apiId string
	row := db.QueryRow("SELECT api_id FROM recordings WHERE id=?", id)

	err := row.Scan(&apiId)
	if ( err == sql.ErrNoRows ) {  //create conference
		return nil
	}
	return &Recording{APIId: apiId}
}



// TODO
func sendLogRoutineEmail(log* LogRoutine, user* User, workspace* Workspace) error {
	return nil
}

func startLogRoutine(log* LogRoutine) (*string, error) {
	var user* User;
	var workspace* Workspace;

    user, err := getUserFromDB(log.UserId)
	if err != nil {
		fmt.Printf("could not get user..")
		return nil, err
	}

	workspace, err = getWorkspaceFromDB(log.WorkspaceId)
	if err != nil {
		fmt.Printf("could not get workspace..")
		return nil, err
	}

	apiId := createAPIID("log")
	stmt, err := db.Prepare("INSERT INTO debugger_logs (`from`, `to`, `title`, `report`, `workspace_id`, `level`, `api_id`) VALUES ( ?, ?, ?, ?, ?, ?, ? )")

	if err != nil {
		fmt.Printf("could not prepare query..")
		return nil, err
	}
	res, err := stmt.Exec(log.From, log.To, log.Title, log.Report, workspace.Id, log.Level, apiId)
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

	err = sendLogRoutineEmail(log, user, workspace)
	if err != nil {
		fmt.Printf("could not send email..")
		return nil, err
	}


	return &logIdStr, err
}
//todo
func checkRouteMatches(from string, to string, prefix string, prepend string, match string) (bool) {
	return true
}
func shouldUseProviderNext(name string, ipPrivate string) (bool, error) {
	return true, nil
}
func checkCIDRMatch(sourceIp string, fullIp string) (bool, error) {
	_, net1, err :=  net.ParseCIDR(sourceIp + "/32")
	if err != nil {
		return false, err
	}
	_, net2, err :=  net.ParseCIDR(fullIp)
	if err != nil {
		return false, err
	}

	return net2.Contains(net1.IP), nil
}
func checkPSTNIPWhitelist(did string, sourceIp string) (bool, error) {
	results, err := db.Query(`SELECT 
	sip_providers_whitelist_ips.ip_address, 
	sip_providers_whitelist_ips.ip_address_range
	FROM sip_providers_whitelist_ips
	INNER JOIN sip_providers ON sip_providers.id = sip_providers_whitelist_ips.provider_id`)
    if err != nil {
		return false, err
	}
    for results.Next() {
		var ipAddr string
		var ipAddrRange string
		err = results.Scan(&ipAddr, &ipAddrRange)
		if err != nil {
			return false, err
		}
		fullIp := ipAddr + ipAddrRange
		match, err := checkCIDRMatch(sourceIp, fullIp) 
		if err != nil {
			return false, err
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
	byo_carriers_ips.range,
	FROM byo_carriers_ips
	INNER JOIN byo_carriers ON byo_carriers.id = byo_carriers_ips.carrier_id`)
    if err != nil {
		return false, err
	}
    for results.Next() {
		var ipAddr string
		var ipAddrRange string
		err = results.Scan(&ipAddr, ipAddrRange)
		if err != nil {
			return false, err
		}
		fullIp := ipAddr + ipAddrRange
		match, err := checkCIDRMatch(sourceIp, fullIp) 
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func finishValidation(number string, didWorkspaceId string) (bool,error) {
	num, err := libphonenumber.Parse(number, "US")
	if err != nil {
		return false, err
	}
	formattedNum := libphonenumber.Format(num, libphonenumber.E164)
	row := db.QueryRow("SELECT id FROM `blocked_numbers` WHERE `workspace_id` = ? AND `number` = ?", didWorkspaceId, formattedNum)
	var id string
	err = row.Scan(&id);
	if ( err == sql.ErrNoRows ) {  //create conference
		return true,nil
	}
	if err != nil {
		return false,err
	}
	return false,nil
}
func checkFreeTrialStatus(plan string, started time.Time) string {
	if plan  == "trial" {
		now := time.Now()
		//make configurable
		expireDays := 10
		expireHours := expireDays * 24
		started.Add(time.Hour * time.Duration(expireHours))
		if started.After( now ) {
			return "expired";
		}
		return "pending-expiry";
	}
	return "not-applicable";
}
func checkIsMakingOutboundCallFirstTime(call Call) {
	var id string
	row := db.QueryRow("SELECT id FROM `calls` WHERE `workspace_id` = ? AND `from` LIKE '?%s' AND `direction = 'outbound'", call.WorkspaceId, call.From, call.Direction)
	err := row.Scan(&id);
	if ( err != sql.ErrNoRows ) {  //create conference
		// all ok
		return
	}
	//send notification
	user, err := getUserFromDB(call.UserId)
	if err != nil {
		panic(err)
	}
	body := `A call was made to ` + call.To + ` for the first time on your account.`;
	sendEmail(user, "First call to destination country", body)
}
func sendEmail(user *User, subject string, body string) {
}
func someLoadBalancingLogic() (*MediaServer,error) {
	results, err := db.Query("SELECT ip_address,private_ip_address FROM media_servers");
    if err != nil {
		return nil,err
	}
    for results.Next() {
		value := MediaServer{};
		err = results.Scan(&value.IpAddress,&value.PrivateIpAddress);
		if err != nil {
			return nil,err
		}
		return &value,nil
	}
	return nil,nil
}
func doVerifyCaller(workspaceId int, number string) (bool, error) {
	var workspace* Workspace;
	workspace, err := getWorkspaceFromDB(workspaceId)
	if err != nil {
		return false, err
	}

	num, err := libphonenumber.Parse(number, "US")
	if err != nil {
		return false, err
	}
	formattedNum := libphonenumber.Format(num, libphonenumber.E164)
	fmt.Printf("looking up number %s", formattedNum)
	var id string
	row := db.QueryRow("SELECT id FROM `did_numbers` WHERE `number` = ? AND `workspace_id` = ?", formattedNum, workspace.Id)
	err = row.Scan(&id);
	if ( err != sql.ErrNoRows ) {  //create conference
		return true, nil
	}
	return false, nil
}

func getQueryVariable(r *http.Request, key string) *string {
	vals := r.URL.Query() // Returns a url.Values, which is a map[string][]string
	results, ok := vals[key] // Note type, not ID. ID wasn't specified anywhere.
	var value *string
	if ok {
		if len(results) >= 1 {
			value = &results[0] // The first `?type=model`
		}
	}
	return value
}
// todo
func uploadS3(folder string, id string, file multipart.File) (error) {
	return nil
}

// todo
func createS3URL(folder string, id string) string {
	return ""
}

func NoContent(w http.ResponseWriter, r *http.Request) {
  // Set up any headers you want here.
  w.WriteHeader(http.StatusNoContent) // send the headers with a 204 response code.
}
func toCents(dollars float64) int {
	result := dollars * 100
	return int( result )
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
	call.APIId = createAPIID("call")

	if call.Direction == "outbound" {
		//check if this is the first time we are making a call to this destination
		go checkIsMakingOutboundCallFirstTime(call)
	}

  // perform a db.Query insert
	stmt, err := db.Prepare("INSERT INTO calls (`from`, `to`, `status`, `direction`, `duration`, `user_id`, `workspace_id`, `started_at`, `api_id`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? )")
	if err != nil {
		handleInternalErr("CreateCall Could not execute query..", err, w);
		return 
	}
	defer stmt.Close()
	fmt.Printf("CreateCall args from=%s, to=%s, status=%s, direction=%s, user_id=%s, workspace_id=%s, started=%s",
		call.From, call.To, call.Status, call.Direction, call.UserId, call.WorkspaceId, now, call.APIId)

	res, err := stmt.Exec(call.From, call.To, call.Status, call.Direction, "8", call.UserId, call.WorkspaceId, now, call.APIId)

		
		if err != nil {
			handleInternalErr("CreateCall Could not execute query", err, w)
		return
	}                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              

	callId, err := res.LastInsertId()
	if err != nil {
			handleInternalErr("CreateCall Could not execute query..", err, w);
		return
	}                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              
    w.Header().Set("X-Call-ID", strconv.FormatInt(callId, 10))
  	json.NewEncoder(w).Encode(&call)
}

func UpdateCall(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  var update CallUpdateReq
   err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		handleInternalErr("UpdateCall Could not decode JSON", err, w)

		return 
	}
	
	
	if ( update.Status == "ended" ) {
		// perform a db.Query insert
		stmt, err := db.Prepare("UPDATE calls SET `status` = ?, `ended_at` = ? WHERE `id` = ?")
		if err != nil {
			fmt.Printf("CreateCall Could not execute query..");
			fmt.Println(err)
  			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
		defer stmt.Close()
		endedAt := time.Now()
		_, err = stmt.Exec(update.Status, endedAt, update.CallId)
		if err != nil {
			handleInternalErr("UpdateCall Could not execute query", err, w)
			return
		}
	}                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              

  	w.WriteHeader(http.StatusNoContent)
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
	err = row.Scan(&id, &name);
	if ( err == sql.ErrNoRows ) {  //create conference
		conference.APIId = createAPIID("conf")
  		// perform a db.Query insert
		stmt, err := db.Prepare("INSERT INTO conferences (`name`, `workspace_id`, `api_id`) VALUES ( ?, ?, ? )");
		if err != nil {
			handleInternalErr("CreateConference 3 Could not execute query..", err, w);
			return 
		}
		defer stmt.Close()
		res, err := stmt.Exec(conference.Name, conference.WorkspaceId, conference.APIId)
		
		if err != nil {
			handleInternalErr("CreateConference 4 Could not execute query", err, w)
			return
		}                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              
		conferenceId, err := res.LastInsertId()
		if err != nil {
			handleInternalErr("CreateConference 5 Could not get ID..", err, w);
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
	rate := lookupBestCallRate(debitReq.Number, debitReq.Type)
	if rate == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	minutes := math.Floor(debitReq.Seconds / 60)
	dollars := minutes * rate.CallRate
	cents := toCents( dollars )
	stmt, err := db.Prepare("INSERT INTO user_debits (`user_id, `cents`, `source`) VALUES ( ?, ?, ? )");
	if err != nil {
		handleInternalErr("CreateDebit Could not execute query..", err, w);
		return 
	}
	_, err = stmt.Exec(debitReq.UserId, cents, debitReq.Source)
	if err != nil {
		handleInternalErr("CreateDebit Could not execute query..", err, w);
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
	if debitReq.Type == "TTS" {
		dollars := calculateTTSCosts(debitReq.Params.Length)
		cents := toCents( dollars )
		source := fmt.Sprintf("API usage - %s", debitReq.Type);
		stmt, err := db.Prepare("INSERT INTO user_debits (`user_id, `cents`, `source`) VALUES ( ?, ?, ? )");
		if err != nil {
			handleInternalErr("CreateDebit Could not execute query..", err, w);
			return 
		}
		_, err = stmt.Exec(debitReq.UserId, cents, source)
		if err != nil {
			handleInternalErr("CreateAPIUsageDebit Could not execute query..", err, w);
			return 
		}
	} else if debitReq.Type == "STT" {
		dollars := calculateSTTCosts(debitReq.Params.RecordingLength)
		cents := toCents( dollars )
		source := fmt.Sprintf("API usage - %s", debitReq.Type);
		stmt, err := db.Prepare("INSERT INTO user_debits (`user_id, `cents`, `source`) VALUES ( ?, ?, ? )");
		if err != nil {
			handleInternalErr("CreateDebit Could not execute query..", err, w);
			return 
		}
		_, err = stmt.Exec(debitReq.UserId, cents, source)
		if err != nil {
			handleInternalErr("CreateAPIUsageDebit Could not execute query..", err, w);
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
		handleInternalErr("CreateLog Could not decode JSON", err, w)
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
	var log* LogRoutine = &LogRoutine{ From: from,
		To: to,
		Level: level,
		Title:logReq.Title,
		Report:logReq.Report,
		FlowId:logReq.FlowId,
		UserId: logReq.UserId,
		WorkspaceId: logReq.WorkspaceId }
	_, err = startLogRoutine( log )
	if err != nil {
		handleInternalErr("CreateLog log routine error", err, w)
		return 
	}
}

func CreateLogSimple(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
	logType := r.FormValue("type")
	level := r.FormValue("level")
	domain := r.FormValue("domain")
	workspace, err := getWorkspaceByDomain( domain )
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
			title = "Caller ID Verify failed..";
			report = "Caller ID Verify failed..";
		}
	var log* LogRoutine = &LogRoutine{ 
		From: "",
		To: "",
		Level: level,
		Title:title,
		Report:report,
		UserId: workspace.CreatorId,
		WorkspaceId: workspace.Id }
	_, err = startLogRoutine( log )
	if err != nil {
		handleInternalErr("CreateLog log routine error", err, w)
		return 
	}

}

func CreateFax(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var fax* Fax
	file, handler, err := r.FormFile("file")
	if err != nil {
		handleInternalErr("CreateFax error occured", err, w)
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
	err = uploadS3("faxes", apiId, file)
	if err != nil {
		handleInternalErr("CreateFax error occured", err, w)
		return
	}
	uri := createS3URL( "faxes", apiId )


	stmt, err := db.Prepare("INSERT INTO faxes (`uri`, `size`, `name`, `user_id`, `call_id`, `workspace_id`, `api_id`) VALUES ( ?, ?, ?, ?, ?, ?, ? )");
	if err != nil {
		handleInternalErr("CreateFax error occured", err, w)
		return
	}
	res, err := stmt.Exec(uri, handler.Size, name, userId, callId, workspaceId, apiId )
	if err != nil {
		handleInternalErr("CreateFax error occured", err, w)
		return
	}

	faxId, err := res.LastInsertId()
	if err != nil {
		handleInternalErr("CreateFax error occured", err, w)
		return
	}

	fax = &Fax{ UserId: userIdInt, WorkspaceId: workspaceIdInt, CallId: callIdInt, Uri: uri }
	w.Header().Set("X-Fax-ID", strconv.FormatInt(faxId, 10))
  	json.NewEncoder(w).Encode(fax)
}

func CreateRecording(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  var recording Recording

   err := json.NewDecoder(r.Body).Decode(&recording)
	if err != nil {
		handleInternalErr("CreateCall Could not decode JSON", err, w)
		return 
	}
	recording.APIId = createAPIID("rec")

  // perform a db.Query insert
	stmt, err := db.Prepare("INSERT INTO recordings (`user_id`, `call_id`, `workspace_id`, `status`, `name`, `uri`, `tag`, `api_id`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)");
	if err != nil {
		handleInternalErr("CreateRecording error.", err, w);
		return 
	}
	res, err := stmt.Exec(recording.UserId, recording.CallId, recording.WorkspaceId, "started", "", "", "", recording.APIId)
	if err != nil {
		handleInternalErr("CreateRecording error.", err, w);
		return 
	}
	recId, err := res.LastInsertId()
	if err != nil {
		handleInternalErr("CreateRecording error.", err, w);
		return
	}
	if recording.Tags != nil {
		for _, v := range *recording.Tags {
			fmt.Printf("adding tag to recording %s\r\n", v)
			stmt, err := db.Prepare("INSERT INTO recording_tags (`recording_id`, `tag`) VALUES (?, ?)");
			if err != nil {
				handleInternalErr("CreateRecording error.", err, w);
			}

			res, err = stmt.Exec(recId, v)
			if err != nil {
				handleInternalErr("CreateRecording error.", err, w);
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
	recordingId := r.FormValue("recording_id")
	recordingIdInt, err := strconv.Atoi(recordingId)
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

	apiId := createAPIID("rec")
	err = uploadS3("recordings", apiId, file)
	if err != nil {
		handleInternalErr("UpdateRecording error occured", err, w)
		return
	}
	uri := createS3URL( "recordings", apiId)
	stmt, err := db.Prepare("UPDATE `recordings` SET `status` = ?, `uri` = ?, `size` = ? WHERE `id` = ?")
	if err != nil {
		handleInternalErr("UpdateRecording error occured", err, w)
		return
	}
	_, err = stmt.Exec(status, uri, handler.Size, recordingIdInt)
	if err != nil {
		handleInternalErr("UpdateRecording error occured", err, w)
		return
	}

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
	_, err  = stmt.Exec("1", update.Text, update.Id)
	if err != nil {
		handleInternalErr("UpdateCall Could not execute query", err, w)
		return
	}

  	w.WriteHeader(http.StatusNoContent)
}

func VerifyCaller(w http.ResponseWriter, r *http.Request) {
	workspaceId := getQueryVariable(r, "workspaceId")
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
	result := VerifyNumber{ Valid: valid }
  	json.NewEncoder(w).Encode(&result)
}
func VerifyCallerByDomain(w http.ResponseWriter, r *http.Request) {
	domain := getQueryVariable(r, "domain")
	number := getQueryVariable(r, "number")

	workspace, err := getWorkspaceByDomain(*domain)
	if err != nil {
		handleInternalErr("VerifyCallerByDomain error occured", err, w)
		return
	}
	valid, err := doVerifyCaller(workspace.Id, *number)
	if err != nil {
		handleInternalErr("VerifyCaller error occured", err, w)
		return
	}

	result := VerifyNumber{ Valid:  valid }
  	json.NewEncoder(w).Encode(&result)
}
func GetUserAssignedIP(w http.ResponseWriter, r *http.Request) {
	server, err := someLoadBalancingLogic()
	if err != nil {
		handleInternalErr("GetUserAssignedIP error occured", err, w)
		return
	}
	if server == nil {
		handleInternalErr("GetUserAssignedIP could not get server", err, w)
	}
	w.Write([]byte(server.PrivateIpAddress));
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
func GetWorkspaceMacros(w http.ResponseWriter, r *http.Request) {
	workspace := getQueryVariable(r, "workspace")
	// Execute the query
	results, err := db.Query("SELECT title, code, compiled_code FROM macro_functions WHERE `workspace_id` = ?", workspace)
    if err != nil {
		handleInternalErr("GetWorkspaceMacros error", err, w)
		return
	}
	values := []MacroFunction{};


    for results.Next() {
		value := MacroFunction{};
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
	var info WorkspaceDIDInfo;
	// Execute the query
	row := db.QueryRow(`SELECT flows.workspace_id, flows.flow_json, did_numbers.number, workspaces.name, workspaces.name AS workspace_name, 
        users.plan,
        workspaces.byo_enabled,
        workspaces.creator_id,
        workspaces.id AS workspace_id,
        workspaces.api_token,
		workspaces.api_secret 
		FROM workspaces
		INNER JOIN did_numbers ON did_numbers.workspace_id = workspaces.id	
		INNER JOIN flows ON flows.workspace_id = flows.id	
		INNER JOIN users ON users.id = workspaces.creator_id
		WHERE did_numbers.api_number = ?	
		`, number);
	err := row.Scan(
			&info.WorkspaceId,
			&info.FlowJSON,
			&info.Number,
			&info.Name,
			&info.WorkspaceName,
			&info.Plan,
			&info.BYOEnabled,
			&info.CreatorId,
			&info.APIToken,
			&info.APISecret )
	if ( err == nil && err != sql.ErrNoRows ) {  
		params, err := getWorkspaceParams(info.WorkspaceId)
		if err != nil {
			handleInternalErr("GetDIDNumberData error", err, w)
			return
		}

		info.WorkspaceParams = params
		json.NewEncoder(w).Encode(&info)
	}
	// Execute the query
	row = db.QueryRow(`SELECT 
		flows.workspace_id, 
		flows.flow_json, 
		byo_did_numbers.number, 
		workspaces.name, 
		workspaces.name AS workspace_name, 
        users.plan,
        workspaces.byo_enabled,
        workspaces.creator_id,
        workspaces.api_token,
		workspaces.api_secret FROM workspaces
		INNER JOIN byo_did_numbers ON byo_did_numbers.workspace_id = workspaces.id	
		INNER JOIN flows ON flows.workspace_id = flows.id	
		INNER JOIN users ON users.id = workspaces.creator_id
		WHERE byo_did_numbers.number = ?	
		`, number);
	err = row.Scan(
			&info.WorkspaceId,
			&info.FlowJSON,
			&info.Number,

			&info.Name,
			&info.WorkspaceName,
			&info.Plan,

			&info.BYOEnabled,
			&info.CreatorId,
			&info.APIToken,

			&info.APISecret )
	if ( err == nil && err != sql.ErrNoRows ) {  
		params, err := getWorkspaceParams(info.WorkspaceId)
		if err != nil {
			handleInternalErr("GetDIDNumberData error", err, w)
			return
		}

		info.WorkspaceParams = params
		json.NewEncoder(w).Encode(&info)
	}

	if err != nil {
		handleInternalErr("GetDIDNumberData error", err, w)
		return
	}

}
func GetPSTNProviderIP(w http.ResponseWriter, r *http.Request) {
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
		results, err := db.Query(`SELECT byo_carriers.name, byo_carriers.ip_address, users.ip_private,
		byo_carriers_routes.prefix, byo_carriers_routes.prepend, byo_carriers_routes.match
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

	    for results.Next() {
			var name string
			var ip string
			var ipPrivate string
			var prefix string
			var prepend string
			var match string
			err = results.Scan(&name, &ip, &ipPrivate, &prefix, &prepend, &match)
			if err != nil {
				handleInternalErr("GetPSTNProviderIP error", err, w)
				return
			}
			if checkRouteMatches(*from, *to, prefix, prepend, match) {
				var number string
				number = prepend + *to
				info := &WorkspacePSTNInfo{ IPAddr: ipPrivate, DID: number }
				json.NewEncoder(w).Encode(&info)
				return
			}
		}
		//w.WriteHeader(http.StatusNotFound)
		return
	}
	results, err := db.Query(`SELECT sip_providers.name, sip_providers_hosts.ip_address
		FROM sip_providers
		INNER JOIN sip_providers_hosts ON sip_providers_hosts.provider_id = sip_providers.id
		WHERE sip_providers.type_of_provider = 'outbound'
		`)
	if err != nil {
		handleInternalErr("GetPSTNProviderIP error", err, w)
		return
	}
	for results.Next() {
		var name string
		var ipAddr string
		err = results.Scan(&name, &ipAddr)
		if err != nil {
			handleInternalErr("GetPSTNProviderIP error", err, w)
			return
		}
		result, err := shouldUseProviderNext(name, ipAddr)
		if err != nil {
			handleInternalErr("GetPSTNProviderIP error", err, w)
			return
		}

		if result {
			var number string
			number = *to
			info := &WorkspacePSTNInfo{ IPAddr: ipAddr, DID: number }
			json.NewEncoder(w).Encode(&info)
			return
		}
	}


}
func IPWhitelistLookup(w http.ResponseWriter, r *http.Request) {
	source := getQueryVariable(r, "ip")
	domain := getQueryVariable(r, "domain")
	workspace, err := getWorkspaceByDomain( *domain )
	if err != nil {
		handleInternalErr("IPWhitelistLookup error occured", err, w)
		return
	}
	results, err := db.Query("SELECT ip, range FROM ip_whitelist WHERE `workspace_id` = ?", workspace.Id)
    if err != nil {
		handleInternalErr("IPWhitelistLookup error", err, w)
		return
	}
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
		fullIp := ip + ipRange
		match,err := checkCIDRMatch(*source, fullIp) 
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
		w.Write([]byte(action));
		return
	}
	if ( err != sql.ErrNoRows ) {  //create conference
		handleInternalErr("GetDIDAcceptOption error occured", err, w)
		return
	}

	row = db.QueryRow(`SELECT did_action FROM byo_did_numbers WHERE did_numbers.number = ?`, did)
	err = row.Scan(&action)
	if ( err != sql.ErrNoRows ) {  //create conference
		handleInternalErr("GetDIDAcceptOption error occured", err, w)
		return
	}
	if err != nil {
		handleInternalErr("GetDIDAcceptOption error occured", err, w)
		return
	}
	w.Write([]byte(action));
}
func GetDIDAssignedIP(w http.ResponseWriter, r *http.Request) {
	server, err := someLoadBalancingLogic()
	if err != nil {
		handleInternalErr("GetUserAssignedIP error occured", err, w)
		return
	}
	if server == nil {
		handleInternalErr("GetUserAssignedIP could not get server", err, w)
	}
	w.Write([]byte(server.PrivateIpAddress));
}
func GetCallerIdToUse(w http.ResponseWriter, r *http.Request) {
	domain := getQueryVariable(r, "domain")
	extension := getQueryVariable(r, "domain")
	workspace, err := getWorkspaceByDomain(*domain)
	if err != nil {
		handleInternalErr("GetCallerIdToUse error", err, w)
		return
	}

	var callerId string
	row := db.QueryRow("SELECT caller_id FROM extensions WHERE workspace_id=? AND username = ?", workspace.Id, extension)

	err = row.Scan(&callerId)
	if ( err == sql.ErrNoRows ) {  //create conference
		w.WriteHeader(http.StatusNotFound)
		return
	}
	info := &CallerIDInfo{ CallerID: callerId }
	json.NewEncoder(w).Encode(&info)
}

func AddPSTNProviderTechPrefix(w http.ResponseWriter, r *http.Request) {

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
	row := db.QueryRow(`SELECT 
		flows.workspace_id, 
		flows.flow_json, 
		extensions.username, 
		workspaces.name,
		workspaces.name AS workspace_name, 
        users.plan,
        workspaces.creator_id,
        workspaces.id AS workspace_id,
        workspaces.api_token,
		workspaces.api_secret
		users.free_trial_started
		FROM workspaces
		INNER JOIN extensions ON extensions.workspace_id = workspaces.id
		INNER JOIN flows ON flows.workspace_id = workspaces.id
		INNER JOIN users ON users.id = workspaces.creator_id
		WHERE extensions.username = ? 
		AND extensions.workspace_id = ? 
		`, extension, workspaceId)
	err := row.Scan(&info.WorkspaceId, &info.FlowJSON, &info.Username, &info.Name, &info.WorkspaceName, &info.Plan,
			&info.CreatorId, &info.Id, &info.APIToken, &info.APISecret, &trialStartedTime)
	if ( err == sql.ErrNoRows ) {  //create conference
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
	if ( err == sql.ErrNoRows ) {  //create conference
		info.FoundCode = false
		json.NewEncoder(w).Encode(&info)
		return
	}

	if ( err != nil ) {
		handleInternalErr("GetCodeFlowInfo error", err, w)
		return

	}
	info.FreeTrialStatus = checkFreeTrialStatus(info.Plan, trialStartedTime)
	info.FoundCode = true
	json.NewEncoder(w).Encode(&info)

}
func IncomingPSTNValidation(w http.ResponseWriter, r *http.Request) {
	did := getQueryVariable(r, "did")
	number := getQueryVariable(r, "number")
	source := getQueryVariable(r, "source")
	// Execute the query
	row := db.QueryRow(`SELECT did_numbers.number, did_numbers.api_number, did_numbers.workspace_id FROM did_numbers WHERE did_numbers.api_number = ?`, did)
	var didNumber string
	var didApiNumber string
	var didWorkspaceId string
	err := row.Scan(&didNumber,
			&didApiNumber,
			&didWorkspaceId)
	if ( err != sql.ErrNoRows ) {  //create conference
		match, err := checkPSTNIPWhitelist(*did, *source) 
		if err != nil {
			handleInternalErr("IncomingPSTNValidation error 1", err, w)
			return
		}

		if match {
			fmt.Printf("Matched incoming DID..")
			valid, err := finishValidation(*number, didWorkspaceId)
			if err != nil {
				handleInternalErr("IncomingPSTNValidation error 2 valid", err, w)
				return
			}
			if !valid {
				fmt.Printf("number not valid..")
				w.WriteHeader(http.StatusInternalServerError) // send the headers with a 204 response code.
				return
			}
		    w.WriteHeader(http.StatusNoContent) // send the headers with a 204 response code.
			return
		}
	}

	if err != nil {
		handleInternalErr("IncomingPSTNValidation error 3", err, w)
		return
	}

	//check BYO DIDs
	row = db.QueryRow(`SELECT byo_did_numbers.number, byo_did_numbers.workspace_id FROM byo_did_numbers WHERE byo_did_numbers.number = ?`, did)
	var byoDidNumber string
	var byoDidWorkspaceId string
	err = row.Scan(&byoDidNumber,
			&byoDidWorkspaceId)
	if ( err != sql.ErrNoRows ) {  //create conference
		match, err := checkBYOPSTNIPWhitelist(*did, *source) 
		if err != nil {
			handleInternalErr("IncomingPSTNValidation error 1", err, w)
			return
		}

		if match {
			fmt.Printf("Matched incoming DID..")
			valid, err := finishValidation(*number, byoDidWorkspaceId)
			if err != nil {
				handleInternalErr("IncomingPSTNValidation error valid", err, w)
				return
			}
			if !valid {
				fmt.Printf("number not valid..")
				w.WriteHeader(http.StatusInternalServerError) // send the headers with a 204 response code.
				return
			}
		    w.WriteHeader(http.StatusNoContent) // send the headers with a 204 response code.
			return
		}
	}


	fmt.Printf("No results were found..")

	w.WriteHeader(http.StatusInternalServerError) // send the headers with a 204 response code.
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
	for results.Next() {
		var ipAddr string
		var ipRange string
		err = results.Scan(&ipAddr, &ipRange)
		if err != nil {
			handleInternalErr("IncomingMediaServerValidation error 2", err, w)
			return
		}
		full := ipAddr + ipRange
		match, err  := checkCIDRMatch(*source, full)
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
	


func main() {
	fmt.Print("starting Lineblocs API server..");
    r := mux.NewRouter()
    // Routes consist of a path and a handler function.
	r.HandleFunc("/call/createCall", CreateCall).Methods("POST");
	r.HandleFunc("/call/updateCall", UpdateCall).Methods("POST");
	r.HandleFunc("/conference/createConference", CreateConference).Methods("POST");
	
	//debits
	r.HandleFunc("/debit/createDebit", CreateDebit).Methods("POST");
	r.HandleFunc("/debit/createAPIUsageDebit", CreateAPIUsageDebit).Methods("POST");

	//logs
	r.HandleFunc("/log/createLog", CreateLog).Methods("POST");
	r.HandleFunc("/log/createLogSimple", CreateLogSimple).Methods("POST");

	//fax
	r.HandleFunc("/fax/createFax", CreateFax).Methods("POST");

	//recording
	r.HandleFunc("/recording/createRecording", CreateRecording).Methods("POST");
	r.HandleFunc("/recording/updateRecording", UpdateRecording).Methods("POST");
	r.HandleFunc("/recording/updateRecordingTranscription", UpdateRecordingTranscription).Methods("POST");


	// user functions
	r.HandleFunc("/user/verifyCaller", VerifyCaller).Methods("GET");
	r.HandleFunc("/user/verifyCallerByDomain", VerifyCallerByDomain).Methods("GET");
	r.HandleFunc("/user/getUserByDomain", GetUserByDomain).Methods("GET");
	r.HandleFunc("/user/getWorkspaceMacros", GetWorkspaceMacros).Methods("GET");
	r.HandleFunc("/user/getDIDNumberData", GetDIDNumberData).Methods("GET");
	r.HandleFunc("/user/getPSTNProviderIP", GetPSTNProviderIP).Methods("GET");
	r.HandleFunc("/user/ipWhitelistLookup", IPWhitelistLookup).Methods("GET");
	r.HandleFunc("/user/getDIDAcceptOption", GetDIDAcceptOption).Methods("GET");
	r.HandleFunc("/user/getDIDAssignedIP", GetDIDAssignedIP).Methods("GET");
	r.HandleFunc("/user/getUserAssignedIP", GetUserAssignedIP).Methods("GET");
	r.HandleFunc("/user/addPSTNProviderTechPrefix", AddPSTNProviderTechPrefix).Methods("GET");
	r.HandleFunc("/user/getCallerIdToUse", GetCallerIdToUse).Methods("GET");
	r.HandleFunc("/user/getExtensionFlowInfo", GetExtensionFlowInfo).Methods("GET");
	r.HandleFunc("/user/getDIDDomain", GetDIDDomain).Methods("GET");
	r.HandleFunc("/user/getCodeFlowInfo", GetCodeFlowInfo).Methods("GET");
	r.HandleFunc("/user/incomingPSTNValidation", IncomingPSTNValidation).Methods("GET");
	r.HandleFunc("/user/incomingMediaServerValidation", IncomingMediaServerValidation).Methods("GET");



	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	var err error
	//db, err = sql.Open("mysql", "lineblocs:lineblocs@lineblocs.ckehyurhpc6m.ca-central-1.rds.amazonaws.com/lineblocs")
	db, err = sql.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/lineblocs?parseTime=true") //add parse time
	if err != nil {
		panic("Could not connect to MySQL");
		return
	}
    // Bind to a port and pass our router in
    log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}