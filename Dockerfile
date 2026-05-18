FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN go build -o /out/feishu-doc-mcp-http-server ./cmd/feishu-doc-mcp-http-server

FROM alpine:3.20
RUN adduser -D -u 10001 appuser
USER appuser
COPY --from=build /out/feishu-doc-mcp-http-server /usr/local/bin/feishu-doc-mcp-http-server
EXPOSE 8080
ENV MCP_HTTP_ADDR=:8080
ENTRYPOINT ["/usr/local/bin/feishu-doc-mcp-http-server"]
