# PromptRails Local

In-memory API emulator for [PromptRails](https://promptrails.ai) — like [LocalStack](https://localstack.cloud) for AWS, but for the PromptRails API.

Develop and test your PromptRails integrations locally without a real backend. All data lives in memory and resets on restart.

## Quick Start

### Docker (recommended)

```bash
docker run -p 8080:8080 promptrails/local
```

### Docker Compose

```yaml
services:
  promptrails-local:
    image: promptrails/local
    ports:
      - "8080:8080"
```

### From Source

```bash
go install github.com/promptrails/promptrails-local@latest
promptrails-local
```

### Binary

Download from [GitHub Releases](https://github.com/promptrails/promptrails-local/releases).

## Usage

Once running, point your SDK to the local emulator:

### Python SDK

```python
from promptrails import PromptRails

client = PromptRails(
    api_key="test-key",          # any value works
    base_url="http://localhost:8080"
)

# List pre-loaded agents
agents = client.agents.list()
print(f"Found {len(agents)} agents")

# Execute an agent (returns simulated response)
result = client.agents.execute("39wNZZu78VawB207IOPonkoP38J", input={"topic": "AI"})
print(result.output)
```

### JavaScript/TypeScript SDK

```typescript
import { PromptRails } from "@promptrails/sdk";

const client = new PromptRails({
  apiKey: "test-key",
  baseUrl: "http://localhost:8080",
});

const agents = await client.agents.list();
console.log(`Found ${agents.length} agents`);
```

### Go SDK

```go
client := promptrails.New("test-key",
    promptrails.WithBaseURL("http://localhost:8080"))

agents, _ := client.Agents.List(ctx, nil)
fmt.Printf("Found %d agents\n", len(agents))
```

### cURL

```bash
# List agents
curl http://localhost:8080/api/v1/agents \
  -H "X-API-Key: test"

# Execute an agent
curl -X POST http://localhost:8080/api/v1/agents/39wNZZu78VawB207IOPonkoP38J/execute \
  -H "X-API-Key: test" \
  -H "Content-Type: application/json" \
  -d '{"input": {"topic": "AI"}}'
```

## API Docs

Interactive API documentation is available at **http://localhost:8080/docs** (powered by [Scalar](https://scalar.com)).

## What's Included

### Pre-loaded Seed Data

The emulator starts with example data so you can immediately test:

- **6 Agents** — simple, chain, multi-agent, and approval-based agents
- **8 Prompts** — with multiple versions and template variables
- **2 Data Sources** — PostgreSQL examples with query templates
- **4 Credentials** — OpenAI, Gemini, PostgreSQL, Linear (masked)
- **47 LLM Models** — Full catalog (OpenAI, Anthropic, Gemini, DeepSeek, xAI, Fireworks, etc.)
- **7 MCP Tools** — Linear, builtin, datasource, and API tools
- **Guardrails, Memories** — Example entries for the Simple Chat Bot agent

### Supported Endpoints

| Resource | CRUD | Execute/Run | Notes |
|----------|------|-------------|-------|
| Agents | Yes | Yes (simulated) | + versions, promote, preview |
| Prompts | Yes | Yes (simulated) | + versions, promote, preview |
| Data Sources | Yes | Yes (mock results) | + versions |
| Executions | Read | Auto-created | From agent execute |
| Credentials | Yes | - | Validation skipped |
| Chat Sessions | Yes | Yes (simulated) | Multi-turn tracking |
| LLM Models | Read | - | From fixture catalog |
| Traces | Read | Auto-created | From executions |
| Scores | Yes | - | + score configs |
| Approvals | Yes | - | approve/reject flow |
| Webhook Triggers | Yes | Yes (hook endpoint) | Token-based |
| MCP Tools | Yes | - | + templates |
| Guardrails | Yes | - | CRUD only |
| Agent Memories | Yes | - | + search |

### Simulated Behavior

- **Agent execution** returns a fake response with simulated token usage, cost, and duration
- **Prompt run** returns simulated content with token metrics
- **Chat messages** generate automatic assistant replies
- **Traces** are auto-created for every execution
- **Data source queries** return mock rows

### Admin Endpoints

```bash
# Reset all data and reload seed
curl -X POST http://localhost:8080/admin/reset

# Reload seed data (additive)
curl -X POST http://localhost:8080/admin/seed

# View store statistics
curl http://localhost:8080/admin/store/stats
```

## Configuration

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `--port` | `PORT` | `8080` | Server port |
| `--seed` | `SEED` | `true` | Load seed data on startup |
| `--fixtures` | `FIXTURES` | `""` | Load additional fixtures from a directory |
| `--log-level` | `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `--cors-origins` | `CORS_ORIGINS` | `*` | CORS allowed origins |

### Custom Fixtures

You can load your own test data from a directory using `--fixtures`:

```bash
# Start with default seed + your custom agents
promptrails-local --fixtures ./my-fixtures/

# Start with only your data (no default seed)
promptrails-local --seed=false --fixtures ./my-fixtures/

# Docker with mounted fixtures
docker run -p 8080:8080 -v ./my-fixtures:/fixtures -e FIXTURES=/fixtures promptrails/local
```

The directory can contain any of these JSON files (all optional):

```
my-fixtures/
  agents.json
  agent_versions.json
  agent_version_prompts.json
  prompts.json
  prompt_versions.json
  data_sources.json
  data_source_versions.json
  credentials.json
  llm_models.json
  mcp_tools.json
  guardrails.json
  memories.json
```

The JSON format matches the PromptRails API — see `internal/seed/data/` for examples.

## Authentication

The emulator accepts any value for `X-API-Key`. No workspace header is needed — all data lives in a single namespace.

## Use Cases

- **Local development** — Build against the PromptRails API without internet access
- **Integration testing** — Use in CI/CD pipelines for automated tests
- **SDK development** — Test SDK changes against a controlled environment
- **Demos** — Showcase PromptRails integrations without a real account

## Not Included (out of scope)

- Authentication (register, login, JWT, SSO)
- Workspace management
- Payments/billing
- Real LLM calls (all responses are simulated)
- Media generation
- Notification channels

## Development

```bash
# Clone and build
git clone https://github.com/promptrails/promptrails-local.git
cd promptrails-local
go build -o promptrails-local .
./promptrails-local

# Run with Docker Compose
docker compose up --build
```

## License

MIT
