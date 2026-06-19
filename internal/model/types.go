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
	Type        string `json:"type"`
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
	Version      string          `json:"version"`
	Config       json.RawMessage `json:"config"`
	InputSchema  json.RawMessage `json:"input_schema"`
	OutputSchema json.RawMessage `json:"output_schema"`
	SetCurrent   bool            `json:"set_current"`
	Message      string          `json:"message"`
	PromptIDs    []struct {
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
	TopK               *int            `json:"top_k"`
	Input              map[string]any  `json:"input"`
	OutputSchema       json.RawMessage `json:"output_schema"`
	Tools              []string        `json:"tools"`
	CredentialID       string          `json:"credential_id"`
	CacheTimeout       int             `json:"cache_timeout"`
	ReasoningEffort    string          `json:"reasoning_effort"`
	WebSearch          bool            `json:"web_search"`
	PromptCaching      bool            `json:"prompt_caching"`
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
	Scores         int `json:"scores"`
	ScoreConfigs   int `json:"score_configs"`
	Approvals      int `json:"approvals"`
	AgentTriggers  int `json:"agent_triggers"`
	MCPTools       int `json:"mcp_tools"`
	Guardrails     int `json:"guardrails"`
	LLMModels      int `json:"llm_models"`
}
