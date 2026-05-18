# ChatGPT Remote MCP Deployment Guide

This guide covers the mainline path for connecting a web MCP client to Feishu/Lark documents through this repository.

## 1. Prepare Feishu/Lark app

Create or reuse a Feishu/Lark internal app and grant the document read/write scopes required by your tenant. Share the target document or parent folder with the app when using tenant access token mode.

## 2. Configure server

Required environment variables:

```bash
export FEISHU_APP_ID="..."
export FEISHU_APP_SECRET="..."
export MCP_SERVER_API_KEY="a-long-random-secret"
export MCP_HTTP_ADDR=":8080"
```

For real writes, set one of:

```bash
export FEISHU_DOC_WRITE_DRY_RUN_DEFAULT=false
```

or pass `dryRun:false` when calling `feishu_doc_create` or `feishu_doc_append`.

## 3. Run locally

```bash
go run ./cmd/feishu-doc-mcp-http-server
```

Health check:

```bash
curl http://localhost:8080/healthz
```

MCP ping:

```bash
curl -s http://localhost:8080/mcp \
  -H 'content-type: application/json' \
  -H "authorization: Bearer $MCP_SERVER_API_KEY" \
  -d '{"jsonrpc":"2.0","id":1,"method":"ping"}'
```

List tools:

```bash
curl -s http://localhost:8080/mcp \
  -H 'content-type: application/json' \
  -H "authorization: Bearer $MCP_SERVER_API_KEY" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'
```

## 4. Deploy for web clients

Deploy the HTTP server behind HTTPS. Configure the web MCP client with:

- MCP URL: `https://your-domain.example/mcp`
- Authorization: bearer token equal to `MCP_SERVER_API_KEY`

## 5. Tools

- `feishu_doc_resolve`: parse Feishu/Lark URLs and tokens.
- `feishu_doc_get_metadata`: fetch document metadata.
- `feishu_doc_read`: read document blocks and Markdown.
- `feishu_doc_create`: create a docx document and optionally seed it with Markdown.
- `feishu_doc_append`: append Markdown to an existing docx document.

## 6. Production hardening still recommended

This mainline implementation supports app/tenant credentials. For multi-user production use, add per-user Feishu OAuth, encrypted token storage, user-level permission checks, and an explicit user confirmation layer for write tools.
