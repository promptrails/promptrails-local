package seed

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
	"go.uber.org/zap"
)

//go:embed data/*.json
var dataFS embed.FS

// Load reads all fixture JSON files and populates the in-memory store.
func Load(s *store.Store, logger *zap.Logger) error {
	now := time.Now()

	// LLM Models
	var llmModels []model.LLMModel
	if err := loadJSON("data/llm_models.json", &llmModels); err != nil {
		return fmt.Errorf("llm_models: %w", err)
	}
	for i := range llmModels {
		llmModels[i].CreatedAt = now
		llmModels[i].UpdatedAt = now
		s.AddLLMModel(llmModels[i])
	}
	logger.Info("loaded llm_models", zap.Int("count", len(llmModels)))

	// Credentials
	var credentials []model.Credential
	if err := loadJSON("data/credentials.json", &credentials); err != nil {
		return fmt.Errorf("credentials: %w", err)
	}
	for i := range credentials {
		credentials[i].CreatedAt = now
		credentials[i].UpdatedAt = now
		s.AddCredential(credentials[i])
	}
	logger.Info("loaded credentials", zap.Int("count", len(credentials)))

	// Prompts
	var prompts []model.Prompt
	if err := loadJSON("data/prompts.json", &prompts); err != nil {
		return fmt.Errorf("prompts: %w", err)
	}
	for i := range prompts {
		prompts[i].CreatedAt = now
		prompts[i].UpdatedAt = now
		s.AddPrompt(prompts[i])
	}
	logger.Info("loaded prompts", zap.Int("count", len(prompts)))

	// Prompt Versions
	var promptVersions []model.PromptVersion
	if err := loadJSON("data/prompt_versions.json", &promptVersions); err != nil {
		return fmt.Errorf("prompt_versions: %w", err)
	}
	for i := range promptVersions {
		promptVersions[i].CreatedAt = now
		s.AddPromptVersion(promptVersions[i])
	}
	logger.Info("loaded prompt_versions", zap.Int("count", len(promptVersions)))

	// Agents
	var agents []model.Agent
	if err := loadJSON("data/agents.json", &agents); err != nil {
		return fmt.Errorf("agents: %w", err)
	}
	for i := range agents {
		agents[i].CreatedAt = now
		agents[i].UpdatedAt = now
		s.AddAgent(agents[i])
	}
	logger.Info("loaded agents", zap.Int("count", len(agents)))

	// Agent Versions
	var agentVersions []model.AgentVersion
	if err := loadJSON("data/agent_versions.json", &agentVersions); err != nil {
		return fmt.Errorf("agent_versions: %w", err)
	}
	for i := range agentVersions {
		agentVersions[i].CreatedAt = now
		s.AddAgentVersion(agentVersions[i])
	}
	logger.Info("loaded agent_versions", zap.Int("count", len(agentVersions)))

	// Agent Version Prompts
	var agentVersionPrompts []model.AgentVersionPrompt
	if err := loadJSON("data/agent_version_prompts.json", &agentVersionPrompts); err != nil {
		return fmt.Errorf("agent_version_prompts: %w", err)
	}
	for i := range agentVersionPrompts {
		s.AddAgentVersionPrompt(agentVersionPrompts[i])
	}
	logger.Info("loaded agent_version_prompts", zap.Int("count", len(agentVersionPrompts)))

	// Data Sources
	var dataSources []model.DataSource
	if err := loadJSON("data/data_sources.json", &dataSources); err != nil {
		return fmt.Errorf("data_sources: %w", err)
	}
	for i := range dataSources {
		dataSources[i].CreatedAt = now
		dataSources[i].UpdatedAt = now
		s.AddDataSource(dataSources[i])
	}
	logger.Info("loaded data_sources", zap.Int("count", len(dataSources)))

	// Data Source Versions
	var dataSourceVersions []model.DataSourceVersion
	if err := loadJSON("data/data_source_versions.json", &dataSourceVersions); err != nil {
		return fmt.Errorf("data_source_versions: %w", err)
	}
	for i := range dataSourceVersions {
		dataSourceVersions[i].CreatedAt = now
		s.AddDataSourceVersion(dataSourceVersions[i])
	}
	logger.Info("loaded data_source_versions", zap.Int("count", len(dataSourceVersions)))

	// MCP Tools
	var mcpTools []model.MCPTool
	if err := loadJSON("data/mcp_tools.json", &mcpTools); err != nil {
		return fmt.Errorf("mcp_tools: %w", err)
	}
	for i := range mcpTools {
		mcpTools[i].CreatedAt = now
		mcpTools[i].UpdatedAt = now
		s.AddMCPTool(mcpTools[i])
	}
	logger.Info("loaded mcp_tools", zap.Int("count", len(mcpTools)))

	// Guardrails
	var guardrails []model.Guardrail
	if err := loadJSON("data/guardrails.json", &guardrails); err != nil {
		return fmt.Errorf("guardrails: %w", err)
	}
	for i := range guardrails {
		guardrails[i].CreatedAt = now
		guardrails[i].UpdatedAt = now
		s.AddGuardrail(guardrails[i])
	}
	logger.Info("loaded guardrails", zap.Int("count", len(guardrails)))

	// Memories
	var memories []model.AgentMemory
	if err := loadJSON("data/memories.json", &memories); err != nil {
		return fmt.Errorf("memories: %w", err)
	}
	for i := range memories {
		memories[i].CreatedAt = now
		memories[i].UpdatedAt = now
		s.AddMemory(memories[i])
	}
	logger.Info("loaded memories", zap.Int("count", len(memories)))

	return nil
}

func loadJSON(path string, dest any) error {
	data, err := dataFS.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

// LoadFromDir loads fixture JSON files from an external directory.
// It looks for the same filenames as the embedded fixtures (agents.json, prompts.json, etc.)
// and loads whichever ones exist. This allows users to provide their own test data.
func LoadFromDir(s *store.Store, dir string, logger *zap.Logger) error {
	now := time.Now()

	loaders := []struct {
		file string
		fn   func(data []byte, now time.Time) (int, error)
	}{
		{"llm_models.json", func(data []byte, now time.Time) (int, error) {
			var items []model.LLMModel
			if err := json.Unmarshal(data, &items); err != nil {
				return 0, err
			}
			for i := range items {
				items[i].CreatedAt = now
				items[i].UpdatedAt = now
				s.AddLLMModel(items[i])
			}
			return len(items), nil
		}},
		{"credentials.json", func(data []byte, now time.Time) (int, error) {
			var items []model.Credential
			if err := json.Unmarshal(data, &items); err != nil {
				return 0, err
			}
			for i := range items {
				items[i].CreatedAt = now
				items[i].UpdatedAt = now
				s.AddCredential(items[i])
			}
			return len(items), nil
		}},
		{"prompts.json", func(data []byte, now time.Time) (int, error) {
			var items []model.Prompt
			if err := json.Unmarshal(data, &items); err != nil {
				return 0, err
			}
			for i := range items {
				items[i].CreatedAt = now
				items[i].UpdatedAt = now
				s.AddPrompt(items[i])
			}
			return len(items), nil
		}},
		{"prompt_versions.json", func(data []byte, now time.Time) (int, error) {
			var items []model.PromptVersion
			if err := json.Unmarshal(data, &items); err != nil {
				return 0, err
			}
			for i := range items {
				items[i].CreatedAt = now
				s.AddPromptVersion(items[i])
			}
			return len(items), nil
		}},
		{"agents.json", func(data []byte, now time.Time) (int, error) {
			var items []model.Agent
			if err := json.Unmarshal(data, &items); err != nil {
				return 0, err
			}
			for i := range items {
				items[i].CreatedAt = now
				items[i].UpdatedAt = now
				s.AddAgent(items[i])
			}
			return len(items), nil
		}},
		{"agent_versions.json", func(data []byte, now time.Time) (int, error) {
			var items []model.AgentVersion
			if err := json.Unmarshal(data, &items); err != nil {
				return 0, err
			}
			for i := range items {
				items[i].CreatedAt = now
				s.AddAgentVersion(items[i])
			}
			return len(items), nil
		}},
		{"agent_version_prompts.json", func(data []byte, now time.Time) (int, error) {
			var items []model.AgentVersionPrompt
			if err := json.Unmarshal(data, &items); err != nil {
				return 0, err
			}
			for i := range items {
				s.AddAgentVersionPrompt(items[i])
			}
			return len(items), nil
		}},
		{"data_sources.json", func(data []byte, now time.Time) (int, error) {
			var items []model.DataSource
			if err := json.Unmarshal(data, &items); err != nil {
				return 0, err
			}
			for i := range items {
				items[i].CreatedAt = now
				items[i].UpdatedAt = now
				s.AddDataSource(items[i])
			}
			return len(items), nil
		}},
		{"data_source_versions.json", func(data []byte, now time.Time) (int, error) {
			var items []model.DataSourceVersion
			if err := json.Unmarshal(data, &items); err != nil {
				return 0, err
			}
			for i := range items {
				items[i].CreatedAt = now
				s.AddDataSourceVersion(items[i])
			}
			return len(items), nil
		}},
		{"mcp_tools.json", func(data []byte, now time.Time) (int, error) {
			var items []model.MCPTool
			if err := json.Unmarshal(data, &items); err != nil {
				return 0, err
			}
			for i := range items {
				items[i].CreatedAt = now
				items[i].UpdatedAt = now
				s.AddMCPTool(items[i])
			}
			return len(items), nil
		}},
		{"guardrails.json", func(data []byte, now time.Time) (int, error) {
			var items []model.Guardrail
			if err := json.Unmarshal(data, &items); err != nil {
				return 0, err
			}
			for i := range items {
				items[i].CreatedAt = now
				items[i].UpdatedAt = now
				s.AddGuardrail(items[i])
			}
			return len(items), nil
		}},
		{"memories.json", func(data []byte, now time.Time) (int, error) {
			var items []model.AgentMemory
			if err := json.Unmarshal(data, &items); err != nil {
				return 0, err
			}
			for i := range items {
				items[i].CreatedAt = now
				items[i].UpdatedAt = now
				s.AddMemory(items[i])
			}
			return len(items), nil
		}},
	}

	for _, l := range loaders {
		path := filepath.Join(dir, l.file)
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("read %s: %w", path, err)
		}
		count, err := l.fn(data, now)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		logger.Info("loaded fixtures from directory", zap.String("file", l.file), zap.Int("count", count))
	}

	return nil
}
