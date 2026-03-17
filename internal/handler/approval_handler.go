package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
)

type ApprovalHandler struct {
	store *store.Store
}

func (h *ApprovalHandler) List(c echo.Context) error {
	p := getPagination(c)
	approvals, total := h.store.ListApprovals(p.Page, p.Limit)
	return listResponse(c, approvals, total, p)
}

func (h *ApprovalHandler) Get(c echo.Context) error {
	approval, ok := h.store.GetApproval(c.Param("approvalId"))
	if !ok {
		return notFound(c, "approval not found")
	}
	return dataResponse(c, http.StatusOK, approval)
}

func (h *ApprovalHandler) Decide(c echo.Context) error {
	approval, ok := h.store.GetApproval(c.Param("approvalId"))
	if !ok {
		return notFound(c, "approval not found")
	}

	if approval.Status != "pending" {
		return badRequest(c, "approval has already been decided")
	}

	var req model.DecideApprovalRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.Decision != "approved" && req.Decision != "rejected" {
		return badRequest(c, "decision must be 'approved' or 'rejected'")
	}

	now := time.Now()
	decidedBy := "local-user"
	approval.Status = req.Decision
	approval.DecidedBy = &decidedBy
	approval.DecidedAt = &now
	if req.Reason != "" {
		approval.Reason = &req.Reason
	}

	h.store.UpdateApproval(approval)

	// Update the associated execution status
	if exec, ok := h.store.GetExecution(approval.ExecutionID); ok {
		if req.Decision == "approved" {
			exec.Status = "completed"
		} else {
			exec.Status = "rejected"
		}
		h.store.UpdateExecution(exec)
	}

	return dataResponse(c, http.StatusOK, approval)
}
