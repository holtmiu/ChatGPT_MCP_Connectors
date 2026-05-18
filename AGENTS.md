# Development Guide for Coding Agents

This repository is intended to be worked on by AI coding agents and humans together.

## Principles

- Keep the MCP server as the primary runtime artifact.
- Keep CLI clients small and focused on local debugging or smoke tests.
- Prefer typed internal packages over ad-hoc JSON handling at command boundaries.
- Keep write operations dry-run by default until permissions, idempotency, and user confirmation are implemented.
- Never log access tokens, refresh tokens, app secrets, or full document contents by default.
- Keep dependencies minimal; add external packages only when they materially reduce risk or complexity.

## Required checks

Run these before committing Go changes:

```bash
gofmt -w $(find . -name '*.go')
go test ./...
```

## Package boundaries

- `cmd/feishu-doc-mcp-server`: local stdio JSON-RPC/MCP server entrypoint.
- `cmd/feishu-doc-mcp-http-server`: remote HTTP JSON-RPC/MCP server entrypoint for web clients.
- `cmd/feishu-doc-cli`: local development CLI.
- `internal/config`: environment-backed configuration.
- `internal/feishu`: Feishu/Lark domain model, resolver, API adapter, service layer, Markdown export, write adapter.
- `internal/mcp`: minimal MCP JSON-RPC transport, HTTP transport, and tool registration.

## Safety defaults

- Write tools must default to dry-run.
- Real write execution must require `dryRun=false` and a deployed server configuration with valid Feishu scopes.
- If an API endpoint or permission model is uncertain, expose it via configuration and document the limitation.
- Return structured errors; do not hide upstream failure causes from developers, but do not leak secrets.
