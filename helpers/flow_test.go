package helpers

import (
	"testing"

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
