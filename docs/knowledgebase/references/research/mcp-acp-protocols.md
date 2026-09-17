---
type: Reference
title: MCP 2025-06-18 and ACP (protocol layer for the bridge)
description: The two protocols the agent bridge composes - MCP for exposing review actions as tools, ACP for driving agents uniformly.
tags:
  - reference
  - research
  - agents
  - mcp
  - acp
  - protocols
generated:
  by: omp-agent/glm-5.3
  at: "2026-09-17T03:36:04Z"
verified:
  by: omp-agent/glm-5.3
  at: "2026-09-17T03:36:04Z"
---

# MCP 2025-06-18 and ACP

## MCP (Model Context Protocol), spec version 2025-06-18

**Resource**:
[modelcontextprotocol.io — Transports](https://modelcontextprotocol.io/specification/2025-06-18/basic/transports)
(**verified verbatim 2026-09-17**)

* Exactly two standard transports: **stdio** (client launches server as
  subprocess; newline-delimited JSON-RPC over stdin/stdout; spec says
  stdio servers retrieve credentials from the environment — no OAuth
  machinery needed) and **Streamable HTTP** (single endpoint, POST +
  optional SSE, `Mcp-Session-Id`, OAuth 2.1 for remote servers with
  Protected Resource Metadata, DCR, PKCE).
* Surface-as-tools precedents: GitHub's official MCP server (~33k stars;
  issues/PRs/review tools; hosted remote with OAuth/PAT + local stdio
  variant) and Linear's (issues/projects/comments; scoped read-only
  endpoint) — the exact template for an ocre review-actions MCP server:
  `list_reviews`, `read_diff`, `post_comment`, `reply`, `approve`,
  `request_changes`, `reference_pr`.

## ACP (Agent Client Protocol)

**Resource**:
[agentclientprotocol.com](https://agentclientprotocol.com/overview/agents)
(**agents list verified verbatim 2026-09-17**)

* JSON-RPC over stdio between a client and a coding agent; official
  schema, stable v1; Rust/TS/Python/Kotlin/Java SDKs; ACP Agent Registry
  for discovery.
* Adoption: Gemini CLI and opencode native; Cursor, Cline, Goose, Junie,
  OpenHands, and ~40 more listed; Copilot CLI in public preview
  (announced 2026-01-28); Claude Code and Codex CLI via official adapter
  shims (zed-industries/claude-agent-acp, agentclientprotocol/codex-acp).
* Precedent for non-editor clients: mobile apps (Happy, Runmote), scheduler
  bridges (Kronos), chat connectors (Telegram/Discord/Slack/Matrix), an
  HTTP stdio bridge, and a remote execution gateway are all listed ACP
  clients — a code-review bridge is an established client type.
* v2 is in draft; pin to stable v1 and monitor.

## Composition

The bridge speaks ACP (or per-runtime headless modes) toward the agent,
and exposes a **local stdio MCP server** toward the same agent's tool
calls — Gemini's ACP mode even defines the client registering its MCP
server at the initialize handshake. No protocol in this composition is
unstandardized.

Used by [agent participation model](/concepts/research/agent-participation.md).
