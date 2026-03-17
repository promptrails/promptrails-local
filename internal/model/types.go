package model

import (
	"encoding/json"
	"time"
)

type Agent struct {
	ID             string          `json:"id"`
	WorkspaceID    string          `json:"workspace_id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Type           string          `json:"type"`
	Status         string          `json:"status"`
	Labels         json.RawMessage `json:"labels"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	CurrentVersion *AgentVersion   `json:"current_version,omitempty"`
}

type AgentVersion struct {
	ID           string               `json:"id"`
	AgentID      string               `json:"agent_id"`
	Version      string               `json:"version"`
	Config       json.RawMessage      `json:"config"`
	InputSchema  json.RawMessage      `json:"input_schema"`
	OutputSchema json.RawMessage      `json:"output_schema"`
	IsCurrent    bool                 `json:"is_current"`
	Message      string               `json:"message"`
	CreatedAt    time.Time            `json:"created_at"`
	Prompts      []AgentVersionPrompt `json:"prompts,omitempty"`
}

type AgentVersionPrompt struct {
	ID              string         `json:"id"`
	AgentVersionID  string         `json:"agent_version_id"`
	PromptVersionID string         `json:"prompt_version_id"`
	Role            string         `json:"role"`
	SortOrder       int            `json:"sort_order"`
	PromptVersion   *PromptVersion `json:"prompt_version,omitempty"`
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

type PromptVersion struct {
	ID                 string          `json:"id"`
	PromptID           string          `json:"prompt_id"`
	Version            string          `json:"version"`
	SystemPrompt       string          `json:"system_prompt"`
	UserPrompt         string          `json:"user_prompt"`
	LLMModelID         *string         `json:"llm_model_id"`
	FallbackLLMModelID *string         `json:"fallback_llm_model_id"`
	Temperature        *float64        `json:"temperature"`
	MaxTokens          *int            `json:"max_tokens"`
	TopP               *float64        `json:"top_p"`
	InputSchema        json.RawMessage `json:"input_schema"`
	OutputSchema       json.RawMessage `json:"output_schema"`
	IsCurrent          bool            `json:"is_current"`
	Message            string          `json:"message"`
	Config             json.RawMessage `json:"config"`
	CacheTimeout       int             `json:"cache_timeout"`
	CreatedAt          time.Time       `json:"created_at"`
	LLMModel           *LLMModel       `json:"llm_model,omitempty"`
	FallbackLLMModel   *LLMModel       `json:"fallback_llm_model,omitempty"`
}

type DataSource struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

type Execution struct {
	ID             string          `json:"id"`
	AgentID        *string         `json:"agent_id"`
	AgentVersionID *string         `json:"agent_version_id"`
	WorkspaceID    string          `json:"workspace_id"`
	UserID         *string         `json:"user_id"`
	SessionID      string          `json:"session_id"`
	Status         string          `json:"status"`
	Input          json.RawMessage `json:"input"`
	Output         json.RawMessage `json:"output"`
	Error          string          `json:"error"`
	Metadata       json.RawMessage `json:"metadata"`
	TokenUsage     json.RawMessage `json:"token_usage"`
	Cost           float64         `json:"cost"`
	DurationMS     *int64          `json:"duration_ms"`
	TraceID        string          `json:"trace_id,omitempty"`
	StartedAt      *time.Time      `json:"started_at"`
	CompletedAt    *time.Time      `json:"completed_at"`
	CreatedAt      time.Time       `json:"created_at"`
	Agent          *Agent          `json:"agent,omitempty"`
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

type Score struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	TraceID     string    `json:"trace_id"`
	SpanID      *string   `json:"span_id,omitempty"`
	Name        string    `json:"name"`
	Value       *float64  `json:"value,omitempty"`
	StringValue *string   `json:"string_value,omitempty"`
	BoolValue   *bool     `json:"bool_value,omitempty"`
	DataType    string    `json:"data_type"`
	Comment     *string   `json:"comment,omitempty"`
	Source      string    `json:"source"`
	ConfigID    *string   `json:"config_id,omitempty"`
	ExecutionID *string   `json:"execution_id,omitempty"`
	AgentID     *string   `json:"agent_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ScoreConfig struct {
	ID          string          `json:"id"`
	WorkspaceID string          `json:"workspace_id"`
	Name        string          `json:"name"`
	DataType    string          `json:"data_type"`
	MinValue    *float64        `json:"min_value,omitempty"`
	MaxValue    *float64        `json:"max_value,omitempty"`
	Categories  json.RawMessage `json:"categories"`
	Description *string         `json:"description,omitempty"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ApprovalRequest struct {
	ID             string          `json:"id"`
	ExecutionID    string          `json:"execution_id"`
	AgentID        string          `json:"agent_id"`
	WorkspaceID    string          `json:"workspace_id"`
	CheckpointName string          `json:"checkpoint_name"`
	Payload        json.RawMessage `json:"payload"`
	Status         string          `json:"status"`
	DecidedBy      *string         `json:"decided_by"`
	DecidedAt      *time.Time      `json:"decided_at"`
	ExpiresAt      *time.Time      `json:"expires_at"`
	Reason         *string         `json:"reason"`
	CreatedAt      time.Time       `json:"created_at"`
}

type WebhookTrigger struct {
	ID          string     `json:"id"`
	WorkspaceID string     `json:"workspace_id"`
	AgentID     string     `json:"agent_id"`
	Name        string     `json:"name"`
	Token       string     `json:"-"`
	TokenPrefix string     `json:"token_prefix"`
	IsActive    bool       `json:"is_active"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Agent       *Agent     `json:"agent,omitempty"`
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

type AgentMemory struct {
	ID             string          `json:"id"`
	WorkspaceID    string          `json:"workspace_id"`
	AgentID        string          `json:"agent_id"`
	Content        string          `json:"content"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	MemoryType     string          `json:"memory_type"`
	Importance     float64         `json:"importance"`
	AccessCount    int             `json:"access_count"`
	LastAccessedAt *time.Time      `json:"last_accessed_at,omitempty"`
	ChatSessionID  *string         `json:"chat_session_id,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type LLMModel struct {
	ID                string    `json:"id"`
	Provider          string    `json:"provider"`
	ModelID           string    `json:"model_id"`
	DisplayName       string    `json:"display_name"`
	InputPrice        *float64  `json:"input_price"`
	OutputPrice       *float64  `json:"output_price"`
	MaxTokens         *int      `json:"max_tokens"`
	SupportsVision    bool      `json:"supports_vision"`
	SupportsTools     bool      `json:"supports_tools"`
	SupportsJSON      bool      `json:"supports_json"`
	SupportsStreaming bool      `json:"supports_streaming"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
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
	Type        string `json:"type"`
}

type UpdateAgentRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

type CreateAgentVersionRequest struct {
	Version      string          `json:"version"`
	Config       json.RawMessage `json:"config"`
	InputSchema  json.RawMessage `json:"input_schema"`
	OutputSchema json.RawMessage `json:"output_schema"`
	SetCurrent   bool            `json:"set_current"`
	Message      string          `json:"message"`
	PromptIDs    []struct {
		PromptVersionID string `json:"prompt_version_id"`
		Role            string `json:"role"`
		SortOrder       int    `json:"sort_order"`
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

type CreatePromptVersionRequest struct {
	Version            string          `json:"version"`
	SystemPrompt       string          `json:"system_prompt"`
	UserPrompt         string          `json:"user_prompt"`
	LLMModelID         *string         `json:"llm_model_id"`
	FallbackLLMModelID *string         `json:"fallback_llm_model_id"`
	Temperature        *float64        `json:"temperature"`
	MaxTokens          *int            `json:"max_tokens"`
	TopP               *float64        `json:"top_p"`
	InputSchema        json.RawMessage `json:"input_schema"`
	OutputSchema       json.RawMessage `json:"output_schema"`
	SetCurrent         bool            `json:"set_current"`
	Message            string          `json:"message"`
	Config             json.RawMessage `json:"config"`
	CacheTimeout       int             `json:"cache_timeout"`
}

type PreviewPromptRequest struct {
	VersionID string         `json:"version_id"`
	Input     map[string]any `json:"input"`
}

type RunPromptRequest struct {
	SystemPrompt       string          `json:"system_prompt"`
	UserPrompt         string          `json:"user_prompt"`
	LLMModelID         string          `json:"llm_model_id"`
	FallbackLLMModelID string          `json:"fallback_llm_model_id"`
	Temperature        *float64        `json:"temperature"`
	MaxTokens          *int            `json:"max_tokens"`
	TopP               *float64        `json:"top_p"`
	Input              map[string]any  `json:"input"`
	OutputSchema       json.RawMessage `json:"output_schema"`
	Tools              []string        `json:"tools"`
	CredentialID       string          `json:"credential_id"`
	CacheTimeout       int             `json:"cache_timeout"`
}

type RunPromptResponse struct {
	Content    string         `json:"content"`
	TokenUsage map[string]int `json:"token_usage"`
	Cost       float64        `json:"cost"`
	DurationMS int64          `json:"duration_ms"`
	Model      string         `json:"model"`
	TraceID    string         `json:"trace_id,omitempty"`
}

type CreateDataSourceRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type UpdateDataSourceRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
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

type CreateScoreRequest struct {
	TraceID     string   `json:"trace_id"`
	SpanID      *string  `json:"span_id,omitempty"`
	Name        string   `json:"name"`
	Value       *float64 `json:"value,omitempty"`
	StringValue *string  `json:"string_value,omitempty"`
	BoolValue   *bool    `json:"bool_value,omitempty"`
	DataType    string   `json:"data_type"`
	Comment     *string  `json:"comment,omitempty"`
	Source      string   `json:"source"`
	ConfigID    *string  `json:"config_id,omitempty"`
	ExecutionID *string  `json:"execution_id,omitempty"`
	AgentID     *string  `json:"agent_id,omitempty"`
}

type UpdateScoreRequest struct {
	Value       *float64 `json:"value,omitempty"`
	StringValue *string  `json:"string_value,omitempty"`
	BoolValue   *bool    `json:"bool_value,omitempty"`
	Comment     *string  `json:"comment,omitempty"`
}

type CreateScoreConfigRequest struct {
	Name        string          `json:"name"`
	DataType    string          `json:"data_type"`
	MinValue    *float64        `json:"min_value,omitempty"`
	MaxValue    *float64        `json:"max_value,omitempty"`
	Categories  json.RawMessage `json:"categories"`
	Description *string         `json:"description,omitempty"`
}

type UpdateScoreConfigRequest struct {
	Name        *string         `json:"name,omitempty"`
	Description *string         `json:"description,omitempty"`
	IsActive    *bool           `json:"is_active,omitempty"`
	MinValue    *float64        `json:"min_value,omitempty"`
	MaxValue    *float64        `json:"max_value,omitempty"`
	Categories  json.RawMessage `json:"categories,omitempty"`
}

type DecideApprovalRequest struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

type CreateWebhookTriggerRequest struct {
	AgentID  string `json:"agent_id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type UpdateWebhookTriggerRequest struct {
	Name     *string `json:"name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
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

type CreateMemoryRequest struct {
	Content    string          `json:"content"`
	MemoryType string          `json:"memory_type"`
	Importance float64         `json:"importance"`
	Metadata   json.RawMessage `json:"metadata"`
}

type SearchMemoryRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
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

type WebhookTriggerCreateResponse struct {
	WebhookTrigger
	FullToken string `json:"token"`
}

type StoreStats struct {
	Agents          int `json:"agents"`
	AgentVersions   int `json:"agent_versions"`
	Prompts         int `json:"prompts"`
	PromptVersions  int `json:"prompt_versions"`
	DataSources     int `json:"data_sources"`
	Executions      int `json:"executions"`
	Credentials     int `json:"credentials"`
	ChatSessions    int `json:"chat_sessions"`
	Traces          int `json:"traces"`
	Scores          int `json:"scores"`
	ScoreConfigs    int `json:"score_configs"`
	Approvals       int `json:"approvals"`
	WebhookTriggers int `json:"webhook_triggers"`
	MCPTools        int `json:"mcp_tools"`
	Guardrails      int `json:"guardrails"`
	Memories        int `json:"memories"`
	LLMModels       int `json:"llm_models"`
}
