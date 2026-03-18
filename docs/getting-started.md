# Getting Started

## Installation

### Docker (recommended)

```bash
docker run -p 8080:8080 bahattincinic/promptrails-local
```

### Docker Compose

```yaml
services:
  promptrails-local:
    image: bahattincinic/promptrails-local
    ports:
      - "8080:8080"
```

Then:

```bash
docker compose up
```

### From source

```bash
go install github.com/promptrails/promptrails-local@latest
promptrails-local
```

### Binary download

Download from [GitHub Releases](https://github.com/promptrails/promptrails-local/releases).

## Verify it's running

```bash
curl http://localhost:8080/health
# {"status":"ok","version":"0.1.0"}
```

## Try the API

```bash
# List pre-loaded agents
curl http://localhost:8080/api/v1/agents -H "X-API-Key: test"

# Execute an agent
curl -X POST http://localhost:8080/api/v1/agents/39wNZZu78VawB207IOPonkoP38J/execute \
  -H "X-API-Key: test" \
  -H "Content-Type: application/json" \
  -d '{"input": {"topic": "AI"}}'

# List prompts
curl http://localhost:8080/api/v1/prompts -H "X-API-Key: test"

# Create a chat session
curl -X POST http://localhost:8080/api/v1/chat/sessions \
  -H "X-API-Key: test" \
  -H "Content-Type: application/json" \
  -d '{"agent_id": "3A1tXOt9iovkA7LEusDSjcKbJQM", "title": "My Chat"}'
```

## Interactive API Docs

Open **http://localhost:8080/docs** in your browser for the Scalar-powered API reference where you can try all endpoints directly.

## Connect your SDK

### Python

```python
from promptrails import PromptRails

client = PromptRails(api_key="test-key", base_url="http://localhost:8080")

agents = client.agents.list()
for agent in agents.data:
    print(f"{agent.name} ({agent.type})")
```

### JavaScript / TypeScript

```typescript
import { PromptRails } from "@promptrails/sdk";

const client = new PromptRails({
  apiKey: "test-key",
  baseUrl: "http://localhost:8080",
});

const agents = await client.agents.list();
agents.data.forEach((a) => console.log(`${a.name} (${a.type})`));
```

### Go

```go
client := promptrails.NewClient("test-key",
    promptrails.WithBaseURL("http://localhost:8080"))

agents, _ := client.Agents.List(ctx, nil)
for _, a := range agents.Data {
    fmt.Printf("%s (%s)\n", a.Name, a.Type)
}
```
