# ChatGPT MCP Connectors

Go implementation of MCP connectors for ChatGPT-oriented tool use.

## Feishu/Lark document connector

This repository contains a Feishu/Lark document connector with:

- a local MCP stdio server: `cmd/feishu-doc-mcp-server`
- a remote HTTP MCP server: `cmd/feishu-doc-mcp-http-server`
- a local CLI client: `cmd/feishu-doc-cli`
- a Feishu/Lark API adapter: `internal/feishu`
- MCP JSON-RPC transport plumbing: `internal/mcp`
- Markdown import/export utilities inside `internal/feishu`

The remote HTTP server is the main path for connecting web MCP clients to Feishu/Lark documents. The stdio server is kept for local MCP clients and development.

## Build and test

```bash
go test ./...
go build ./cmd/feishu-doc-mcp-server
go build ./cmd/feishu-doc-mcp-http-server
go build ./cmd/feishu-doc-cli
```

## Run remote HTTP MCP server

```bash
export FEISHU_APP_ID="..."
export FEISHU_APP_SECRET="..."
export MCP_SERVER_API_KEY="change-me"
go run ./cmd/feishu-doc-mcp-http-server
```

Default address: `:8080`.

Endpoints:

- `GET /healthz`
- `POST /mcp`
- `OPTIONS /mcp`

Deploy this server behind HTTPS and configure your web MCP client with the public `/mcp` URL.

## Run local stdio MCP server

```bash
go run ./cmd/feishu-doc-mcp-server
```

The server implements:

- `initialize`
- `notifications/initialized`
- `ping`
- `tools/list`
- `tools/call`

Tools:

- `feishu_doc_resolve`
- `feishu_doc_get_metadata`
- `feishu_doc_read`
- `feishu_doc_create`
- `feishu_doc_append`

## Run CLI

```bash
go run ./cmd/feishu-doc-cli resolve "https://..."
go run ./cmd/feishu-doc-cli metadata "https://..."
go run ./cmd/feishu-doc-cli read "https://..."
go run ./cmd/feishu-doc-cli create "New title" "# Hello"
go run ./cmd/feishu-doc-cli append "https://..." "## Added from CLI"
```

## Write behavior

Write tools are dry-run by default. To execute real writes, grant the Feishu/Lark app docx read/write scopes, share target docs or folders with the app, and set `FEISHU_DOC_WRITE_DRY_RUN_DEFAULT=false` or pass `dryRun:false` in the tool arguments.

See `doc/chatgpt-remote-mcp-deploy.md` for deployment details.
