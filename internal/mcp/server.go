package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type Handler interface {
	Tools() []Tool
	CallTool(ctx context.Context, name string, args json.RawMessage) (any, error)
}

type Server struct {
	handler Handler
	name    string
	version string
}

func NewServer(name, version string, handler Handler) *Server {
	return &Server{handler: handler, name: name, version: version}
}

func (s *Server) Serve(ctx context.Context, in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	encoder := json.NewEncoder(out)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			_ = encoder.Encode(Response{JSONRPC: "2.0", Error: &Error{Code: -32700, Message: "parse error", Data: err.Error()}})
			continue
		}
		if len(req.ID) == 0 {
			_ = s.HandleNotification(ctx, req)
			continue
		}
		resp := s.HandleRequest(ctx, req)
		if err := encoder.Encode(resp); err != nil {
			return fmt.Errorf("encode response: %w", err)
		}
	}
	return scanner.Err()
}

func (s *Server) HandleNotification(ctx context.Context, req Request) error {
	switch req.Method {
	case "notifications/initialized", "$/cancelRequest":
		return nil
	default:
		return nil
	}
}

func (s *Server) HandleRequest(ctx context.Context, req Request) Response {
	resp := Response{JSONRPC: "2.0", ID: req.ID}
	switch req.Method {
	case "initialize":
		resp.Result = map[string]any{
			"protocolVersion": "2024-11-05",
			"serverInfo":      map[string]any{"name": s.name, "version": s.version},
			"capabilities":    map[string]any{"tools": map[string]any{}},
		}
	case "ping":
		resp.Result = map[string]any{}
	case "tools/list":
		resp.Result = map[string]any{"tools": s.handler.Tools()}
	case "tools/call":
		var params ToolCallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			resp.Error = &Error{Code: -32602, Message: "invalid params", Data: err.Error()}
			return resp
		}
		result, err := s.handler.CallTool(ctx, params.Name, params.Arguments)
		if err != nil {
			resp.Result = ToolCallResult{IsError: true, Content: []ToolContent{{Type: "text", Text: err.Error()}}}
			return resp
		}
		raw, _ := json.MarshalIndent(result, "", "  ")
		resp.Result = ToolCallResult{Content: []ToolContent{{Type: "text", Text: string(raw)}}}
	default:
		resp.Error = &Error{Code: -32601, Message: "method not found"}
	}
	return resp
}
