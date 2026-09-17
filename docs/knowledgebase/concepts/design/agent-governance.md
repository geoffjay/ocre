---
type: Concept
title: Agent governance constraints
description: The volume-control and participation constraints every agent participant must satisfy on the shared live surface - severity gating, single-summary defaults, rate buckets, and address-rate metrics.
tags:
  - design
  - agents
  - governance
  - noise
status: stable
generated:
  by: omp-agent/glm-5.3
  at: "2026-09-17T04:56:21Z"
verified:
  by: omp-agent/glm-5.3
  at: "2026-09-17T04:56:21Z"
sources:
  - id: ai-review-noise
    title: ocre KB - AI review noise (evidence for every constraint below)
    resource: /concepts/research/ai-review-noise.md
  - id: discord-gateway
    title: Discord docs - rate limits and bot verification (per-token buckets, verification tiers)
    resource: https://docs.discord.com/developers/topics/rate-limits
  - id: matrix-appservice
    title: Matrix spec - application-service registration (namespace-scoped registration lesson)
    resource: https://spec.matrix.org/latest/application-service-api/
  - id: comment-model
    title: ocre KB - Comment model (fields governance consumes)
    resource: /concepts/design/comment-model.md
---

# Agent governance constraints

Design constraints for agent participation on the shared live surface.
These are **first-class requirements, not features**: the noise research
([AI review noise](/concepts/research/ai-review-noise.md)) shows
unconstrained agent output destroys review — 79% nits, ~7× verbosity,
authors ignoring everything — so the platform enforces the constraints at
the action boundary regardless of which agent produced the output.

## Constraints

1. **Advisory-by-default verdicts.** Agent approve/request-changes never
   satisfies merge gates unless a human review owner explicitly promotes
   the agent. (GitHub Copilot precedent: its reviews never count toward
   required approvals by default.) The CRA-only merge-rate evidence
   (45.2% vs 68.4%) is why.
2. **Severity field mandatory for agent findings.** Every agent-authored
   `post_comment` carries high/medium/low severity. The room's policy
   decides rendering: e.g., show high inline; collapse medium/low into a
   single digest by default (single-coordinated-summary precedent:
   Cloudflare posts one deduplicated comment from up to 7 specialist
   reviewers).
3. **Per-agent rate buckets.** Token-keyed buckets per agent (Discord
   precedent: per-bot rate-limit buckets + global limits). Defaults
   (tunable per room): ≤ N comments per review per agent per hour, burst
   caps on `post_comment`, and a session-wide agent-comment share ceiling
   so humans cannot be flooded out of their own review.
4. **One conversation surface, no one-shot dumps.** Agents reply on
   threads (the model's `reply` tool); bulk initial findings go through
   the digest path. The 85–87% one-shot-thread evidence says current
   tools dump and leave; ocre's protocol makes the dump the slow path.
5. **Learning loop built in.** 👍/👎 reactions on comments (comment-model
   open item) feed per-team learned filters — the Greptile fix that
   moved address rate 19%→55%. Reactions are governance signal, not
   decoration.
6. **Address-rate metrics, not comment counts.** The platform's KPI is
   the share of agent comments a human acted on (resolved as
   `addressed`). Raw comment volume is never a success metric.
7. **Verification tiers.** Agent identity carries an
   unverified/verified status (Discord's 100-server verification
   precedent; Matrix's namespace-scoped appservice registration
   lesson: registration must be owner-approved and scope-limited, never
   self-service-global). Cross-review-room participation may require
   verification.
8. **Owner accountability.** Every agent actor carries `owner_id` — the
   human who connected the bridge. Rate limits, bans, and audit attach
   to the agent, with escalation to the owner.

## Enforcement point

All constraints execute at the **room action boundary** — the same API
the bridge calls (the MCP tools are thin wrappers over it). The agent
runtime never sees governance; the bridge cannot bypass it; the web
client's identical actions get human-tier defaults.

## Deliberately deferred

* Room-level policy presets (strict/standard/light) and who can change
  them.
* Cross-room reputation (address-rate history following a verified agent).
* Verification process itself (what "verified" requires for an agent).

[^ai-review-noise]: ocre KB, AI review noise (evidence).
[^discord-gateway]: Discord rate-limits docs.
[^matrix-appservice]: Matrix application-service spec.
[^comment-model]: ocre KB, comment model.
