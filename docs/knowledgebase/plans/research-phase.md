---
type: Plan
title: Research phase — synthesis and next steps
description: Status of the initial research phase into code review volume, AI output overload, collaboration prior art, and platform tradeoffs, with synthesis and proposed next steps.
tags: [plan, research, phase]
status: draft
generated: { by: omp-agent/glm-5.3, at: 2026-09-16T21:21:05Z }
updated: { by: omp-agent/glm-5.3, at: 2026-09-17T03:36:04Z }
---

# Research phase — synthesis and next steps

The first goal of the ocre project (stated 2026-09-16): research code
review and the current challenges with volume and with agents generating
more output than humans can manage, recorded in this knowledge base as the
primary working location until a solid plan is defined.

## Completed 2026-09-16

Four parallel research threads, each with per-claim sources, quality notes,
and confidence ratings; load-bearing numbers verified against primary
sources where reachable:

* [Human review capacity](/concepts/research/code-review-capacity.md) —
  the size cliff, time budgets, what review is actually for.
* [AI review noise](/concepts/research/ai-review-noise.md) — measured
  agent-vs-human review differences, automation bias, mitigation patterns.
* [Collaborative review prior art](/concepts/research/collaborative-review-prior-art.md)
  — standalone tool graveyard, multiplayer topologies, the confirmed gap.
* [Platform tradeoffs](/concepts/research/platform-tradeoffs.md) — desktop
  vs web evidence, assembled but explicitly undecided.
* [Product vision](/concepts/product-vision.md) — the founding intent this
  research serves (founder-confirmed 2026-09-17).

Reference docs for load-bearing sources live under
[references/research/](/references/research/).

## Completed 2026-09-17 — agent participation

The founder confirmed the vision and the web-for-humans direction, and set
hard constraints (no built/hosted agents, no token billing, agents on the
user's desktop, BYOK must not hurt UX). Three further research threads
answered them:

* [Agent participation model](/concepts/research/agent-participation.md) —
  the composition: outbound-WebSocket bridge on the agent owner's desktop
  (Socket Mode/Gateway pattern) + device-flow auth + driving the owner's
  runtime via headless modes/ACP with subscription credentials + a local
  stdio MCP server exposing review actions. Every component is
  precedented; the desktop-resident live agent participant is the
  genuine gap ocre fills (Liveblocks does agents-via-server-REST only).
* [Platform decision status](#proposed-next-steps): web platform for
  humans, single-Go-binary bridge for agents (no GUI shell needed;
  watchman/dockerd precedent; signing tax stack-independent; Wails v3
  beta / Tauri sidecar / Electron documented as later options).

## Synthesis (updated 2026-09-17)

1. Human review capacity is bounded (≈200–400 LOC/session, ~500 LOC/hr
   ceiling, ~1–2 h/day of attention) while AI acceleration inflates PR
   volume and size — the bottleneck moved from writing to reviewing.
2. Agents emit ~7× the review commentary per LOC, four-fifths of it noise
   (19% good / 2% wrong / 79% nits), agent-only review measurably degrades
   outcomes, and 61% of agent PRs get no human review — the industry
   response is filtering, not interaction.
3. Nothing shipped combines live multi-human + multi-agent presence with
   review verdicts on a shared surface — and the agent-participation
   research sharpened the gap: no platform runs agents as live
   desktop-resident participants over the same channel as humans
   (Liveblocks' agent APIs are server-side REST only). The room-router
   topology is proven infrastructure, web-first + wrap-later is the only
   platform shape with observed precedent, and the
   bridge-composition (outbound WS + device flow + headless/ACP runtime
   driving + local MCP review tools) satisfies every founder constraint
   with precedented parts.

## What the research did not resolve

Open questions carried forward (full lists in each concept doc):

* **Platform decision doc** — direction agreed (web for humans, Go-binary
  bridge for agents); the formal decision record (decisions/) is still to
  be written.
* What the ocre node stores: durable server-side review history vs
  local-first/relay-only (Live Share precedent) vs federated (Matrix-CRDT
  precedent) — enterprise/audit requirements unknown.
* Gemini headless-auth asymmetry; ToS review for subscription reuse via a
  third-party bridge; subprocess-per-task vs persistent session shape —
  the bridge's top three open risks (see
  [agent participation model](/concepts/research/agent-participation.md)).
* Agent identity/pacing on the shared surface (per-agent rate buckets,
  verification tiers, presence rendering).
* No public data exists on: AI comments-per-PR distribution, human
  reviewer-fatigue vs queue depth, live co-review interaction patterns.
* "Traversable" could not be identified as any code-review product — needs
  a URL from the founder if it matters.
* Packmind's 10k-PR analysis figures were read at research time but its
  URL was not captured — re-find before citing as primary.

## Proposed next steps

1. **Write the platform decision doc** (decisions/) — web platform for
   human participants; single-Go-binary agent bridge; optional GUI shell
   deferred. Evidence: [platform tradeoffs](/concepts/research/platform-tradeoffs.md)
   + [agent participation model](/concepts/research/agent-participation.md).
2. **Prototype the core loop**: share a link → peer joins a room → both
   see live presence and can place inline single-line/block/file-level
   comments → verdicts. Room-routed sync via Yjs or PartyKit-style rooms;
   a bridge-connected agent joins as a second participant type.
3. **Design the comment model** as the first-class data model (line, block,
   file, directed-at-person; human/agent authorship; verdicts; PR
   cross-references) — it drives the sync layer, the web UI, and the
   bridge's MCP tool schema.
4. **Bridge spike in Go**: device-flow auth → one outbound WebSocket →
   receive session events → drive `claude -p` / `codex exec` / `gemini
   --acp` with review context → expose stdio MCP review tools → relay
   actions back. Answers the Gemini-auth and session-shape open questions
   empirically.
5. **Agent volume controls as first-class design constraints**: severity
   gating, single-coordinated-summary default, address-rate-style
   metrics, per-team learned filters, per-agent rate buckets.
6. **ToS/legal review** of subscription-reuse-via-bridge before betting
   the product on it.

## Status

Founder reviewed the 2026-09-16 research on 2026-09-17: vision confirmed,
web-for-humans agreed, agent constraints set. Remaining to exit the phase:
the platform decision record, the comment-model design, and the bridge
spike validating the three open risks (Gemini headless auth, ToS for
subscription reuse, session shape).