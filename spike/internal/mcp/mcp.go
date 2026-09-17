// Package mcp implements a minimal JSON-RPC 2.0 server for the MCP
// stdio transport (spec 2025-06-18: client launches server as
// subprocess, newline-delimited JSON-RPC over stdin/stdout) plus a
// WebSocket dial helper shared by the bridge.
package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"

	"github.com/coder/websocket"
)

// Conn aliases the websocket connection the bridge listens on.
type Conn = websocket.Conn

// Dial opens the bridge's outbound WebSocket to a room.
func Dial(url string) (*Conn, error) {
	ctx := context.Background()
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ---------- minimal MCP stdio server ----------

type Tool struct {
	Name        string
	Description string
	InputSchema map[string]any
	Handler     func(args map[string]any) (any, error)
}

type Server struct {
	mu    sync.Mutex
	tools map[string]Tool
}

func NewServer() *Server {
	return &Server{tools: map[string]Tool{}}
}

// Obj builds a JSON-Schema object schema from a name→type map.
func Obj(props map[string]any) map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": props,
	}
}

func (s *Server) Tool(name, desc string, schema map[string]any, h func(map[string]any) (any, error)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tools[name] = Tool{Name: name, Description: desc, InputSchema: schema, Handler: h}
}

func (s *Server) ToolNames() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	names := make([]string, 0, len(s.tools))
	for n := range s.tools {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// ServeStdio runs the MCP initialize/tools-list/tools-call flow over the
// given reader/writer per the 2025-06-18 stdio transport.
func (s *Server) ServeStdio(r io.Reader, w io.Writer) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 1024*1024), 16*1024*1024)
	enc := json.NewEncoder(w)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var req struct {
			JSONRPC string          `json:"jsonrpc"`
			ID      json.RawMessage `json:"id"`
			Method  string          `json:"method"`
			Params  struct {
				Name      string          `json:"name"`
				Arguments json.RawMessage `json:"arguments"`
			} `json:"params"`
		}
		if err := json.Unmarshal(line, &req); err != nil {
			continue
		}
		if req.JSONRPC != "2.0" {
			continue
		}
		var result any
		switch req.Method {
		case "initialize":
			result = map[string]any{
				"protocolVersion": "2025-06-18",
				"capabilities":    map[string]any{"tools": map[string]any{"listChanged": false}},
				"serverInfo":      map[string]any{"name": "ocre-bridge", "version": "0.1.0"},
			}
		case "notifications/initialized":
			continue // notification: no response
		case "tools/list":
			s.mu.Lock()
			tools := make([]map[string]any, 0, len(s.tools))
			for _, t := range s.tools {
				tools = append(tools, map[string]any{
					"name": t.Name, "description": t.Description, "inputSchema": t.InputSchema,
				})
			}
			s.mu.Unlock()
			result = map[string]any{"tools": tools}
		case "tools/call":
			s.mu.Lock()
			tool, ok := s.tools[req.Params.Name]
			s.mu.Unlock()
			if !ok {
				writeErr(enc, req.ID, -32602, fmt.Sprintf("unknown tool: %s", req.Params.Name))
				continue
			}
			args := map[string]any{}
			if len(req.Params.Arguments) > 0 {
				json.Unmarshal(req.Params.Arguments, &args)
			}
			out, err := tool.Handler(args)
			if err != nil {
				// Tool errors are results with isError, not JSON-RPC errors.
				result = map[string]any{
					"content": []map[string]any{{"type": "text", "text": err.Error()}},
					"isError": true,
				}
			} else {
				text, _ := json.Marshal(out)
				result = map[string]any{
					"content": []map[string]any{{"type": "text", "text": string(text)}},
				}
			}
		default:
			if len(req.ID) > 0 {
				writeErr(enc, req.ID, -32601, "method not found: "+req.Method)
			}
			continue
		}
		if len(req.ID) == 0 {
			continue // notification
		}
		enc.Encode(map[string]any{"jsonrpc": "2.0", "id": json.RawMessage(req.ID), "result": result})
	}
}

func writeErr(enc *json.Encoder, id json.RawMessage, code int, msg string) {
	enc.Encode(map[string]any{
		"jsonrpc": "2.0",
		"id":      json.RawMessage(id),
		"error":   map[string]any{"code": code, "message": msg},
	})
}

var _ = os.Stdin // referenced by callers via ServeStdio(os.Stdin, os.Stdout)