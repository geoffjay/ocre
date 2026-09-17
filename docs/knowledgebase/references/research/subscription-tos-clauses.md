---
type: Reference
title: Subscription ToS clauses for bridge-driven runtimes
description: Verbatim ToS/policy text assembled for legal review - the clauses governing a third-party local bridge driving Claude Code / Codex CLI with the user's consumer subscription.
tags:
  - reference
  - tos
  - legal
  - agents
  - byok
status: draft
generated:
  by: omp-agent/glm-5.3
  at: "2026-09-17T16:23:37Z"
sources:
  - id: anthropic-consumer-tos
    title: Anthropic Consumer Terms of Service
    resource: https://www.anthropic.com/legal/consumer-terms
    last_modified: 2025-10-08
  - id: anthropic-agent-sdk
    title: Claude Code Agent SDK overview (third-party auth note)
    resource: https://code.claude.com/docs/en/agent-sdk/overview
  - id: openai-terms
    title: OpenAI Terms of Use
    resource: https://openai.com/policies/terms-of-use
  - id: google-tos
    title: Google Terms of Service / Gemini Additional Terms
    resource: https://policies.google.com/terms
---

# Subscription ToS clauses for bridge-driven runtimes

Research-assembled for legal review (**not legal conclusions**).
**Status 2026-09-17: the founder read these clauses and accepted the
risk.** Building on subscription-auth driving proceeds; the ACP-only
interface is the recorded fallback if a clause becomes a real concern.
This doc stays as the input for any future formal legal review.

The question: may ocre's bridge — a third-party local tool — drive the
user's own Claude Code / Codex CLI session using the user's consumer
subscription, and relay other participants' review prompts through it?

## Restrict-side clauses (verbatim)

**Anthropic Consumer ToS §3 (2025-10-08)** — the central restriction
(quoted verbatim; ellipses ours, [n] ours):

> "You may not access or use, or help another person to access or use,
> our Services in the following ways: [...] 7. Except when you are
> accessing our Services via an Anthropic API Key or where we otherwise
> explicitly permit it, to access the Services through automated or
> non-human means, whether through a bot, script, or otherwise."

Facially, a bridge spawning `claude -p` is "automated means" with a
consumer OAuth login (not an API key). Whether Anthropic's own headless
docs + `setup-token` constitute "explicit permission" is the dispositive
question.

**Anthropic Consumer ToS §2** — account sharing:

> "You may not share your Account login information, Anthropic API key,
> or Account credentials with anyone else. You also may not make your
> Account available to anyone else. You are responsible for all activity
> occurring under your Account"

The bridge never moves credentials off the owner's machine — but relaying
*other participants' prompts* through one subscriber's account is the
"make your Account available to anyone else" exposure.

**Anthropic Consumer ToS §3(2)** — reselling:

> "2. To develop any products or services that compete with our Services,
> including to develop or train any artificial intelligence or machine
> learning algorithms or models resell the Services."

ocre is not a Claude competitor, but *paid hosted review rooms powered by
users' subscriptions* is the resell-shaped risk to design around (e.g.
ocre's product value = the collaboration surface, not the model capacity).

**Anthropic Agent SDK docs** — third-party auth note (from
code.claude.com/docs/en/agent-sdk/overview): third-party developers may
not offer claude.ai login / rate limits for their products without prior
approval; API-key auth is the directed path. This is the most on-point
vendor statement aimed at exactly ocre's pattern.

**OpenAI Terms of Use** — the extraction clause: no programmatic
extraction of Content for redistribution, no automated access outside the
API; ChatGPT-plan credentials with `codex exec` sit in the same
consumer-vs-API ambiguity (their Terms govern ChatGPT; the CLI is
documented as a ChatGPT-plan feature, which is a permit-side signal).

**Google** — headline only: Gemini CLI headless with consumer Google
account is the least documented path; Gemini Additional Terms restrict
programmatic access outside provided APIs.

## Permit-side signals (vendor-documented mechanisms)

* `claude setup-token` — Anthropic ships a first-party mechanism to mint a
  subscription OAuth token for non-interactive use (docs frame it for
  environments without a browser).
* Headless `-p` mode + `--output-format json` + Agent SDK — documented,
  recommended for scripts and CI; "SDK gives you the same tools, agent
  loop... available as a CLI for scripts and CI/CD".
* `codex exec` — "reuses saved CLI authentication by default" (ChatGPT
  plan login); `codex login --with-access-token` explicitly intended for
  "trusted, non-interactive Codex local workflows".
* These are strong evidence the vendors *intend* programmatic local
  driving — but none addresses the third-party-product question, and the
  Consumer ToS automation clause has no matching consumer-side carve-out.

## Risk classification (for the lawyer)

| Question | Clause tension | Class |
|---|---|---|
| Bridge drives owner's own runtime, owner's own prompts | SDK/headless docs suggest permitted | LOW |
| Owner's agent reviews *other participants' PRs* in ocre rooms | ToS §2 "make Account available to others" | MEDIUM-HIGH |
| ocre charges for rooms while review runs on user subscriptions | ToS §3(2) reselling shape | MEDIUM |
| ocre advertises "use your Claude subscription" in its product | Agent SDK third-party auth note | MEDIUM-HIGH |

## Mitigations to design toward

1. The agent works only on reviews its owner is a participant of — the
   owner's own reviews or rooms they were invited to — keeping the
   account-use owner-scoped.
2. ocre's pricing never prices model capacity (already a founder
   constraint); price the collaboration surface.
3. Ask Anthropic/OpenAI for the third-party-product clarification
   (partner/developer programs exist); document the fallback path
   (per-owner API-key mode for commercial deployments).
4. Legal review of the exact quotes above before public launch; this doc
   is the input, not the answer.
