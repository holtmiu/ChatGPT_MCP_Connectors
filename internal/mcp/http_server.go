package mcp

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTPServer struct {
	server *Server
	apiKey string
}

func NewHTTPServer(name, version string, handler Handler, apiKey string) *HTTPServer {
	return &HTTPServer{server: NewServer(name, version, handler), apiKey: apiKey}
}

func (h *HTTPServer) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.handleHealthz)
	mux.HandleFunc("/mcp", h.handleMCP)
	return withCORS(mux)
}

func (h *HTTPServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "time": time.Now().UTC().Format(time.RFC3339)})
}

func (h *HTTPServer) handleMCP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "use POST for JSON-RPC MCP requests"})
		return
	}
	if !h.authorized(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "missing or invalid bearer token"})
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 16*1024*1024))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "empty JSON-RPC body"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()
	if body[0] == '[' {
		h.handleBatch(ctx, w, body)
		return
	}
	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSON(w, http.StatusOK, Response{JSONRPC: "2.0", Error: &Error{Code: -32700, Message: "parse error", Data: err.Error()}})
		return
	}
	if len(req.ID) == 0 {
		_ = h.server.HandleNotification(ctx, req)
		w.WriteHeader(http.StatusAccepted)
		return
	}
	writeJSON(w, http.StatusOK, h.server.HandleRequest(ctx, req))
}

func (h *HTTPServer) handleBatch(ctx context.Context, w http.ResponseWriter, body []byte) {
	var reqs []Request
	if err := json.Unmarshal(body, &reqs); err != nil {
		writeJSON(w, http.StatusOK, Response{JSONRPC: "2.0", Error: &Error{Code: -32700, Message: "parse error", Data: err.Error()}})
		return
	}
	responses := make([]Response, 0, len(reqs))
	for _, req := range reqs {
		if len(req.ID) == 0 {
			_ = h.server.HandleNotification(ctx, req)
			continue
		}
		responses = append(responses, h.server.HandleRequest(ctx, req))
	}
	if len(responses) == 0 {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	writeJSON(w, http.StatusOK, responses)
}

func (h *HTTPServer) authorized(r *http.Request) bool {
	if h.apiKey == "" {
		return true
	}
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return false
	}
	got := strings.TrimSpace(strings.TrimPrefix(auth, prefix))
	return subtle.ConstantTimeCompare([]byte(got), []byte(h.apiKey)) == 1
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, MCP-Protocol-Version")
		next.ServeHTTP(w, r)
	})
}
