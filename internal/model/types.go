package model

import (
	"encoding/json"
	"time"
)

type Agent struct {
	ID          string          `json:"id"`
	WorkspaceID string          `json:"workspace_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Type        string          `json:"type"`
	Status      string          `json:"status"`
	Labels      json.RawMessage `json:"labels"`
	// MaskingEnabled overrides the workspace PII masking policy for this
	// agent. nil = inherit; true/false = explicit override. Matches the
	// shape served by the real API.
	MaskingEnabled *bool         `json:"masking_enabled,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	CurrentVersion *AgentVersion `json:"current_version,omitempty"`
}

// AgentVersion is a committed version of an agent's runtime behavior. In the
// Agent/Workflow (v2) model the version OWNS the model + sampling
// (model_config), run budget, approval policy, cache TTL, version-scoped
// vfs/masking overrides, and the attached MCP tools. Prompt versions are pure
// content and carry no model.
type AgentVersion struct {
	ID             string               `json:"id"`
	AgentID        string               `json:"agent_id"`
	Version        string               `json:"version"`
	Message        string               `json:"message"`
	Config         json.RawMessage      `json:"config"`
	InputSchema    json.RawMessage      `json:"input_schema"`
	OutputSchema   json.RawMessage      `json:"output_schema"`
	IsCurrent      bool                 `json:"is_current"`
	ModelConfig    *ModelConfig         `json:"model_config"`
	RunBudget      *RunBudget           `json:"run_budget"`
	ApprovalPolicy *ApprovalPolicy      `json:"approval_policy"`
	CacheTimeout   *int                 `json:"cache_timeout"`
	VFSEnabled     *bool                `json:"vfs_enabled"`
	MaskingEnabled *bool                `json:"masking_enabled"`
	Tools          []AgentVersionTool   `json:"tools,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
	Prompts        []AgentVersionPrompt `json:"prompts,omitempty"`
}

// AgentVersionTool is an MCP tool attached to an agent version with per-tool
// approval / retry policy.
type AgentVersionTool struct {
	ID               string `json:"id"`
	MCPToolID        string `json:"mcp_tool_id"`
	RequiresApproval bool   `json:"requires_approval"`
	NoRetry          bool   `json:"no_retry"`
	SortOrder        int    `json:"sort_order"`
}

// ModelConfig is the version-scoped model + sampling ownership. Each field is
// optional; unset sampling inherits the provider/model default.
type ModelConfig struct {
	ModelID         *string  `json:"model_id,omitempty"`
	FallbackModelID *string  `json:"fallback_model_id,omitempty"`
	Temperature     *float64 `json:"temperature,omitempty"`
	TopP            *float64 `json:"top_p,omitempty"`
	TopK            *int     `json:"top_k,omitempty"`
	MaxTokens       *int     `json:"max_tokens,omitempty"`
}

// RunBudget bounds the whole execution tree, enforced at the root.
type RunBudget struct {
	MaxCost        *float64 `json:"max_cost,omitempty"`
	MaxTotalTokens *int     `json:"max_total_tokens,omitempty"`
	MaxToolCalls   *int     `json:"max_tool_calls,omitempty"`
	MaxChildren    *int     `json:"max_children,omitempty"`
	MaxDepth       *int     `json:"max_depth,omitempty"`
}

// ApprovalPolicy declares who may approve/deny a parked, approval-gated call.
type ApprovalPolicy struct {
	Mode      string   `json:"mode"`
	MemberIDs []string `json:"member_ids,omitempty"`
}

// GuardrailSpec is a version-scoped guardrail attached to an agent version.
type GuardrailSpec struct {
	ID          string          `json:"id,omitempty"`
	Type        string          `json:"type"`
	ScannerType string          `json:"scanner_type"`
	Action      string          `json:"action"`
	Config      json.RawMessage `json:"config,omitempty"`
	IsActive    bool            `json:"is_active"`
	SortOrder   int             `json:"sort_order"`
}

type AgentVersionPrompt struct {
	ID             string  `json:"id"`
	AgentVersionID string  `json:"agent_version_id"`
	PromptID       string  `json:"prompt_id"`
	Role           string  `json:"role"`
	SortOrder      int     `json:"sort_order"`
	Prompt         *Prompt `json:"prompt,omitempty"`
}

type Prompt struct {
	ID             string         `json:"id"`
	WorkspaceID    string         `json:"workspace_id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Status         string         `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	CurrentVersion *PromptVersion `json:"current_version,omitempty"`
}

// PromptVersion is PURE CONTENT in the Agent/Workflow (v2) model: system/user
// templates plus the declared input schema. Model, fallback, sampling, output
// schema, tools, and cache TTL now live on the agent version.
type PromptVersion struct {
	ID           string          `json:"id"`
	PromptID     string          `json:"prompt_id"`
	Version      string          `json:"version"`
	SystemPrompt string          `json:"system_prompt"`
	UserPrompt   string          `json:"user_prompt"`
	InputSchema  json.RawMessage `json:"input_schema"`
	IsCurrent    bool            `json:"is_current"`
	Message      string          `json:"message"`
	Config       json.RawMessage `json:"config"`
	CreatedAt    time.Time       `json:"created_at"`
}

type DataSource struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	// MaskingEnabled overrides the workspace PII masking policy for data
	// flowing out of this source. nil = inherit; true/false = explicit.
	MaskingEnabled *bool `json:"masking_enabled,omitempty"`
	// MaskingRules carries the per-field always/never masking overrides
	// for columns returned by this source. Stored verbatim so the
	// emulator round-trips whatever shape the real API expects without
	// re-implementing the validation.
	MaskingRules json.RawMessage `json:"masking_rules,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type DataSourceVersion struct {
	ID               string          `json:"id"`
	DataSourceID     string          `json:"data_source_id"`
	Version          string          `json:"version"`
	CredentialID     *string         `json:"credential_id"`
	ConnectionConfig json.RawMessage `json:"connection_config"`
	QueryTemplate    string          `json:"query_template"`
	Parameters       json.RawMessage `json:"parameters"`
	IsCurrent        bool            `json:"is_current"`
	Message          string          `json:"message"`
	CacheTimeout     int             `json:"cache_timeout"`
	OutputFormat     string          `json:"output_format"`
	CreatedAt        time.Time       `json:"created_at"`
}

// Execution is one agent run. In the Agent/Workflow (v2) model an execution
// may be part of a tree: parent_execution_id links a sub-agent delegation,
// handoff continuation, or workflow agent-node run to its parent, and
// children[] holds the direct (single-get) or full (tree endpoint) descendants.
type Execution struct {
	ID                string          `json:"id"`
	AgentID           *string         `json:"agent_id"`
	AgentVersionID    *string         `json:"agent_version_id"`
	ParentExecutionID *string         `json:"parent_execution_id"`
	WorkspaceID       string          `json:"workspace_id"`
	UserID            *string         `json:"user_id"`
	SessionID         string          `json:"session_id"`
	Status            string          `json:"status"`
	Input             json.RawMessage `json:"input"`
	Output            json.RawMessage `json:"output"`
	Error             string          `json:"error"`
	Metadata          json.RawMessage `json:"metadata"`
	TokenUsage        json.RawMessage `json:"token_usage"`
	Cost              float64         `json:"cost"`
	DurationMS        *int64          `json:"duration_ms"`
	TraceID           string          `json:"trace_id,omitempty"`
	ApprovalExpiresAt *time.Time      `json:"approval_expires_at,omitempty"`
	StartedAt         *time.Time      `json:"started_at"`
	CompletedAt       *time.Time      `json:"completed_at"`
	CreatedAt         time.Time       `json:"created_at"`
	Agent             *Agent          `json:"agent,omitempty"`
	Children          []Execution     `json:"children,omitempty"`
}

type Credential struct {
	ID            string    `json:"id"`
	WorkspaceID   string    `json:"workspace_id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Category      string    `json:"category"`
	Description   string    `json:"description"`
	MaskedContent string    `json:"masked_content"`
	IsDefault     bool      `json:"is_default"`
	SchemaType    string    `json:"schema_type"`
	IsValid       bool      `json:"is_valid"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ChatSession struct {
	ID          string          `json:"id"`
	WorkspaceID string          `json:"workspace_id"`
	AgentID     string          `json:"agent_id"`
	UserID      *string         `json:"user_id"`
	Title       string          `json:"title"`
	Metadata    json.RawMessage `json:"metadata"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Agent       *Agent          `json:"agent,omitempty"`
}

type ChatMessage struct {
	ID               string          `json:"id"`
	SessionID        string          `json:"session_id"`
	Role             string          `json:"role"`
	Content          string          `json:"content"`
	Metadata         json.RawMessage `json:"metadata"`
	ToolCalls        json.RawMessage `json:"tool_calls,omitempty"`
	ToolResults      json.RawMessage `json:"tool_results,omitempty"`
	Model            string          `json:"model,omitempty"`
	Cost             *float64        `json:"cost,omitempty"`
	TokenCount       *int            `json:"token_count"`
	PromptTokens     *int            `json:"prompt_tokens,omitempty"`
	CompletionTokens *int            `json:"completion_tokens,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
}

type Trace struct {
	ID           string          `json:"id"`
	WorkspaceID  string          `json:"workspace_id"`
	TraceID      string          `json:"trace_id"`
	SpanID       string          `json:"span_id"`
	ParentSpanID string          `json:"parent_span_id"`
	Name         string          `json:"name"`
	Kind         string          `json:"kind"`
	Status       string          `json:"status"`
	Level        string          `json:"level"`
	Input        json.RawMessage `json:"input"`
	Output       json.RawMessage `json:"output"`
	Attributes   json.RawMessage `json:"attributes"`
	Tags         json.RawMessage `json:"tags"`
	TokenUsage   json.RawMessage `json:"token_usage"`
	Cost         *float64        `json:"cost"`
	DurationMS   *int            `json:"duration_ms"`
	ErrorMessage string          `json:"error_message,omitempty"`
	ModelName    string          `json:"model_name,omitempty"`
	AgentID      *string         `json:"agent_id,omitempty"`
	ExecutionID  *string         `json:"execution_id,omitempty"`
	SessionID    string          `json:"session_id,omitempty"`
	StartedAt    time.Time       `json:"started_at"`
	EndedAt      *time.Time      `json:"ended_at"`
	CreatedAt    time.Time       `json:"created_at"`
}

// TraceSummary is the aggregate metering rollup returned by GET
// /traces/summary (the v2 replacement for the removed cost summary).
type TraceSummary struct {
	TotalTraces    int64   `json:"total_traces"`
	TotalTokens    int64   `json:"total_tokens"`
	TotalCost      float64 `json:"total_cost"`
	AvgDurationMS  float64 `json:"avg_duration_ms"`
	ErrorCount     int64   `json:"error_count"`
	UniqueModels   int64   `json:"unique_models"`
	UniqueSessions int64   `json:"unique_sessions"`
}

type AgentTrigger struct {
	ID           string                 `json:"id"`
	WorkspaceID  string                 `json:"workspace_id"`
	AgentID      string                 `json:"agent_id"`
	Name         string                 `json:"name"`
	Token        string                 `json:"-"`
	TokenPrefix  string                 `json:"token_prefix"`
	Source       string                 `json:"source"`
	SourceConfig map[string]interface{} `json:"source_config,omitempty"`
	ReplyConfig  map[string]interface{} `json:"reply_config,omitempty"`
	IsActive     bool                   `json:"is_active"`
	LastUsedAt   *time.Time             `json:"last_used_at"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Agent        *Agent                 `json:"agent,omitempty"`
}

// AgentVFSFile is a single entry in an agent's Virtual Filesystem.
type AgentVFSFile struct {
	ID             string                 `json:"id"`
	WorkspaceID    string                 `json:"workspace_id"`
	AgentID        string                 `json:"agent_id"`
	Path           string                 `json:"path"`
	ParentPath     string                 `json:"parent_path"`
	Name           string                 `json:"name"`
	IsDir          bool                   `json:"is_dir"`
	Content        string                 `json:"content,omitempty"`
	SizeBytes      int64                  `json:"size_bytes"`
	MimeType       string                 `json:"mime_type,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	LastWriterKind string                 `json:"last_writer_kind"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type MCPTool struct {
	ID           string          `json:"id"`
	WorkspaceID  string          `json:"workspace_id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Type         string          `json:"type"`
	Config       json.RawMessage `json:"config"`
	Schema       json.RawMessage `json:"schema"`
	CredentialID *string         `json:"credential_id"`
	TemplateID   *string         `json:"template_id"`
	Status       string          `json:"status"`
	IsActive     bool            `json:"is_active"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type Guardrail struct {
	ID          string          `json:"id"`
	AgentID     string          `json:"agent_id"`
	Type        string          `json:"type"`
	ScannerType string          `json:"scanner_type"`
	Action      string          `json:"action"`
	Config      json.RawMessage `json:"config"`
	IsActive    bool            `json:"is_active"`
	SortOrder   int             `json:"sort_order"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type LLMModel struct {
	ID                    string     `json:"id"`
	Provider              string     `json:"provider"`
	ModelID               string     `json:"model_id"`
	DisplayName           string     `json:"display_name"`
	InputPrice            *float64   `json:"input_price"`
	OutputPrice           *float64   `json:"output_price"`
	CachedInputPrice      *float64   `json:"cached_input_price"`
	MaxTokens             *int       `json:"max_tokens"`
	SupportsVision        bool       `json:"supports_vision"`
	SupportsTools         bool       `json:"supports_tools"`
	SupportsJSON          bool       `json:"supports_json"`
	SupportsStreaming     bool       `json:"supports_streaming"`
	SupportsTemperature   bool       `json:"supports_temperature"`
	SupportsTopP          bool       `json:"supports_top_p"`
	SupportsTopK          bool       `json:"supports_top_k"`
	SupportsReasoning     bool       `json:"supports_reasoning"`
	SupportsWebSearch     bool       `json:"supports_web_search"`
	SupportsPromptCaching bool       `json:"supports_prompt_caching"`
	IsActive              bool       `json:"is_active"`
	IsDeprecated          bool       `json:"is_deprecated"`
	DeprecatedAt          *time.Time `json:"deprecated_at"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// AvailableModelEntry is one model in the available-models response, grouped
// by provider.
type AvailableModelEntry struct {
	ID                    string   `json:"id"`
	ModelID               string   `json:"model_id"`
	DisplayName           string   `json:"display_name"`
	MaxTokens             *int     `json:"max_tokens"`
	SupportsVision        bool     `json:"supports_vision"`
	SupportsTools         bool     `json:"supports_tools"`
	SupportsJSON          bool     `json:"supports_json"`
	SupportsTemperature   bool     `json:"supports_temperature"`
	SupportsTopP          bool     `json:"supports_top_p"`
	SupportsTopK          bool     `json:"supports_top_k"`
	SupportsReasoning     bool     `json:"supports_reasoning"`
	SupportsWebSearch     bool     `json:"supports_web_search"`
	SupportsPromptCaching bool     `json:"supports_prompt_caching"`
	InputPrice            *float64 `json:"input_price"`
	OutputPrice           *float64 `json:"output_price"`
	IsDeprecated          bool     `json:"is_deprecated"`
}

// AvailableModelGroup groups available models by provider.
type AvailableModelGroup struct {
	Provider string                `json:"provider"`
	Models   []AvailableModelEntry `json:"models"`
}

type MCPTemplate struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	Type        string          `json:"type"`
	Category    string          `json:"category"`
	Config      json.RawMessage `json:"config"`
	Schema      json.RawMessage `json:"schema"`
	IconURL     string          `json:"icon_url"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// Request types

type CreateAgentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Type is the agent kind: "agent" (prompt + tools + optional sub-agents)
	// or "workflow" (deterministic DAG). These are the only valid v2 types.
	Type string `json:"type"`
}

type UpdateAgentRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	// MaskingEnabled is tri-state: omitted = leave unchanged, JSON null
	// = clear the override (back to inherit), true/false = set. Stored as
	// json.RawMessage because Go's encoding/json collapses "field
	// omitted" and "explicit null" into the same nil for *bool.
	MaskingEnabled json.RawMessage `json:"masking_enabled,omitempty"`
}

type CreateAgentVersionRequest struct {
	Version        string             `json:"version"`
	Config         json.RawMessage    `json:"config"`
	InputSchema    json.RawMessage    `json:"input_schema"`
	OutputSchema   json.RawMessage    `json:"output_schema"`
	SetCurrent     bool               `json:"set_current"`
	Message        string             `json:"message"`
	ModelConfig    *ModelConfig       `json:"model_config"`
	RunBudget      *RunBudget         `json:"run_budget"`
	ApprovalPolicy *ApprovalPolicy    `json:"approval_policy"`
	CacheTimeout   *int               `json:"cache_timeout"`
	VFSEnabled     *bool              `json:"vfs_enabled"`
	MaskingEnabled *bool              `json:"masking_enabled"`
	Tools          []AgentVersionTool `json:"tools"`
	PromptIDs      []struct {
		PromptID  string `json:"prompt_id"`
		Role      string `json:"role"`
		SortOrder int    `json:"sort_order"`
	} `json:"prompt_ids"`
}

type ExecuteAgentRequest struct {
	Input     map[string]any `json:"input"`
	SessionID string         `json:"session_id"`
	UserID    string         `json:"user_id"`
	Stream    bool           `json:"stream"`
	VersionID string         `json:"version_id"`
	Sync      bool           `json:"sync"`
}

type PreviewAgentRequest struct {
	VersionID string         `json:"version_id"`
	Input     map[string]any `json:"input"`
}

type CreatePromptRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdatePromptRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

// CreatePromptVersionRequest creates a content-only prompt version (v2 model).
// Model/sampling/tools/output-schema/cache belong to the agent version.
type CreatePromptVersionRequest struct {
	Version      string          `json:"version"`
	SystemPrompt string          `json:"system_prompt"`
	UserPrompt   string          `json:"user_prompt"`
	InputSchema  json.RawMessage `json:"input_schema"`
	SetCurrent   bool            `json:"set_current"`
	Message      string          `json:"message"`
	Config       json.RawMessage `json:"config"`
}

type PreviewPromptRequest struct {
	VersionID string         `json:"version_id"`
	Input     map[string]any `json:"input"`
}

// PreviewPromptResponse is the dry-run render returned by prompt preview — no
// LLM is called, so there is no token usage / cost.
type PreviewPromptResponse struct {
	SystemPrompt string         `json:"system_prompt"`
	UserPrompt   string         `json:"user_prompt"`
	Input        map[string]any `json:"input"`
}

type CreateDataSourceRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type UpdateDataSourceRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	// MaskingEnabled tri-state — see UpdateAgentRequest.MaskingEnabled.
	MaskingEnabled json.RawMessage `json:"masking_enabled,omitempty"`
	// MaskingRules carries per-field always/never overrides; stored
	// verbatim by the emulator.
	MaskingRules json.RawMessage `json:"masking_rules,omitempty"`
}

type CreateDataSourceVersionRequest struct {
	Version          string          `json:"version"`
	CredentialID     *string         `json:"credential_id"`
	ConnectionConfig json.RawMessage `json:"connection_config"`
	QueryTemplate    string          `json:"query_template"`
	Parameters       json.RawMessage `json:"parameters"`
	SetCurrent       bool            `json:"set_current"`
	CacheTimeout     int             `json:"cache_timeout"`
	OutputFormat     string          `json:"output_format"`
	Message          string          `json:"message"`
}

type CreateCredentialRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Value       string `json:"value"`
	IsDefault   bool   `json:"is_default"`
	SchemaType  string `json:"schema_type"`
}

type UpdateCredentialRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Value       *string `json:"value,omitempty"`
	IsDefault   *bool   `json:"is_default,omitempty"`
}

type CreateChatSessionRequest struct {
	AgentID  string          `json:"agent_id"`
	Title    string          `json:"title"`
	Metadata json.RawMessage `json:"metadata"`
}

type SendMessageRequest struct {
	Content string `json:"content"`
}

type SendMessageResponse struct {
	UserMessage      *ChatMessage `json:"user_message"`
	AssistantMessage *ChatMessage `json:"assistant_message"`
	ExecutionID      string       `json:"execution_id,omitempty"`
}

// PlaygroundRequest runs a persisted agent version with an ephemeral prompt
// snapshot. The override is recorded on the execution but does not create or
// promote prompt/agent versions.
type PlaygroundRequest struct {
	Input          map[string]any `json:"input"`
	VersionID      string         `json:"version_id"`
	PromptOverride struct {
		SystemPrompt string          `json:"system_prompt"`
		UserPrompt   string          `json:"user_prompt"`
		InputSchema  json.RawMessage `json:"input_schema"`
	} `json:"prompt_override"`
}

type CreateAgentTriggerRequest struct {
	AgentID        string                 `json:"agent_id"`
	Name           string                 `json:"name"`
	Source         string                 `json:"source,omitempty"`
	SourceConfig   map[string]interface{} `json:"source_config,omitempty"`
	ReplyConfig    map[string]interface{} `json:"reply_config,omitempty"`
	GenerateSecret bool                   `json:"generate_secret,omitempty"`
	IsActive       bool                   `json:"is_active"`
}

type UpdateAgentTriggerRequest struct {
	Name         *string                `json:"name,omitempty"`
	IsActive     *bool                  `json:"is_active,omitempty"`
	Source       *string                `json:"source,omitempty"`
	SourceConfig map[string]interface{} `json:"source_config,omitempty"`
	ReplyConfig  map[string]interface{} `json:"reply_config,omitempty"`
}

type CreateMCPToolRequest struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Type         string          `json:"type"`
	Config       json.RawMessage `json:"config"`
	Schema       json.RawMessage `json:"schema"`
	CredentialID *string         `json:"credential_id"`
	TemplateID   *string         `json:"template_id"`
}

type UpdateMCPToolRequest struct {
	Name        *string         `json:"name,omitempty"`
	Description *string         `json:"description,omitempty"`
	Config      json.RawMessage `json:"config,omitempty"`
	Schema      json.RawMessage `json:"schema,omitempty"`
	IsActive    *bool           `json:"is_active,omitempty"`
}

type CreateGuardrailRequest struct {
	Type        string          `json:"type"`
	ScannerType string          `json:"scanner_type"`
	Action      string          `json:"action"`
	Config      json.RawMessage `json:"config"`
	IsActive    bool            `json:"is_active"`
	SortOrder   int             `json:"sort_order"`
}

type UpdateGuardrailRequest struct {
	Action   *string         `json:"action,omitempty"`
	Config   json.RawMessage `json:"config,omitempty"`
	IsActive *bool           `json:"is_active,omitempty"`
}

// Response types

type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

type ListResponse struct {
	Data any             `json:"data"`
	Meta *PaginationMeta `json:"meta,omitempty"`
}

type DataResponse struct {
	Data any `json:"data"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type AgentTriggerCreateResponse struct {
	AgentTrigger
	FullToken string `json:"token"`
}

type StoreStats struct {
	Agents         int `json:"agents"`
	AgentVersions  int `json:"agent_versions"`
	Prompts        int `json:"prompts"`
	PromptVersions int `json:"prompt_versions"`
	DataSources    int `json:"data_sources"`
	Executions     int `json:"executions"`
	Credentials    int `json:"credentials"`
	ChatSessions   int `json:"chat_sessions"`
	Traces         int `json:"traces"`
	AgentTriggers  int `json:"agent_triggers"`
	MCPTools       int `json:"mcp_tools"`
	Guardrails     int `json:"guardrails"`
	LLMModels      int `json:"llm_models"`
}
