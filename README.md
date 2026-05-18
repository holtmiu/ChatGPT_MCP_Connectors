# ChatGPT MCP Connectors

Go implementation of MCP connectors for ChatGPT-oriented tool use.

## Feishu/Lark document connector

This repository currently contains a Feishu/Lark document connector scaffold with:

- an MCP stdio server: `cmd/feishu-doc-mcp-server`
- a local CLI client: `cmd/feishu-doc-cli`
- a Feishu/Lark API adapter: `internal/feishu`
- MCP JSON-RPC transport plumbing: `internal/mcp`
- Markdown export utilities inside `internal/feishu`

The server exposes read-first document tools. Write operations are intentionally dry-run first and require explicit enablement before real mutations are added.

## Build

```bash
go test ./...
go build ./cmd/feishu-doc-mcp-server
go build ./cmd/feishu-doc-cli
```

## Run MCP server

```bash
export FEISHU_APP_ID="..."
export FEISHU_APP_SECRET="..."
go run ./cmd/feishu-doc-mcp-server
```

The server speaks JSON-RPC over stdio and implements a minimal MCP-compatible tool surface:

- `initialize`
- `tools/list`
- `tools/call`

Tools:

- `feishu_doc_resolve`
- `feishu_doc_get_metadata`
- `feishu_doc_read`
- `feishu_doc_append`

## Run CLI

```bash
go run ./cmd/feishu-doc-cli resolve "https://..."
go run ./cmd/feishu-doc-cli metadata "https://..."
go run ./cmd/feishu-doc-cli read "https://..."
```

## Configuration

| Env var | Default | Description |
| --- | --- | --- |
| `FEISHU_PROVIDER` | `feishu` | `feishu` or `lark`. |
| `FEISHU_BASE_URL` | provider default | OpenAPI base URL. |
| `FEISHU_APP_ID` | empty | App ID for tenant-token auth. |
| `FEISHU_APP_SECRET` | empty | App secret for tenant-token auth. |
| `FEISHU_TENANT_ACCESS_TOKEN` | empty | Optional pre-provisioned tenant access token. |
| `FEISHU_API_TIMEOUT_MS` | `15000` | HTTP timeout. |
| `FEISHU_API_MAX_RETRIES` | `3` | Retry count for retryable API calls. |
| `FEISHU_DOC_MAX_BLOCKS` | `3000` | Safety cap for document reads. |
| `FEISHU_DOC_MAX_DEPTH` | `20` | Safety cap for recursive block reads. |
| `FEISHU_DOC_WRITE_DRY_RUN_DEFAULT` | `true` | Default dry-run for write tools. |
| `FEISHU_DOCX_METADATA_PATH_TEMPLATE` | `/open-apis/docx/v1/documents/%s` | Metadata endpoint template. |
| `FEISHU_DOCX_CHILDREN_PATH_TEMPLATE` | `/open-apis/docx/v1/documents/%s/blocks/%s/children` | Children endpoint template. |

Endpoint templates are configurable because Feishu/Lark document APIs evolve and tenant configurations may differ.
