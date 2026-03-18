# Configuration

## Command-line flags & environment variables

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `--port` | `PORT` | `8080` | Server port |
| `--seed` | `SEED` | `true` | Load built-in seed data on startup |
| `--fixtures` | `FIXTURES` | `""` | Load additional fixtures from a directory |
| `--log-level` | `LOG_LEVEL` | `info` | Log level (`debug`, `info`, `warn`, `error`) |
| `--cors-origins` | `CORS_ORIGINS` | `*` | CORS allowed origins |

## Examples

```bash
# Run on a different port
docker run -p 9090:9090 -e PORT=9090 bahattincinic/promptrails-local

# Run without seed data (empty database)
docker run -p 8080:8080 -e SEED=false bahattincinic/promptrails-local

# Run with debug logging
docker run -p 8080:8080 -e LOG_LEVEL=debug bahattincinic/promptrails-local

# Run with custom fixtures
docker run -p 8080:8080 \
  -v ./my-data:/fixtures \
  -e FIXTURES=/fixtures \
  bahattincinic/promptrails-local

# Combine: no default seed, only your fixtures
docker run -p 8080:8080 \
  -v ./my-data:/fixtures \
  -e SEED=false \
  -e FIXTURES=/fixtures \
  bahattincinic/promptrails-local
```

## Authentication

The emulator accepts **any value** for the `X-API-Key` header. It only checks that the header exists for compatibility with the SDKs. No real authentication is performed.

There is no workspace concept in the emulator. All data lives in a single flat namespace.

## Admin endpoints

These endpoints are specific to the emulator and not part of the real PromptRails API.

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/admin/reset` | Clear all data and reload seed fixtures |
| `POST` | `/admin/seed` | Reload seed fixtures (additive, does not clear existing data) |
| `GET` | `/admin/store/stats` | Return item counts for every resource type |

```bash
# Reset between test runs
curl -X POST http://localhost:8080/admin/reset

# Check how much data is loaded
curl http://localhost:8080/admin/store/stats
```
