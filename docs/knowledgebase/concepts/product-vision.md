---
type: Concept
title: ocre product vision
description: What ocre is — a standalone, interactive, collaborative code review system where humans and AI agents review together.
tags: [vision, product, code-review, collaboration]
status: stable
generated: { by: omp-agent/glm-5.3, at: 2026-09-16T21:21:05Z }
verified: { by: human:geoff, at: 2026-09-17T03:32:36Z }
---

# ocre product vision

ocre is a **standalone, interactive, collaborative code review system for
humans and AI agents**. This doc records the founding vision as stated by
the project founder on 2026-09-16 and **confirmed by the founder on
2026-09-17**.

## Motivation

"Code review is broken." Two converging failure modes drive the project:

1. **Human review capacity is exceeded by volume.** Effectiveness collapses
   beyond a few hundred LOC per review session and reviewers have roughly
   1–2 hours/day of review attention — see
   [human review capacity](/concepts/research/code-review-capacity.md).
2. **AI agents generate more review output than humans can manage.** Agent
   reviews are ~7× more verbose per line of code, most of that output goes
   unaddressed, and agent-authored PRs increasingly escape human review
   entirely — see
   [AI review noise](/concepts/research/ai-review-noise.md).

## Product shape

* **Standalone system.** ocre is its own product, not a feature embedded in
  a forge (GitHub/GitLab), though forge integration is expected.
* **Shared session topology.** A user joins a router/node/gateway their
  peers are connected to; everyone in a review works on one live, shared
  surface. Production multiplayer systems implement exactly this shape as
  URL-addressed rooms — see
  [prior art](/concepts/research/collaborative-review-prior-art.md).
* **Link-driven contribution.** Users create and share links to reviews of
  pull requests; visiting a link makes the visitor a live contributor to
  that review.
* **First-class comment model.** Comments can be:
  * inline on a single line,
  * inline on a block/range of lines,
  * on a file as a whole,
  * directed at a specific individual,
  * authored by humans or by AI agents — agents are review participants,
    not one-shot bots bolted onto the side.
* **Review verdicts.** PRs can be approved; changes can be requested.
* **Cross-references.** Reviews can reference other PRs/reviews.

## Confirmed direction (2026-09-17)

The founder confirmed the vision and agreed with the platform-research
conclusion: **web platform for human participants** — the web multiplayer
stack is commodity infrastructure, while desktop collaboration is bespoke
engineering (see
[platform tradeoffs](/concepts/research/platform-tradeoffs.md)).

## Hard constraints on agent participation

Stated by the founder 2026-09-17:

1. **ocre does not build or host code review agents.**
2. **ocre does not bill for AI token use.** No metering, no markup, no
   pass-through billing of model inference.
3. **Agents run on the user's own desktop**, on the user's own compute,
   keys, and subscriptions (BYO everything).
4. **BYO agents must not mean a poor experience.** Building custom agents is
   a massive project; connecting an existing agent runtime must be easy.

Implication: the likely shape is a **web platform for humans + a desktop
client/CLI (an "agent bridge") that runs on the agent owner's machine**,
connects their agent runtime to live review sessions, and relays between
the two. This is a candidate architecture, not a decision — see
[agent participation model](/concepts/research/agent-participation.md)
for the research.

## Deliberately open questions

* Agent bridge form factor: CLI/daemon, Wails/Tauri/Electron GUI, or
  both; what the minimum viable desktop footprint is.
* Who operates the router/node/gateway (hosted rooms, self-hosted relay, or
  local-first hybrid), and what lives in durable server-side history vs
  client-side replicas.
* Forge integration depth: read-only overlay vs write-back of
  comments/verdicts.
* Agent identity, pacing, and volume constraints on a shared live surface.

See the [research phase plan](/plans/research-phase.md) for current status,
synthesis, and next steps.