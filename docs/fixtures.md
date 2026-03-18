# Fixtures

The emulator can load data from JSON fixture files. There are two sources:

1. **Built-in seed data** — Embedded in the binary, loaded by default (`--seed=true`)
2. **Custom fixtures** — Loaded from a directory you provide (`--fixtures ./path/`)

Both can be used together. Built-in seed loads first, then custom fixtures are added on top.

## Supported files

Place any of these JSON files in your fixtures directory. All files are optional — only the ones present will be loaded.

| File | Description |
|------|-------------|
| `agents.json` | Agent definitions |
| `agent_versions.json` | Agent version configurations |
| `agent_version_prompts.json` | Links between agent versions and prompt versions |
| `prompts.json` | Prompt templates |
| `prompt_versions.json` | Prompt version definitions with system/user prompts |
| `data_sources.json` | Data source definitions |
| `data_source_versions.json` | Data source version configurations with query templates |
| `credentials.json` | Credential entries (no real secrets needed) |
| `llm_models.json` | LLM model catalog entries |
| `mcp_tools.json` | MCP tool definitions |
| `guardrails.json` | Agent guardrail configurations |
| `memories.json` | Agent memory entries |

## File formats

Each file is a JSON array of objects. The `id` field should be a KSUID string (27 characters). If you don't provide an `id`, the emulator will not generate one — always include IDs in your fixtures.

The `workspace_id` field is accepted but ignored for filtering purposes (the emulator has no workspace concept).

### agents.json

```json
[
  {
    "id": "YOUR_KSUID_HERE",
    "workspace_id": "any-value",
    "name": "My Agent",
    "description": "Description of the agent",
    "type": "simple",
    "status": "active",
    "labels": []
  }
]
```

**Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | yes | KSUID (27 chars) |
| `workspace_id` | string | yes | Any value |
| `name` | string | yes | Display name |
| `description` | string | no | Description |
| `type` | string | yes | `simple`, `chain`, `multi_agent`, `workflow`, `composite` |
| `status` | string | yes | `active` or `archived` |
| `labels` | array | no | String array of labels |

### agent_versions.json

```json
[
  {
    "id": "VERSION_KSUID",
    "agent_id": "AGENT_KSUID",
    "version": "1.0.0",
    "config": {
      "prompt_id": "PROMPT_KSUID"
    },
    "input_schema": {
      "type": "object",
      "required": ["topic"],
      "properties": {
        "topic": { "type": "string" }
      }
    },
    "output_schema": null,
    "is_current": true,
    "message": "Initial version"
  }
]
```

**Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | yes | KSUID |
| `agent_id` | string | yes | References `agents.json` id |
| `version` | string | yes | Semver string (e.g. `1.0.0`) |
| `config` | object | yes | Agent configuration (prompt_id, temperature, etc.) |
| `input_schema` | object | no | JSON Schema for input validation |
| `output_schema` | object | no | JSON Schema for output |
| `is_current` | boolean | yes | Whether this is the active version |
| `message` | string | no | Version message / changelog |

### agent_version_prompts.json

Links agent versions to prompt versions. Required for chain and multi-agent types.

```json
[
  {
    "id": "AVP_KSUID",
    "agent_version_id": "VERSION_KSUID",
    "prompt_version_id": "PROMPT_VERSION_KSUID",
    "role": "main",
    "sort_order": 0
  }
]
```

### prompts.json

```json
[
  {
    "id": "PROMPT_KSUID",
    "workspace_id": "any-value",
    "name": "My Prompt",
    "description": "Generates a summary",
    "status": "active"
  }
]
```

### prompt_versions.json

```json
[
  {
    "id": "PV_KSUID",
    "prompt_id": "PROMPT_KSUID",
    "version": "1.0.0",
    "system_prompt": "You are a helpful assistant.",
    "user_prompt": "Summarize the following: {{ text }}",
    "llm_model_id": "2N8Dx8a2a1qUspBKbP02vqP33iL",
    "temperature": 0.7,
    "max_tokens": 1024,
    "top_p": null,
    "input_schema": {
      "parameters": [
        {
          "name": "text",
          "type": "string",
          "required": true,
          "description": "Text to summarize"
        }
      ]
    },
    "output_schema": null,
    "is_current": true,
    "config": {},
    "message": "Initial version",
    "cache_timeout": 0
  }
]
```

**Key fields:**

| Field | Type | Description |
|-------|------|-------------|
| `system_prompt` | string | System prompt text (supports `{{ variable }}` syntax) |
| `user_prompt` | string | User prompt text (supports `{{ variable }}` syntax) |
| `llm_model_id` | string | References a model from `llm_models.json` |
| `temperature` | number | 0.0 - 1.0 |
| `max_tokens` | number | Max output tokens |
| `input_schema.parameters` | array | Template variable definitions |

### data_sources.json

```json
[
  {
    "id": "DS_KSUID",
    "workspace_id": "any-value",
    "name": "Users",
    "description": "User database",
    "type": "postgresql",
    "status": "active"
  }
]
```

**Supported types:** `postgresql`, `mysql`, `bigquery`, `snowflake`, `redshift`, `mssql`, `clickhouse`, `static_file`

### data_source_versions.json

```json
[
  {
    "id": "DSV_KSUID",
    "data_source_id": "DS_KSUID",
    "version": "1.0.0",
    "credential_id": "CRED_KSUID",
    "connection_config": {},
    "query_template": "SELECT * FROM users WHERE role = '{{ role }}'",
    "parameters": [
      {
        "name": "role",
        "type": "string",
        "required": false,
        "description": "Filter by role"
      }
    ],
    "is_current": true,
    "cache_timeout": 3600,
    "output_format": "json",
    "message": "Initial version"
  }
]
```

### credentials.json

```json
[
  {
    "id": "CRED_KSUID",
    "workspace_id": "any-value",
    "name": "OpenAI",
    "type": "openai",
    "category": "llm",
    "description": "",
    "masked_content": "sk-...abcd",
    "is_default": true,
    "schema_type": "",
    "is_valid": true
  }
]
```

**Categories:** `llm`, `data_warehouse`, `tools`, `speech`, `image`, `video`

**LLM types:** `openai`, `anthropic`, `gemini`, `deepseek`, `fireworks`, `xai`, `openrouter`, `together_ai`, `mistral`, `cohere`, `groq`

> Note: The emulator does not use real credentials. `masked_content` is what gets returned in API responses. No `encrypted_value` is needed.

### llm_models.json

```json
[
  {
    "id": "MODEL_KSUID",
    "provider": "openai",
    "model_id": "gpt-4o",
    "display_name": "GPT-4o",
    "input_price": 2.50,
    "output_price": 10.00,
    "max_tokens": 16384,
    "supports_vision": true,
    "supports_tools": true,
    "supports_json": true,
    "supports_streaming": true,
    "is_active": true
  }
]
```

### mcp_tools.json

```json
[
  {
    "id": "TOOL_KSUID",
    "workspace_id": "any-value",
    "name": "Time",
    "description": "Returns current time",
    "type": "builtin",
    "config": {
      "function": "current_time"
    },
    "schema": null,
    "credential_id": null,
    "template_id": null,
    "status": "active",
    "is_active": true
  }
]
```

**Tool types:** `builtin`, `remote_mcp`, `datasource`, `api`

### guardrails.json

```json
[
  {
    "id": "GUARD_KSUID",
    "agent_id": "AGENT_KSUID",
    "type": "output",
    "scanner_type": "pii",
    "action": "redact",
    "config": {
      "pii_types": ["email", "phone", "ssn"]
    },
    "is_active": true,
    "sort_order": 0
  }
]
```

### memories.json

```json
[
  {
    "id": "MEM_KSUID",
    "workspace_id": "any-value",
    "agent_id": "AGENT_KSUID",
    "content": "The user prefers concise answers",
    "memory_type": "fact",
    "importance": 0.9,
    "access_count": 0
  }
]
```

**Memory types:** `fact`, `event`, `context`, `procedure`, `semantic`

## Generating KSUIDs

If you need to generate KSUID values for your fixtures:

```bash
# Go
go run github.com/segmentio/ksuid/cmd/ksuid@latest

# Python
pip install ksuid
python -c "from ksuid import Ksuid; print(str(Ksuid()))"

# Node.js
npx ksuid
```

## Built-in seed data

The default seed includes:

| Resource | Count | Notes |
|----------|-------|-------|
| Agents | 6 | simple, chain, multi_agent, approval |
| Agent Versions | 15 | Multiple versions per agent |
| Prompts | 8 | Various prompt templates |
| Prompt Versions | 12 | With system/user prompts and schemas |
| Data Sources | 2 | PostgreSQL examples |
| Credentials | 4 | OpenAI, Gemini, PostgreSQL, Linear |
| LLM Models | 42 | Full catalog across all providers |
| MCP Tools | 6 | Builtin, datasource, remote, API |
| Guardrails | 1 | PII redaction example |
| Memories | 4 | Fact, procedure, semantic types |

See `internal/seed/data/` in the repository for the exact fixture files.
