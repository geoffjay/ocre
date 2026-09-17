---
type: Concept
title: Core-loop spike report
description: What the core-loop prototype (roomd + ocre-bridge) proved, what it verified, and its known limitations - the plan's step 2/4 validation.
tags: [spike, prototype, verification, bridge, core-loop]
status: stable
generated: { by: omp-agent/glm-5.3, at: 2026-09-17T16:23:37Z }
verified: { by: omp-agent/glm-5.3, at: 2026-09-17T16:23:37Z }
sources:
  - id: comment-model
    title: ocre KB - Comment model (the wire schema the spike implements)
    resource: /concepts/design/comment-model.md
  - id: agent-governance
    title: ocre KB - Agent governance constraints (enforced server-side in the spike)
    resource: /concepts/design/agent-governance.md
  - id: agent-participation
    title: ocre KB - Agent participation model (the composition the spike validates)
    resource: /concepts/research/agent-participation.md
---

# Core-loop spike report

The prototype (`spike/` in the repo) implements the plan's step 2 (core
loop) and step 4 (Go bridge spike) in minimal form. Everything below was
**executed and verified on 2026-09-17**, not designed on paper.

## What was built

* **`spike/cmd/roomd`** — Go server: rooms over WebSocket (coder/websocket)
  + REST action API, embedding a single-page web client
  (`spike/cmd/roomd/web/index.html`). Implements the comment-model wire
  contract (Comment/Anchor/Verdict/Event JSON), governance checks at the
  action boundary, presence, and the `agent.mentioned` event.
* **`spike/cmd/ocre-bridge`** — the desktop agent bridge: one outbound
  WebSocket to the room (reconnect loop), listens for mentions, drives the
  owner's `claude -p` headlessly with review context, parses `[severity]`
  findings, posts them back as threaded replies; exposes the review verbs
  as MCP tools over stdio (`spike/internal/mcp`), protocol 2025-06-18.

## Verified end-to-end (2026-09-17)

1. **Create-review → share-link → join**: `POST /api/reviews` returns the
   share link; opening it in a browser joins as a live human participant.
2. **Live UI**: diff renders per file with +/- lines; clicking a line
   opens an inline comment form; presence renders humans + `🤖 agent`
   correctly after the idempotent-join fix.
3. **Human → agent loop (the core interactive path)**: human comments on
   the buggy line, addressed-to the agent → room broadcasts
   `agent.mentioned` → bridge invokes real `claude -p` → agent posts a
   threaded reply into the room → **the reply renders in every open
   browser without reload**. The planted bug
   (`for i := 1; i <= len(nums); i++` → index-out-of-range) was correctly
   found both runs: `[high]` off-by-one panic finding + `[low]` no-test and
   revert-to-range findings + `request_changes` verdict.
4. **Runtime auth**: headless `claude -p` worked only after the bridge
   strips ambient auth-override vars (`ANTHROPIC_API_KEY`,
   `CLAUDE_CODE_USE_FOUNDRY`, `ANTHROPIC_FOUNDRY_*`) and lets the owner's
   configured auth (subscription OAuth or, on this workstation, the
   corporate gateway via `apiKeyHelper`/settings env) act — confirming the
   auth research's "stray key silently overrides" caveat is a real
   implementation trap, handled in `driveRuntime`.
5. **Governance at the boundary**: agent verdicts recorded with
   `advisory: true` (rendered "(advisory)" in UI); agent comment without
   severity → HTTP 400; 6th agent comment inside the 10s window → HTTP 429.
   Enforcement is server-side; the bridge cannot bypass it.
6. **MCP stdio surface**: `initialize` handshake (protocol 2025-06-18,
   `ocre-bridge v0.1.0`), `tools/list` returns the five comment-model
   verbs (`read_diff`, `post_comment`, `reply`, `approve`,
   `request_changes`).

## Known limitations (spike scope, recorded deliberately)

* No persistence: rooms are in-memory; a `roomd` restart loses state, and
  WS events are not replayed on reconnect (the mention must be re-sent).
  The comment model's commit-sha anchor stability is therefore untested.
* No real auth: query-token stand-ins for the OAuth device flow; the agent
  token is the room's single agent identity (no per-owner agents yet).
* The runtime is driven by prompt, not by MCP tool calls: the agent
  receives the review context inline and the bridge posts its findings.
  The MCP server exists and speaks the protocol, but `claude -p` does not
  call the tools yet — the next iteration wires tool-use
  (`--mcp-config`) so the agent posts findings itself via `post_comment`.
* Presence is join/leave only (no cursors); verdicts have no merge-gate
  semantics; no cross-PR references in the spike.

## Answers this spike gives the open questions

* **Bridge viability**: confirmed — the whole loop runs with ~1.1k lines
  of Go, one outbound WS, zero inbound endpoints.
* **Subprocess-per-task shape**: `claude -p` per mention worked fine
  (3–40s per review); latency acceptable for interactive review. No
  persistent-session complexity needed yet.
* **Workstation-gateway caveat**: headless auth requires env sanitization
  (above) — a real deployment consideration for any user with corporate
  gateway config, now handled and documented in the code.
* **Gemini headless auth**: still untested (no gemini CLI on this
  workstation) — remains open.

[^comment-model]: ocre KB, comment model.
[^agent-governance]: ocre KB, agent governance constraints.
[^agent-participation]: ocre KB, agent participation model.