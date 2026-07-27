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
    {"name": "Traces", "description": "Execution traces and metering summary"},
    {"name": "Agent Triggers", "description": "Webhook-based agent triggers"},
    {"name": "MCP Tools", "description": "Model Context Protocol tools"},
    {"name": "Guardrails", "description": "Agent guardrails"},
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
                  "type": {"type": "string", "enum": ["agent", "workflow"]}
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
        "summary": "Preview agent execution (dry-run, no LLM call)",
        "operationId": "previewAgent",
        "parameters": [
          {"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {"200": {"description": "Preview result"}}
      }
    },
    "/api/v1/agents/{agentId}/playground": {
      "get": {
        "tags": ["Agents"],
        "summary": "Get the agent's current version content to pre-fill the playground",
        "operationId": "getAgentPlayground",
        "parameters": [
          {"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {"200": {"description": "Agent with current version"}, "404": {"description": "Not found"}}
      },
      "post": {
        "tags": ["Executions"],
        "summary": "Execute an agent with temporary prompt content",
        "operationId": "executeAgentPlayground",
        "parameters": [
          {"name": "agentId", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["input", "prompt_override"],
                "properties": {
                  "input": {"type": "object"},
                  "version_id": {"type": "string"},
                  "prompt_override": {
                    "type": "object",
                    "properties": {
                      "system_prompt": {"type": "string"},
                      "user_prompt": {"type": "string"},
                      "input_schema": {"type": "object"}
                    }
                  }
                }
              }
            }
          }
        },
        "responses": {"201": {"description": "Playground execution created"}, "404": {"description": "Not found"}}
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
      "post": {"tags": ["Prompts"], "summary": "Preview prompt rendering (dry-run, no LLM call)", "operationId": "previewPrompt", "parameters": [{"name": "promptId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Preview"}}}
    },
    "/api/v1/executions": {
      "get": {"tags": ["Executions"], "summary": "List executions", "operationId": "listExecutions", "parameters": [{"name": "page", "in": "query", "schema": {"type": "integer"}}, {"name": "limit", "in": "query", "schema": {"type": "integer"}}, {"name": "agent_id", "in": "query", "schema": {"type": "string"}}, {"name": "session_id", "in": "query", "schema": {"type": "string"}}, {"name": "status", "in": "query", "schema": {"type": "string", "enum": ["pending", "running", "completed", "failed", "cancelled", "rejected", "waiting_approval", "cancel_requested"]}}], "responses": {"200": {"description": "Paginated executions"}}}
    },
    "/api/v1/executions/approval-inbox": {
      "get": {"tags": ["Executions"], "summary": "List executions awaiting approval", "operationId": "listApprovalInbox", "parameters": [{"name": "page", "in": "query", "schema": {"type": "integer"}}, {"name": "limit", "in": "query", "schema": {"type": "integer"}}], "responses": {"200": {"description": "Paginated executions parked at waiting_approval"}}}
    },
    "/api/v1/executions/{executionId}": {
      "get": {"tags": ["Executions"], "summary": "Get execution (fills one level of children)", "operationId": "getExecution", "parameters": [{"name": "executionId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Execution details"}}}
    },
    "/api/v1/executions/{executionId}/tree": {
      "get": {"tags": ["Executions"], "summary": "Get an execution with its full descendant tree", "operationId": "getExecutionTree", "parameters": [{"name": "executionId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Execution with populated children[]"}, "404": {"description": "Not found"}}}
    },
    "/api/v1/executions/{executionId}/cancel": {
      "post": {"tags": ["Executions"], "summary": "Request cancellation of a running execution", "operationId": "cancelExecution", "parameters": [{"name": "executionId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Execution updated to cancel_requested"}, "404": {"description": "Not found"}}}
    },
    "/api/v1/executions/{executionId}/approve": {
      "post": {"tags": ["Executions"], "summary": "Approve a run parked at waiting_approval", "operationId": "approveExecution", "parameters": [{"name": "executionId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Execution resumed with the approved call"}, "400": {"description": "Not awaiting approval"}, "404": {"description": "Not found"}}}
    },
    "/api/v1/executions/{executionId}/deny": {
      "post": {"tags": ["Executions"], "summary": "Deny a run parked at waiting_approval", "operationId": "denyExecution", "parameters": [{"name": "executionId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Execution resumed with the denial injected"}, "400": {"description": "Not awaiting approval"}, "404": {"description": "Not found"}}}
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
    "/api/v1/traces/summary": {
      "get": {"tags": ["Traces"], "summary": "Get trace summary statistics", "operationId": "getTraceSummary", "responses": {"200": {"description": "Aggregate metering statistics (total_traces, total_tokens, total_cost, avg_duration_ms, error_count, unique_models, unique_sessions)"}}}
    },
    "/api/v1/traces/{traceId}": {
      "get": {"tags": ["Traces"], "summary": "Get trace", "operationId": "getTrace", "parameters": [{"name": "traceId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Trace details"}}}
    },
    "/api/v1/triggers": {
      "get": {"tags": ["Agent Triggers"], "summary": "List triggers", "operationId": "listAgentTriggers", "responses": {"200": {"description": "List"}}},
      "post": {"tags": ["Agent Triggers"], "summary": "Create trigger", "operationId": "createAgentTrigger", "responses": {"201": {"description": "Created with full token"}}}
    },
    "/api/v1/triggers/{triggerId}": {
      "get": {"tags": ["Agent Triggers"], "summary": "Get trigger", "operationId": "getAgentTrigger", "parameters": [{"name": "triggerId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Details"}}},
      "patch": {"tags": ["Agent Triggers"], "summary": "Update trigger", "operationId": "updateAgentTrigger", "parameters": [{"name": "triggerId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "Updated"}}},
      "delete": {"tags": ["Agent Triggers"], "summary": "Delete trigger", "operationId": "deleteAgentTrigger", "parameters": [{"name": "triggerId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"204": {"description": "Deleted"}}}
    },
    "/hooks/{token}": {
      "post": {"tags": ["Agent Triggers"], "summary": "Execute agent via webhook", "operationId": "hookExecute", "parameters": [{"name": "token", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"201": {"description": "Execution created"}}}
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
    "/api/v1/llm-models": {
      "get": {"tags": ["LLM Models"], "summary": "List LLM models", "operationId": "listLLMModels", "responses": {"200": {"description": "List of models"}}}
    },
    "/api/v1/llm-models/available": {
      "get": {"tags": ["LLM Models"], "summary": "List available LLM models grouped by provider", "operationId": "listAvailableLLMModels", "responses": {"200": {"description": "Models grouped by provider"}}}
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
