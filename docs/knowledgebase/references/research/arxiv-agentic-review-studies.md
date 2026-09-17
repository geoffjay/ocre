---
type: Reference
title: 2026 arXiv empirical studies of agentic code review
description: The AIDev-dataset studies quantifying agent review volume, human oversight gaps, and CRA-only outcomes.
tags: [reference, research, ai, code-review, agents]
generated: { by: omp-agent/glm-5.3, at: 2026-09-16T21:21:05Z }
verified: { by: omp-agent/glm-5.3, at: 2026-09-16T21:21:05Z }
---

# 2026 arXiv empirical studies of agentic code review

The strongest quantitative evidence base for ocre's premise. All mine the
AIDev dataset (agent-authored vs human-authored GitHub PRs) or large GitHub
review corpora; all preprints (not yet peer-reviewed at capture time).

## Human-AI Synergy in Agentic Code Review — arXiv 2603.15911 (2026-03-16)

Zhong, Noei, Zou, Adams. 278,790 inline review conversations, 300 OSS
projects, 2022–2025, 16 identified AI review bots.
AI reviews: **29.6 vs 4.1 tokens per LOC** (~7× human verbosity); >95% Code
Improvement/Defect Detection (humans add Understanding 31%, Testing,
Knowledge Transfer); **85–87% of AI-initiated threads end after one
comment**; suggestion adoption **16.6% vs 56.5%** (agents vs humans);
conversations ending at AI responses → higher rejection (7.1–25.8% vs
0.9–7.8%). Full text re-read 2026-09-16, figures confirmed verbatim.

## These Aren't the Reviews You're Looking For — arXiv 2605.02273 (2026-05-04)

Duma et al. 33,596 agent-authored PRs in ≥100-star repos: **61.38%
received no recorded review**; of reviewed ones, 58.77% agent-only /
10.14% human-only / 31.09% mixed; **71.58% of review comments on agent PRs
were agent-authored**; 25.92% of human comments were agent-steering
commands (vs 1.63% on human PRs). Full text re-read 2026-09-16, table
figures confirmed.

## From Industry Claims to Empirical Reality — arXiv 2604.03196 (2026-04-03)

Chowdhury et al. 3,109 commented AIDev PRs: CRA-only reviewed PRs merged at
**45.20% vs 68.37%** (human-only); abandonment 34.88% vs 21.60%; **60.2%
of closed CRA-only PRs in the 0–30% comment-signal range**; 12 of 13 CRAs
below 60% average signal. Conclusion: CRAs should augment, not replace,
human reviewers. Abstract + tables re-read 2026-09-16, figures confirmed.

## Code Review is a Conversation — arXiv 2607.22095 (2026-07)

Vision paper arguing review assistants should be conversational (ask
grounded questions, respond to explanations, know when to abstain) rather
than one-shot commenters. Research-pass read; anchors ocre's interactive
framing.

Used by [ai-review-noise](/concepts/research/ai-review-noise.md).