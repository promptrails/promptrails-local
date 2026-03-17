package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
	"go.uber.org/zap"
)

// Handlers aggregates all sub-handlers for the emulator.
type Handlers struct {
	Agent          *AgentHandler
	Prompt         *PromptHandler
	Execution      *ExecutionHandler
	DataSource     *DataSourceHandler
	Credential     *CredentialHandler
	Chat           *ChatHandler
	Trace          *TraceHandler
	Score          *ScoreHandler
	Approval       *ApprovalHandler
	WebhookTrigger *WebhookTriggerHandler
	MCPTool        *MCPToolHandler
	Guardrail      *GuardrailHandler
	Memory         *MemoryHandler
	LLMModel       *LLMModelHandler
	Admin          *AdminHandler

	version string
}

// New creates all sub-handlers wired to the given store.
func New(s *store.Store, logger *zap.Logger, version string) *Handlers {
	return &Handlers{
		Agent:          &AgentHandler{store: s, logger: logger},
		Prompt:         &PromptHandler{store: s, logger: logger},
		Execution:      &ExecutionHandler{store: s},
		DataSource:     &DataSourceHandler{store: s},
		Credential:     &CredentialHandler{store: s},
		Chat:           &ChatHandler{store: s},
		Trace:          &TraceHandler{store: s},
		Score:          &ScoreHandler{store: s},
		Approval:       &ApprovalHandler{store: s},
		WebhookTrigger: &WebhookTriggerHandler{store: s},
		MCPTool:        &MCPToolHandler{store: s},
		Guardrail:      &GuardrailHandler{store: s},
		Memory:         &MemoryHandler{store: s},
		LLMModel:       &LLMModelHandler{store: s},
		Admin:          &AdminHandler{store: s, logger: logger},
		version:        version,
	}
}

// ---------- Health ----------

// Health returns a simple health check response.
func (h *Handlers) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"version": h.version,
	})
}

// ---------- Admin delegations ----------

func (h *Handlers) AdminReset(c echo.Context) error { return h.Admin.Reset(c) }
func (h *Handlers) AdminSeed(c echo.Context) error  { return h.Admin.Seed(c) }
func (h *Handlers) AdminStats(c echo.Context) error { return h.Admin.Stats(c) }

// ---------- Agent delegations ----------

func (h *Handlers) ListAgents(c echo.Context) error          { return h.Agent.List(c) }
func (h *Handlers) CreateAgent(c echo.Context) error         { return h.Agent.Create(c) }
func (h *Handlers) GetAgent(c echo.Context) error            { return h.Agent.Get(c) }
func (h *Handlers) UpdateAgent(c echo.Context) error         { return h.Agent.Update(c) }
func (h *Handlers) DeleteAgent(c echo.Context) error         { return h.Agent.Delete(c) }
func (h *Handlers) ListAgentVersions(c echo.Context) error   { return h.Agent.ListVersions(c) }
func (h *Handlers) CreateAgentVersion(c echo.Context) error  { return h.Agent.CreateVersion(c) }
func (h *Handlers) PromoteAgentVersion(c echo.Context) error { return h.Agent.PromoteVersion(c) }
func (h *Handlers) ExecuteAgent(c echo.Context) error        { return h.Agent.Execute(c) }
func (h *Handlers) PreviewAgent(c echo.Context) error        { return h.Agent.Preview(c) }

// ---------- Prompt delegations ----------

func (h *Handlers) ListPrompts(c echo.Context) error          { return h.Prompt.List(c) }
func (h *Handlers) CreatePrompt(c echo.Context) error         { return h.Prompt.Create(c) }
func (h *Handlers) GetPrompt(c echo.Context) error            { return h.Prompt.Get(c) }
func (h *Handlers) UpdatePrompt(c echo.Context) error         { return h.Prompt.Update(c) }
func (h *Handlers) DeletePrompt(c echo.Context) error         { return h.Prompt.Delete(c) }
func (h *Handlers) ListPromptVersions(c echo.Context) error   { return h.Prompt.ListVersions(c) }
func (h *Handlers) CreatePromptVersion(c echo.Context) error  { return h.Prompt.CreateVersion(c) }
func (h *Handlers) PromotePromptVersion(c echo.Context) error { return h.Prompt.PromoteVersion(c) }
func (h *Handlers) PreviewPrompt(c echo.Context) error        { return h.Prompt.Preview(c) }
func (h *Handlers) RunPrompt(c echo.Context) error            { return h.Prompt.Run(c) }

// ---------- Execution delegations ----------

func (h *Handlers) ListExecutions(c echo.Context) error     { return h.Execution.List(c) }
func (h *Handlers) GetExecution(c echo.Context) error       { return h.Execution.Get(c) }
func (h *Handlers) GetPendingApproval(c echo.Context) error { return h.Execution.GetPendingApproval(c) }

// ---------- Data Source delegations ----------

func (h *Handlers) ListDataSources(c echo.Context) error        { return h.DataSource.List(c) }
func (h *Handlers) CreateDataSource(c echo.Context) error       { return h.DataSource.Create(c) }
func (h *Handlers) GetDataSource(c echo.Context) error          { return h.DataSource.Get(c) }
func (h *Handlers) UpdateDataSource(c echo.Context) error       { return h.DataSource.Update(c) }
func (h *Handlers) DeleteDataSource(c echo.Context) error       { return h.DataSource.Delete(c) }
func (h *Handlers) ListDataSourceVersions(c echo.Context) error { return h.DataSource.ListVersions(c) }
func (h *Handlers) CreateDataSourceVersion(c echo.Context) error {
	return h.DataSource.CreateVersion(c)
}
func (h *Handlers) QueryDataSource(c echo.Context) error { return h.DataSource.Query(c) }

// ---------- Credential delegations ----------

func (h *Handlers) ListCredentials(c echo.Context) error  { return h.Credential.List(c) }
func (h *Handlers) CreateCredential(c echo.Context) error { return h.Credential.Create(c) }
func (h *Handlers) GetCredential(c echo.Context) error    { return h.Credential.Get(c) }
func (h *Handlers) UpdateCredential(c echo.Context) error { return h.Credential.Update(c) }
func (h *Handlers) DeleteCredential(c echo.Context) error { return h.Credential.Delete(c) }

// ---------- Chat delegations ----------

func (h *Handlers) ListChatSessions(c echo.Context) error  { return h.Chat.ListSessions(c) }
func (h *Handlers) CreateChatSession(c echo.Context) error { return h.Chat.CreateSession(c) }
func (h *Handlers) GetChatSession(c echo.Context) error    { return h.Chat.GetSession(c) }
func (h *Handlers) DeleteChatSession(c echo.Context) error { return h.Chat.DeleteSession(c) }
func (h *Handlers) ListChatMessages(c echo.Context) error  { return h.Chat.ListMessages(c) }
func (h *Handlers) SendChatMessage(c echo.Context) error   { return h.Chat.SendMessage(c) }

// ---------- Trace delegations ----------

func (h *Handlers) ListTraces(c echo.Context) error { return h.Trace.List(c) }
func (h *Handlers) GetTrace(c echo.Context) error   { return h.Trace.Get(c) }

// ---------- Score delegations ----------

func (h *Handlers) ListScores(c echo.Context) error        { return h.Score.List(c) }
func (h *Handlers) CreateScore(c echo.Context) error       { return h.Score.Create(c) }
func (h *Handlers) GetScore(c echo.Context) error          { return h.Score.Get(c) }
func (h *Handlers) UpdateScore(c echo.Context) error       { return h.Score.Update(c) }
func (h *Handlers) DeleteScore(c echo.Context) error       { return h.Score.Delete(c) }
func (h *Handlers) ListScoreConfigs(c echo.Context) error  { return h.Score.ListConfigs(c) }
func (h *Handlers) CreateScoreConfig(c echo.Context) error { return h.Score.CreateConfig(c) }
func (h *Handlers) GetScoreConfig(c echo.Context) error    { return h.Score.GetConfig(c) }
func (h *Handlers) UpdateScoreConfig(c echo.Context) error { return h.Score.UpdateConfig(c) }
func (h *Handlers) DeleteScoreConfig(c echo.Context) error { return h.Score.DeleteConfig(c) }

// ---------- Approval delegations ----------

func (h *Handlers) ListApprovals(c echo.Context) error  { return h.Approval.List(c) }
func (h *Handlers) GetApproval(c echo.Context) error    { return h.Approval.Get(c) }
func (h *Handlers) DecideApproval(c echo.Context) error { return h.Approval.Decide(c) }

// ---------- Webhook Trigger delegations ----------

func (h *Handlers) ListWebhookTriggers(c echo.Context) error  { return h.WebhookTrigger.List(c) }
func (h *Handlers) CreateWebhookTrigger(c echo.Context) error { return h.WebhookTrigger.Create(c) }
func (h *Handlers) GetWebhookTrigger(c echo.Context) error    { return h.WebhookTrigger.Get(c) }
func (h *Handlers) UpdateWebhookTrigger(c echo.Context) error { return h.WebhookTrigger.Update(c) }
func (h *Handlers) DeleteWebhookTrigger(c echo.Context) error { return h.WebhookTrigger.Delete(c) }
func (h *Handlers) HookExecute(c echo.Context) error          { return h.WebhookTrigger.Hook(c) }

// ---------- MCP Tool delegations ----------

func (h *Handlers) ListMCPTools(c echo.Context) error     { return h.MCPTool.List(c) }
func (h *Handlers) CreateMCPTool(c echo.Context) error    { return h.MCPTool.Create(c) }
func (h *Handlers) GetMCPTool(c echo.Context) error       { return h.MCPTool.Get(c) }
func (h *Handlers) UpdateMCPTool(c echo.Context) error    { return h.MCPTool.Update(c) }
func (h *Handlers) DeleteMCPTool(c echo.Context) error    { return h.MCPTool.Delete(c) }
func (h *Handlers) ListMCPTemplates(c echo.Context) error { return h.MCPTool.ListTemplates(c) }
func (h *Handlers) GetMCPTemplate(c echo.Context) error   { return h.MCPTool.GetTemplate(c) }

// ---------- Guardrail delegations ----------

func (h *Handlers) ListGuardrails(c echo.Context) error  { return h.Guardrail.List(c) }
func (h *Handlers) CreateGuardrail(c echo.Context) error { return h.Guardrail.Create(c) }
func (h *Handlers) UpdateGuardrail(c echo.Context) error { return h.Guardrail.Update(c) }
func (h *Handlers) DeleteGuardrail(c echo.Context) error { return h.Guardrail.Delete(c) }

// ---------- Memory delegations ----------

func (h *Handlers) ListMemories(c echo.Context) error      { return h.Memory.List(c) }
func (h *Handlers) CreateMemory(c echo.Context) error      { return h.Memory.Create(c) }
func (h *Handlers) SearchMemories(c echo.Context) error    { return h.Memory.Search(c) }
func (h *Handlers) GetMemory(c echo.Context) error         { return h.Memory.Get(c) }
func (h *Handlers) DeleteMemory(c echo.Context) error      { return h.Memory.Delete(c) }
func (h *Handlers) DeleteAllMemories(c echo.Context) error { return h.Memory.DeleteAll(c) }

// ---------- LLM Model delegations ----------

func (h *Handlers) ListLLMModels(c echo.Context) error { return h.LLMModel.List(c) }

// ---------- helpers ----------

// getWorkspaceID returns a constant workspace ID used when creating new resources
// so the response JSON has a workspace_id field. The emulator uses a single flat namespace.
func getWorkspaceID() string {
	return "2N8Dx8a2a1qUspBKbP02vqP29eH"
}

// pagination holds page/limit values parsed from query parameters.
type pagination struct {
	Page  int
	Limit int
}

// getPagination reads page and limit query params with sensible defaults.
func getPagination(c echo.Context) pagination {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return pagination{Page: page, Limit: limit}
}

// dataResponse wraps a single item in the standard envelope.
func dataResponse(c echo.Context, status int, data any) error {
	return c.JSON(status, model.DataResponse{Data: data})
}

// listResponse wraps a slice with pagination metadata.
func listResponse(c echo.Context, data any, total int, p pagination) error {
	totalPages := int(math.Ceil(float64(total) / float64(p.Limit)))
	return c.JSON(http.StatusOK, model.ListResponse{
		Data: data,
		Meta: &model.PaginationMeta{
			Total:      total,
			Page:       p.Page,
			Limit:      p.Limit,
			TotalPages: totalPages,
		},
	})
}

// notFound returns a 404 error response.
func notFound(c echo.Context, msg string) error {
	return c.JSON(http.StatusNotFound, model.ErrorResponse{
		Error: model.ErrorDetail{Code: "not_found", Message: msg},
	})
}

// badRequest returns a 400 error response.
func badRequest(c echo.Context, msg string) error {
	return c.JSON(http.StatusBadRequest, model.ErrorResponse{
		Error: model.ErrorDetail{Code: "bad_request", Message: msg},
	})
}
