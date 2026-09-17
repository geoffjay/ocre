---
type: Reference
title: Agent runtime programmatic surfaces (2026)
description: How Claude Code, Codex CLI, Gemini CLI, opencode, and ACP agents can be driven by an external local process using the user's own subscription credentials.
tags:
  - reference
  - research
  - agents
  - byok
  - acp
generated:
  by: omp-agent/glm-5.3
  at: "2026-09-17T03:36:04Z"
verified:
  by: omp-agent/glm-5.3
  at: "2026-09-17T03:36:04Z"
---

# Agent runtime programmatic surfaces (2026)

The load-bearing table behind ocre's agent bridge: every major runtime is
drivable today, and the good-UX auth path is subscription reuse, not
API-key paste.

## Claude Code

* Headless: `claude -p` — stdin/arg prompt, exit codes, `--output-format
  json|stream-json`, caller-supplied `--json-schema` (result in
  `structured_output`), `--allowedTools`/`--permission-mode`,
  `--mcp-config`; Agent SDK (TS/Python) for full programmatic control.
  **Verified verbatim** at
  [code.claude.com/docs/en/headless](https://code.claude.com/docs/en/headless).
* Auth: subscription OAuth by default; `claude setup-token` mints a
  one-year browserless subscription token. **Caveats**: `--bare` mode
  never reads OAuth credentials (API key only); a stray
  `ANTHROPIC_API_KEY` env var silently overrides the subscription in
  non-interactive mode — the bridge must avoid both.
  [code.claude.com/docs/en/authentication](https://code.claude.com/docs/en/authentication)

## Codex CLI

* Headless: `codex exec` — JSONL event stream (`--json`),
  `--output-schema`, `-o/--output-last-message`, `codex exec resume`,
  read-only sandbox default; SDK (TS `@openai/codex-sdk` / Python) over
  the local app-server. [learn.chatgpt.com/docs/non-interactive-mode](https://learn.chatgpt.com/docs/non-interactive-mode)
* Auth: ChatGPT browser login default; "codex exec reuses saved CLI
  authentication by default"; `--with-access-token` for non-interactive
  enterprise workflows. [learn.chatgpt.com/docs/auth](https://learn.chatgpt.com/docs/auth)

## Gemini CLI

* Headless: `gemini -p` with `--output-format json`/JSONL, defined exit
  codes. **Native ACP**: `gemini --acp` (JSON-RPC over stdio; the client
  can register its own MCP server during initialize).
  [ACP mode docs](https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/acp-mode.md)
* Auth asymmetry (open question): docs recommend API key / Vertex for
  headless; cached Google OAuth works interactively; unattended behavior
  unverified.

## opencode

* `opencode serve`: HTTP server with OpenAPI spec — sessions, messages
  (sync/async), aborts, permission responses, diffs, provider auth
  endpoints; `opencode acp` for stdio. BYO provider keys at the model
  layer. [opencode.ai/docs/server](https://opencode.ai/docs/server/)

## ACP (Agent Client Protocol)

* JSON-RPC over stdio, stable v1 (v2 draft); SDKs in Rust/TS/Python/Kotlin/
  Java; 40+ agents listed (Gemini CLI, opencode, Cursor, Cline, Goose,
  Copilot CLI in preview natively; Claude Code and Codex via official
  adapters — zed-industries/claude-agent-acp, agentclientprotocol/codex-acp).
  **Verified verbatim** at
  [agentclientprotocol.com/overview/agents](https://agentclientprotocol.com/overview/agents)
* The clients directory lists non-editor bridges (mobile apps, schedulers,
  chat connectors, HTTP gateways) as an established client category — a
  review bridge is precedented, not novel.

## BYOK UX evidence

Negative: Cursor restricted then killed BYOK (late 2025) amid backlash;
agent features never worked with API keys; commentators call key-paste
"a conversion tax" (secondary sources — flagged). Positive: Claude Code,
Codex, Gemini CLI all default to subscription login; Claude docs warn
about accidental API-key billing — vendors treat subscription reuse as
the expected good UX.

Used by [agent participation model](/concepts/research/agent-participation.md).
