// roomd is the ocre core-loop spike server: rooms over WebSocket with a
// REST action API, enforcing the agent governance constraints from
// docs/knowledgebase/concepts/design/agent-governance.md.
//
// Wire contract per docs/knowledgebase/concepts/design/comment-model.md:
// WS events out (comment.created, verdict.recorded, presence.*, agent.mentioned),
// REST actions in. Comments are discrete append-mostly objects; verdicts
// are review-level facets. Agent actors authenticate with the room's
// agent token (stand-in for the OAuth device flow per the decision doc).
package main

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
)

//go:embed web/index.html
var webFS embed.FS

// ---------- model (mirrors concepts/design/comment-model.md) ----------

type Actor struct {
	Kind        string `json:"kind"` // human | agent
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Runtime     string `json:"runtime,omitempty"`  // agent only
	OwnerID     string `json:"owner_id,omitempty"` // agent only
	Verified    bool   `json:"verified"`
}

type Anchor struct {
	Kind      string `json:"kind"` // line | block | file | review
	FilePath  string `json:"file_path,omitempty"`
	CommitSHA string `json:"commit_sha,omitempty"`
	LineStart int    `json:"line_start,omitempty"`
	LineEnd   *int   `json:"line_end,omitempty"` // block anchors; null for single-line
}

type Comment struct {
	ID               string  `json:"id"`
	ReviewID         string  `json:"review_id"`
	ParentID         *string `json:"parent_id"`
	Author           Actor   `json:"author"`
	Anchor           Anchor  `json:"anchor"`
	Body             string  `json:"body"`
	AddressedTo      *string `json:"addressed_to,omitempty"`
	Severity         *string `json:"severity,omitempty"` // required for agents
	CreatedAt        string  `json:"created_at"`
	ResolvedBy       *string `json:"resolved_by,omitempty"`
	ResolvedAt       *string `json:"resolved_at,omitempty"`
	ResolutionReason *string `json:"resolution_reason,omitempty"`
}

type Verdict struct {
	Type string `json:"type"` // approve | request_changes
	By   Actor  `json:"by"`
	At   string `json:"at"`
	// Advisory is the governance marker: agent verdicts never gate merges.
	Advisory bool `json:"advisory"`
}

type DiffLine struct {
	N    int    `json:"n"` // new-file line number (0 for pure deletions)
	OldN int    `json:"old_n"`
	Type string `json:"type"` // ctx | add | del
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
}

// ---------- room ----------

type client struct {
	conn  *websocket.Conn
	actor Actor
	token string
}

type room struct {
	mu         sync.Mutex
	id         string
	agentToken string
	review     Review
	byToken    map[string]*Actor // join/agent token -> actor
	clients    map[*client]bool
	comments   []*Comment
	verdicts   []Verdict
	agentRates map[string][]time.Time // agent actor id -> recent post times
}

var (
	roomsMu sync.Mutex
	rooms   = map[string]*room{}
)

func newID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// sampleReview carries a subtle real bug so an agent reviewer has
// something to find: Sum indexes nums[1..len] inclusive → out of range.
func sampleReview() Review {
	sumNew := []string{
		"package calc",
		"",
		"// Sum returns the total of nums.",
		"func Sum(nums []int) int {",
		"\ttotal := 0",
		"\tfor i := 1; i <= len(nums); i++ {",
		"\t\ttotal += nums[i]",
		"\t}",
		"\treturn total",
		"}",
	}
	lines := make([]DiffLine, 0, len(sumNew)+2)
	lines = append(lines,
		DiffLine{OldN: 6, Type: "del", Text: "\tfor _, n := range nums {"},
		DiffLine{OldN: 7, Type: "del", Text: "\t\ttotal += n"},
	)
	for i, t := range sumNew {
		lines = append(lines, DiffLine{N: i + 1, OldN: i + 1, Type: "add", Text: t})
	}
	return Review{
		ID:    "rev-" + newID()[:6],
		Title: "calc: switch Sum to indexed loop",
		SHA:   "spike0123",
		Files: []DiffFile{
			{Path: "calc/sum.go", Lines: lines},
			{Path: "calc/README.md", Lines: []DiffLine{
				{N: 1, Type: "add", Text: "# calc"},
				{N: 2, Type: "add", Text: ""},
				{N: 3, Type: "add", Text: "Tiny arithmetic helpers."},
			}},
		},
	}
}

// join registers an actor and mints its token. Human joins are
// idempotent by display name: reconnects (and the same person in two
// tabs) map to one actor, so presence doesn't duplicate.
func (r *room) join(kind, name, runtime string) (Actor, string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if kind == "agent" {
		if a, ok := r.byToken[r.agentToken]; ok {
			return *a, r.agentToken
		}
		a := &Actor{Kind: "agent", ID: "agent", DisplayName: name, Runtime: runtime, OwnerID: "owner"}
		r.byToken[r.agentToken] = a
		return *a, r.agentToken
	}
	for tok, a := range r.byToken {
		if a.Kind == "human" && a.DisplayName == name {
			return *a, tok
		}
	}
	tok := newID()
	id := "h-" + tok[:6]
	a := &Actor{Kind: "human", ID: id, DisplayName: name}
	r.byToken[tok] = a
	return *a, tok
}


func newRoom() *room {
	return &room{
		id:         newID()[:10],
		agentToken: newID(),
		review:     sampleReview(),
		byToken:    map[string]*Actor{},
		clients:    map[*client]bool{},
		comments:   []*Comment{},
		agentRates: map[string][]time.Time{},
	}
}


func (r *room) broadcast(ev Event) {
	data, _ := json.Marshal(ev)
	r.mu.Lock()
	conns := make([]*client, 0, len(r.clients))
	for c := range r.clients {
		conns = append(conns, c)
	}
	r.mu.Unlock()
	for _, c := range conns {
		go func(c *client) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			c.conn.Write(ctx, websocket.MessageText, data)
		}(c)
	}
}

func (r *room) actorSlice() []Actor {
	out := make([]Actor, 0, len(r.byToken))
	seen := map[string]bool{}
	for _, a := range r.byToken {
		if !seen[a.ID] {
			seen[a.ID] = true
			out = append(out, *a)
		}
	}
	return out
}


// ---------- governance (concepts/design/agent-governance.md) ----------

const (
	agentRateWindow = 10 * time.Second
	agentRateMax    = 5 // agent comments per window (spike default)
)

func (r *room) governanceCheck(c *Comment) (int, string) {
	if c.Author.Kind != "agent" {
		return 0, ""
	}
	// #2: severity is mandatory on agent findings.
	if c.Severity == nil || (*c.Severity != "high" && *c.Severity != "medium" && *c.Severity != "low") {
		return http.StatusBadRequest, `{"error":"agent comments require severity: high|medium|low"}`
	}
	// #3: per-agent rate bucket.
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	times := r.agentRates[c.Author.ID]
	kept := times[:0]
	for _, t := range times {
		if now.Sub(t) < agentRateWindow {
			kept = append(kept, t)
		}
	}
	if len(kept) >= agentRateMax {
		r.agentRates[c.Author.ID] = kept
		return http.StatusTooManyRequests,
			fmt.Sprintf(`{"error":"rate limit: max %d agent comments per %s"}`, agentRateMax, agentRateWindow)
	}
	r.agentRates[c.Author.ID] = append(kept, now)
	return 0, ""
}

// ---------- server ----------

func main() {
	port := flag.Int("port", 8080, "listen port")
	flag.Parse()

	mux := http.NewServeMux()

	// Create review → share link (humans) + agent token (bridge).
	mux.HandleFunc("POST /api/reviews", func(w http.ResponseWriter, req *http.Request) {
		roomsMu.Lock()
		r := newRoom()
		rooms[r.id] = r
		roomsMu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"review_id":   r.review.ID,
			"room":        r.id,
			"share_link":  fmt.Sprintf("http://localhost:%d/r/%s", *port, r.id),
			"agent_token": r.agentToken,
		})
	})

	mux.HandleFunc("GET /r/{room}", func(w http.ResponseWriter, req *http.Request) {
		if getRoom(req) == nil {
			http.NotFound(w, req)
			return
		}
		data, _ := webFS.ReadFile("web/index.html")
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	})

	mux.HandleFunc("GET /api/rooms/{room}/review", func(w http.ResponseWriter, req *http.Request) {
		r := getRoom(req)
		if r == nil {
			http.NotFound(w, req)
			return
		}
		r.mu.Lock()
		defer r.mu.Unlock()
		json.NewEncoder(w).Encode(map[string]any{
			"review":   r.review,
			"comments": r.comments,
			"verdicts": r.verdicts,
			"actors":   r.actorSlice(),
		})
	})

	// WebSocket: join + live events. Query token stands in for the OAuth
	// device flow (decision doc); agents must hold the room's agent token.
	mux.HandleFunc("GET /api/rooms/{room}/ws", func(w http.ResponseWriter, req *http.Request) {
		r := getRoom(req)
		if r == nil {
			http.NotFound(w, req)
			return
		}
		q := req.URL.Query()
		kind, name, runtime := q.Get("t"), q.Get("name"), q.Get("runtime")
		if kind == "agent" && q.Get("token") != r.agentToken {
			http.Error(w, "invalid agent token", http.StatusUnauthorized)
			return
		}
		if name == "" {
			name = "Anonymous"
		}
		actor, tok := r.join(kind, name, runtime)
		conn, err := websocket.Accept(w, req, nil)
		if err != nil {
			return
		}
		c := &client{conn: conn, actor: actor, token: tok}
		r.mu.Lock()
		r.clients[c] = true
		state, _ := json.Marshal(map[string]any{
			"type":     "state",
			"review":   r.review,
			"comments": r.comments,
			"verdicts": r.verdicts,
			"actors":   r.actorSlice(),
			"you":      actor,
			"token":    tok,
		})
		r.mu.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		conn.Write(ctx, websocket.MessageText, state)
		cancel()
		r.broadcast(Event{Type: "presence.joined", Actor: &actor, TS: now()})

		go func() {
			defer func() {
				r.mu.Lock()
				delete(r.clients, c)
				r.mu.Unlock()
				r.broadcast(Event{Type: "presence.left", Actor: &actor, TS: now()})
				conn.CloseNow()
			}()
			for {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				_, _, err := conn.Read(ctx)
				cancel()
				if err != nil {
					return
				}
			}
		}()
	})

	// REST actions. authActor resolves the participant from its token.
	mux.HandleFunc("POST /api/rooms/{room}/comments", func(w http.ResponseWriter, req *http.Request) {
		r := getRoom(req)
		if r == nil {
			http.NotFound(w, req)
			return
		}
		actor := authActor(r, req)
		if actor == nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		var in struct {
			Anchor      Anchor  `json:"anchor"`
			Body        string  `json:"body"`
			ParentID    *string `json:"parent_id"`
			AddressedTo *string `json:"addressed_to"`
			Severity    *string `json:"severity"`
		}
		if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
			http.Error(w, `{"error":"bad json"}`, http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(in.Body) == "" {
			http.Error(w, `{"error":"body required"}`, http.StatusBadRequest)
			return
		}
		if in.Anchor.Kind == "" {
			in.Anchor.Kind = "review"
		}
		if in.Anchor.Kind != "review" {
			if in.Anchor.CommitSHA == "" {
				in.Anchor.CommitSHA = r.review.SHA
			}
		}
		c := &Comment{
			ID: newID(), ReviewID: r.review.ID, ParentID: in.ParentID,
			Author: *actor, Anchor: in.Anchor, Body: in.Body,
			AddressedTo: in.AddressedTo, Severity: in.Severity,
			CreatedAt: now(),
		}
		if code, msg := r.governanceCheck(c); code != 0 {
			http.Error(w, msg, code)
			return
		}
		r.mu.Lock()
		r.comments = append(r.comments, c)
		r.mu.Unlock()
		r.broadcast(Event{Type: "comment.created", Comment: c, TS: now()})

		// Comment directed at an agent → the mention event the bridge
		// listens for (the human↔agent interaction loop).
		if in.AddressedTo != nil {
			r.mu.Lock()
			var target *Actor
			for _, a := range r.byToken {
				if a.ID == *in.AddressedTo {
					target = a
				}
			}
			r.mu.Unlock()
			if target != nil && target.Kind == "agent" {
				r.broadcast(Event{Type: "agent.mentioned", Actor: actor, Question: in.Body, CommentID: c.ID, TS: now()})
			}
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(c)
	})

	mux.HandleFunc("POST /api/rooms/{room}/comments/{cid}/resolve", func(w http.ResponseWriter, req *http.Request) {
		r := getRoom(req)
		if r == nil {
			http.NotFound(w, req)
			return
		}
		actor := authActor(r, req)
		if actor == nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		var in struct {
			Reason string `json:"reason"`
		}
		json.NewDecoder(req.Body).Decode(&in)
		if in.Reason == "" {
			in.Reason = "addressed"
		}
		r.mu.Lock()
		var found *Comment
		for _, c := range r.comments {
			if c.ID == req.PathValue("cid") {
				found = c
			}
		}
		if found != nil {
			found.ResolvedBy = new(actor.ID)
			found.ResolvedAt = new(now())
			found.ResolutionReason = new(in.Reason)
		}
		r.mu.Unlock()
		if found == nil {
			http.NotFound(w, req)
			return
		}
		r.broadcast(Event{Type: "comment.resolved", Comment: found, TS: now()})
		json.NewEncoder(w).Encode(found)
	})

	mux.HandleFunc("POST /api/rooms/{room}/verdicts", func(w http.ResponseWriter, req *http.Request) {
		r := getRoom(req)
		if r == nil {
			http.NotFound(w, req)
			return
		}
		actor := authActor(r, req)
		if actor == nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusBadRequest)
			return
		}
		var in struct {
			Type string `json:"type"`
		}
		if err := json.NewDecoder(req.Body).Decode(&in); err != nil ||
			(in.Type != "approve" && in.Type != "request_changes") {
			http.Error(w, `{"error":"type must be approve|request_changes"}`, http.StatusBadRequest)
			return
		}
		// Governance #1: agent verdicts are advisory, never gate merges.
		v := Verdict{Type: in.Type, By: *actor, At: now(), Advisory: actor.Kind == "agent"}
		r.mu.Lock()
		r.verdicts = append(r.verdicts, v)
		r.mu.Unlock()
		r.broadcast(Event{Type: "verdict.recorded", Verdict: &v, TS: now()})
		json.NewEncoder(w).Encode(v)
	})

	log.Printf("roomd listening on :%d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func getRoom(req *http.Request) *room {
	roomsMu.Lock()
	defer roomsMu.Unlock()
	return rooms[req.PathValue("room")]
}

// authActor resolves the acting participant from its token: humans use
// the join token minted at WS connect; the bridge uses the agent token.
func authActor(r *room, req *http.Request) *Actor {
	tok := req.Header.Get("X-ocre-Token")
	if tok == "" {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if a, ok := r.byToken[tok]; ok {
		cp := *a
		return &cp
	}
	return nil
}

func now() string { return time.Now().UTC().Format(time.RFC3339) }