package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

/*
Input: Debit model
Todo : Create new user_debit and store to db
Output: If success return NoContent else return err
*/
func (h *Handler) CreateDebit(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "CreateDebit is called...")

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

	// Get Call Rate depends number and type
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

/*
Input: DebitAPI model
Todo : Calculate cents based on debit type and create user_debit
Output: If success return NoContent else return err
*/
func (h *Handler) CreateAPIUsageDebit(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "CreateAPIUsageDebit is called...\r\n")

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
