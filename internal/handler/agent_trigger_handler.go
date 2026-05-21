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

type AgentTriggerHandler struct {
	store *store.Store
}

func (h *AgentTriggerHandler) List(c echo.Context) error {
	p := getPagination(c)
	triggers, total := h.store.ListAgentTriggers(p.Page, p.Limit)
	return listResponse(c, triggers, total, p)
}

func (h *AgentTriggerHandler) Create(c echo.Context) error {
	wid := getWorkspaceID()
	var req model.CreateAgentTriggerRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.AgentID == "" || req.Name == "" {
		return badRequest(c, "agent_id and name are required")
	}

	token := ksuid.New().String()
	now := time.Now()
	source := req.Source
	if source == "" {
		source = "generic"
	}

	trigger := model.AgentTrigger{
		ID:           ksuid.New().String(),
		WorkspaceID:  wid,
		AgentID:      req.AgentID,
		Name:         req.Name,
		Token:        token,
		TokenPrefix:  token[:8],
		Source:       source,
		SourceConfig: req.SourceConfig,
		ReplyConfig:  req.ReplyConfig,
		IsActive:     req.IsActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	h.store.CreateAgentTrigger(trigger)

	// Return the full token only on creation
	resp := model.AgentTriggerCreateResponse{
		AgentTrigger: trigger,
		FullToken:    token,
	}
	return dataResponse(c, http.StatusCreated, resp)
}

func (h *AgentTriggerHandler) Get(c echo.Context) error {
	trigger, ok := h.store.GetAgentTrigger(c.Param("triggerId"))
	if !ok {
		return notFound(c, "agent trigger not found")
	}
	return dataResponse(c, http.StatusOK, trigger)
}

func (h *AgentTriggerHandler) Update(c echo.Context) error {
	trigger, ok := h.store.GetAgentTrigger(c.Param("triggerId"))
	if !ok {
		return notFound(c, "agent trigger not found")
	}

	var req model.UpdateAgentTriggerRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if req.Name != nil {
		trigger.Name = *req.Name
	}
	if req.IsActive != nil {
		trigger.IsActive = *req.IsActive
	}
	if req.Source != nil {
		trigger.Source = *req.Source
	}
	if req.SourceConfig != nil {
		trigger.SourceConfig = req.SourceConfig
	}
	if req.ReplyConfig != nil {
		trigger.ReplyConfig = req.ReplyConfig
	}
	trigger.UpdatedAt = time.Now()

	h.store.UpdateAgentTrigger(trigger)
	return dataResponse(c, http.StatusOK, trigger)
}

func (h *AgentTriggerHandler) Delete(c echo.Context) error {
	if !h.store.DeleteAgentTrigger(c.Param("triggerId")) {
		return notFound(c, "agent trigger not found")
	}
	return c.NoContent(http.StatusNoContent)
}

// Hook is the public endpoint that receives webhook calls and executes the agent.
func (h *AgentTriggerHandler) Hook(c echo.Context) error {
	token := c.Param("token")
	trigger, ok := h.store.GetAgentTriggerByToken(token)
	if !ok {
		return notFound(c, "invalid webhook token")
	}

	if !trigger.IsActive {
		return badRequest(c, "agent trigger is disabled")
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
	h.store.UpdateAgentTrigger(trigger)

	return dataResponse(c, http.StatusCreated, execution)
}
