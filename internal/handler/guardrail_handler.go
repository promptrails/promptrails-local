package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
	"github.com/segmentio/ksuid"
)

type GuardrailHandler struct {
	store *store.Store
}

func (h *GuardrailHandler) List(c echo.Context) error {
	agentID := c.Param("agentId")
	if _, ok := h.store.GetAgent(agentID); !ok {
		return notFound(c, "agent not found")
	}
	guardrails := h.store.ListGuardrails(agentID)
	return dataResponse(c, http.StatusOK, guardrails)
}

func (h *GuardrailHandler) Create(c echo.Context) error {
	agentID := c.Param("agentId")
	if _, ok := h.store.GetAgent(agentID); !ok {
		return notFound(c, "agent not found")
	}

	var req model.CreateGuardrailRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	now := time.Now()
	guardrail := model.Guardrail{
		ID:          ksuid.New().String(),
		AgentID:     agentID,
		Type:        req.Type,
		ScannerType: req.ScannerType,
		Action:      req.Action,
		Config:      req.Config,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	h.store.CreateGuardrail(guardrail)
	return dataResponse(c, http.StatusCreated, guardrail)
}

func (h *GuardrailHandler) Update(c echo.Context) error {
	guardrail, ok := h.store.GetGuardrail(c.Param("guardrailId"))
	if !ok {
		return notFound(c, "guardrail not found")
	}

	var req model.UpdateGuardrailRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if req.Action != nil {
		guardrail.Action = *req.Action
	}
	if req.Config != nil {
		guardrail.Config = req.Config
	}
	if req.IsActive != nil {
		guardrail.IsActive = *req.IsActive
	}
	guardrail.UpdatedAt = time.Now()

	h.store.UpdateGuardrail(guardrail)
	return dataResponse(c, http.StatusOK, guardrail)
}

func (h *GuardrailHandler) Delete(c echo.Context) error {
	if !h.store.DeleteGuardrail(c.Param("guardrailId")) {
		return notFound(c, "guardrail not found")
	}
	return c.NoContent(http.StatusNoContent)
}
