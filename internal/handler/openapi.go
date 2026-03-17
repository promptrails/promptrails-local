package handler

// OpenAPISpec is the OpenAPI 3.0 specification for the PromptRails Local emulator
var OpenAPISpec = `{
  "openapi": "3.0.3",
  "info": {
    "title": "PromptRails Local Emulator",
    "description": "In-memory API emulator for PromptRails. Use this to develop and test against the PromptRails API without a real backend.",
    "version": "0.1.0",
    "contact": {
      "name": "PromptRails",
      "url": "https://github.com/promptrails"
    },
    "license": {
      "name": "MIT"
    }
  },
  "servers": [
    {
      "url": "http://localhost:8080",
      "description": "Local emulator"
    }
  ],
  "tags": [
    {"name": "Agents", "description": "Agent management and execution"},
    {"name": "Prompts", "description": "Prompt template management"},
    {"name": "Executions", "description": "Execution history"},
    {"name": "Data Sources", "description": "External data source management"},
    {"name": "Credentials", "description": "Credential management"},
    {"name": "Chat", "description": "Chat sessions and messages"},
    {"name": "Traces", "description": "Execution traces"},
    {"name": "Scores", "description": "Evaluation scores"},
    {"name": "Approvals", "description": "Human-in-the-loop approvals"},
    {"name": "Webhook Triggers", "description": "Webhook-based agent triggers"},
    {"name": "MCP Tools", "description": "Model Context Protocol tools"},
    {"name": "Guardrails", "description": "Agent guardrails"},
    {"name": "Memories", "description": "Agent memory management"},
    {"name": "LLM Models", "description": "Available LLM models"},
    {"name": "Admin", "description": "Emulator administration"}
  ],
  "paths": {
    "/health": {
      "get": {
        "summary": "Health check",
        "operationId": "healthCheck",
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": {"type": "string"},
                    "version": {"type": "string"},
                    "service": {"type": "string"}
                  }
                }
              }
            }
          }
        }
      }
    },
    "/api/v1/agents": {
      "get": {
        "tags": ["Agents"],
        "summary": "List agents",
        "operationId": "listAgents",
        "parameters": [
          {"name": "page", "in": "query", "schema": {"type": "integer", "default": 1}},
          {"name": "limit", "in": "query", "schema": {"type": "integer", "default": 20}},
          {"name": "type", "in": "query", "schema": {"type": "string", "enum": ["simple", "chain", "multi_agent", "workflow", "composite"]}},
          {"name": "status", "in": "query", "schema": {"type": "string", "enum": ["active", "archived"]}},
          {"name": "name", "in": "query", "schema": {"type": "string"}}
        ],
        "responses": {
          "200": {"description": "Paginated list of agents"}
        }
      },
      "post": {
        "tags": ["Agents"],
        "summary": "Create agent",
        "operationId": "createAgent",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["name", "type"],
                "properties": {
                  "name": {"type": "string"},
                  "description": {"type": "string"},
                  "type": {"type": "string", "enum": ["simple", "chain", "multi_agent", "workflow", "composite"]}
                }
              }
            }
          }
        },
        "responses": {
          "201": {"description": "Agent created"}
        }
      }
    },
    "/api/v1/agents/{agentId}": {
      "get": {
        "tags": ["Agents"],
        "summary": "Get agent",
        "operationId": "getAgent",
        "parameters": [
          {"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {
          "200": {"description": "Agent details"},
          "404": {"description": "Not found"}
        }
      },
      "patch": {
        "tags": ["Agents"],
        "summary": "Update agent",
        "operationId": "updateAgent",
        "parameters": [
          {"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "name": {"type": "string"},
                  "description": {"type": "string"},
                  "status": {"type": "string"}
                }
              }
            }
          }
        },
        "responses": {
          "200": {"description": "Agent updated"}
        }
      },
      "delete": {
        "tags": ["Agents"],
        "summary": "Delete agent",
        "operationId": "deleteAgent",
        "parameters": [
          {"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {
          "204": {"description": "Deleted"}
        }
      }
    },
    "/api/v1/agents/{agentId}/versions": {
      "get": {
        "tags": ["Agents"],
        "summary": "List agent versions",
        "operationId": "listAgentVersions",
        "parameters": [
          {"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {"200": {"description": "List of versions"}}
      },
      "post": {
        "tags": ["Agents"],
        "summary": "Create agent version",
        "operationId": "createAgentVersion",
        "parameters": [
          {"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "version": {"type": "string"},
                  "config": {"type": "object"},
                  "input_schema": {"type": "object"},
                  "output_schema": {"type": "object"},
                  "set_current": {"type": "boolean"},
                  "message": {"type": "string"},
                  "prompt_ids": {"type": "array", "items": {"type": "object"}}
                }
              }
            }
          }
        },
        "responses": {"201": {"description": "Version created"}}
      }
    },
    "/api/v1/agents/{agentId}/versions/{versionId}/promote": {
      "put": {
        "tags": ["Agents"],
        "summary": "Promote agent version",
        "operationId": "promoteAgentVersion",
        "parameters": [
          {"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}},
          {"name": "versionId", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {"200": {"description": "Version promoted"}}
      }
    },
    "/api/v1/agents/{agentId}/execute": {
      "post": {
        "tags": ["Agents"],
        "summary": "Execute agent (simulated)",
        "operationId": "executeAgent",
        "parameters": [
          {"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "input": {"type": "object"},
                  "session_id": {"type": "string"},
                  "user_id": {"type": "string"},
                  "sync": {"type": "boolean"}
                }
              }
            }
          }
        },
        "responses": {"201": {"description": "Execution created"}}
      }
    },
    "/api/v1/agents/{agentId}/preview": {
      "post": {
        "tags": ["Agents"],
        "summary": "Preview agent",
        "operationId": "previewAgent",
        "parameters": [
          {"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {"200": {"description": "Preview result"}}
      }
    },
    "/api/v1/prompts": {
      "get": {
        "tags": ["Prompts"],
        "summary": "List prompts",
        "operationId": "listPrompts",
        "parameters": [
          {"name": "page", "in": "query", "schema": {"type": "integer"}},
          {"name": "limit", "in": "query", "schema": {"type": "integer"}}
        ],
        "responses": {"200": {"description": "Paginated list of prompts"}}
      },
      "post": {
        "tags": ["Prompts"],
        "summary": "Create prompt",
        "operationId": "createPrompt",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["name"],
                "properties": {
                  "name": {"type": "string"},
                  "description": {"type": "string"}
                }
              }
            }
          }
        },
        "responses": {"201": {"description": "Prompt created"}}
      }
    },
    "/api/v1/prompts/{promptId}": {
      "get": {"tags": ["Prompts"], "summary": "Get prompt", "operationId": "getPrompt", "parameters": [{"name": "promptId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Prompt details"}}},
      "patch": {"tags": ["Prompts"], "summary": "Update prompt", "operationId": "updatePrompt", "parameters": [{"name": "promptId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Updated"}}},
      "delete": {"tags": ["Prompts"], "summary": "Delete prompt", "operationId": "deletePrompt", "parameters": [{"name": "promptId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"204": {"description": "Deleted"}}}
    },
    "/api/v1/prompts/{promptId}/versions": {
      "get": {"tags": ["Prompts"], "summary": "List prompt versions", "operationId": "listPromptVersions", "parameters": [{"name": "promptId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Versions"}}},
      "post": {"tags": ["Prompts"], "summary": "Create prompt version", "operationId": "createPromptVersion", "parameters": [{"name": "promptId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"201": {"description": "Created"}}}
    },
    "/api/v1/prompts/{promptId}/versions/{versionId}/promote": {
      "put": {"tags": ["Prompts"], "summary": "Promote prompt version", "operationId": "promotePromptVersion", "parameters": [{"name": "promptId", "in": "path", "required": true, "schema": {"type": "string"}}, {"name": "versionId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Promoted"}}}
    },
    "/api/v1/prompts/{promptId}/preview": {
      "post": {"tags": ["Prompts"], "summary": "Preview prompt", "operationId": "previewPrompt", "parameters": [{"name": "promptId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Preview"}}}
    },
    "/api/v1/prompts/{promptId}/run": {
      "post": {"tags": ["Prompts"], "summary": "Run prompt (simulated)", "operationId": "runPrompt", "parameters": [{"name": "promptId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Run result"}}}
    },
    "/api/v1/executions": {
      "get": {"tags": ["Executions"], "summary": "List executions", "operationId": "listExecutions", "parameters": [{"name": "page", "in": "query", "schema": {"type": "integer"}}, {"name": "limit", "in": "query", "schema": {"type": "integer"}}, {"name": "agent_id", "in": "query", "schema": {"type": "string"}}, {"name": "session_id", "in": "query", "schema": {"type": "string"}}, {"name": "status", "in": "query", "schema": {"type": "string"}}], "responses": {"200": {"description": "Paginated executions"}}}
    },
    "/api/v1/executions/{executionId}": {
      "get": {"tags": ["Executions"], "summary": "Get execution", "operationId": "getExecution", "parameters": [{"name": "executionId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Execution details"}}}
    },
    "/api/v1/executions/{executionId}/pending-approval": {
      "get": {"tags": ["Executions"], "summary": "Get pending approval for execution", "operationId": "getPendingApproval", "parameters": [{"name": "executionId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Approval"}}}
    },
    "/api/v1/data-sources": {
      "get": {"tags": ["Data Sources"], "summary": "List data sources", "operationId": "listDataSources", "responses": {"200": {"description": "List"}}},
      "post": {"tags": ["Data Sources"], "summary": "Create data source", "operationId": "createDataSource", "responses": {"201": {"description": "Created"}}}
    },
    "/api/v1/data-sources/{dataSourceId}": {
      "get": {"tags": ["Data Sources"], "summary": "Get data source", "operationId": "getDataSource", "parameters": [{"name": "dataSourceId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Details"}}},
      "patch": {"tags": ["Data Sources"], "summary": "Update data source", "operationId": "updateDataSource", "parameters": [{"name": "dataSourceId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Updated"}}},
      "delete": {"tags": ["Data Sources"], "summary": "Delete data source", "operationId": "deleteDataSource", "parameters": [{"name": "dataSourceId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"204": {"description": "Deleted"}}}
    },
    "/api/v1/data-sources/{dataSourceId}/versions": {
      "get": {"tags": ["Data Sources"], "summary": "List versions", "operationId": "listDataSourceVersions", "parameters": [{"name": "dataSourceId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Versions"}}},
      "post": {"tags": ["Data Sources"], "summary": "Create version", "operationId": "createDataSourceVersion", "parameters": [{"name": "dataSourceId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"201": {"description": "Created"}}}
    },
    "/api/v1/data-sources/{dataSourceId}/query": {
      "post": {"tags": ["Data Sources"], "summary": "Execute query (mock)", "operationId": "queryDataSource", "parameters": [{"name": "dataSourceId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Mock query results"}}}
    },
    "/api/v1/credentials": {
      "get": {"tags": ["Credentials"], "summary": "List credentials", "operationId": "listCredentials", "responses": {"200": {"description": "List"}}},
      "post": {"tags": ["Credentials"], "summary": "Create credential", "operationId": "createCredential", "responses": {"201": {"description": "Created"}}}
    },
    "/api/v1/credentials/{credentialId}": {
      "get": {"tags": ["Credentials"], "summary": "Get credential", "operationId": "getCredential", "parameters": [{"name": "credentialId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Details"}}},
      "patch": {"tags": ["Credentials"], "summary": "Update credential", "operationId": "updateCredential", "parameters": [{"name": "credentialId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Updated"}}},
      "delete": {"tags": ["Credentials"], "summary": "Delete credential", "operationId": "deleteCredential", "parameters": [{"name": "credentialId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"204": {"description": "Deleted"}}}
    },
    "/api/v1/chat/sessions": {
      "get": {"tags": ["Chat"], "summary": "List sessions", "operationId": "listChatSessions", "responses": {"200": {"description": "List"}}},
      "post": {"tags": ["Chat"], "summary": "Create session", "operationId": "createChatSession", "responses": {"201": {"description": "Created"}}}
    },
    "/api/v1/chat/sessions/{sessionId}": {
      "get": {"tags": ["Chat"], "summary": "Get session", "operationId": "getChatSession", "parameters": [{"name": "sessionId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Details"}}},
      "delete": {"tags": ["Chat"], "summary": "Delete session", "operationId": "deleteChatSession", "parameters": [{"name": "sessionId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"204": {"description": "Deleted"}}}
    },
    "/api/v1/chat/sessions/{sessionId}/messages": {
      "get": {"tags": ["Chat"], "summary": "List messages", "operationId": "listChatMessages", "parameters": [{"name": "sessionId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Messages"}}},
      "post": {"tags": ["Chat"], "summary": "Send message (simulated)", "operationId": "sendChatMessage", "parameters": [{"name": "sessionId", "in": "path", "required": true, "schema": {"type": "string"}}], "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "required": ["content"], "properties": {"content": {"type": "string"}}}}}}, "responses": {"201": {"description": "Message sent"}}}
    },
    "/api/v1/traces": {
      "get": {"tags": ["Traces"], "summary": "List traces", "operationId": "listTraces", "responses": {"200": {"description": "List"}}}
    },
    "/api/v1/traces/{traceId}": {
      "get": {"tags": ["Traces"], "summary": "Get trace", "operationId": "getTrace", "parameters": [{"name": "traceId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Trace details"}}}
    },
    "/api/v1/scores": {
      "get": {"tags": ["Scores"], "summary": "List scores", "operationId": "listScores", "responses": {"200": {"description": "List"}}},
      "post": {"tags": ["Scores"], "summary": "Create score", "operationId": "createScore", "responses": {"201": {"description": "Created"}}}
    },
    "/api/v1/scores/{scoreId}": {
      "get": {"tags": ["Scores"], "summary": "Get score", "operationId": "getScore", "parameters": [{"name": "scoreId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Details"}}},
      "patch": {"tags": ["Scores"], "summary": "Update score", "operationId": "updateScore", "parameters": [{"name": "scoreId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Updated"}}},
      "delete": {"tags": ["Scores"], "summary": "Delete score", "operationId": "deleteScore", "parameters": [{"name": "scoreId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"204": {"description": "Deleted"}}}
    },
    "/api/v1/score-configs": {
      "get": {"tags": ["Scores"], "summary": "List score configs", "operationId": "listScoreConfigs", "responses": {"200": {"description": "List"}}},
      "post": {"tags": ["Scores"], "summary": "Create score config", "operationId": "createScoreConfig", "responses": {"201": {"description": "Created"}}}
    },
    "/api/v1/score-configs/{configId}": {
      "get": {"tags": ["Scores"], "summary": "Get score config", "operationId": "getScoreConfig", "parameters": [{"name": "configId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Details"}}},
      "patch": {"tags": ["Scores"], "summary": "Update score config", "operationId": "updateScoreConfig", "parameters": [{"name": "configId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Updated"}}},
      "delete": {"tags": ["Scores"], "summary": "Delete score config", "operationId": "deleteScoreConfig", "parameters": [{"name": "configId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"204": {"description": "Deleted"}}}
    },
    "/api/v1/approvals": {
      "get": {"tags": ["Approvals"], "summary": "List approvals", "operationId": "listApprovals", "responses": {"200": {"description": "List"}}}
    },
    "/api/v1/approvals/{approvalId}": {
      "get": {"tags": ["Approvals"], "summary": "Get approval", "operationId": "getApproval", "parameters": [{"name": "approvalId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Details"}}}
    },
    "/api/v1/approvals/{approvalId}/decide": {
      "post": {"tags": ["Approvals"], "summary": "Decide approval", "operationId": "decideApproval", "parameters": [{"name": "approvalId", "in": "path", "required": true, "schema": {"type": "string"}}], "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "required": ["decision"], "properties": {"decision": {"type": "string", "enum": ["approved", "rejected"]}, "reason": {"type": "string"}}}}}}, "responses": {"200": {"description": "Decided"}}}
    },
    "/api/v1/webhook-triggers": {
      "get": {"tags": ["Webhook Triggers"], "summary": "List triggers", "operationId": "listWebhookTriggers", "responses": {"200": {"description": "List"}}},
      "post": {"tags": ["Webhook Triggers"], "summary": "Create trigger", "operationId": "createWebhookTrigger", "responses": {"201": {"description": "Created with full token"}}}
    },
    "/api/v1/webhook-triggers/{triggerId}": {
      "get": {"tags": ["Webhook Triggers"], "summary": "Get trigger", "operationId": "getWebhookTrigger", "parameters": [{"name": "triggerId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Details"}}},
      "patch": {"tags": ["Webhook Triggers"], "summary": "Update trigger", "operationId": "updateWebhookTrigger", "parameters": [{"name": "triggerId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Updated"}}},
      "delete": {"tags": ["Webhook Triggers"], "summary": "Delete trigger", "operationId": "deleteWebhookTrigger", "parameters": [{"name": "triggerId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"204": {"description": "Deleted"}}}
    },
    "/hooks/{token}": {
      "post": {"tags": ["Webhook Triggers"], "summary": "Execute agent via webhook", "operationId": "hookExecute", "parameters": [{"name": "token", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"201": {"description": "Execution created"}}}
    },
    "/api/v1/mcp-tools": {
      "get": {"tags": ["MCP Tools"], "summary": "List MCP tools", "operationId": "listMCPTools", "responses": {"200": {"description": "List"}}},
      "post": {"tags": ["MCP Tools"], "summary": "Create MCP tool", "operationId": "createMCPTool", "responses": {"201": {"description": "Created"}}}
    },
    "/api/v1/mcp-tools/{toolId}": {
      "get": {"tags": ["MCP Tools"], "summary": "Get tool", "operationId": "getMCPTool", "parameters": [{"name": "toolId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Details"}}},
      "patch": {"tags": ["MCP Tools"], "summary": "Update tool", "operationId": "updateMCPTool", "parameters": [{"name": "toolId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Updated"}}},
      "delete": {"tags": ["MCP Tools"], "summary": "Delete tool", "operationId": "deleteMCPTool", "parameters": [{"name": "toolId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"204": {"description": "Deleted"}}}
    },
    "/api/v1/mcp-templates": {
      "get": {"tags": ["MCP Tools"], "summary": "List MCP templates", "operationId": "listMCPTemplates", "responses": {"200": {"description": "List"}}}
    },
    "/api/v1/mcp-templates/{templateId}": {
      "get": {"tags": ["MCP Tools"], "summary": "Get MCP template", "operationId": "getMCPTemplate", "parameters": [{"name": "templateId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Details"}}}
    },
    "/api/v1/agents/{agentId}/guardrails": {
      "get": {"tags": ["Guardrails"], "summary": "List guardrails", "operationId": "listGuardrails", "parameters": [{"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "List"}}},
      "post": {"tags": ["Guardrails"], "summary": "Create guardrail", "operationId": "createGuardrail", "parameters": [{"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"201": {"description": "Created"}}}
    },
    "/api/v1/guardrails/{guardrailId}": {
      "patch": {"tags": ["Guardrails"], "summary": "Update guardrail", "operationId": "updateGuardrail", "parameters": [{"name": "guardrailId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Updated"}}},
      "delete": {"tags": ["Guardrails"], "summary": "Delete guardrail", "operationId": "deleteGuardrail", "parameters": [{"name": "guardrailId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"204": {"description": "Deleted"}}}
    },
    "/api/v1/agents/{agentId}/memories": {
      "get": {"tags": ["Memories"], "summary": "List memories", "operationId": "listMemories", "parameters": [{"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "List"}}},
      "post": {"tags": ["Memories"], "summary": "Create memory", "operationId": "createMemory", "parameters": [{"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"201": {"description": "Created"}}},
      "delete": {"tags": ["Memories"], "summary": "Delete all memories", "operationId": "deleteAllMemories", "parameters": [{"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Deleted"}}}
    },
    "/api/v1/agents/{agentId}/memories/search": {
      "post": {"tags": ["Memories"], "summary": "Search memories", "operationId": "searchMemories", "parameters": [{"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Results"}}}
    },
    "/api/v1/memories/{memoryId}": {
      "get": {"tags": ["Memories"], "summary": "Get memory", "operationId": "getMemory", "parameters": [{"name": "memoryId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Details"}}},
      "delete": {"tags": ["Memories"], "summary": "Delete memory", "operationId": "deleteMemory", "parameters": [{"name": "memoryId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"204": {"description": "Deleted"}}}
    },
    "/api/v1/llm-models": {
      "get": {"tags": ["LLM Models"], "summary": "List available LLM models", "operationId": "listLLMModels", "responses": {"200": {"description": "List of models"}}}
    },
    "/admin/reset": {
      "post": {"tags": ["Admin"], "summary": "Reset all data and reload seed", "operationId": "adminReset", "responses": {"200": {"description": "Reset complete"}}}
    },
    "/admin/seed": {
      "post": {"tags": ["Admin"], "summary": "Reload seed data", "operationId": "adminSeed", "responses": {"200": {"description": "Seed complete"}}}
    },
    "/admin/store/stats": {
      "get": {"tags": ["Admin"], "summary": "Get store statistics", "operationId": "adminStats", "responses": {"200": {"description": "Stats"}}}
    }
  }
}`
