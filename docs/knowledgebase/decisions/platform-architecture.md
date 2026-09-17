---
type: Decision
title: Platform architecture - web for humans, Go agent bridge, no GUI shell
description: The founding platform decision - web platform for human participants, a single-Go-binary agent bridge on the agent owner's desktop, GUI shells deferred.
tags: [decision, platform, architecture, web, bridge, go]
status: stable
generated: { by: omp-agent/glm-5.3, at: 2026-09-17T04:56:21Z }
verified: { by: human:geoff, at: 2026-09-17T04:52:04Z }
sources:
  - id: platform-tradeoffs
    title: ocre KB - Desktop vs web platform tradeoffs (evidence base)
    resource: /concepts/research/platform-tradeoffs.md
  - id: agent-participation
    title: ocre KB - Agent participation model (bridge composition)
    resource: /concepts/research/agent-participation.md
  - id: collab-prior-art
    title: ocre KB - Collaborative review prior art (room topologies)
    resource: /concepts/research/collaborative-review-prior-art.md
  - id: figma-multiplayer
    title: Figma - How Figma's multiplayer technology works
    resource: https://www.figma.com/blog/how-figmas-multiplayer-technology-works/
  - id: partykit-docs
    title: PartyKit - How PartyKit works (same-id room routing)
    resource: https://docs.partykit.io/how-partykit-works/
  - id: slack-socket-mode
    title: Slack - Using Socket Mode (outbound bot connection)
    resource: https://docs.slack.dev/apis/events-api/using-socket-mode/
---

# Platform architecture - web for humans, Go agent bridge, no GUI shell

**Decision date**: 2026-09-17. **Status**: decided; founder-verified.

## Decision

1. **Web platform for human participants.** Humans join live review
   sessions in the browser via share links. No desktop application is
   required to be a contributor.
2. **Agents participate via a desktop agent bridge**: a single static Go
   binary (CLI + daemon in one executable) running on the agent owner's
   machine. It connects out to the review rooms over one long-lived
   WebSocket per bridge (Socket Mode / Discord Gateway pattern),
   authenticates via OAuth device flow, and drives the owner's own agent
   runtimes (Claude Code, Codex CLI, Gemini CLI, opencode, any ACP agent)
   with the owner's own subscription credentials, exposing review actions
   to the runtime as a local stdio MCP server.
3. **No GUI shell for v1.** The bridge is daemon-first (watchman /
   dockerd precedent). Wails v3 / Tauri sidecar / Electron are recorded
   escape hatches if a human-facing desktop app is ever needed — the
   wrap-later path stays open.
4. **ocre does not build or host agents and does not bill for AI
   tokens** (founder constraints, 2026-09-17). Model inference happens
   entirely inside the owner's runtime on their credentials.

## Evidence (why this and not the alternatives)

* **Web owns the adoption loop.** Every product with ocre's
   share-a-link → instant-live-contributor shape (Figma, Google Docs,
   Miro, CodeSandbox) is browser-first; no counterexample exists.
   Desktop-first has no observed precedent for adding zero-install link
   sharing later. See [platform tradeoffs](/concepts/research/platform-tradeoffs.md).
* **Desktop multiplayer is bespoke; web multiplayer is commodity.**
   Figma's per-document server + property-LWW sync, PartyKit's same-id
   Durable-Object rooms, and the Yjs ecosystem mean live presence and
   shared state in a browser is a weeks-scale integration; Zed spent
   years building the native equivalent. Founder explicitly weighed in
   agreement on 2026-09-17.
* **Large-diff rendering is a solved-in-practice web constraint**
   (GitHub's virtualized Files Changed; Figma's GPU rendering) — it
   costs deliberate engineering, not a platform change.
* **The bridge satisfies every founder constraint with precedented
   parts** (outbound WS, device flow, headless/ACP runtime driving,
   local MCP review tools): see
   [agent participation model](/concepts/research/agent-participation.md).
   The Liveblocks gap — no platform runs agents as desktop-resident live
   participants — is ocre's differentiation.

## Consequences

* Server-side review state (comments, verdicts, session history) is the
  **provisional storage choice for v1**; the Live Share relay-only model
  and Matrix-style federation are recorded escape hatches, not chosen
  paths. Revisit if enterprise data-residency demands arise.
* Two implementation surfaces to maintain: web app + bridge. Both are
  thin over a shared room protocol — the bridge consumes the same
  event-stream + action API the web client uses.
* ocre's server carries only review payloads (JSON), never model I/O;
  token-cost exposure to ocre is zero by construction.
* Cross-platform distribution cost for the bridge = macOS
  signing+notarization + Windows signing (stack-independent tax).

## What this decision deliberately does not decide

* Room protocol details (Yjs vs plain JSON events vs other) — decided by
  the core-loop prototype and comment-model work.
* Forge integration depth (read-only overlay vs write-back).
* Agent identity/pacing specifics — designed with the volume controls.

Evidence base: [platform tradeoffs](/concepts/research/platform-tradeoffs.md),
[agent participation model](/concepts/research/agent-participation.md),
[collaborative review prior art](/concepts/research/collaborative-review-prior-art.md).

[^platform-tradeoffs]: ocre KB, platform tradeoffs (assembled evidence).
[^agent-participation]: ocre KB, agent participation model.
[^collab-prior-art]: ocre KB, collaborative review prior art.
[^figma-multiplayer]: Figma engineering blog.
[^partykit-docs]: PartyKit docs.
[^slack-socket-mode]: Slack developer docs.