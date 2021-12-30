package helpers
import (
	"fmt"
	"database/sql"
	"strconv"
	"errors"
	//"encoding/json"
	"reflect"
)

type FlowInfo struct {
	FlowId          int               `json:"flow_id"`
	FlowJSON        string            `json:"flow_json"`

}
type Vertice struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type CellConnection struct {
	Id string `json:"id"`
	Port string `json:"port"`
}
type GraphCell struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Source CellConnection `json:"source"`
	Target CellConnection `json:"target"`
	Vertices []Vertice `json:"vertices"`
}

type Link struct {
	Link *GraphCell
	Source *Cell
	Target *Cell
}
type Cell struct {
	Cell *GraphCell
	Model *Model
	SourceLinks []*Link
	TargetLinks []*Link
	EventVars map[string]string
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
	Value map[string] string
}
type ModelLink struct {
	Type string `json:"type"`
	Condition string `json:"condition"`
	Value string `json:"value"`
	Cell string `json:"cell"`
}
type Model struct {
	Id string
	Name string
	Data map[string] ModelData
	Links []*ModelLink `json:"links"`
}

type UnparsedModel struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Data map[string]interface{} `json:"data"`
	Links []*ModelLink `json:"links"`

}
type Graph struct {
	Cells []*GraphCell `json:"cells"`
}
type FlowVars struct {
	Graph Graph `json:"graph"`
	Models []UnparsedModel `json:"models"`
}

type FlowDIDData struct {
	//FlowJson FlowVars `json:"flow_json"`
	FlowId int `json:"flow_id"`
	WorkspaceId int `json:"workspace_id"`
	WorkspaceName string `json:"workspace_name"`
	CreatorId int `json:"creator_id"`
	FlowJson string `json:"flow_json"`
	Plan string `json:"plan"`
}
func findCellInFlow(id string, flow *Flow) (*Cell) {

	var cellToFind *GraphCell
	for _, cell := range flow.Vars.Graph.Cells {
		if cell.Id == id {
			cellToFind = cell
		}
	}
	if cellToFind == nil {
		// could not find
	}
	cell := Cell{ Cell: cellToFind, EventVars: make( map[string]string ) }
	return &cell
}
	
func createCellData(cell *Cell, flow *Flow) {
	var model Model = Model{
		Id: "",
		Data: make(map[string] ModelData) }
	sourceLinks := make( []*Link, 0 )
	targetLinks := make( []*Link, 0 )
	for _, item := range flow.Vars.Models {
		if (item.Id == cell.Cell.Id) {
			//unparsedModel := item
			var modelData map[string]interface{}
			modelData = item.Data
			model.Name =  item.Name
			model.Links =  item.Links
			//json.Unmarshal([]byte(unparsedModel.Data), &modelData)

			for key, v := range modelData {
				//var item ModelData
				//item = ModelData{}
 				typeOfValue := fmt.Sprintf("%s", reflect.TypeOf(v))

				fmt.Printf("parsing type %s\r\n", typeOfValue)
					fmt.Printf("setting key: %s\r\n", key)
				switch ; typeOfValue {
					case "[]string":
						// it's an array
						value := v.([]string)
						model.Data[key]= ModelDataArr{Value: value}
						//item.ValueArr = v
						//item.IsArray = true

					case "map[string]string":
						// it's an object
						fmt.Println("converting obj")
						value := v.(map[string]string)
						model.Data[key]=ModelDataObj{Value: value}
						//item.ValueObj = v
						//item.IsObj = true
					case "string":
						// it's something else
						//item.ValueStr = v.(string)
						//item.IsStr = true
						value := v.(string)
						model.Data[key]=ModelDataStr{Value: value}
					case "boolean":
						// it's something else
						//item.ValueBool = v.(bool)
						//item.IsBool = true
						value := v.(bool)
						model.Data[key]=ModelDataBool{Value: value}
					}

			}
		}
	}	

	cell.Model = &model

	for _, item := range flow.Vars.Graph.Cells {
		if item.Type == "devs.FlowLink" {
			fmt.Printf("createCellData processing link %s\r\n", item.Type)
			if item.Source.Id == cell.Cell.Id {
				fmt.Printf("createCellData adding target link %s\r\n", item.Target.Id)
				destCell := addCellToFlow( item.Target.Id, flow )
				link := &Link{
					Link: item,
					Source: cell,
					Target: destCell }
				sourceLinks = append( sourceLinks, link )
			} else if item.Target.Id == cell.Cell.Id {
				fmt.Printf("createCellData adding source link %s\r\n", item.Target.Id)
				srcCell := addCellToFlow( item.Target.Id, flow )
				link := &Link{
					Link: item,
					Source: srcCell,
					Target: cell }
				targetLinks = append( targetLinks, link )
			}

		}
	}
	cell.SourceLinks = sourceLinks
	cell.TargetLinks = targetLinks
}
func addCellToFlow(id string, flow *Flow) (*Cell) {

	for _, cell := range flow.Cells {
		if cell.Cell.Id == id {
			return cell
		}
	}

	cellInFlow := findCellInFlow(id, flow)

	fmt.Printf("adding cell %s", cellInFlow.Cell.Id)
	flow.Cells = append(flow.Cells, cellInFlow)
	createCellData(cellInFlow, flow)
	return cellInFlow
}

type FlowContext struct {
	DbConn *sql.DB
}


type RoutablePSTNProvider struct {
	IPAddr string
}

type FlowResponse struct {
	Providers []*RoutablePSTNProvider
	Link *Link
}
type BaseManager interface {
	Process() (*FlowResponse, error)
}
type Manager struct {
	Ctx *FlowContext 
}
type CallCapacityManager struct {
	Manager
}
func NewCallCapacityManager(ctx *FlowContext) (*CallCapacityManager) {
	return &CallCapacityManager{};
}

func (man *CallCapacityManager) Process() (*FlowResponse, error) {
	return nil, nil
}


type LowCostManager struct {
	Manager
}

func NewLowCostManager(ctx *FlowContext) (*LowCostManager) {
	return &LowCostManager{}
}

func (man *LowCostManager) Process() (*FlowResponse, error) {
	return nil, nil
}


type HighCostManager struct {
	Manager
}

func NewHighCostManager(ctx *FlowContext) (*HighCostManager) {
	return &HighCostManager{}
}

func (man *HighCostManager) Process() (*FlowResponse, error) {
	return nil, nil
}

type LocationCheckManager struct {
	Manager
}

func NewLocationCheckManager(ctx *FlowContext) (*LocationCheckManager) {
	return &LocationCheckManager{}
}

func (man *LocationCheckManager) Process() (*FlowResponse, error) {
	return nil, nil
}

type SortServersManager struct {
	Manager
}

func NewSortServersManager(ctx *FlowContext) (*SortServersManager) {
	return &SortServersManager{}
}

func (man *SortServersManager) Process() (*FlowResponse, error) {
	return nil, nil
}





type UserPriorityManager struct {
	Manager
}

func NewUserPriorityManager(ctx *FlowContext) (*UserPriorityManager) {
	return &UserPriorityManager{}
}

func (man *UserPriorityManager) Process() (*FlowResponse, error) {
	return nil, nil
}

type EndRoutingManager struct {
	Manager
}

func NewEndRoutingManager(ctx *FlowContext) (*EndRoutingManager) {
	return &EndRoutingManager{}
}

func (man *EndRoutingManager) Process() (*FlowResponse, error) {
	return nil, nil
}



type NoRoutingManager struct {
	Manager
}

func NewNoRoutingManager(ctx *FlowContext) (*NoRoutingManager) {
	return &NoRoutingManager{}
}

func (man *NoRoutingManager) Process() (*FlowResponse, error) {
	return nil, nil
}







func ProcessFlow( ctx *FlowContext, flow *Flow, cell *Cell) ([]*RoutablePSTNProvider, error) {
	fmt.Println("source link count: " + strconv.Itoa( len( cell.SourceLinks )))
	fmt.Println("target link count: " + strconv.Itoa( len( cell.TargetLinks )))
	// execute it
	var mngr BaseManager
	providers:=make([]*RoutablePSTNProvider,0)
	switch ; cell.Cell.Type {
		case "devs.LaunchModel":
			for _, link := range cell.SourceLinks {
				return ProcessFlow( ctx, flow, link.Target )
			}
			return providers,nil
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
		case "devs.NoRoutingModel":
			mngr = NewNoRoutingManager(ctx)
		default:
			fmt.Println("unknown type of cell..")
			return nil,errors.New("unknown type of cell..")
	}
	var resp *FlowResponse
	var err error
	resp, err = mngr.Process()

	if err != nil {
		return nil, err
	}

	if resp.Link == nil {
		return resp.Providers, nil
	}
	next := resp.Link
	return ProcessFlow( ctx, flow, next.Target )
}
func NewFlow(id int,vars *FlowVars) (*Flow) {
	flow := &Flow{FlowId: id, Vars: vars}
	fmt.Printf("number of cells %d\r\n", len(flow.Vars.Graph.Cells))
	// create cells from flow.Vars
	for _, cell := range flow.Vars.Graph.Cells {
		fmt.Printf("processing %s\r\n", cell.Type)
		// creating a cell	
		if cell != nil {
			if cell.Type != "devs.FlowLink" {
				addCellToFlow( cell.Id, flow )
			}
		}
	}
	return flow
}


type Flow struct {
	Exten string
	CallerId string
	Cells []*Cell
	Models []*Model
	Vars *FlowVars
	FlowId int
}

type Runner struct {
	Cancelled bool
}

