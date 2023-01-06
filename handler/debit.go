package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

func (h *Handler) CreateDebit(c echo.Context) error {
	var debit model.Debit

	if err := c.Bind(&debit); err != nil {
		return utils.HandleInternalErr("CreateDebit 1 Could not decode JSON", err, c)
	}
	if err := c.Validate(&debit); err != nil {
		return utils.HandleInternalErr("CreateDebit 2 Could not decode JSON", err, c)
	}

	workspace, err := h.callStore.GetWorkspaceFromDB(debit.WorkspaceId)
	if err != nil {
		return utils.HandleInternalErr("Could not get workspace..", err, c)
	}
	rate := utils.LookupBestCallRate(debit.Number, debit.Type)
	if rate == nil {
		return c.NoContent(http.StatusNotFound)
	}
	debit.PlanSnapshot = workspace.Plan
	err = h.debitStore.CreateDebit(rate, &debit)
	if err != nil {
		return utils.HandleInternalErr("CreateDebit Could not execute query..", err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) CreateAPIUsageDebit(c echo.Context) error {
	var debitApi model.DebitAPI

	if err := c.Bind(&debitApi); err != nil {
		return utils.HandleInternalErr("CreateDebit 1 Could not decode JSON", err, c)
	}
	if err := c.Validate(&debitApi); err != nil {
		return utils.HandleInternalErr("CreateDebit 2 Could not decode JSON", err, c)
	}
	workspace, err := h.callStore.GetWorkspaceFromDB(debitApi.WorkspaceId)
	if err != nil {
		return utils.HandleInternalErr("Could not get workspace..", err, c)
	}

	err = h.debitStore.CreateAPIUsageDebit(workspace, &debitApi)

	if err != nil {
		return utils.HandleInternalErr("CreateDebit Could not execute query..", err, c)
	}

	return c.NoContent(http.StatusNoContent)
}
