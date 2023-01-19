package helpers

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"

	//"encoding/json"
	"reflect"

	"github.com/sirupsen/logrus"
	"lineblocs.com/api/utils"
)

type FlowInfo struct {
	FlowId   int    `json:"flow_id"`
	FlowJSON string `json:"flow_json"`
}
type Vertice struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type CellConnection struct {
	Id   string `json:"id"`
	Port string `json:"port"`
}
type GraphCell struct {
	Id       string         `json:"id"`
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	Source   CellConnection `json:"source"`
	Target   CellConnection `json:"target"`
	Vertices []Vertice      `json:"vertices"`
}

type Link struct {
	Link   *GraphCell
	Source *Cell
	Target *Cell
}
type Cell struct {
	Cell        *GraphCell
	Model       *Model
	SourceLinks []*Link
	TargetLinks []*Link
	EventVars   map[string]string
}

type ModelData interface {
}

type ModelDataStr struct {
	Value string
}
type ModelDataBool struct {
	Value bool
}
type ModelDataArr struct {
	Value []string
}
type ModelDataObj struct {
	Value map[string]string
}
type ModelLink struct {
	Type      string `json:"type"`
	Condition string `json:"condition"`
	Value     string `json:"value"`
	Cell      string `json:"cell"`
}
type Model struct {
	Id    string
	Name  string
	Data  map[string]ModelData
	Links []*ModelLink `json:"links"`
}

type UnparsedModel struct {
	Id    string                 `json:"id"`
	Name  string                 `json:"name"`
	Data  map[string]interface{} `json:"data"`
	Links []*ModelLink           `json:"links"`
}
type Graph struct {
	Cells []*GraphCell `json:"cells"`
}
type FlowVars struct {
	Graph  Graph           `json:"graph"`
	Models []UnparsedModel `json:"models"`
}

type FlowDIDData struct {
	//FlowJson FlowVars `json:"flow_json"`
	FlowId        int    `json:"flow_id"`
	WorkspaceId   int    `json:"workspace_id"`
	WorkspaceName string `json:"workspace_name"`
	CreatorId     int    `json:"creator_id"`
	FlowJson      string `json:"flow_json"`
	Plan          string `json:"plan"`
}

func findCellInFlow(id string, flow *Flow) *Cell {

	var cellToFind *GraphCell
	for _, cell := range flow.Vars.Graph.Cells {
		if cell.Id == id {
			cellToFind = cell
		}
	}
	if cellToFind == nil {
		// could not find
	}
	cell := Cell{Cell: cellToFind, EventVars: make(map[string]string)}
	return &cell
}

func findLinkByName(links []*Link, direction string, tag string) (*Link, error) {
	utils.Log(logrus.InfoLevel, "FindLinkByName called...")
	for _, link := range links {
		utils.Log(logrus.InfoLevel, "FindLinkByName checking source port: "+link.Link.Source.Port)
		utils.Log(logrus.InfoLevel, "FindLinkByName checking target port: "+link.Link.Target.Port)
		if direction == "source" {
			utils.Log(logrus.InfoLevel, "FindLinkByName checking link: "+link.Source.Cell.Name)
			if link.Link.Source.Port == tag {
				return link, nil
			}
		} else if direction == "target" {

			utils.Log(logrus.InfoLevel, "FindLinkByName checking link: "+link.Target.Cell.Name)
			if link.Link.Target.Port == tag {
				return link, nil
			}
		}
	}
	return nil, errors.New("Could not find link")
}

func createCellData(cell *Cell, flow *Flow) {
	var model Model = Model{
		Id:   "",
		Data: make(map[string]ModelData)}
	sourceLinks := make([]*Link, 0)
	targetLinks := make([]*Link, 0)
	for _, item := range flow.Vars.Models {
		if item.Id == cell.Cell.Id {
			//unparsedModel := item
			var modelData map[string]interface{}
			modelData = item.Data
			model.Name = item.Name
			model.Links = item.Links
			//json.Unmarshal([]byte(unparsedModel.Data), &modelData)

			for key, v := range modelData {
				//var item ModelData
				//item = ModelData{}
				typeOfValue := fmt.Sprintf("%s", reflect.TypeOf(v))

				utils.Log(logrus.InfoLevel, fmt.Sprintf("parsing type %s\r\n", typeOfValue))
				utils.Log(logrus.InfoLevel, fmt.Sprintf("setting key: %s\r\n", key))
				switch typeOfValue {
				case "[]string":
					// it's an array
					value := v.([]string)
					model.Data[key] = ModelDataArr{Value: value}
					//item.ValueArr = v
					//item.IsArray = true

				case "map[string]string":
					// it's an object
					utils.Log(logrus.InfoLevel, "converting obj")
					value := v.(map[string]string)
					model.Data[key] = ModelDataObj{Value: value}
					//item.ValueObj = v
					//item.IsObj = true
				case "string":
					// it's something else
					//item.ValueStr = v.(string)
					//item.IsStr = true
					value := v.(string)
					model.Data[key] = ModelDataStr{Value: value}
				case "boolean":
					// it's something else
					//item.ValueBool = v.(bool)
					//item.IsBool = true
					value := v.(bool)
					model.Data[key] = ModelDataBool{Value: value}
				}

			}
		}
	}

	cell.Model = &model

	for _, item := range flow.Vars.Graph.Cells {
		if item.Type == "devs.FlowLink" {
			utils.Log(logrus.InfoLevel, fmt.Sprintf("createCellData processing link %s\r\n", item.Type))
			if item.Source.Id == cell.Cell.Id {
				utils.Log(logrus.InfoLevel, fmt.Sprintf("createCellData adding target link %s\r\n", item.Target.Id))
				destCell := addCellToFlow(item.Target.Id, flow)
				link := &Link{
					Link:   item,
					Source: cell,
					Target: destCell}
				sourceLinks = append(sourceLinks, link)
			} else if item.Target.Id == cell.Cell.Id {
				utils.Log(logrus.InfoLevel, fmt.Sprintf("createCellData adding source link %s\r\n", item.Target.Id))
				srcCell := addCellToFlow(item.Target.Id, flow)
				link := &Link{
					Link:   item,
					Source: srcCell,
					Target: cell}
				targetLinks = append(targetLinks, link)
			}

		}
	}
	cell.SourceLinks = sourceLinks
	cell.TargetLinks = targetLinks
}
func addCellToFlow(id string, flow *Flow) *Cell {

	for _, cell := range flow.Cells {
		if cell.Cell.Id == id {
			return cell
		}
	}

	cellInFlow := findCellInFlow(id, flow)

	utils.Log(logrus.InfoLevel, fmt.Sprintf("adding cell %s", cellInFlow.Cell.Id))
	flow.Cells = append(flow.Cells, cellInFlow)
	createCellData(cellInFlow, flow)
	return cellInFlow
}

type FlowContext struct {
	DbConn    *sql.DB
	Cell      *Cell
	Data      map[string]string
	Providers []*RoutablePSTNProvider
}

type RoutablePSTNProvider struct {
	Id    int
	Name  string
	Hosts []RoutableHost
	Data  map[string]int
}

type RoutableHost struct {
	Priority int
	IPAddr   string
	Prefix   string
}

type FlowResponse struct {
	Providers []*RoutablePSTNProvider
	Link      *Link
}
type BaseManager interface {
	Process() (*FlowResponse, error)
}
type Manager struct {
	Ctx *FlowContext
}
type CallCapacityManager struct {
	*Manager
}

func NewCallCapacityManager(ctx *FlowContext) *CallCapacityManager {
	return &CallCapacityManager{&Manager{Ctx: ctx}}
}

func (man *CallCapacityManager) Process() (*FlowResponse, error) {
	db := man.Ctx.DbConn
	cell := man.Ctx.Cell
	//providers := make( []*RoutablePSTNProvider,0 )
	providers := man.Ctx.Providers

	outLink, _ := findLinkByName(cell.TargetLinks, "source", "Out")
	noMatchLink, _ := findLinkByName(cell.TargetLinks, "source", "No match")
	// lookup by country
	results, err := db.Query(`SELECT sip_providers.provider_id, 
sip_providers_hosts.name,
sip_providers_hosts.ip_address,
sip_providers_hosts.priority,
sip_providers.active_channels,
sip_providers_hosts.priority_prefixes
FROM sip_providers_hosts
INNER JOIN sip_providers_call_rates ON sip_providers_call_rates.provider_id = sip_providers_hosts.provider_id
INNER JOIN sip_providers ON sip_providers.id = sip_providers_hosts.provider_id
INNER JOIN sip_countries ON sip_countries.id = sip_providers_call_rates.country_id
WHERE sip_countries.country_code= ?`, man.Ctx.Data["dest_code"])

	defer results.Close()

	if err != nil {
		return nil, err
	}

	for results.Next() {
		var providerId int
		var name string
		var ipAddr string
		var priority int
		var channels int
		var prefixes string
		// for each row, scan the result into our tag composite object
		err = results.Scan(&providerId, &name, &ipAddr, &priority, &channels, &prefixes)
		if err != nil {
			return nil, err
		}
		provider := createOrUseExistingProvider(providers, providerId)
		host := RoutableHost{
			Prefix:   prefixes,
			Priority: priority,
			IPAddr:   ipAddr}
		provider.Data["channels"] = channels
		provider.Hosts = append(provider.Hosts, host)
		providers = append(providers, provider)
	}
	// sort based on costs
	sort.SliceStable(providers, func(i, j int) bool {
		return providers[i].Data["channels"] < providers[j].Data["channels"]
	})

	return createFlowResponse(providers, outLink, noMatchLink), nil
}

type LowCostManager struct {
	*Manager
}

func NewLowCostManager(ctx *FlowContext) *LowCostManager {
	return &LowCostManager{&Manager{Ctx: ctx}}
}
func createOrUseExistingProvider(providers []*RoutablePSTNProvider, providerId int) *RoutablePSTNProvider {
	for _, value := range providers {
		if value.Id == providerId {
			return value
		}
	}

	// create new one
	return &RoutablePSTNProvider{Id: providerId, Hosts: make([]RoutableHost, 0), Data: make(map[string]int)}
}
func createFlowResponse(providers []*RoutablePSTNProvider, outLink, noMatchLink *Link) *FlowResponse {
	var link *Link = outLink

	if len(providers) == 0 {
		link = noMatchLink
	}
	resp := FlowResponse{
		Providers: providers,
		Link:      link}
	return &resp
}

func (man *LowCostManager) Process() (*FlowResponse, error) {
	db := man.Ctx.DbConn
	cell := man.Ctx.Cell
	//providers := make( []*RoutablePSTNProvider,0 )
	providers := man.Ctx.Providers

	outLink, _ := findLinkByName(cell.TargetLinks, "source", "Out")
	noMatchLink, _ := findLinkByName(cell.TargetLinks, "source", "No match")
	// lookup by country
	results, err := db.Query(`SELECT sip_providers_call_rates.provider_id, 
sip_providers_hosts.name,
sip_providers_hosts.ip_address,
sip_providers_hosts.priority,
sip_providers_hosts.priority,
sip_providers_call_rates.rate,
FROM sip_providers_hosts
INNER JOIN sip_providers_call_rates ON sip_providers_call_rates.provider_id = sip_providers_hosts.provider_id
INNER JOIN sip_providers ON sip_providers.id = sip_providers_hosts.provider_id
INNER JOIN sip_countries ON sip_countries.id = sip_providers_call_rates.country_id
WHERE sip_countries.country_code= ?`, man.Ctx.Data["dest_code"])

	defer results.Close()

	if err != nil {
		return nil, err
	}

	for results.Next() {
		var providerId int
		var name string
		var ipAddr string
		var priority int
		var rate int
		// for each row, scan the result into our tag composite object
		err = results.Scan(&providerId, &name, &ipAddr, &priority, &rate)
		if err != nil {
			return nil, err
		}
		provider := createOrUseExistingProvider(providers, providerId)
		host := RoutableHost{
			Priority: priority,
			IPAddr:   ipAddr}
		provider.Data["cost"] = rate
		provider.Hosts = append(provider.Hosts, host)
		providers = append(providers, provider)
	}
	// sort based on costs
	sort.SliceStable(providers, func(i, j int) bool {
		return providers[i].Data["cost"] < providers[j].Data["cost"]
	})

	return createFlowResponse(providers, outLink, noMatchLink), nil
}

type HighCostManager struct {
	*Manager
}

func NewHighCostManager(ctx *FlowContext) *HighCostManager {
	return &HighCostManager{&Manager{Ctx: ctx}}
}

func (man *HighCostManager) Process() (*FlowResponse, error) {
	return nil, nil
}

type LocationCheckManager struct {
	*Manager
}

func NewLocationCheckManager(ctx *FlowContext) *LocationCheckManager {
	return &LocationCheckManager{&Manager{Ctx: ctx}}
}

func (man *LocationCheckManager) Process() (*FlowResponse, error) {
	return nil, nil
}

type SortServersManager struct {
	*Manager
}

func NewSortServersManager(ctx *FlowContext) *SortServersManager {
	return &SortServersManager{&Manager{Ctx: ctx}}
}

func (man *SortServersManager) Process() (*FlowResponse, error) {
	return nil, nil
}

type UserPriorityManager struct {
	*Manager
}

func NewUserPriorityManager(ctx *FlowContext) *UserPriorityManager {
	return &UserPriorityManager{&Manager{Ctx: ctx}}
}

func (man *UserPriorityManager) Process() (*FlowResponse, error) {
	return nil, nil
}

type EndRoutingManager struct {
	*Manager
}

func NewEndRoutingManager(ctx *FlowContext) *EndRoutingManager {
	return &EndRoutingManager{&Manager{Ctx: ctx}}
}

func (man *EndRoutingManager) Process() (*FlowResponse, error) {
	return nil, nil
}

type NoRoutingManager struct {
	*Manager
}

func NewNoRoutingManager(ctx *FlowContext) *NoRoutingManager {
	return &NoRoutingManager{&Manager{Ctx: ctx}}
}

func (man *NoRoutingManager) Process() (*FlowResponse, error) {
	return nil, nil
}

func ProcessFlow(flow *Flow, cell *Cell, providers []*RoutablePSTNProvider, data map[string]string, db *sql.DB) ([]*RoutablePSTNProvider, error) {
	utils.Log(logrus.InfoLevel, "source link count: "+strconv.Itoa(len(cell.SourceLinks)))
	utils.Log(logrus.InfoLevel, "target link count: "+strconv.Itoa(len(cell.TargetLinks)))
	// execute it
	var mngr BaseManager
	var isFinished bool = false

	ctx := &FlowContext{
		DbConn: db,
		Data:   data,
		Cell:   cell}
	switch cell.Cell.Type {
	case "devs.LaunchModel":
		for _, link := range cell.SourceLinks {
			return ProcessFlow(flow, link.Target, providers, data, db)
		}
		return providers, nil
	case "devs.CallCapacityModel":
		mngr = NewCallCapacityManager(ctx)
	case "devs.LowCostModel":
		mngr = NewLowCostManager(ctx)
	case "devs.HighCostModel":
		mngr = NewHighCostManager(ctx)
	case "devs.LocationCheckModel":
		mngr = NewLocationCheckManager(ctx)
	case "devs.UserPriorityModel":
		mngr = NewUserPriorityManager(ctx)
	case "devs.SortServersModel":
		mngr = NewSortServersManager(ctx)
	case "devs.EndRoutingModel":
		mngr = NewEndRoutingManager(ctx)
		isFinished = true
	case "devs.NoRoutingModel":
		mngr = NewNoRoutingManager(ctx)
		isFinished = true
	default:
		utils.Log(logrus.InfoLevel, "unknown type of cell..")
		return nil, errors.New("unknown type of cell..")
	}
	var resp *FlowResponse
	var err error
	resp, err = mngr.Process()

	if err != nil {
		return nil, err
	}

	if resp.Link == nil || isFinished {
		return resp.Providers, nil
	}
	next := resp.Link
	return ProcessFlow(flow, next.Target, providers, data, db)
}
func StartProcessingFlow(flow *Flow, cell *Cell, data map[string]string, db *sql.DB) ([]*RoutablePSTNProvider, error) {

	emptyProviders := make([]*RoutablePSTNProvider, 0)
	providers, err := ProcessFlow(flow, flow.Cells[0], emptyProviders, data, db)
	return providers, err
}

func NewFlow(id int, vars *FlowVars) *Flow {
	flow := &Flow{FlowId: id, Vars: vars}
	utils.Log(logrus.InfoLevel, fmt.Sprintf("number of cells %d\r\n", len(flow.Vars.Graph.Cells)))

	// create cells from flow.Vars
	for _, cell := range flow.Vars.Graph.Cells {
		utils.Log(logrus.InfoLevel, fmt.Sprintf("processing %s\r\n", cell.Type))
		// creating a cell
		if cell != nil {
			if cell.Type != "devs.FlowLink" {
				addCellToFlow(cell.Id, flow)
			}
		}
	}
	return flow
}

type Flow struct {
	Exten    string
	CallerId string
	Cells    []*Cell
	Models   []*Model
	Vars     *FlowVars
	FlowId   int
}

type Runner struct {
	Cancelled bool
}
