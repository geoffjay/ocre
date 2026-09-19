---
type: Concept
title: Shared substrate - a common library for ocre and lore
description: Proposal for extracting the domain-free parts of the ocre stack into a common Go library that lore (the knowledge-building sibling) can reuse - bridge, room protocol, MCP, auth; the boundary and the deferred ruling.
tags:
  - design
  - framework
  - substrate
  - lore
  - reuse
status: draft
generated:
  by: omp-agent/glm-5.3
  at: "2026-09-19T16:11:35Z"
---

# Shared substrate - a common library for ocre and lore

This is a **proposal, not a decision.** The founder raised the question
on 2026-09-19: ocre (code review) and lore (knowledge building) share
almost all their platform research and likely much of their code. Is
there a common library that keeps each app focused on its domain?

The answer proposed here: **yes, but extract on the second need, not
the first.** The library is justified when lore's spike exists, not
before. The rest of this doc draws the boundary so both apps build
toward it from day one.

## The family thesis

ocre and lore are two applications of one platform pattern:

* **ocre**: humans and agents review code together on a live shared
  surface. The record is a review (comments, verdicts).
* **lore**: humans and agents build knowledge together on a live shared
  surface. The record is an OKF bundle.

Both inherit the [platform decision](/decisions/platform-architecture.md)
topology: web for humans, a Go agent bridge on the owner's desktop,
governance at the action boundary, no agent hosting, no token billing.
The domains differ; the platform does not.

## The boundary

Divide the stack by one question: **would this code change if the
artifact being discussed were knowledge instead of code?** If no, it
belongs in the substrate. If yes, it belongs in the app.

| Layer | Question it answers | Belongs |
|---|---|---|
| Agent bridge core: outbound WebSocket, reconnect, heartbeat, device-flow auth, self-update | How does a desktop process join a room? | substrate |
| Local stdio MCP server: JSON-RPC framing, tool registration, schema helpers | How does a runtime discover and call tools? | substrate |
| Room protocol: actor, presence, event envelope, action API contract | Who is here and what can they do? | substrate |
| Runtime driving: headless CLI invocation, ACP, env sanitization | How does the bridge invoke the owner's runtime? | substrate |
| Governance engine: rate buckets, digest policy hooks, identity tiers | How is agent volume controlled? | substrate |
| **ocre only:** comment anchors on diffs, verdicts, review lifecycle, forge integration | What is being reviewed? | app |
| **lore only:** OKF document serving, contribution model, edit sync, OKF trust fields | What is being built? | app |

The spike code already reflects this split: `spike/internal/mcp` is
domain-free (a minimal JSON-RPC/MCP stdio server with a websocket dial
helper), while the comment/diff/verdict types live in `roomd` and the
bridge mirrors them. The spike is the extraction inventory.

## What the substrate is not

* **Not a runtime or an agent.** The substrate never invokes a model.
  It carries bytes between the room and the owner's runtime. This
  preserves the no-hosting and no-token-billing constraints in the
  [vision](/concepts/product-vision.md) by construction.
* **Not a product.** No user-facing surface. Apps own their web UIs.
* **Not a protocol standard.** The room protocol is a contract between
  the app server and its clients. If it stabilizes, publishing it is a
  separate decision.

## Extraction path

The [build phase](/plans/build-phase.md) M1 extracts `internal/room`
from the spike into the repo-root Go module. Two options from there:

* **Option A - extract now:** create the substrate repo first, make
  ocre's M1 depend on it.
* **Option B - extract on second use (proposed):** ocre's M1 stays in
  the `ocre` module as `internal/room`. lore's spike extracts from it
  (or from `spike/`) when the second consumer exists. The extraction
  boundary is documented here; the code stays app-internal until two
  consumers force the interface to stabilize.

Option B is proposed because premature extraction before a second
consumer is the classic abstraction-from-one-case failure. The spike
inventory above is the contract both apps build toward. When the second
consumer exists, the library moves to its own repository (the ruling
names it; see the [lore research phase](../../../lore/plans/research-phase.md)).

## What is deliberately deferred

* The substrate repository location and name.
* Whether the web client shares a TS/JS room-protocol package with the
  substrate (the protocol types exist in Go today; a TS mirror would
  serve both apps' front ends).
* Versioning policy between the substrate and the apps.

These are decided by the framework ruling in the
[lore research phase](../../../lore/plans/research-phase.md), step 5.
