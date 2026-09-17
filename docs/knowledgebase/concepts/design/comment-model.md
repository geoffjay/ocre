---
type: Concept
title: Comment model
description: The first-class data model for ocre review artifacts - comment anchors, authorship, verdicts, and cross-references - which drives the sync layer, web UI, and the bridge's MCP tool schema.
tags: [design, comment-model, data-model, mcp]
status: stable
generated: { by: omp-agent/glm-5.3, at: 2026-09-17T04:56:21Z }
verified: { by: omp-agent/glm-5.3, at: 2026-09-17T04:56:21Z }
sources:
  - id: vision
    title: ocre KB - Product vision (comment/verdict requirements)
    resource: /concepts/product-vision.md
  - id: github-review-docs
    title: GitHub docs - About pull request reviews (verdict vocabulary)
    resource: https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/about-pull-request-reviews
  - id: ai-review-noise
    title: ocre KB - AI review noise (severity/signal evidence for governance fields)
    resource: /concepts/research/ai-review-noise.md
---

# Comment model

The comment model is ocre's first-class data model. Everything else —
room sync, web UI, bridge MCP tools, volume governance — is derived from
it. Design rule: **the model is written once here; the prototype's wire
schema, the UI's rendering, and the bridge's MCP tool schema must all be
traceable to this doc.**

## Schema

A review is a session over a change-set (a PR). Every artifact is a
comment; verdicts and cross-references are comment facets.

```yaml
Comment:
  id:            uuid                  # immutable, client- or server-minted
  review_id:     uuid                  # the review session this belongs to
  parent_id:     uuid | null           # null = top-level; set = reply (one level of nesting for replies; threads are linear)
  author:        Actor
  anchor:        Anchor                # where the comment lives
  body:          string                # markdown
  addressed_to:  ActorId | null        # directed-at-a-specific-individual
  severity:      enum | null           # for agent-authored findings: high | medium | low
  created_at:    iso8601
  updated_at:    iso8601
  resolved:      { by: ActorId, at: iso8601 } | null
  resolution_reason: enum | null      # addressed | wont_fix | incorrect | duplicate

Actor:
  kind:          human | agent
  id:            string                # platform user id or bridge-agent id
  display_name:  string
  # agent actors additionally carry: runtime (claude_code|codex|gemini|opencode|acp:...),
  # owner_id (the human who connected the bridge), and verification tier.

Anchor:
  kind:          line | block | file | review
  # line/block/file anchors are stable against re-pushes via:
  file_path:     string                # required for line|block|file
  commit_sha:    string                # the diff the anchor was made on
  line_start:    int                   # line anchors
  line_end:      int | null            # block anchors; null for single-line
  # review-level anchor: kind=review, no file fields

Verdict:                             # a facet of a review-level comment
  type:          approve | request_changes | comment
  by:            Actor
  at:            iso8601
  # agent verdicts are advisory by default (cannot satisfy merge gates unless the
  # review owner promotes the agent - matches GitHub Copilot's default posture)

CrossReference:                      # a facet of any comment body
  target:        { kind: review, id: uuid } | { kind: pr, forge: github|gitlab|..., id: string }
  syntax:        "review:<uuid>" | "pr:<forge>:<id>"   # parsed from body text; rendered as links
```

## Invariants

1. **Anchors are stable across re-pushes**: an anchor records the
   `commit_sha` it was made on; the server computes and stores the
   mapping forward (or marks the comment "moved/outdated") — never
   mutates the original anchor.
2. **Replies are linear threads** (`parent_id` chain), not trees.
3. **Every artifact is a comment** — verdicts are review-level comments
   with a Verdict facet; cross-references parse from body text. One
   storage type, one sync type.
4. **Agent authorship is always visible**: `author.kind == agent`
   comments render with runtime + owner, never impersonating a human.
5. **Severity is advisory input to governance, not display rank.** The
   platform's volume controls (see
   [governance constraints](/concepts/design/agent-governance.md)) may
   gate, collapse, or batch agent comments by severity and
   confidence — the model carries the fields; policy consumes them.

## MCP tool schema (bridge → runtime)

The bridge's local stdio MCP server exposes the model as tools — this is
the review verbs list the runtime sees:

| Tool | Maps to |
|---|---|
| `list_reviews` | sessions the connected agent's owner is in |
| `read_diff` | the change-set for a review (file filter optional) |
| `post_comment` | Comment with anchor: line / block / file / review |
| `reply` | Comment with parent_id |
| `approve` / `request_changes` | Verdict facet |
| `resolve_comment` | resolved + resolution_reason |
| `reference` | parse/emit cross-references in body text |

## Wire contract (room protocol)

The prototype uses a minimal JSON event stream over WebSocket:

```
event: comment.created | comment.updated | comment.resolved |
       verdict.recorded | presence.joined | presence.left | review.state
payload: the Comment/Verdict JSON above (plus presence: actor + cursor state)
```

Comments are discrete append-mostly objects with server-assigned order
and last-writer-wins field updates — the CRDT-easy shape per the
[prior art](/concepts/research/collaborative-review-prior-art.md) (Figma
property-LWW precedent). If the web client later adopts a Yjs document
for live cursors, comments still sync through this event contract;
cursors/awareness ride a separate ephemeral channel (presence split per
Phoenix Presence precedent).

## Open items (deliberately deferred)

* Editing windows / edit history policy for comments.
* Reaction model (👍/👎 on comments — the Greptile learning loop's
  signal source, worth designing early).
* Reviewer roles (required reviewers, owner promotion of agent
  verdicts) — needed before merge gating, not before the core loop.

[^vision]: ocre KB, product vision.
[^github-review-docs]: GitHub docs (verdict vocabulary precedent).
[^ai-review-noise]: ocre KB, AI review noise (governance evidence).