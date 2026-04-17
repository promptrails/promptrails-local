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
	"go.uber.org/zap"
)

type AgentHandler struct {
	store  *store.Store
	logger *zap.Logger
}

func (h *AgentHandler) List(c echo.Context) error {
	p := getPagination(c)

	filters := store.AgentFilters{
		Type:   c.QueryParam("type"),
		Status: c.QueryParam("status"),
		Name:   c.QueryParam("name"),
	}

	agents, total := h.store.ListAgents(filters, p.Page, p.Limit)
	return listResponse(c, agents, total, p)
}

func (h *AgentHandler) Create(c echo.Context) error {
	wid := getWorkspaceID()
	var req model.CreateAgentRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.Name == "" || req.Type == "" {
		return badRequest(c, "name and type are required")
	}

	now := time.Now()
	agent := model.Agent{
		ID:          ksuid.New().String(),
		WorkspaceID: wid,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Status:      "active",
		Labels:      json.RawMessage("[]"),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	h.store.CreateAgent(agent)
	return dataResponse(c, http.StatusCreated, agent)
}

func (h *AgentHandler) Get(c echo.Context) error {
	agent, ok := h.store.GetAgent(c.Param("agentId"))
	if !ok {
		return notFound(c, "agent not found")
	}
	return dataResponse(c, http.StatusOK, agent)
}

func (h *AgentHandler) Update(c echo.Context) error {
	agent, ok := h.store.GetAgent(c.Param("agentId"))
	if !ok {
		return notFound(c, "agent not found")
	}

	var req model.UpdateAgentRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if req.Name != nil {
		agent.Name = *req.Name
	}
	if req.Description != nil {
		agent.Description = *req.Description
	}
	if req.Status != nil {
		agent.Status = *req.Status
	}
	agent.UpdatedAt = time.Now()

	h.store.UpdateAgent(agent)
	return dataResponse(c, http.StatusOK, agent)
}

func (h *AgentHandler) Delete(c echo.Context) error {
	if !h.store.DeleteAgent(c.Param("agentId")) {
		return notFound(c, "agent not found")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *AgentHandler) ListVersions(c echo.Context) error {
	agentID := c.Param("agentId")
	if _, ok := h.store.GetAgent(agentID); !ok {
		return notFound(c, "agent not found")
	}
	versions := h.store.ListAgentVersions(agentID)
	return dataResponse(c, http.StatusOK, versions)
}

func (h *AgentHandler) CreateVersion(c echo.Context) error {
	agentID := c.Param("agentId")
	if _, ok := h.store.GetAgent(agentID); !ok {
		return notFound(c, "agent not found")
	}

	var req model.CreateAgentVersionRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	version := model.AgentVersion{
		ID:           ksuid.New().String(),
		AgentID:      agentID,
		Version:      req.Version,
		Config:       req.Config,
		InputSchema:  req.InputSchema,
		OutputSchema: req.OutputSchema,
		IsCurrent:    req.SetCurrent,
		Message:      req.Message,
		CreatedAt:    time.Now(),
	}

	// Build prompt associations. Links by prompt_id so executions auto-follow
	// the prompt's current version rather than pinning a specific version.
	for _, pid := range req.PromptIDs {
		version.Prompts = append(version.Prompts, model.AgentVersionPrompt{
			ID:             ksuid.New().String(),
			AgentVersionID: version.ID,
			PromptID:       pid.PromptID,
			Role:           pid.Role,
			SortOrder:      pid.SortOrder,
		})
	}

	if req.SetCurrent {
		h.store.DemoteAgentVersions(agentID)
	}

	h.store.CreateAgentVersion(version)
	return dataResponse(c, http.StatusCreated, version)
}

func (h *AgentHandler) PromoteVersion(c echo.Context) error {
	agentID := c.Param("agentId")
	versionID := c.Param("versionId")

	version, ok := h.store.GetAgentVersion(versionID)
	if !ok || version.AgentID != agentID {
		return notFound(c, "version not found")
	}

	h.store.DemoteAgentVersions(agentID)
	version.IsCurrent = true
	h.store.UpdateAgentVersion(version)
	return dataResponse(c, http.StatusOK, version)
}

func (h *AgentHandler) Preview(c echo.Context) error {
	agentID := c.Param("agentId")
	agent, ok := h.store.GetAgent(agentID)
	if !ok {
		return notFound(c, "agent not found")
	}

	var req model.PreviewAgentRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	result := fake.GenerateExecutionOutput(agent.Name, req.Input)
	return dataResponse(c, http.StatusOK, result)
}

func (h *AgentHandler) Execute(c echo.Context) error {
	agentID := c.Param("agentId")
	wid := getWorkspaceID()
	agent, ok := h.store.GetAgent(agentID)
	if !ok {
		return notFound(c, "agent not found")
	}

	var req model.ExecuteAgentRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = ksuid.New().String()
	}

	// Generate fake output and trace
	output := fake.GenerateExecutionOutput(agent.Name, req.Input)
	inputJSON, _ := json.Marshal(req.Input)
	outputJSON, _ := json.Marshal(output)

	now := time.Now()
	durationMS := int64(250)
	execution := model.Execution{
		ID:          ksuid.New().String(),
		AgentID:     &agentID,
		WorkspaceID: wid,
		SessionID:   sessionID,
		Status:      "completed",
		Input:       inputJSON,
		Output:      outputJSON,
		TokenUsage:  json.RawMessage(`{"prompt_tokens":150,"completion_tokens":80,"total_tokens":230}`),
		Cost:        0.0023,
		DurationMS:  &durationMS,
		TraceID:     ksuid.New().String(),
		StartedAt:   &now,
		CompletedAt: &now,
		CreatedAt:   now,
	}

	if req.UserID != "" {
		execution.UserID = &req.UserID
	}
	if req.VersionID != "" {
		execution.AgentVersionID = &req.VersionID
	}

	h.store.CreateExecution(execution)

	// Create trace spans for this execution
	traces := fake.CreateExecutionTrace(execution, agent.Name)
	for _, t := range traces {
		h.store.CreateTrace(t)
	}

	return dataResponse(c, http.StatusCreated, execution)
}
