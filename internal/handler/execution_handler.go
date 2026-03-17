package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/store"
)

type ExecutionHandler struct {
	store *store.Store
}

func (h *ExecutionHandler) List(c echo.Context) error {
	p := getPagination(c)

	filters := store.ExecutionFilters{
		AgentID:   c.QueryParam("agent_id"),
		SessionID: c.QueryParam("session_id"),
		Status:    c.QueryParam("status"),
	}

	executions, total := h.store.ListExecutions(filters, p.Page, p.Limit)
	return listResponse(c, executions, total, p)
}

func (h *ExecutionHandler) Get(c echo.Context) error {
	execution, ok := h.store.GetExecution(c.Param("executionId"))
	if !ok {
		return notFound(c, "execution not found")
	}
	return dataResponse(c, http.StatusOK, execution)
}

func (h *ExecutionHandler) GetPendingApproval(c echo.Context) error {
	executionID := c.Param("executionId")
	if _, ok := h.store.GetExecution(executionID); !ok {
		return notFound(c, "execution not found")
	}

	approval, ok := h.store.GetApprovalByExecutionID(executionID)
	if !ok {
		return notFound(c, "no pending approval for this execution")
	}
	return dataResponse(c, http.StatusOK, approval)
}
