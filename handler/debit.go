package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

func (h *Handler) CreateDebit(c echo.Context) error {
	var debit model.Debit

	if err := c.Bind(&debit); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := c.Validate(&debit); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	workspace, err := h.callStore.GetWorkspaceFromDB(debit.WorkspaceId)
	if err != nil {
		fmt.Printf("could not get workspace..")
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	rate := utils.LookupBestCallRate(debit.Number, debit.Type)
	if rate == nil {
		return c.JSON(http.StatusNotFound, utils.NewError(err))
	}
	debit.PlanSnapshot = workspace.Plan
	err = h.debitStore.CreateDebit(rate, &debit)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) CreateAPIUsageDebit(c echo.Context) error {
	var debitApi model.DebitAPI

	if err := c.Bind(&debitApi); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := c.Validate(&debitApi); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	workspace, err := h.callStore.GetWorkspaceFromDB(debitApi.WorkspaceId)
	if err != nil {
		fmt.Printf("could not get workspace..")
		return err
	}

	err = h.debitStore.CreateAPIUsageDebit(workspace, &debitApi)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.NoContent(http.StatusNoContent)
}
