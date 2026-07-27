# API Endpoints

All endpoints are under `/api/v1`. The `X-API-Key` header is accepted but any value works.

For the interactive API reference with try-it-out functionality, visit **http://localhost:8080/docs**.

## Agents

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/agents` | List agents (supports `?type=` (`agent`\|`workflow`), `?status=`, `?name=` filters) |
| `POST` | `/agents` | Create agent (`type` must be `agent` or `workflow`) |
| `GET` | `/agents/:agentId` | Get agent with current version |
| `PATCH` | `/agents/:agentId` | Update agent |
| `DELETE` | `/agents/:agentId` | Delete agent |
| `GET` | `/agents/:agentId/versions` | List agent versions |
| `POST` | `/agents/:agentId/versions` | Create agent version (owns `model_config`, `run_budget`, `approval_policy`, `cache_timeout`, `tools`) |
| `PUT` | `/agents/:agentId/versions/:versionId/promote` | Promote version to current |
| `POST` | `/agents/:agentId/execute` | Execute agent (simulated response) |
| `POST` | `/agents/:agentId/preview` | Preview agent execution (dry-run) |
| `GET` | `/agents/:agentId/playground` | Get current version content to pre-fill the playground |
| `POST` | `/agents/:agentId/playground` | Execute an agent with temporary prompt content |

## Prompts

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/prompts` | List prompts |
| `POST` | `/prompts` | Create prompt |
| `GET` | `/prompts/:promptId` | Get prompt with current version |
| `PATCH` | `/prompts/:promptId` | Update prompt |
| `DELETE` | `/prompts/:promptId` | Delete prompt |
| `GET` | `/prompts/:promptId/versions` | List prompt versions |
| `POST` | `/prompts/:promptId/versions` | Create prompt version |
| `PUT` | `/prompts/:promptId/versions/:versionId/promote` | Promote version |
| `POST` | `/prompts/:promptId/preview` | Preview rendered prompt (dry-run, no LLM call) |

Prompt versions are **content-only** in API v2 — model, sampling, tools, output
schema, and cache TTL live on the agent version, not the prompt version.

## Executions

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/executions` | List executions (supports `?agent_id=`, `?session_id=`, `?status=`) |
| `GET` | `/executions/approval-inbox` | List executions parked at `waiting_approval` |
| `GET` | `/executions/:executionId` | Get execution details (fills one level of `children`) |
| `GET` | `/executions/:executionId/tree` | Get execution with its full descendant tree |
| `POST` | `/executions/:executionId/cancel` | Request cancellation (→ `cancel_requested`) |
| `POST` | `/executions/:executionId/approve` | Approve a run parked at `waiting_approval` |
| `POST` | `/executions/:executionId/deny` | Deny a run parked at `waiting_approval` |

Statuses: `pending`, `running`, `completed`, `failed`, `cancelled`, `rejected`,
`waiting_approval`, `cancel_requested`.

## Data Sources

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/data-sources` | List data sources |
| `POST` | `/data-sources` | Create data source |
| `GET` | `/data-sources/:dataSourceId` | Get data source |
| `PATCH` | `/data-sources/:dataSourceId` | Update data source |
| `DELETE` | `/data-sources/:dataSourceId` | Delete data source |
| `GET` | `/data-sources/:dataSourceId/versions` | List versions |
| `POST` | `/data-sources/:dataSourceId/versions` | Create version |
| `POST` | `/data-sources/:dataSourceId/query` | Execute query (returns mock results) |

## Credentials

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/credentials` | List credentials |
| `POST` | `/credentials` | Create credential |
| `GET` | `/credentials/:credentialId` | Get credential |
| `PATCH` | `/credentials/:credentialId` | Update credential |
| `DELETE` | `/credentials/:credentialId` | Delete credential |

## Chat

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/chat/sessions` | List chat sessions |
| `POST` | `/chat/sessions` | Create chat session |
| `GET` | `/chat/sessions/:sessionId` | Get session |
| `DELETE` | `/chat/sessions/:sessionId` | Delete session |
| `GET` | `/chat/sessions/:sessionId/messages` | List messages |
| `POST` | `/chat/sessions/:sessionId/messages` | Send message (simulated reply) |

## Traces

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/traces` | List root traces |
| `GET` | `/traces/summary` | Aggregate metering rollup (traces, tokens, cost, duration, errors) |
| `GET` | `/traces/:traceId` | Get trace with spans |

## Agent Triggers

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/triggers` | List triggers |
| `POST` | `/triggers` | Create trigger; accepts `source` (`generic` / `slack` / `telegram` / `whatsapp` / `teams` / `schedule`), `source_config`, `reply_config` |
| `GET` | `/triggers/:triggerId` | Get trigger |
| `PATCH` | `/triggers/:triggerId` | Update trigger |
| `DELETE` | `/triggers/:triggerId` | Delete trigger |

### Public inbound endpoints (no auth)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/hooks/:token` | Generic webhook trigger |
| `POST` | `/api/v1/hooks/:token` | Generic webhook trigger (versioned path) |
| `POST` | `/api/v1/hooks/slack/:token` | Slack Events API |
| `POST` | `/api/v1/hooks/telegram/:token` | Telegram bot update |
| `POST` | `/api/v1/hooks/teams/:token` | Microsoft Teams message |
| `POST` `GET` | `/api/v1/hooks/whatsapp/:token` | WhatsApp Cloud API (POST for updates, GET for verify) |

The emulator dispatches every source uniformly — signature verification and channel-specific reply behavior are platform-only.

## Agent Virtual Filesystem

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/agents/:agentId/vfs` | List directory entries (query: `path`, `recursive`) |
| `GET` | `/agents/:agentId/vfs/file` | Read file content (query: `path`) |
| `PUT` | `/agents/:agentId/vfs/file` | Write file (`path`, `content`, `mode=overwrite|append`) |
| `GET` | `/agents/:agentId/vfs/stat` | Metadata for a single path |
| `POST` | `/agents/:agentId/vfs/mkdir` | Create directory (parents auto-created) |
| `POST` | `/agents/:agentId/vfs/move` | Rename or move (`from`, `to`) |
| `POST` | `/agents/:agentId/vfs/copy` | Copy file or subtree |
| `DELETE` | `/agents/:agentId/vfs` | Delete (query: `path`, `recursive`) |
| `GET` | `/agents/:agentId/vfs/grep` | Content search (query: `q`, `path`, `limit`) |
| `GET` | `/agents/:agentId/vfs/glob` | Name/path pattern search (query: `pattern`, `path`) |
| `GET` | `/agents/:agentId/vfs/usage` | Total bytes used by this agent |

## MCP Tools

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/mcp-tools` | List MCP tools |
| `POST` | `/mcp-tools` | Create MCP tool |
| `GET` | `/mcp-tools/:toolId` | Get tool |
| `PATCH` | `/mcp-tools/:toolId` | Update tool |
| `DELETE` | `/mcp-tools/:toolId` | Delete tool |
| `GET` | `/mcp-templates` | List MCP templates |
| `GET` | `/mcp-templates/:templateId` | Get template |

## Guardrails

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/agents/:agentId/guardrails` | List guardrails for agent |
| `POST` | `/agents/:agentId/guardrails` | Create guardrail |
| `PATCH` | `/guardrails/:guardrailId` | Update guardrail |
| `DELETE` | `/guardrails/:guardrailId` | Delete guardrail |

## LLM Models

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/llm-models` | List available LLM models |
| `GET` | `/llm-models/available` | List models grouped by provider |

## Response format

All responses use a standard envelope:

```json
// Single item
{"data": { ... }}

// Paginated list
{"data": [ ... ], "meta": {"total": 42, "page": 1, "limit": 20, "total_pages": 3}}

// Error
{"error": {"code": "not_found", "message": "agent not found"}}
```

## Pagination

List endpoints support `?page=` and `?limit=` query parameters. Default: page 1, limit 20, max 100.

## Simulated behavior

- **Agent execute** returns fake output with simulated token usage, cost, and duration; a version with an approval policy or approval-gated tool parks at `waiting_approval`
- **Prompt preview** renders the content-only templates (dry-run, no LLM call)
- **Chat send message** generates an automatic assistant reply
- **Traces** are auto-created for every execution
- **Data source query** returns mock rows
