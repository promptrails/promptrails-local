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
	if req.Type != "agent" && req.Type != "workflow" {
		return badRequest(c, "type must be 'agent' or 'workflow'")
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
	applyMaskingEnabled(req.MaskingEnabled, &agent.MaskingEnabled)
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
		ID:             ksuid.New().String(),
		AgentID:        agentID,
		Version:        req.Version,
		Config:         req.Config,
		InputSchema:    req.InputSchema,
		OutputSchema:   req.OutputSchema,
		IsCurrent:      req.SetCurrent,
		Message:        req.Message,
		ModelConfig:    req.ModelConfig,
		RunBudget:      req.RunBudget,
		ApprovalPolicy: req.ApprovalPolicy,
		CacheTimeout:   req.CacheTimeout,
		VFSEnabled:     req.VFSEnabled,
		MaskingEnabled: req.MaskingEnabled,
		CreatedAt:      time.Now(),
	}

	// Attach MCP tools with per-tool approval/retry policy, assigning IDs.
	for _, t := range req.Tools {
		tool := t
		if tool.ID == "" {
			tool.ID = ksuid.New().String()
		}
		version.Tools = append(version.Tools, tool)
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

	version := h.resolveVersion(agent, req.VersionID)
	execution := h.newExecution(agent, sessionID, req.UserID, req.VersionID, req.Input, versionRequiresApproval(version))
	h.store.CreateExecution(execution)

	// Create trace spans for this execution.
	traces := fake.CreateExecutionTrace(execution, agent.Name)
	for _, t := range traces {
		h.store.CreateTrace(t)
	}

	return dataResponse(c, http.StatusCreated, execution)
}

// Playground runs a persisted agent version with an ephemeral prompt snapshot.
// The override is recorded on the execution metadata but does not create or
// promote any prompt/agent version.
func (h *AgentHandler) Playground(c echo.Context) error {
	agentID := c.Param("agentId")
	agent, ok := h.store.GetAgent(agentID)
	if !ok {
		return notFound(c, "agent not found")
	}

	var req model.PlaygroundRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	sessionID := ksuid.New().String()
	version := h.resolveVersion(agent, req.VersionID)
	execution := h.newExecution(agent, sessionID, "", req.VersionID, req.Input, versionRequiresApproval(version))

	// Record the ephemeral prompt override on the execution metadata.
	meta, _ := json.Marshal(map[string]any{
		"playground":      true,
		"prompt_override": req.PromptOverride,
	})
	execution.Metadata = meta

	h.store.CreateExecution(execution)
	traces := fake.CreateExecutionTrace(execution, agent.Name)
	for _, t := range traces {
		h.store.CreateTrace(t)
	}
	return dataResponse(c, http.StatusCreated, execution)
}

// PlaygroundInfo returns the agent's current version content so a client can
// pre-fill the playground editor before running an ephemeral prompt.
func (h *AgentHandler) PlaygroundInfo(c echo.Context) error {
	agent, ok := h.store.GetAgent(c.Param("agentId"))
	if !ok {
		return notFound(c, "agent not found")
	}
	return dataResponse(c, http.StatusOK, map[string]any{
		"agent":           agent,
		"current_version": agent.CurrentVersion,
	})
}

// resolveVersion returns the agent version to run: the requested version, or
// the agent's current version when none is given.
func (h *AgentHandler) resolveVersion(agent model.Agent, versionID string) *model.AgentVersion {
	if versionID != "" {
		if v, ok := h.store.GetAgentVersion(versionID); ok {
			return &v
		}
	}
	return agent.CurrentVersion
}

// newExecution builds a simulated execution. When the version parks at an
// approval-gated call the run stops at waiting_approval instead of completing.
func (h *AgentHandler) newExecution(agent model.Agent, sessionID, userID, versionID string, input map[string]any, park bool) model.Execution {
	inputJSON, _ := json.Marshal(input)
	now := time.Now()

	execution := model.Execution{
		ID:          ksuid.New().String(),
		AgentID:     &agent.ID,
		WorkspaceID: agent.WorkspaceID,
		SessionID:   sessionID,
		Input:       inputJSON,
		TraceID:     ksuid.New().String(),
		StartedAt:   &now,
		CreatedAt:   now,
	}
	if userID != "" {
		execution.UserID = &userID
	}
	if versionID != "" {
		execution.AgentVersionID = &versionID
	}

	if park {
		expires := now.Add(24 * time.Hour)
		execution.Status = "waiting_approval"
		execution.ApprovalExpiresAt = &expires
		return execution
	}

	output := fake.GenerateExecutionOutput(agent.Name, input)
	outputJSON, _ := json.Marshal(output)
	durationMS := int64(250)
	execution.Status = "completed"
	execution.Output = outputJSON
	execution.TokenUsage = json.RawMessage(`{"prompt_tokens":150,"completion_tokens":80,"total_tokens":230}`)
	execution.Cost = 0.0023
	execution.DurationMS = &durationMS
	execution.CompletedAt = &now
	return execution
}

// versionRequiresApproval reports whether a run of this version parks at an
// approval-gated call: either an explicit approval policy or a tool/sub-agent
// flagged requires_approval.
func versionRequiresApproval(v *model.AgentVersion) bool {
	if v == nil {
		return false
	}
	if v.ApprovalPolicy != nil {
		return true
	}
	for _, t := range v.Tools {
		if t.RequiresApproval {
			return true
		}
	}
	return false
}
