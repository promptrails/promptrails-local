package handler

import (
	"net/http"
	"time"

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
	// Single-get fills one level of children.
	if tree, ok := h.store.ExecutionTree(execution.ID); ok {
		execution.Children = tree.Children
	}
	return dataResponse(c, http.StatusOK, execution)
}

// ApprovalInbox lists the runs parked at waiting_approval. The real API filters
// by the caller's approval-policy membership; the single-namespace emulator
// returns every parked run.
func (h *ExecutionHandler) ApprovalInbox(c echo.Context) error {
	p := getPagination(c)
	executions, total := h.store.ListExecutions(store.ExecutionFilters{Status: "waiting_approval"}, p.Page, p.Limit)
	return listResponse(c, executions, total, p)
}

// Tree returns the execution with its full descendant tree populated.
func (h *ExecutionHandler) Tree(c echo.Context) error {
	execution, ok := h.store.ExecutionTree(c.Param("executionId"))
	if !ok {
		return notFound(c, "execution not found")
	}
	return dataResponse(c, http.StatusOK, execution)
}

// Cancel requests cooperative cancellation of a running execution, moving it to
// cancel_requested (the runner would finalize it as cancelled).
func (h *ExecutionHandler) Cancel(c echo.Context) error {
	execution, ok := h.store.GetExecution(c.Param("executionId"))
	if !ok {
		return notFound(c, "execution not found")
	}
	if execution.Status == "completed" || execution.Status == "failed" || execution.Status == "cancelled" || execution.Status == "rejected" {
		return badRequest(c, "execution has already finished")
	}
	execution.Status = "cancel_requested"
	h.store.UpdateExecution(execution)
	return dataResponse(c, http.StatusOK, execution)
}

// Approve resumes a run parked at waiting_approval and completes the approved
// call. The decision is single-use.
func (h *ExecutionHandler) Approve(c echo.Context) error {
	return h.decide(c, true)
}

// Deny resumes a parked run with the denial injected; the run finishes rejected.
func (h *ExecutionHandler) Deny(c echo.Context) error {
	return h.decide(c, false)
}

func (h *ExecutionHandler) decide(c echo.Context, approve bool) error {
	execution, ok := h.store.GetExecution(c.Param("executionId"))
	if !ok {
		return notFound(c, "execution not found")
	}
	if execution.Status != "waiting_approval" {
		return badRequest(c, "execution is not awaiting approval")
	}

	now := time.Now()
	execution.ApprovalExpiresAt = nil
	execution.CompletedAt = &now
	durationMS := int64(250)
	execution.DurationMS = &durationMS
	if approve {
		execution.Status = "completed"
		execution.Output = mustJSON(map[string]any{"content": "Approved and completed.", "type": "text"})
		execution.TokenUsage = mustJSON(map[string]any{"prompt_tokens": 150, "completion_tokens": 80, "total_tokens": 230})
		execution.Cost = 0.0023
	} else {
		execution.Status = "rejected"
		execution.Error = "approval denied"
	}
	h.store.UpdateExecution(execution)
	return dataResponse(c, http.StatusOK, execution)
}
