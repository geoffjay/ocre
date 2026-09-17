---
type: Reference
title: Greptile — How to Make LLMs Shut Up (2024-12-18)
description: First-party postmortem on AI review comment noise; what failed (prompting, LLM-as-judge) and what worked (per-team embedding-cluster filters).
resource: https://www.greptile.com/blog/make-llms-shut-up
tags: [reference, research, ai, code-review, noise]
generated: { by: omp-agent/glm-5.3, at: 2026-09-16T21:21:05Z }
verified: { by: omp-agent/glm-5.3, at: 2026-09-16T21:21:05Z }
---

# Greptile — How to Make LLMs Shut Up

**Resource**: [greptile.com/blog/make-llms-shut-up](https://www.greptile.com/blog/make-llms-shut-up)
(Daksh Gupta, co-founder; 2024-12-18).

The canonical AI-review-noise case study, first-party and candid about
failures.

## What it established

* At launch the bot left up to 10 comments on a 20-change PR; authors
  "started ignoring all of them".
* Self-analysis of comment quality: **~19% good, 2% flat-out incorrect,
  79% nits.**
* Failed: prompting (couldn't cut nits without cutting critical comments;
  few-shot made it worse) and LLM-as-judge severity rating ("nearly
  random", plus latency).
* Worked: per-team embedding-cluster filtering over past
  addressed/upvoted/downvoted comments (block new comments similar to ≥3
  downvoted ones). **Address rate rose 19% → 55+% within two weeks.**
* Key learning recorded: nits are subjective and team-specific — the fix
  must be learned per team, not prompted globally.

## Verification

Re-read verbatim 2026-09-16 (published 2024-12-18 confirmed from page
metadata); all figures confirmed. "Address rate" (share of comments the
author actually acts on, detected from subsequent commits) is the de facto
industry noise KPI. Used by
[ai-review-noise](/concepts/research/ai-review-noise.md).