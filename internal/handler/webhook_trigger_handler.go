package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/fake"
	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
	"github.com/segmentio/ksuid"
)

type WebhookTriggerHandler struct {
	store *store.Store
}

func (h *WebhookTriggerHandler) List(c echo.Context) error {
	p := getPagination(c)
	triggers, total := h.store.ListWebhookTriggers(p.Page, p.Limit)
	return listResponse(c, triggers, total, p)
}

func (h *WebhookTriggerHandler) Create(c echo.Context) error {
	wid := getWorkspaceID()
	var req model.CreateWebhookTriggerRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.AgentID == "" || req.Name == "" {
		return badRequest(c, "agent_id and name are required")
	}

	token := ksuid.New().String()
	now := time.Now()

	trigger := model.WebhookTrigger{
		ID:          ksuid.New().String(),
		WorkspaceID: wid,
		AgentID:     req.AgentID,
		Name:        req.Name,
		Token:       token,
		TokenPrefix: token[:8],
		IsActive:    req.IsActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	h.store.CreateWebhookTrigger(trigger)

	// Return the full token only on creation
	resp := model.WebhookTriggerCreateResponse{
		WebhookTrigger: trigger,
		FullToken:      token,
	}
	return dataResponse(c, http.StatusCreated, resp)
}

func (h *WebhookTriggerHandler) Get(c echo.Context) error {
	trigger, ok := h.store.GetWebhookTrigger(c.Param("triggerId"))
	if !ok {
		return notFound(c, "webhook trigger not found")
	}
	return dataResponse(c, http.StatusOK, trigger)
}

func (h *WebhookTriggerHandler) Update(c echo.Context) error {
	trigger, ok := h.store.GetWebhookTrigger(c.Param("triggerId"))
	if !ok {
		return notFound(c, "webhook trigger not found")
	}

	var req model.UpdateWebhookTriggerRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if req.Name != nil {
		trigger.Name = *req.Name
	}
	if req.IsActive != nil {
		trigger.IsActive = *req.IsActive
	}
	trigger.UpdatedAt = time.Now()

	h.store.UpdateWebhookTrigger(trigger)
	return dataResponse(c, http.StatusOK, trigger)
}

func (h *WebhookTriggerHandler) Delete(c echo.Context) error {
	if !h.store.DeleteWebhookTrigger(c.Param("triggerId")) {
		return notFound(c, "webhook trigger not found")
	}
	return c.NoContent(http.StatusNoContent)
}

// Hook is the public endpoint that receives webhook calls and executes the agent.
func (h *WebhookTriggerHandler) Hook(c echo.Context) error {
	token := c.Param("token")
	trigger, ok := h.store.GetWebhookTriggerByToken(token)
	if !ok {
		return notFound(c, "invalid webhook token")
	}

	if !trigger.IsActive {
		return badRequest(c, "webhook trigger is disabled")
	}

	// Parse incoming body as input
	var input map[string]any
	if err := c.Bind(&input); err != nil {
		input = map[string]any{}
	}

	// Look up the agent name for fake output
	agentName := "Agent"
	if agent, ok := h.store.GetAgent(trigger.AgentID); ok {
		agentName = agent.Name
	}

	output := fake.GenerateExecutionOutput(agentName, input)
	inputJSON, _ := json.Marshal(input)
	outputJSON, _ := json.Marshal(output)

	now := time.Now()
	durationMS := int64(200)
	execution := model.Execution{
		ID:          ksuid.New().String(),
		AgentID:     &trigger.AgentID,
		WorkspaceID: trigger.WorkspaceID,
		SessionID:   ksuid.New().String(),
		Status:      "completed",
		Input:       inputJSON,
		Output:      outputJSON,
		TokenUsage:  json.RawMessage(`{"prompt_tokens":120,"completion_tokens":60,"total_tokens":180}`),
		Cost:        0.0018,
		DurationMS:  &durationMS,
		TraceID:     ksuid.New().String(),
		StartedAt:   &now,
		CompletedAt: &now,
		CreatedAt:   now,
	}
	h.store.CreateExecution(execution)

	// Update last used timestamp
	trigger.LastUsedAt = &now
	h.store.UpdateWebhookTrigger(trigger)

	return dataResponse(c, http.StatusCreated, execution)
}
