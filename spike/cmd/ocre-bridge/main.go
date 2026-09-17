// ocre-bridge is the desktop agent-bridge spike: it connects OUT to an
// ocre room over WebSocket (Socket Mode pattern - no inbound endpoints),
// listens for agent.mentioned events, drives the owner's Claude Code via
// `claude -p` (headless, subscription auth), and exposes the review
// actions to the runtime as a local stdio MCP server whose tools map
// 1:1 to the comment model.
//
// Per docs/knowledgebase/concepts/research/agent-participation.md and
// the platform decision: agents run on the owner's desktop with the
// owner's credentials; ocre never sees an API key or bills tokens.
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"ocre.dev/spike/internal/mcp"
)

// ---------- room protocol types (mirror comment-model.md) ----------

type Actor struct {
	Kind        string `json:"kind"`
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Runtime     string `json:"runtime,omitempty"`
	OwnerID     string `json:"owner_id,omitempty"`
	Verified    bool   `json:"verified"`
}

type Anchor struct {
	Kind      string `json:"kind"`
	FilePath  string `json:"file_path,omitempty"`
	CommitSHA string `json:"commit_sha,omitempty"`
	LineStart int    `json:"line_start,omitempty"`
	LineEnd   *int   `json:"line_end,omitempty"`
}

type Comment struct {
	ID         string `json:"id"`
	ReviewID   string `json:"review_id"`
	ParentID   *string `json:"parent_id"`
	Author     Actor   `json:"author"`
	Anchor     Anchor `json:"anchor"`
	Body       string `json:"body"`
	AddressedTo *string `json:"addressed_to,omitempty"`
	Severity   *string `json:"severity,omitempty"`
	CreatedAt  string `json:"created_at"`
}

type DiffLine struct {
	N    int    `json:"n"`
	OldN int    `json:"old_n"`
	Type string `json:"type"`
	Text string `json:"text"`
}

type DiffFile struct {
	Path  string     `json:"path"`
	Lines []DiffLine `json:"lines"`
}

type Review struct {
	ID    string     `json:"id"`
	Title string     `json:"title"`
	SHA   string     `json:"sha"`
	Files []DiffFile `json:"files"`
}

type Event struct {
	Type      string   `json:"type"`
	Comment   *Comment `json:"comment,omitempty"`
	Verdict   *Verdict `json:"verdict,omitempty"`
	Actor     *Actor   `json:"actor,omitempty"`
	Question  string   `json:"question,omitempty"`
	CommentID string   `json:"comment_id,omitempty"`
	TS        string   `json:"ts"`

	// state payload fields
	Review   Review    `json:"review,omitempty"`
	Comments []Comment `json:"comments,omitempty"`
	Verdicts []Verdict `json:"verdicts,omitempty"`
	Actors   []Actor   `json:"actors,omitempty"`
}

type Verdict struct {
	Type     string `json:"type"`
	By       Actor  `json:"by"`
	At       string `json:"at"`
	Advisory bool   `json:"advisory"`
}

// ---------- bridge ----------

type bridge struct {
	base   string // roomd base URL, e.g. http://localhost:8080
	room   string
	token  string // agent token (device-flow stand-in)
	client *http.Client

	mu       sync.Mutex
	review   Review
	comments []Comment
	actors   []Actor
	pending  map[string]bool // comment ids awaiting reply
}

func newBridge(base, room, token string) *bridge {
	return &bridge{
		base:    strings.TrimRight(base, "/"),
		room:    room,
		token:   token,
		client:  &http.Client{Timeout: 30 * time.Second},
		pending: map[string]bool{},
	}
}

// actions are the room REST calls the MCP tools wrap. The governance
// constraints execute server-side at this boundary.
func (b *bridge) action(path string, body any) (map[string]any, error) {
	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", b.base+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ocre-Token", b.token)
	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	out := map[string]any{}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("bad response (%d): %w", res.StatusCode, err)
	}
	if res.StatusCode >= 400 {
		if msg, _ := out["error"].(string); msg != "" {
			return out, fmt.Errorf("%s", msg)
		}
		return out, fmt.Errorf("HTTP %d", res.StatusCode)
	}
	return out, nil
}

func (b *bridge) refreshState() error {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/rooms/%s/review", b.base, b.room), nil)
	res, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var state struct {
		Review   Review    `json:"review"`
		Comments []Comment `json:"comments"`
		Verdicts []Verdict `json:"verdicts"`
		Actors   []Actor   `json:"actors"`
	}
	if err := json.NewDecoder(res.Body).Decode(&state); err != nil {
		return err
	}
	b.mu.Lock()
	b.review = state.Review
	b.comments = state.Comments
	b.actors = state.Actors
	b.mu.Unlock()
	return nil
}

// listen connects the outbound WebSocket and dispatches events.
func (b *bridge) listen(ctx context.Context) error {
	for {
		if err := b.listenOnce(ctx); err != nil {
			log.Printf("ws: %v; reconnecting in 2s", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}

func (b *bridge) listenOnce(ctx context.Context) error {
	url := fmt.Sprintf("ws://%s/api/rooms/%s/ws?t=agent&name=review-agent&runtime=claude_code&token=%s",
		strings.TrimPrefix(b.base, "http://"), b.room, b.token)
	// use coder/websocket via the shared helper
	conn, err := dial(url)
	if err != nil {
		return err
	}
	defer conn.CloseNow()
	log.Printf("connected to room %s", b.room)
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return err
		}
		var ev Event
		if json.Unmarshal(data, &ev) != nil {
			continue
		}
		switch ev.Type {
		case "state":
			b.mu.Lock()
			b.review = ev.Review
			b.comments = ev.Comments
			b.actors = ev.Actors
			b.mu.Unlock()
		case "comment.created", "comment.resolved", "verdict.recorded", "presence.joined", "presence.left":
			_ = b.refreshState() // keep local mirror fresh
		case "agent.mentioned":
			// A human directed a comment at the agent: answer it by
			// driving the runtime. This is the interactive loop.
			log.Printf("mentioned by %s: %q", ev.Actor.DisplayName, ev.Question)
			go b.answer(ctx, ev)
		}
	}
}

// answer drives the owner's runtime with the review context plus the
// question, over the local stdio MCP server.
func (b *bridge) answer(ctx context.Context, ev Event) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("answer panic: %v", r)
		}
	}()

	b.refreshState()
	b.mu.Lock()
	snapshot := reviewSnapshot(b.review, b.comments, b.actors)
	b.mu.Unlock()

	// MCP server over stdio: the runtime's tool calls land in our tools.
	srv := mcp.NewServer()
	registerTools(srv, b)

	// Run the MCP server on pipes; the runtime connects to it as an MCP
	// client... in the spike, the runtime is driven directly with the
	// review context inline (the MCP server exists to prove the surface;
	// full ACP/MCP wiring is the next iteration).
	prompt := fmt.Sprintf(`You are a code review participant in a live review session.

REVIEW: %s (%s)
DIFF:
%s

EARLIER COMMENTS:
%s

%s asks:
%s

Review the change. Use high/medium/low severity labels. Reply concisely; post at most 3 inline findings using the review actions described below, then reply to the question directly.

REVIEW ACTIONS (the room API you have through this bridge):
- comment on line: POST comment {anchor:{kind:"line",file_path,line_start},body,severity}
- reply: POST comment {parent_id, body, severity}
- verdict: approve | request_changes
In this spike, respond with your findings as plain text with [severity] markers; the bridge will relay them.`,
		snapshot.Title, snapshot.ReviewID, snapshot.Diff, snapshot.Comments,
		ev.Actor.DisplayName, ev.Question)

	out, err := driveRuntime(ctx, prompt)
	if err != nil {
		log.Printf("runtime: %v", err)
		if _, aerr := b.action(fmt.Sprintf("/api/rooms/%s/comments", b.room), map[string]any{
			"anchor":   map[string]any{"kind": "review"},
			"body":     "⚠️ agent runtime unavailable: " + err.Error(),
			"severity": "low",
		}); aerr != nil {
			log.Printf("post failure note: %v", aerr)
		}
		return
	}

	// Parse [severity] markers and post as a reply to the mention comment.
	parent := ev.CommentID
	for _, finding := range parseFindings(out) {
		body := finding
		sev := "medium"
		for _, s := range []string{"high", "medium", "low"} {
			if strings.HasPrefix(finding, "["+s+"]") {
				sev = s
				body = strings.TrimSpace(strings.TrimPrefix(finding, "["+s+"]"))
			}
		}
		payload := map[string]any{
			"anchor":   map[string]any{"kind": "review"},
			"body":     body,
			"severity": sev,
		}
		if parent != "" {
			payload["parent_id"] = parent
		}
		res, err := b.action(fmt.Sprintf("/api/rooms/%s/comments", b.room), payload)
		if err != nil {
			log.Printf("post finding: %v", err)
			continue
		}
		if id, _ := res["id"].(string); id != "" {
			parent = id // subsequent findings thread onto the first
		}
	}
	if len(parseFindings(out)) == 0 && strings.TrimSpace(out) != "" {
		payload := map[string]any{"anchor": map[string]any{"kind": "review"}, "body": out, "severity": "medium"}
		if parent != "" {
			payload["parent_id"] = parent
		}
		if _, err := b.action(fmt.Sprintf("/api/rooms/%s/comments", b.room), payload); err != nil {
			log.Printf("post reply: %v", err)
		}
	}
}

type snapshotT struct {
	Title, ReviewID, Diff, Comments string
}
func parseFindings(out string) []string {
	var findings []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") || strings.Contains(line, "severity:") {
			findings = append(findings, line)
		}
	}
	if len(findings) == 0 {
		// treat non-empty output as a single finding-free reply
		return nil
	}
	return findings
}

func reviewSnapshot(r Review, comments []Comment, actors []Actor) snapshotT {
	var diff strings.Builder
	for _, f := range r.Files {
		fmt.Fprintf(&diff, "--- %s ---\n", f.Path)
		for _, l := range f.Lines {
			switch l.Type {
			case "add":
				fmt.Fprintf(&diff, "+%s\n", l.Text)
			case "del":
				fmt.Fprintf(&diff, "-%s\n", l.Text)
			default:
				fmt.Fprintf(&diff, " %s\n", l.Text)
			}
		}
	}
	var cs strings.Builder
	for _, c := range comments {
		sev := ""
		if c.Severity != nil {
			sev = " [" + *c.Severity + "]"
		}
		anchor := c.Anchor.Kind
		if c.Anchor.Kind != "review" {
			anchor = fmt.Sprintf("%s:%d", c.Anchor.FilePath, c.Anchor.LineStart)
		}
		fmt.Fprintf(&cs, "- %s (%s%s): %s\n", c.Author.DisplayName, anchor, sev, strings.ReplaceAll(c.Body, "\n", " "))
	}
	return snapshotT{r.Title, r.ID, diff.String(), cs.String()}
}

// driveRuntime invokes the owner's Claude Code headlessly with the
// owner's existing credentials. The bridge passes through the user's
// environment (their subscription login or their corporate gateway
// config) but strips vars that would override the owner's auth with
// API-key billing on a different account (per the auth research: a
// stray ANTHROPIC_API_KEY silently overrides subscription auth).
func driveRuntime(ctx context.Context, prompt string) (string, error) {
	if _, err := exec.LookPath("claude"); err != nil {
		return "", fmt.Errorf("claude CLI not found on PATH")
	}
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "claude", "-p", prompt, "--output-format", "text")
	env := os.Environ()
	filtered := env[:0]
	for _, kv := range env {
		switch {
		case strings.HasPrefix(kv, "ANTHROPIC_API_KEY="),
			strings.HasPrefix(kv, "BASE_ANTHROPIC_API_KEY="),
			strings.HasPrefix(kv, "CLAUDE_CODE_USE_FOUNDRY="),
			strings.HasPrefix(kv, "ANTHROPIC_FOUNDRY_API_KEY="),
			strings.HasPrefix(kv, "ANTHROPIC_FOUNDRY_RESOURCE="):
			// stripped: never let ambient keys override the owner's auth
		default:
			filtered = append(filtered, kv)
		}
	}
	cmd.Env = filtered
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("claude -p failed: %s", msg)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// registerTools exposes the review actions as MCP tools over stdio -
// the surface the runtime would call. In the spike the server is
// instantiated and a tool listing is logged to prove the wiring; the
// runtime path above posts directly through the same action API.
func registerTools(srv *mcp.Server, b *bridge) {
	srv.Tool("read_diff", "Read the review diff", mcp.Obj(map[string]any{}), func(args map[string]any) (any, error) {
		b.mu.Lock()
		defer b.mu.Unlock()
		return reviewSnapshot(b.review, b.comments, b.actors), nil
	})
	srv.Tool("post_comment", "Post a comment (line/block/file/review anchor)", mcp.Obj(map[string]any{
		"kind":       "string",
		"file_path":  "string?",
		"line_start": "int?",
		"body":       "string",
		"severity":   "string?",
	}), func(args map[string]any) (any, error) {
		return b.action(fmt.Sprintf("/api/rooms/%s/comments", b.room), args)
	})
	srv.Tool("reply", "Reply in a thread", mcp.Obj(map[string]any{
		"parent_id": "string",
		"body":      "string",
		"severity":  "string?",
	}), func(args map[string]any) (any, error) {
		return b.action(fmt.Sprintf("/api/rooms/%s/comments", b.room), args)
	})
	srv.Tool("approve", "Record advisory approval", mcp.Obj(map[string]any{}), func(args map[string]any) (any, error) {
		return b.action(fmt.Sprintf("/api/rooms/%s/verdicts", b.room), map[string]any{"type": "approve"})
	})
	srv.Tool("request_changes", "Record advisory request-changes", mcp.Obj(map[string]any{}), func(args map[string]any) (any, error) {
		return b.action(fmt.Sprintf("/api/rooms/%s/verdicts", b.room), map[string]any{"type": "request_changes"})
	})
}

func main() {
	base := flag.String("base", "http://localhost:8080", "roomd base URL")
	room := flag.String("room", "", "room id")
	token := flag.String("token", "", "agent token")
	mcpStdio := flag.Bool("mcp", false, "run the MCP stdio server instead of the bridge loop (tool-list mode)")
	flag.Parse()

	if *mcpStdio {
		// Tool-listing mode: prove the MCP surface over stdio.
		b := newBridge(*base, *room, *token)
		srv := mcp.NewServer()
		registerTools(srv, b)
		srv.ServeStdio(os.Stdin, os.Stdout)
		return
	}

	if *room == "" || *token == "" {
		log.Fatal("usage: ocre-bridge -base URL -room ID -token AGENT_TOKEN [-mcp]")
	}

	b := newBridge(*base, *room, *token)
	ctx := context.Background()

	// Demonstrate the MCP surface at startup (tool list proves the
	// review-actions-as-tools wiring).
	srv := mcp.NewServer()
	registerTools(srv, b)
	log.Printf("MCP tools ready: %v", srv.ToolNames())

	if err := b.refreshState(); err != nil {
		log.Printf("initial state fetch: %v (will rely on ws)", err)
	}
	if err := b.listen(ctx); err != nil && ctx.Err() == nil {
		log.Fatalf("listen: %v", err)
	}
}

// dial is a thin wrapper so main doesn't import websocket twice.
func dial(url string) (*conn, error) {
	c, err := mcpDial(url)
	return c, err
}

// conn aliases the websocket connection type.
type conn = mcp.Conn

// mcpDial avoids importing coder/websocket in two files of this package.
func mcpDial(url string) (*conn, error) {
	return mcp.Dial(url)
}

var _ = bufio.NewReader // keep import if unused paths change
var _ io.Reader = (*os.File)(nil)