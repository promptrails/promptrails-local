# API Endpoints

All endpoints are under `/api/v1`. The `X-API-Key` header is accepted but any value works.

For the interactive API reference with try-it-out functionality, visit **http://localhost:8080/docs**.

## Agents

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/agents` | List agents (supports `?type=`, `?status=`, `?name=` filters) |
| `POST` | `/agents` | Create agent |
| `GET` | `/agents/:agentId` | Get agent with current version |
| `PATCH` | `/agents/:agentId` | Update agent |
| `DELETE` | `/agents/:agentId` | Delete agent |
| `GET` | `/agents/:agentId/versions` | List agent versions |
| `POST` | `/agents/:agentId/versions` | Create agent version |
| `PUT` | `/agents/:agentId/versions/:versionId/promote` | Promote version to current |
| `POST` | `/agents/:agentId/execute` | Execute agent (simulated response) |
| `POST` | `/agents/:agentId/preview` | Preview agent |

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
| `POST` | `/prompts/:promptId/preview` | Preview rendered prompt |
| `POST` | `/prompts/:promptId/run` | Run prompt (simulated response) |

## Executions

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/executions` | List executions (supports `?agent_id=`, `?session_id=`, `?status=`) |
| `GET` | `/executions/:executionId` | Get execution details |
| `GET` | `/executions/:executionId/pending-approval` | Get pending approval for execution |

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
| `GET` | `/traces/:traceId` | Get trace with spans |

## Scores

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/scores` | List scores |
| `POST` | `/scores` | Create score |
| `GET` | `/scores/:scoreId` | Get score |
| `PATCH` | `/scores/:scoreId` | Update score |
| `DELETE` | `/scores/:scoreId` | Delete score |
| `GET` | `/score-configs` | List score configs |
| `POST` | `/score-configs` | Create score config |
| `GET` | `/score-configs/:configId` | Get score config |
| `PATCH` | `/score-configs/:configId` | Update score config |
| `DELETE` | `/score-configs/:configId` | Delete score config |

## Approvals

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/approvals` | List approval requests |
| `GET` | `/approvals/:approvalId` | Get approval |
| `POST` | `/approvals/:approvalId/decide` | Approve or reject (`{"decision": "approved"}`) |

## Webhook Triggers

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/webhook-triggers` | List triggers |
| `POST` | `/webhook-triggers` | Create trigger (returns full token once) |
| `GET` | `/webhook-triggers/:triggerId` | Get trigger |
| `PATCH` | `/webhook-triggers/:triggerId` | Update trigger |
| `DELETE` | `/webhook-triggers/:triggerId` | Delete trigger |
| `POST` | `/hooks/:token` | Execute agent via webhook (public, no auth) |

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

## Memories

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/agents/:agentId/memories` | List memories |
| `POST` | `/agents/:agentId/memories` | Create memory |
| `POST` | `/agents/:agentId/memories/search` | Search memories |
| `DELETE` | `/agents/:agentId/memories` | Delete all memories for agent |
| `GET` | `/memories/:memoryId` | Get memory |
| `DELETE` | `/memories/:memoryId` | Delete memory |

## LLM Models

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/llm-models` | List available LLM models |

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

- **Agent execute** returns fake output with simulated token usage, cost, and duration
- **Prompt run** returns simulated content with token metrics
- **Chat send message** generates an automatic assistant reply
- **Traces** are auto-created for every execution
- **Data source query** returns mock rows
