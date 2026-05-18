# ChatGPT MCP Connectors

面向 ChatGPT / MCP 工具调用场景的 Go 语言连接器实现。

当前仓库主要实现了一个 **飞书 / Lark 文档 MCP Connector**，用于让支持 Remote MCP 的客户端通过工具调用读取、创建和追加飞书文档内容。

## 项目状态

当前版本定位为：

- 可本地开发和调试的 MCP stdio server。
- 可部署到服务器并通过 HTTPS 暴露的 Remote MCP HTTP server。
- 面向飞书 / Lark 文档的只读和写入主链路 MVP。

需要注意：当前实现使用 **飞书应用 / tenant credential** 模式，适合自用、内部应用或单租户测试。多用户生产环境建议后续补充飞书用户 OAuth、加密 token 存储、用户级权限校验和写操作二次确认。

## 功能概览

仓库包含：

- 本地 MCP stdio 服务端：`cmd/feishu-doc-mcp-server`
- 远程 HTTP MCP 服务端：`cmd/feishu-doc-mcp-http-server`
- 本地调试 CLI：`cmd/feishu-doc-cli`
- 飞书 / Lark API 适配层：`internal/feishu`
- MCP JSON-RPC transport：`internal/mcp`
- Markdown 导入 / 导出工具：`internal/feishu`

远程 HTTP MCP server 是给 ChatGPT 网页端或其他 Web MCP 客户端连接的主路径；stdio server 主要用于本地 MCP 客户端和开发调试。

## 支持的 MCP 工具

| 工具名 | 说明 |
| --- | --- |
| `feishu_doc_resolve` | 解析飞书 / Lark 文档 URL 或 token，返回标准化文档身份。 |
| `feishu_doc_get_metadata` | 获取飞书 / Lark Docx 文档元信息。 |
| `feishu_doc_read` | 读取文档块结构，并导出 Markdown。 |
| `feishu_doc_create` | 创建飞书 / Lark Docx 文档，可选写入初始 Markdown 内容。 |
| `feishu_doc_append` | 向已有飞书 / Lark Docx 文档追加 Markdown 内容。 |

## 重要概念：不是把 GitHub 仓库地址填到 GPT 里

GPT 端不能直接使用：

```text
git@github.com:holtmiu/ChatGPT_MCP_Connectors.git
```

这个地址只是源码仓库地址。

你需要先把本仓库代码部署成一个可访问的 HTTP 服务，然后在 GPT / Web MCP 客户端里填写部署后的 MCP 地址，例如：

```text
https://your-domain.example/mcp
```

## 准备飞书 / Lark 应用

1. 在飞书开放平台或 Lark 开放平台创建内部应用。
2. 给应用开通文档读取和写入所需权限。
3. 如果使用 tenant credential 模式，需要把目标文档或父文件夹共享给该应用。
4. 获取 `FEISHU_APP_ID` 和 `FEISHU_APP_SECRET`。

## 本地构建与测试

```bash
go test ./...
go build ./cmd/feishu-doc-mcp-server
go build ./cmd/feishu-doc-mcp-http-server
go build ./cmd/feishu-doc-cli
```

## 启动远程 HTTP MCP Server

设置环境变量：

```bash
export FEISHU_APP_ID="你的飞书 App ID"
export FEISHU_APP_SECRET="你的飞书 App Secret"
export MCP_SERVER_API_KEY="一个长随机字符串"
export MCP_HTTP_ADDR=":8080"
```

启动服务：

```bash
go run ./cmd/feishu-doc-mcp-http-server
```

默认监听地址：

```text
:8080
```

HTTP 端点：

| 端点 | 方法 | 说明 |
| --- | --- | --- |
| `/healthz` | `GET` | 健康检查。 |
| `/mcp` | `POST` | JSON-RPC MCP 调用入口。 |
| `/mcp` | `OPTIONS` | CORS 预检。 |

如果设置了 `MCP_SERVER_API_KEY`，调用 `/mcp` 时需要带：

```text
Authorization: Bearer <MCP_SERVER_API_KEY>
```

## 本地验证

健康检查：

```bash
curl http://localhost:8080/healthz
```

MCP ping：

```bash
curl -s http://localhost:8080/mcp \
  -H 'content-type: application/json' \
  -H "authorization: Bearer $MCP_SERVER_API_KEY" \
  -d '{"jsonrpc":"2.0","id":1,"method":"ping"}'
```

列出工具：

```bash
curl -s http://localhost:8080/mcp \
  -H 'content-type: application/json' \
  -H "authorization: Bearer $MCP_SERVER_API_KEY" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'
```

## 使用 Docker 运行

构建镜像：

```bash
docker build -t feishu-doc-mcp .
```

运行容器：

```bash
docker run -p 8080:8080 \
  -e FEISHU_APP_ID="你的飞书 App ID" \
  -e FEISHU_APP_SECRET="你的飞书 App Secret" \
  -e MCP_SERVER_API_KEY="一个长随机字符串" \
  -e MCP_HTTP_ADDR=":8080" \
  feishu-doc-mcp
```

## 部署给 GPT / Web MCP 客户端使用

ChatGPT 网页端或其他 Web MCP 客户端通常不能访问你的本地 `localhost:8080`。

你需要把服务部署到公网，并放在 HTTPS 后面，例如：

```text
https://feishu-mcp.example.com/mcp
```

在 GPT / Web MCP 客户端中配置：

```text
MCP URL: https://feishu-mcp.example.com/mcp
Authorization: Bearer <MCP_SERVER_API_KEY>
```

## 写入行为

写入工具默认是 dry-run，即默认不真正修改飞书文档。

要执行真实写入，需要满足：

1. 飞书 / Lark 应用具备文档读写权限。
2. 目标文档或父文件夹已共享给应用。
3. 设置：

```bash
export FEISHU_DOC_WRITE_DRY_RUN_DEFAULT=false
```

或者在工具调用参数里显式传：

```json
{
  "dryRun": false
}
```

## 示例：读取文档

```json
{
  "name": "feishu_doc_read",
  "arguments": {
    "input": "https://example.feishu.cn/docx/xxxx",
    "format": "both"
  }
}
```

## 示例：追加内容

```json
{
  "name": "feishu_doc_append",
  "arguments": {
    "input": "https://example.feishu.cn/docx/xxxx",
    "markdown": "## 来自 ChatGPT 的追加内容\n\n这是一段通过 MCP 写入的内容。",
    "dryRun": false
  }
}
```

## 示例：创建文档

```json
{
  "name": "feishu_doc_create",
  "arguments": {
    "title": "ChatGPT 生成的文档",
    "markdown": "# 标题\n\n这是通过 MCP 创建的飞书文档。",
    "dryRun": false
  }
}
```

## 常用环境变量

| 环境变量 | 默认值 | 说明 |
| --- | --- | --- |
| `MCP_HTTP_ADDR` | `:8080` | 远程 MCP HTTP server 监听地址。 |
| `MCP_SERVER_API_KEY` | 空 | `/mcp` 的可选 Bearer token。生产环境建议设置。 |
| `FEISHU_PROVIDER` | `feishu` | `feishu` 或 `lark`。 |
| `FEISHU_BASE_URL` | 按 provider 自动设置 | 飞书 / Lark OpenAPI 基础地址。 |
| `FEISHU_APP_ID` | 空 | 飞书 / Lark 应用 ID。 |
| `FEISHU_APP_SECRET` | 空 | 飞书 / Lark 应用密钥。 |
| `FEISHU_TENANT_ACCESS_TOKEN` | 空 | 可选：直接使用预置 tenant access token。 |
| `FEISHU_API_TIMEOUT_MS` | `15000` | API 请求超时时间。 |
| `FEISHU_API_MAX_RETRIES` | `3` | 可重试请求的最大重试次数。 |
| `FEISHU_DOC_MAX_BLOCKS` | `3000` | 单次读取文档的最大块数量。 |
| `FEISHU_DOC_MAX_DEPTH` | `20` | 文档块递归读取最大深度。 |
| `FEISHU_DOC_WRITE_DRY_RUN_DEFAULT` | `true` | 写入工具是否默认 dry-run。 |

## 安全建议

- 生产环境必须通过 HTTPS 暴露 `/mcp`。
- 建议设置 `MCP_SERVER_API_KEY`，不要裸露 `/mcp`。
- 不要把 `FEISHU_APP_SECRET` 提交到仓库。
- 多用户生产环境建议使用飞书用户 OAuth，而不是共享 tenant token。
- 写操作建议保留二次确认或审计日志。

## 进一步文档

- 远程 MCP 部署说明：`doc/chatgpt-remote-mcp-deploy.md`
- 飞书文档模块 SDD：`doc/feishu-doc-module-sdd-spec.md`
