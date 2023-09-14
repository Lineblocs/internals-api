package helpers

import (
	"testing"

	helpers "github.com/Lineblocs/go-helpers"
	"github.com/stretchr/testify/assert"
)

func TestFindCellInFlow(t *testing.T) {
	t.Run("Should Found", func(t *testing.T) {
		graphCells := []GraphCell{
			{
				Id: "cell1",
			},
			{
				Id: "cell2",
			},
		}

		flow := &Flow{
			Vars: &FlowVars{
				Graph: Graph{
					Cells: make([]*GraphCell, len(graphCells)),
				},
			},
		}

		for i, cell := range graphCells {
			flow.Vars.Graph.Cells[i] = &cell
		}

		cell := findCellInFlow("cell2", flow)

		assert.NotNil(t, cell)
		assert.Equal(t, "cell2", cell.Cell.Id)
	})

	t.Run("Should NotFound", func(t *testing.T) {
		graphCells := []GraphCell{
			{
				Id: "cell1",
			},
			{
				Id: "cell3",
			},
		}

		flow := &Flow{
			Vars: &FlowVars{
				Graph: Graph{
					Cells: make([]*GraphCell, len(graphCells)),
				},
			},
		}

		for i, cell := range graphCells {
			flow.Vars.Graph.Cells[i] = &cell
		}

		cell := findCellInFlow("cell2", flow)

		expCell := &Cell{
			EventVars: make(map[string]string),
		}

		assert.Equal(t, cell, expCell)
	})
}

func TestNewCallCapacityManager(t *testing.T) {
	t.Run("Should create a new call capacity manager", func(t *testing.T) {
		capacity := NewCallCapacityManager(nil)
		assert.NotNil(t, capacity)
	})
}

func TestNewLowCostManager(t *testing.T) {
	t.Run("Should create a new low cost manager", func(t *testing.T) {
		capacity := NewLowCostManager(nil)
		assert.NotNil(t, capacity)
	})
}

func TestNewHighCostManager(t *testing.T) {
	t.Run("Should create a new high cost manager", func(t *testing.T) {
		capacity := NewHighCostManager(nil)
		assert.NotNil(t, capacity)
	})
}

func TestNewHighCostManager_Process(t *testing.T) {
	t.Run("Should process a high cost manager", func(t *testing.T) {
		capacity := NewHighCostManager(nil)
		_, error := capacity.Process()
		assert.NoError(t, error)
	})
}

func TestNewLocationCheckManager(t *testing.T) {
	t.Run("Should create a new location check manager", func(t *testing.T) {
		capacity := NewLocationCheckManager(nil)
		assert.NotNil(t, capacity)
	})
}

func TestNewLocationCheckManager_Process(t *testing.T) {
	t.Run("Should process a new location check manager", func(t *testing.T) {
		capacity := NewLocationCheckManager(nil)
		_, error := capacity.Process()
		assert.NoError(t, error)
	})
}

func TestNewSortServersManager(t *testing.T) {
	t.Run("Should create a new SortServersManager", func(t *testing.T) {
		ctx := &FlowContext{}
		sortServersManager := NewSortServersManager(ctx)
		assert.NotNil(t, sortServersManager)
	})
}

func TestNewUserPriorityManager(t *testing.T) {
	t.Run("Should create a new UserPriorityManager", func(t *testing.T) {
		ctx := &FlowContext{}
		userPriorityManager := NewUserPriorityManager(ctx)
		assert.NotNil(t, userPriorityManager)
	})
}

func TestNewEndRoutingManager(t *testing.T) {
	t.Run("Should create a new EndRoutingManager", func(t *testing.T) {
		ctx := &FlowContext{}
		endRoutingManager := NewEndRoutingManager(ctx)
		assert.NotNil(t, endRoutingManager)
	})
}

func TestNewNoRoutingManager(t *testing.T) {
	t.Run("Should create a new NoRoutingManager", func(t *testing.T) {
		ctx := &FlowContext{}
		noRoutingManager := NewNoRoutingManager(ctx)
		assert.NotNil(t, noRoutingManager)
	})
}

func TestSortServersManager_Process(t *testing.T) {
	t.Run("Should process a SortServersManager", func(t *testing.T) {
		ctx := &FlowContext{}
		sortServersManager := NewSortServersManager(ctx)
		_, err := sortServersManager.Process()
		assert.NoError(t, err)
	})
}

func TestUserPriorityManager_Process(t *testing.T) {
	t.Run("Should process a UserPriorityManager", func(t *testing.T) {
		ctx := &FlowContext{}
		userPriorityManager := NewUserPriorityManager(ctx)
		_, err := userPriorityManager.Process()
		assert.NoError(t, err)
	})
}

func TestEndRoutingManager_Process(t *testing.T) {
	t.Run("Should process an EndRoutingManager", func(t *testing.T) {
		ctx := &FlowContext{}
		endRoutingManager := NewEndRoutingManager(ctx)
		_, err := endRoutingManager.Process()
		assert.NoError(t, err)
	})
}

func TestNoRoutingManager_Process(t *testing.T) {
	t.Run("Should process a NoRoutingManager", func(t *testing.T) {
		ctx := &FlowContext{}
		noRoutingManager := NewNoRoutingManager(ctx)
		_, err := noRoutingManager.Process()
		assert.NoError(t, err)
	})
}

func TestCreateOrUseExistingProvider(t *testing.T) {
	t.Run("CreateNewProvider", func(t *testing.T) {

		existingProvider := &RoutablePSTNProvider{Id: 1}
		providers := []*RoutablePSTNProvider{existingProvider}

		newProviderId := 2
		newProvider := createOrUseExistingProvider(providers, newProviderId)

		assert.NotNil(t, newProvider)
		assert.Equal(t, newProviderId, newProvider.Id)
		assert.Len(t, providers, 1)
	})

	t.Run("UseExistingProvider", func(t *testing.T) {

		existingProvider := &RoutablePSTNProvider{Id: 1}
		providers := []*RoutablePSTNProvider{existingProvider}

		existingProviderId := 1
		existingProviderResult := createOrUseExistingProvider(providers, existingProviderId)

		assert.NotNil(t, existingProviderResult)
		assert.Equal(t, existingProvider, existingProviderResult)
		assert.Len(t, providers, 1)
	})
}

func TestCreateFlowResponse(t *testing.T) {
	t.Run("ProvidersNotEmpty", func(t *testing.T) {

		providers := []*RoutablePSTNProvider{}
		outLink := &Link{}
		noMatchLink := &Link{}

		response := createFlowResponse(providers, outLink, noMatchLink)

		assert.NotNil(t, response)
		assert.Equal(t, providers, response.Providers)
		assert.Equal(t, outLink, response.Link)
	})

	t.Run("ProvidersEmpty", func(t *testing.T) {

		providers := []*RoutablePSTNProvider{}
		outLink := &Link{}
		noMatchLink := &Link{}

		response := createFlowResponse(providers, outLink, noMatchLink)

		assert.NotNil(t, response)
		assert.Empty(t, response.Providers)
		assert.Equal(t, noMatchLink, response.Link)
	})
}

func TestCreateCellData(t *testing.T) {
	t.Run("CreateCellData", func(t *testing.T) {

		flow := &Flow{
			Vars: &FlowVars{
				Graph: Graph{
					Cells: []*GraphCell{},
				},
				Models: []UnparsedModel{
					{
						Id: "model1",
					},
					{
						Id: "model2",
					}},
			},
		}

		cell := &Cell{
			Cell: &GraphCell{
				Id: "cell1",
			},
		}

		createCellData(cell, flow)

		assert.NotNil(t, cell.Model)

	})

}

func TestAddCellToFlow(t *testing.T) {

	helpers.InitLogrus("stdout")

	t.Run("Should return the cell", func(t *testing.T) {

		flow := &Flow{
			Exten:    "",
			CallerId: "",
			Cells: []*Cell{
				{
					Cell: &GraphCell{
						Id: "cell1",
					},
				},
			},
			Models: []*Model{},
			Vars: &FlowVars{
				Graph:  Graph{},
				Models: []UnparsedModel{},
			},
			FlowId: 0,
		}

		cellID := "cell1"
		cell := addCellToFlow(cellID, flow)
		assert.Equal(t, cell.Cell.Id, cellID)
	})

	t.Run("Should create cell data", func(t *testing.T) {

		graphCells := []GraphCell{
			{
				Id: "cell1",
			},
			{
				Id: "cell2",
			},
		}

		flow := &Flow{
			Vars: &FlowVars{
				Graph: Graph{
					Cells: make([]*GraphCell, len(graphCells)),
				},
			},
		}

		for i, cell := range graphCells {
			flow.Vars.Graph.Cells[i] = &cell
		}

		cellID := "cell2"
		cell := addCellToFlow(cellID, flow)
		assert.Equal(t, cell.Cell.Id, cellID)
	})
}

func TestFindLinkByName_SourceDirection(t *testing.T) {

	helpers.InitLogrus("stdout")

	t.Run("FoundSourceLink", func(t *testing.T) {

		sourceLink := &Link{
			Link: &GraphCell{
				Source: CellConnection{Port: "tag"},
				Name:   "link1",
			},
			Source: &Cell{
				Cell: &GraphCell{
					Id:   "cell1",
					Name: "cell1",
				},
			},
			Target: &Cell{},
		}
		links := []*Link{sourceLink}

		direction := "source"
		tag := "tag"
		result, err := findLinkByName(links, direction, tag)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, sourceLink, result)
	})

	t.Run("NotFoundSourceLink", func(t *testing.T) {

		links := []*Link{}

		direction := "source"
		tag := "tag"
		result, err := findLinkByName(links, direction, tag)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestFindLinkByName_TargetDirection(t *testing.T) {

	helpers.InitLogrus("stdout")

	t.Run("FoundTargetLink", func(t *testing.T) {

		targetLink := &Link{
			Link: &GraphCell{
				Target: CellConnection{Port: "tag"},
			},
			Source: &Cell{},
			Target: &Cell{
				Cell: &GraphCell{
					Id:   "cell1",
					Name: "cell1",
				},
			},
		}
		links := []*Link{targetLink}

		direction := "target"
		tag := "tag"
		result, err := findLinkByName(links, direction, tag)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, targetLink, result)
	})

	t.Run("NotFoundTargetLink", func(t *testing.T) {

		links := []*Link{}

		direction := "target"
		tag := "tag"
		result, err := findLinkByName(links, direction, tag)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
