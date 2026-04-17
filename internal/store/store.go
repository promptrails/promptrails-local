package store

import (
	"sort"
	"strings"
	"sync"

	"github.com/promptrails/promptrails-local/internal/model"
)

// ---------- filter types ----------

// AgentFilters holds optional filters for listing agents.
type AgentFilters struct {
	Type   string
	Status string
	Name   string
}

// ExecutionFilters holds optional filters for listing executions.
type ExecutionFilters struct {
	AgentID   string
	SessionID string
	Status    string
}

// ---------- Store ----------

// Store is a thread-safe in-memory store for all PromptRails resources.
type Store struct {
	mu sync.RWMutex

	agents              map[string]model.Agent
	agentVersions       map[string]model.AgentVersion
	agentVersionPrompts map[string]model.AgentVersionPrompt
	prompts             map[string]model.Prompt
	promptVersions      map[string]model.PromptVersion
	dataSources         map[string]model.DataSource
	dataSourceVersions  map[string]model.DataSourceVersion
	executions          map[string]model.Execution
	credentials         map[string]model.Credential
	chatSessions        map[string]model.ChatSession
	chatMessages        map[string]model.ChatMessage
	traces              map[string]model.Trace
	scores              map[string]model.Score
	scoreConfigs        map[string]model.ScoreConfig
	approvals           map[string]model.ApprovalRequest
	webhookTriggers     map[string]model.WebhookTrigger
	mcpTools            map[string]model.MCPTool
	guardrails          map[string]model.Guardrail
	memories            map[string]model.AgentMemory
	llmModels           map[string]model.LLMModel
	mcpTemplates        map[string]model.MCPTemplate
}

// New creates a new empty Store.
func New() *Store {
	s := &Store{}
	s.initMaps()
	return s
}

func (s *Store) initMaps() {
	s.agents = make(map[string]model.Agent)
	s.agentVersions = make(map[string]model.AgentVersion)
	s.agentVersionPrompts = make(map[string]model.AgentVersionPrompt)
	s.prompts = make(map[string]model.Prompt)
	s.promptVersions = make(map[string]model.PromptVersion)
	s.dataSources = make(map[string]model.DataSource)
	s.dataSourceVersions = make(map[string]model.DataSourceVersion)
	s.executions = make(map[string]model.Execution)
	s.credentials = make(map[string]model.Credential)
	s.chatSessions = make(map[string]model.ChatSession)
	s.chatMessages = make(map[string]model.ChatMessage)
	s.traces = make(map[string]model.Trace)
	s.scores = make(map[string]model.Score)
	s.scoreConfigs = make(map[string]model.ScoreConfig)
	s.approvals = make(map[string]model.ApprovalRequest)
	s.webhookTriggers = make(map[string]model.WebhookTrigger)
	s.mcpTools = make(map[string]model.MCPTool)
	s.guardrails = make(map[string]model.Guardrail)
	s.memories = make(map[string]model.AgentMemory)
	s.llmModels = make(map[string]model.LLMModel)
	s.mcpTemplates = make(map[string]model.MCPTemplate)
}

// Reset clears all data and re-initialises internal maps.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.initMaps()
}

// Stats returns counts for every entity type.
func (s *Store) Stats() model.StoreStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return model.StoreStats{
		Agents:          len(s.agents),
		AgentVersions:   len(s.agentVersions),
		Prompts:         len(s.prompts),
		PromptVersions:  len(s.promptVersions),
		DataSources:     len(s.dataSources),
		Executions:      len(s.executions),
		Credentials:     len(s.credentials),
		ChatSessions:    len(s.chatSessions),
		Traces:          len(s.traces),
		Scores:          len(s.scores),
		ScoreConfigs:    len(s.scoreConfigs),
		Approvals:       len(s.approvals),
		WebhookTriggers: len(s.webhookTriggers),
		MCPTools:        len(s.mcpTools),
		Guardrails:      len(s.guardrails),
		Memories:        len(s.memories),
		LLMModels:       len(s.llmModels),
	}
}

// ============================================================
// Agents
// ============================================================

// AddAgent inserts an agent (used by seed loader).
func (s *Store) AddAgent(a model.Agent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agents[a.ID] = a
}

// CreateAgent inserts a new agent.
func (s *Store) CreateAgent(a model.Agent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agents[a.ID] = a
}

// GetAgent returns an agent by ID.
func (s *Store) GetAgent(id string) (model.Agent, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.agents[id]
	if ok {
		a.CurrentVersion = s.currentAgentVersion(a.ID)
	}
	return a, ok
}

// UpdateAgent replaces an agent in place.
func (s *Store) UpdateAgent(a model.Agent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agents[a.ID] = a
}

// DeleteAgent removes an agent and returns whether it existed.
func (s *Store) DeleteAgent(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.agents[id]; !ok {
		return false
	}
	delete(s.agents, id)
	return true
}

// ListAgents returns a filtered, paginated list of agents.
func (s *Store) ListAgents(f AgentFilters, page, limit int) ([]model.Agent, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.Agent
	for _, a := range s.agents {
		if f.Type != "" && a.Type != f.Type {
			continue
		}
		if f.Status != "" && a.Status != f.Status {
			continue
		}
		if f.Name != "" && !strings.Contains(strings.ToLower(a.Name), strings.ToLower(f.Name)) {
			continue
		}
		a.CurrentVersion = s.currentAgentVersion(a.ID)
		all = append(all, a)
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	return paginate(all, page, limit)
}

// currentAgentVersion returns the current version pointer (must be called under at least RLock).
func (s *Store) currentAgentVersion(agentID string) *model.AgentVersion {
	for _, v := range s.agentVersions {
		if v.AgentID == agentID && v.IsCurrent {
			cv := v
			cv.Prompts = s.agentVersionPromptsFor(cv.ID)
			return &cv
		}
	}
	return nil
}

// ============================================================
// Agent Versions
// ============================================================

// AddAgentVersion inserts an agent version (seed loader).
func (s *Store) AddAgentVersion(v model.AgentVersion) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agentVersions[v.ID] = v
}

// CreateAgentVersion inserts a new agent version and its prompts.
func (s *Store) CreateAgentVersion(v model.AgentVersion) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range v.Prompts {
		s.agentVersionPrompts[p.ID] = p
	}
	s.agentVersions[v.ID] = v
}

// GetAgentVersion returns an agent version by ID.
func (s *Store) GetAgentVersion(id string) (model.AgentVersion, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.agentVersions[id]
	if ok {
		v.Prompts = s.agentVersionPromptsFor(v.ID)
	}
	return v, ok
}

// UpdateAgentVersion replaces an agent version in place.
func (s *Store) UpdateAgentVersion(v model.AgentVersion) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agentVersions[v.ID] = v
}

// ListAgentVersions returns all versions for a given agent, newest first.
func (s *Store) ListAgentVersions(agentID string) []model.AgentVersion {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []model.AgentVersion
	for _, v := range s.agentVersions {
		if v.AgentID == agentID {
			v.Prompts = s.agentVersionPromptsFor(v.ID)
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

// DemoteAgentVersions sets IsCurrent=false for all versions of the given agent.
func (s *Store) DemoteAgentVersions(agentID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, v := range s.agentVersions {
		if v.AgentID == agentID && v.IsCurrent {
			v.IsCurrent = false
			s.agentVersions[id] = v
		}
	}
}

// ============================================================
// Agent Version Prompts
// ============================================================

// AddAgentVersionPrompt inserts an agent-version-prompt association (seed loader).
func (s *Store) AddAgentVersionPrompt(p model.AgentVersionPrompt) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agentVersionPrompts[p.ID] = p
}

// agentVersionPromptsFor returns prompts for a version (must be called under at least RLock).
// The link is by prompt_id, and the prompt's current version is attached on each
// read so agents always use the latest published prompt content.
func (s *Store) agentVersionPromptsFor(versionID string) []model.AgentVersionPrompt {
	var out []model.AgentVersionPrompt
	for _, p := range s.agentVersionPrompts {
		if p.AgentVersionID == versionID {
			cp := p
			if prompt, ok := s.prompts[cp.PromptID]; ok {
				pc := prompt
				pc.CurrentVersion = s.currentPromptVersion(pc.ID)
				cp.Prompt = &pc
			}
			out = append(out, cp)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].SortOrder < out[j].SortOrder
	})
	return out
}

// ============================================================
// Prompts
// ============================================================

// AddPrompt inserts a prompt (seed loader).
func (s *Store) AddPrompt(p model.Prompt) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prompts[p.ID] = p
}

// CreatePrompt inserts a new prompt.
func (s *Store) CreatePrompt(p model.Prompt) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prompts[p.ID] = p
}

// GetPrompt returns a prompt by ID.
func (s *Store) GetPrompt(id string) (model.Prompt, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.prompts[id]
	if ok {
		p.CurrentVersion = s.currentPromptVersion(p.ID)
	}
	return p, ok
}

// UpdatePrompt replaces a prompt in place.
func (s *Store) UpdatePrompt(p model.Prompt) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prompts[p.ID] = p
}

// DeletePrompt removes a prompt and returns whether it existed.
func (s *Store) DeletePrompt(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.prompts[id]; !ok {
		return false
	}
	delete(s.prompts, id)
	return true
}

// ListPrompts returns a paginated list of all prompts.
func (s *Store) ListPrompts(page, limit int) ([]model.Prompt, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.Prompt
	for _, p := range s.prompts {
		p.CurrentVersion = s.currentPromptVersion(p.ID)
		all = append(all, p)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// currentPromptVersion returns the current version pointer (must be called under at least RLock).
func (s *Store) currentPromptVersion(promptID string) *model.PromptVersion {
	for _, v := range s.promptVersions {
		if v.PromptID == promptID && v.IsCurrent {
			cv := v
			cv.LLMModel = s.llmModelPtr(cv.LLMModelID)
			cv.FallbackLLMModel = s.llmModelPtr(cv.FallbackLLMModelID)
			return &cv
		}
	}
	return nil
}

// ============================================================
// Prompt Versions
// ============================================================

// AddPromptVersion inserts a prompt version (seed loader).
func (s *Store) AddPromptVersion(v model.PromptVersion) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.promptVersions[v.ID] = v
}

// CreatePromptVersion inserts a new prompt version.
func (s *Store) CreatePromptVersion(v model.PromptVersion) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.promptVersions[v.ID] = v
}

// GetPromptVersion returns a prompt version by ID.
func (s *Store) GetPromptVersion(id string) (model.PromptVersion, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.promptVersions[id]
	if ok {
		v.LLMModel = s.llmModelPtr(v.LLMModelID)
		v.FallbackLLMModel = s.llmModelPtr(v.FallbackLLMModelID)
	}
	return v, ok
}

// UpdatePromptVersion replaces a prompt version in place.
func (s *Store) UpdatePromptVersion(v model.PromptVersion) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.promptVersions[v.ID] = v
}

// ListPromptVersions returns all versions for a given prompt, newest first.
func (s *Store) ListPromptVersions(promptID string) []model.PromptVersion {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []model.PromptVersion
	for _, v := range s.promptVersions {
		if v.PromptID == promptID {
			v.LLMModel = s.llmModelPtr(v.LLMModelID)
			v.FallbackLLMModel = s.llmModelPtr(v.FallbackLLMModelID)
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

// DemotePromptVersions sets IsCurrent=false for all versions of the given prompt.
func (s *Store) DemotePromptVersions(promptID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, v := range s.promptVersions {
		if v.PromptID == promptID && v.IsCurrent {
			v.IsCurrent = false
			s.promptVersions[id] = v
		}
	}
}

// ============================================================
// Executions
// ============================================================

// CreateExecution inserts a new execution.
func (s *Store) CreateExecution(e model.Execution) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.executions[e.ID] = e
}

// GetExecution returns an execution by ID.
func (s *Store) GetExecution(id string) (model.Execution, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.executions[id]
	if ok && e.AgentID != nil {
		if a, aok := s.agents[*e.AgentID]; aok {
			e.Agent = &a
		}
	}
	return e, ok
}

// UpdateExecution replaces an execution in place.
func (s *Store) UpdateExecution(e model.Execution) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.executions[e.ID] = e
}

// ListExecutions returns a filtered, paginated list of executions.
func (s *Store) ListExecutions(f ExecutionFilters, page, limit int) ([]model.Execution, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.Execution
	for _, e := range s.executions {
		if f.AgentID != "" && (e.AgentID == nil || *e.AgentID != f.AgentID) {
			continue
		}
		if f.SessionID != "" && e.SessionID != f.SessionID {
			continue
		}
		if f.Status != "" && e.Status != f.Status {
			continue
		}
		if e.AgentID != nil {
			if a, ok := s.agents[*e.AgentID]; ok {
				e.Agent = &a
			}
		}
		all = append(all, e)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// ============================================================
// Data Sources
// ============================================================

// AddDataSource inserts a data source (seed loader).
func (s *Store) AddDataSource(ds model.DataSource) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dataSources[ds.ID] = ds
}

// CreateDataSource inserts a new data source.
func (s *Store) CreateDataSource(ds model.DataSource) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dataSources[ds.ID] = ds
}

// GetDataSource returns a data source by ID.
func (s *Store) GetDataSource(id string) (model.DataSource, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ds, ok := s.dataSources[id]
	return ds, ok
}

// UpdateDataSource replaces a data source in place.
func (s *Store) UpdateDataSource(ds model.DataSource) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dataSources[ds.ID] = ds
}

// DeleteDataSource removes a data source and returns whether it existed.
func (s *Store) DeleteDataSource(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dataSources[id]; !ok {
		return false
	}
	delete(s.dataSources, id)
	return true
}

// ListDataSources returns a paginated list of all data sources.
func (s *Store) ListDataSources(page, limit int) ([]model.DataSource, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.DataSource
	for _, ds := range s.dataSources {
		all = append(all, ds)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// ============================================================
// Data Source Versions
// ============================================================

// AddDataSourceVersion inserts a data source version (seed loader).
func (s *Store) AddDataSourceVersion(v model.DataSourceVersion) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dataSourceVersions[v.ID] = v
}

// CreateDataSourceVersion inserts a new data source version.
func (s *Store) CreateDataSourceVersion(v model.DataSourceVersion) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dataSourceVersions[v.ID] = v
}

// GetDataSourceVersion returns a data source version by ID.
func (s *Store) GetDataSourceVersion(id string) (model.DataSourceVersion, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.dataSourceVersions[id]
	return v, ok
}

// UpdateDataSourceVersion replaces a data source version in place.
func (s *Store) UpdateDataSourceVersion(v model.DataSourceVersion) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dataSourceVersions[v.ID] = v
}

// ListDataSourceVersions returns all versions for a given data source, newest first.
func (s *Store) ListDataSourceVersions(dataSourceID string) []model.DataSourceVersion {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []model.DataSourceVersion
	for _, v := range s.dataSourceVersions {
		if v.DataSourceID == dataSourceID {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

// DemoteDataSourceVersions sets IsCurrent=false for all versions of the given data source.
func (s *Store) DemoteDataSourceVersions(dataSourceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, v := range s.dataSourceVersions {
		if v.DataSourceID == dataSourceID && v.IsCurrent {
			v.IsCurrent = false
			s.dataSourceVersions[id] = v
		}
	}
}

// ============================================================
// Credentials
// ============================================================

// AddCredential inserts a credential (seed loader).
func (s *Store) AddCredential(c model.Credential) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.credentials[c.ID] = c
}

// CreateCredential inserts a new credential.
func (s *Store) CreateCredential(c model.Credential) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.credentials[c.ID] = c
}

// GetCredential returns a credential by ID.
func (s *Store) GetCredential(id string) (model.Credential, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.credentials[id]
	return c, ok
}

// UpdateCredential replaces a credential in place.
func (s *Store) UpdateCredential(c model.Credential) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.credentials[c.ID] = c
}

// DeleteCredential removes a credential and returns whether it existed.
func (s *Store) DeleteCredential(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.credentials[id]; !ok {
		return false
	}
	delete(s.credentials, id)
	return true
}

// ListCredentials returns a paginated list of all credentials.
func (s *Store) ListCredentials(page, limit int) ([]model.Credential, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.Credential
	for _, c := range s.credentials {
		all = append(all, c)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// ============================================================
// Chat Sessions
// ============================================================

// CreateChatSession inserts a new chat session.
func (s *Store) CreateChatSession(cs model.ChatSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.chatSessions[cs.ID] = cs
}

// GetChatSession returns a chat session by ID.
func (s *Store) GetChatSession(id string) (model.ChatSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cs, ok := s.chatSessions[id]
	if ok {
		if a, aok := s.agents[cs.AgentID]; aok {
			cs.Agent = &a
		}
	}
	return cs, ok
}

// DeleteChatSession removes a chat session and returns whether it existed.
func (s *Store) DeleteChatSession(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.chatSessions[id]; !ok {
		return false
	}
	delete(s.chatSessions, id)
	// Also delete associated messages
	for mid, m := range s.chatMessages {
		if m.SessionID == id {
			delete(s.chatMessages, mid)
		}
	}
	return true
}

// ListChatSessions returns a paginated list of all chat sessions.
func (s *Store) ListChatSessions(page, limit int) ([]model.ChatSession, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.ChatSession
	for _, cs := range s.chatSessions {
		if a, ok := s.agents[cs.AgentID]; ok {
			cs.Agent = &a
		}
		all = append(all, cs)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// ============================================================
// Chat Messages
// ============================================================

// CreateChatMessage inserts a new chat message.
func (s *Store) CreateChatMessage(m model.ChatMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.chatMessages[m.ID] = m
}

// ListChatMessages returns a paginated list of messages for a session.
func (s *Store) ListChatMessages(sessionID string, page, limit int) ([]model.ChatMessage, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.ChatMessage
	for _, m := range s.chatMessages {
		if m.SessionID == sessionID {
			all = append(all, m)
		}
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.Before(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// ============================================================
// Traces
// ============================================================

// AddTrace inserts a trace (used by fake LLM).
func (s *Store) AddTrace(t model.Trace) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.traces[t.ID] = t
}

// CreateTrace inserts a new trace.
func (s *Store) CreateTrace(t model.Trace) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.traces[t.ID] = t
}

// GetTrace returns a trace by ID.
func (s *Store) GetTrace(id string) (model.Trace, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.traces[id]
	return t, ok
}

// ListTraces returns a paginated list of all traces.
func (s *Store) ListTraces(page, limit int) ([]model.Trace, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.Trace
	for _, t := range s.traces {
		all = append(all, t)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// ============================================================
// Scores
// ============================================================

// CreateScore inserts a new score.
func (s *Store) CreateScore(sc model.Score) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scores[sc.ID] = sc
}

// GetScore returns a score by ID.
func (s *Store) GetScore(id string) (model.Score, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sc, ok := s.scores[id]
	return sc, ok
}

// UpdateScore replaces a score in place.
func (s *Store) UpdateScore(sc model.Score) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scores[sc.ID] = sc
}

// DeleteScore removes a score and returns whether it existed.
func (s *Store) DeleteScore(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.scores[id]; !ok {
		return false
	}
	delete(s.scores, id)
	return true
}

// ListScores returns a paginated list of all scores.
func (s *Store) ListScores(page, limit int) ([]model.Score, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.Score
	for _, sc := range s.scores {
		all = append(all, sc)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// ============================================================
// Score Configs
// ============================================================

// CreateScoreConfig inserts a new score config.
func (s *Store) CreateScoreConfig(cfg model.ScoreConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scoreConfigs[cfg.ID] = cfg
}

// GetScoreConfig returns a score config by ID.
func (s *Store) GetScoreConfig(id string) (model.ScoreConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cfg, ok := s.scoreConfigs[id]
	return cfg, ok
}

// UpdateScoreConfig replaces a score config in place.
func (s *Store) UpdateScoreConfig(cfg model.ScoreConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scoreConfigs[cfg.ID] = cfg
}

// DeleteScoreConfig removes a score config and returns whether it existed.
func (s *Store) DeleteScoreConfig(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.scoreConfigs[id]; !ok {
		return false
	}
	delete(s.scoreConfigs, id)
	return true
}

// ListScoreConfigs returns a paginated list of all score configs.
func (s *Store) ListScoreConfigs(page, limit int) ([]model.ScoreConfig, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.ScoreConfig
	for _, cfg := range s.scoreConfigs {
		all = append(all, cfg)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// ============================================================
// Approvals
// ============================================================

// CreateApproval inserts a new approval request.
func (s *Store) CreateApproval(a model.ApprovalRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.approvals[a.ID] = a
}

// GetApproval returns an approval request by ID.
func (s *Store) GetApproval(id string) (model.ApprovalRequest, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.approvals[id]
	return a, ok
}

// UpdateApproval replaces an approval request in place.
func (s *Store) UpdateApproval(a model.ApprovalRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.approvals[a.ID] = a
}

// GetApprovalByExecutionID returns the approval request for an execution.
func (s *Store) GetApprovalByExecutionID(executionID string) (model.ApprovalRequest, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.approvals {
		if a.ExecutionID == executionID {
			return a, true
		}
	}
	return model.ApprovalRequest{}, false
}

// ListApprovals returns a paginated list of all approvals.
func (s *Store) ListApprovals(page, limit int) ([]model.ApprovalRequest, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.ApprovalRequest
	for _, a := range s.approvals {
		all = append(all, a)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// ============================================================
// Webhook Triggers
// ============================================================

// CreateWebhookTrigger inserts a new webhook trigger.
func (s *Store) CreateWebhookTrigger(wt model.WebhookTrigger) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.webhookTriggers[wt.ID] = wt
}

// GetWebhookTrigger returns a webhook trigger by ID.
func (s *Store) GetWebhookTrigger(id string) (model.WebhookTrigger, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	wt, ok := s.webhookTriggers[id]
	if ok {
		if a, aok := s.agents[wt.AgentID]; aok {
			wt.Agent = &a
		}
	}
	return wt, ok
}

// GetWebhookTriggerByToken returns a webhook trigger by its token.
func (s *Store) GetWebhookTriggerByToken(token string) (model.WebhookTrigger, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, wt := range s.webhookTriggers {
		if wt.Token == token {
			if a, ok := s.agents[wt.AgentID]; ok {
				wt.Agent = &a
			}
			return wt, true
		}
	}
	return model.WebhookTrigger{}, false
}

// UpdateWebhookTrigger replaces a webhook trigger in place.
func (s *Store) UpdateWebhookTrigger(wt model.WebhookTrigger) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.webhookTriggers[wt.ID] = wt
}

// DeleteWebhookTrigger removes a webhook trigger and returns whether it existed.
func (s *Store) DeleteWebhookTrigger(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.webhookTriggers[id]; !ok {
		return false
	}
	delete(s.webhookTriggers, id)
	return true
}

// ListWebhookTriggers returns a paginated list of all webhook triggers.
func (s *Store) ListWebhookTriggers(page, limit int) ([]model.WebhookTrigger, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.WebhookTrigger
	for _, wt := range s.webhookTriggers {
		if a, ok := s.agents[wt.AgentID]; ok {
			wt.Agent = &a
		}
		all = append(all, wt)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// ============================================================
// MCP Tools
// ============================================================

// AddMCPTool inserts an MCP tool (seed loader).
func (s *Store) AddMCPTool(t model.MCPTool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mcpTools[t.ID] = t
}

// CreateMCPTool inserts a new MCP tool.
func (s *Store) CreateMCPTool(t model.MCPTool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mcpTools[t.ID] = t
}

// GetMCPTool returns an MCP tool by ID.
func (s *Store) GetMCPTool(id string) (model.MCPTool, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.mcpTools[id]
	return t, ok
}

// UpdateMCPTool replaces an MCP tool in place.
func (s *Store) UpdateMCPTool(t model.MCPTool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mcpTools[t.ID] = t
}

// DeleteMCPTool removes an MCP tool and returns whether it existed.
func (s *Store) DeleteMCPTool(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.mcpTools[id]; !ok {
		return false
	}
	delete(s.mcpTools, id)
	return true
}

// ListMCPTools returns a paginated list of all MCP tools.
func (s *Store) ListMCPTools(page, limit int) ([]model.MCPTool, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.MCPTool
	for _, t := range s.mcpTools {
		all = append(all, t)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// ============================================================
// MCP Templates
// ============================================================

// AddMCPTemplate inserts an MCP template (seed loader).
func (s *Store) AddMCPTemplate(t model.MCPTemplate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mcpTemplates[t.ID] = t
}

// ListMCPTemplates returns all MCP templates.
func (s *Store) ListMCPTemplates() []model.MCPTemplate {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]model.MCPTemplate, 0, len(s.mcpTemplates))
	for _, t := range s.mcpTemplates {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

// GetMCPTemplate returns an MCP template by ID.
func (s *Store) GetMCPTemplate(id string) (model.MCPTemplate, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.mcpTemplates[id]
	return t, ok
}

// ============================================================
// Guardrails
// ============================================================

// AddGuardrail inserts a guardrail (seed loader).
func (s *Store) AddGuardrail(g model.Guardrail) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.guardrails[g.ID] = g
}

// CreateGuardrail inserts a new guardrail.
func (s *Store) CreateGuardrail(g model.Guardrail) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.guardrails[g.ID] = g
}

// GetGuardrail returns a guardrail by ID.
func (s *Store) GetGuardrail(id string) (model.Guardrail, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g, ok := s.guardrails[id]
	return g, ok
}

// UpdateGuardrail replaces a guardrail in place.
func (s *Store) UpdateGuardrail(g model.Guardrail) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.guardrails[g.ID] = g
}

// DeleteGuardrail removes a guardrail and returns whether it existed.
func (s *Store) DeleteGuardrail(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.guardrails[id]; !ok {
		return false
	}
	delete(s.guardrails, id)
	return true
}

// ListGuardrails returns all guardrails for a given agent, sorted by sort_order.
func (s *Store) ListGuardrails(agentID string) []model.Guardrail {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []model.Guardrail
	for _, g := range s.guardrails {
		if g.AgentID == agentID {
			out = append(out, g)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].SortOrder < out[j].SortOrder
	})
	return out
}

// ============================================================
// Memories
// ============================================================

// AddMemory inserts a memory (seed loader).
func (s *Store) AddMemory(m model.AgentMemory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.memories[m.ID] = m
}

// CreateMemory inserts a new memory.
func (s *Store) CreateMemory(m model.AgentMemory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.memories[m.ID] = m
}

// GetMemory returns a memory by ID.
func (s *Store) GetMemory(id string) (model.AgentMemory, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.memories[id]
	return m, ok
}

// DeleteMemory removes a memory and returns whether it existed.
func (s *Store) DeleteMemory(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.memories[id]; !ok {
		return false
	}
	delete(s.memories, id)
	return true
}

// DeleteAllMemories removes all memories for a given agent, returning the count.
func (s *Store) DeleteAllMemories(agentID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	for id, m := range s.memories {
		if m.AgentID == agentID {
			delete(s.memories, id)
			count++
		}
	}
	return count
}

// ListMemories returns a paginated list of memories for an agent.
func (s *Store) ListMemories(agentID string, page, limit int) ([]model.AgentMemory, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []model.AgentMemory
	for _, m := range s.memories {
		if m.AgentID != agentID {
			continue
		}
		all = append(all, m)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return paginate(all, page, limit)
}

// ============================================================
// LLM Models
// ============================================================

// AddLLMModel inserts an LLM model (seed loader).
func (s *Store) AddLLMModel(m model.LLMModel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.llmModels[m.ID] = m
}

// ListLLMModels returns all LLM models.
func (s *Store) ListLLMModels() []model.LLMModel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]model.LLMModel, 0, len(s.llmModels))
	for _, m := range s.llmModels {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].DisplayName < out[j].DisplayName
	})
	return out
}

// GetLLMModel returns an LLM model by ID.
func (s *Store) GetLLMModel(id string) (model.LLMModel, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.llmModels[id]
	return m, ok
}

// llmModelPtr returns a pointer to an LLM model or nil (must be called under at least RLock).
func (s *Store) llmModelPtr(id *string) *model.LLMModel {
	if id == nil {
		return nil
	}
	if m, ok := s.llmModels[*id]; ok {
		return &m
	}
	return nil
}

// ============================================================
// Pagination helper
// ============================================================

// paginate applies page/limit to a sorted slice and returns the page plus total count.
func paginate[T any](items []T, page, limit int) ([]T, int) {
	total := len(items)
	if total == 0 {
		return []T{}, 0
	}

	start := (page - 1) * limit
	if start >= total {
		return []T{}, total
	}

	end := start + limit
	if end > total {
		end = total
	}
	return items[start:end], total
}
