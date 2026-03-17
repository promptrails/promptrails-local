package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/promptrails/promptrails-local/internal/handler"
	"github.com/promptrails/promptrails-local/internal/store"
	"go.uber.org/zap"
)

const scalarHTML = `<!DOCTYPE html>
<html>
<head>
  <title>PromptRails Local - API Docs</title>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
  <script id="api-reference" data-url="/openapi.json"></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`

// New creates and configures the Echo server with all routes.
func New(s *store.Store, logger *zap.Logger, corsOrigins string, version string) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{corsOrigins},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	h := handler.New(s, logger, version)

	// Health & docs
	e.GET("/health", h.Health)
	e.GET("/docs", func(c echo.Context) error {
		return c.HTML(http.StatusOK, scalarHTML)
	})
	e.GET("/openapi.json", func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, []byte(handler.OpenAPISpec))
	})

	// Admin routes
	e.POST("/admin/reset", h.AdminReset)
	e.POST("/admin/seed", h.AdminSeed)
	e.GET("/admin/store/stats", h.AdminStats)

	// Webhook hooks (no workspace middleware)
	e.POST("/hooks/:token", h.HookExecute)

	// API v1 routes
	api := e.Group("/api/v1")

	// Agents
	api.GET("/agents", h.ListAgents)
	api.POST("/agents", h.CreateAgent)
	api.GET("/agents/:agentId", h.GetAgent)
	api.PATCH("/agents/:agentId", h.UpdateAgent)
	api.DELETE("/agents/:agentId", h.DeleteAgent)
	api.GET("/agents/:agentId/versions", h.ListAgentVersions)
	api.POST("/agents/:agentId/versions", h.CreateAgentVersion)
	api.PUT("/agents/:agentId/versions/:versionId/promote", h.PromoteAgentVersion)
	api.POST("/agents/:agentId/execute", h.ExecuteAgent)
	api.POST("/agents/:agentId/preview", h.PreviewAgent)

	// Prompts
	api.GET("/prompts", h.ListPrompts)
	api.POST("/prompts", h.CreatePrompt)
	api.GET("/prompts/:promptId", h.GetPrompt)
	api.PATCH("/prompts/:promptId", h.UpdatePrompt)
	api.DELETE("/prompts/:promptId", h.DeletePrompt)
	api.GET("/prompts/:promptId/versions", h.ListPromptVersions)
	api.POST("/prompts/:promptId/versions", h.CreatePromptVersion)
	api.PUT("/prompts/:promptId/versions/:versionId/promote", h.PromotePromptVersion)
	api.POST("/prompts/:promptId/preview", h.PreviewPrompt)
	api.POST("/prompts/:promptId/run", h.RunPrompt)

	// Executions
	api.GET("/executions", h.ListExecutions)
	api.GET("/executions/:executionId", h.GetExecution)
	api.GET("/executions/:executionId/pending-approval", h.GetPendingApproval)

	// Data Sources
	api.GET("/data-sources", h.ListDataSources)
	api.POST("/data-sources", h.CreateDataSource)
	api.GET("/data-sources/:dataSourceId", h.GetDataSource)
	api.PATCH("/data-sources/:dataSourceId", h.UpdateDataSource)
	api.DELETE("/data-sources/:dataSourceId", h.DeleteDataSource)
	api.GET("/data-sources/:dataSourceId/versions", h.ListDataSourceVersions)
	api.POST("/data-sources/:dataSourceId/versions", h.CreateDataSourceVersion)
	api.POST("/data-sources/:dataSourceId/query", h.QueryDataSource)

	// Credentials
	api.GET("/credentials", h.ListCredentials)
	api.POST("/credentials", h.CreateCredential)
	api.GET("/credentials/:credentialId", h.GetCredential)
	api.PATCH("/credentials/:credentialId", h.UpdateCredential)
	api.DELETE("/credentials/:credentialId", h.DeleteCredential)

	// Chat
	api.GET("/chat/sessions", h.ListChatSessions)
	api.POST("/chat/sessions", h.CreateChatSession)
	api.GET("/chat/sessions/:sessionId", h.GetChatSession)
	api.DELETE("/chat/sessions/:sessionId", h.DeleteChatSession)
	api.GET("/chat/sessions/:sessionId/messages", h.ListChatMessages)
	api.POST("/chat/sessions/:sessionId/messages", h.SendChatMessage)

	// Traces
	api.GET("/traces", h.ListTraces)
	api.GET("/traces/:traceId", h.GetTrace)

	// Scores
	api.GET("/scores", h.ListScores)
	api.POST("/scores", h.CreateScore)
	api.GET("/scores/:scoreId", h.GetScore)
	api.PATCH("/scores/:scoreId", h.UpdateScore)
	api.DELETE("/scores/:scoreId", h.DeleteScore)

	// Score Configs
	api.GET("/score-configs", h.ListScoreConfigs)
	api.POST("/score-configs", h.CreateScoreConfig)
	api.GET("/score-configs/:configId", h.GetScoreConfig)
	api.PATCH("/score-configs/:configId", h.UpdateScoreConfig)
	api.DELETE("/score-configs/:configId", h.DeleteScoreConfig)

	// Approvals
	api.GET("/approvals", h.ListApprovals)
	api.GET("/approvals/:approvalId", h.GetApproval)
	api.POST("/approvals/:approvalId/decide", h.DecideApproval)

	// Webhook Triggers
	api.GET("/webhook-triggers", h.ListWebhookTriggers)
	api.POST("/webhook-triggers", h.CreateWebhookTrigger)
	api.GET("/webhook-triggers/:triggerId", h.GetWebhookTrigger)
	api.PATCH("/webhook-triggers/:triggerId", h.UpdateWebhookTrigger)
	api.DELETE("/webhook-triggers/:triggerId", h.DeleteWebhookTrigger)

	// MCP Tools
	api.GET("/mcp-tools", h.ListMCPTools)
	api.POST("/mcp-tools", h.CreateMCPTool)
	api.GET("/mcp-tools/:toolId", h.GetMCPTool)
	api.PATCH("/mcp-tools/:toolId", h.UpdateMCPTool)
	api.DELETE("/mcp-tools/:toolId", h.DeleteMCPTool)

	// MCP Templates
	api.GET("/mcp-templates", h.ListMCPTemplates)
	api.GET("/mcp-templates/:templateId", h.GetMCPTemplate)

	// Guardrails
	api.GET("/agents/:agentId/guardrails", h.ListGuardrails)
	api.POST("/agents/:agentId/guardrails", h.CreateGuardrail)
	api.PATCH("/guardrails/:guardrailId", h.UpdateGuardrail)
	api.DELETE("/guardrails/:guardrailId", h.DeleteGuardrail)

	// Memories
	api.GET("/agents/:agentId/memories", h.ListMemories)
	api.POST("/agents/:agentId/memories", h.CreateMemory)
	api.DELETE("/agents/:agentId/memories", h.DeleteAllMemories)
	api.POST("/agents/:agentId/memories/search", h.SearchMemories)
	api.GET("/memories/:memoryId", h.GetMemory)
	api.DELETE("/memories/:memoryId", h.DeleteMemory)

	// LLM Models
	api.GET("/llm-models", h.ListLLMModels)

	return e
}
