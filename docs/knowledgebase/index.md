---
okf_version: "0.2"
---

# ocre knowledge base

This is the working knowledge base for the ocre project, conforming to the
[Open Knowledge Format (OKF) v0.2](https://github.com/GoogleCloudPlatform/knowledge-catalog/blob/main/okf/SPEC.md).

It consolidates working knowledge about the project: what ocre is, how it is
structured, decisions and their rationale, recurring patterns, and plans.
It is authored by people and agents and meant to be read by both.

## For agents (policy)

This section is the single source of truth for how agents should use this knowledge
base. Tooling injects it into context automatically, so it does not depend on
`CLAUDE.md`/`AGENTS.md` being picked up.

Wired agents: **Claude Code** via a `SessionStart` hook (`.claude/hooks/kb-inject.py`,
with `PostToolUse`/`Stop` reminders in `.claude/hooks/kb-reminder.py`), **opencode**
via the `instructions` array in `.opencode/opencode.jsonc`, and **oh-my-pi** via the
`.omp/extensions/kb-hooks.ts` extension.

**Consult before acting.** Before working on a task, scan the entries below and
read any concept/decision/pattern doc relevant to what you are about to change.
Prefer the recorded decision or pattern over re-deriving one. This index is the
map; read the specific doc on demand rather than guessing.

**Update after acting.** Update the knowledge base when a change would make an
existing entry wrong or leave a new fact unrecorded. In particular:

* A new architectural decision or a change to startup/threading →
  add or update a [decision](decisions/index.md) and relevant concept docs.
* A new recurring convention → add a [pattern](patterns/index.md).
* A new concept or architectural understanding → add a [concept](concepts/index.md).
* A forward-looking plan or roadmap item → add a [plan](plans/index.md).
* An external source or spec referenced by the KB → add a [reference](references/index.md).

Concept docs require YAML frontmatter with a `type` field; `index.md` and
`log.md` are reserved. When you add a doc, add a one-line pointer to the matching
category index below and a line to [`log.md`](log.md). If you deliberately decide
*not* to record a change, that is fine — the policy is judgement, not a mandate
to touch the KB on every edit.

## Concepts

* [ocre product vision](concepts/product-vision.md) - what ocre is: a standalone, interactive, collaborative code review system for humans and AI agents.
* [Human code review capacity](concepts/research/code-review-capacity.md) - empirical evidence on human review limits and the volume problem.
* [AI review noise](concepts/research/ai-review-noise.md) - evidence that agent review output exceeds human capacity, and mitigation patterns.
* [Collaborative review prior art](concepts/research/collaborative-review-prior-art.md) - standalone tool history, multiplayer topologies, and the confirmed live-review gap.
* [Platform tradeoffs](concepts/research/platform-tradeoffs.md) - desktop vs web evidence, assembled but not decided.
* [Agent participation model](concepts/research/agent-participation.md) - how user-owned agents become live review participants: desktop bridge, BYO runtime and credentials, MCP/ACP composition.
* [Comment model](concepts/design/comment-model.md) - the first-class data model: comment anchors, authorship, verdicts, cross-references; drives sync, UI, and the bridge's MCP tools.
* [Agent governance constraints](concepts/design/agent-governance.md) - volume-control requirements for agent participants: advisory verdicts, severity, rate buckets, address-rate metrics.
* [Core-loop spike report](concepts/design/core-loop-spike.md) - what the roomd + ocre-bridge prototype proved end-to-end, and its known limitations.
* [Shared substrate](concepts/design/shared-substrate.md) - proposal for the common library ocre and lore share: bridge, room protocol, MCP, governance; extraction deferred to the second consumer.

## Decisions

* [Platform architecture](decisions/platform-architecture.md) - web for humans, Go agent bridge on the owner's desktop, no GUI shell for v1; founder-verified.

## Patterns

* _(empty — add pattern docs here)_

## Sibling projects

* [lore seed bundle](../lore/index.md) - the seed knowledge base for lore, ocre's knowledge-building sibling. It lives at `docs/lore/` temporarily and moves to its own repository. The ocre KB site does not include it.

## Plans

* [Build phase](plans/build-phase.md) - milestones M1 to M5 from the verified spike to v1; founder rulings recorded, research phase closed.
* [lore research phase](../lore/plans/research-phase.md) - the plan that turns inherited ocre research into lore's foundation: OKF serving, prior art, data model, edit sync, framework ruling, spike.

## References

* [OKF spec](references/okf-spec.md) - pointer to the Open Knowledge Format v0.2 specification.
* [Subscription ToS clauses](references/research/subscription-tos-clauses.md) - verbatim ToS text for legal review of bridge-driven subscription auth.
