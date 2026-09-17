---
type: Reference
title: Microsoft code review empirical studies (Bacchelli & Bird 2013; Czerwonka 2015)
description: ICSE/IEEE studies establishing what code review is actually used for at Microsoft.
tags: [reference, research, code-review]
generated: { by: omp-agent/glm-5.3, at: 2026-09-16T21:21:05Z }
---

# Microsoft code review empirical studies

## Bacchelli & Bird — Expectations, Outcomes, and Challenges of Modern Code Review (ICSE 2013)

**Resource**:
[Microsoft Research PDF](https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/ICSE202013-codereview.pdf)

873 developers + 165 managers surveyed, 17 interviews, 570 CodeFlow review
comments hand-classified. Established: defect-related comments ≈ one-eighth
(~14%) of review comments, mostly micro/superficial; largest category is
code improvement; **44% of developers rank defect-finding as the #1 review
motivation** — the expectations-vs-outcomes gap.

## Czerwonka, Greiler, Tilford, Bird — Code Reviews Do Not Find Bugs (IEEE 2015)

**Resource**:
[Microsoft Research PDF](https://www.microsoft.com/en-us/research/wp-content/uploads/2015/05/PID3556473.pdf)

Windows review-comment analysis: ~15% of reviewer comments indicate a
possible defect, "much less a blocking defect". Interpretation caveat
(recorded): this measures comment *composition*, not defect-detection rate.

## Verification

Research-pass reads of the primary PDFs during the 2026-09-16 research
phase; figures corroborated across the two studies and secondary
literature. Used by
[code-review-capacity](/concepts/research/code-review-capacity.md) — the
"review is communication, not defect-finding" design anchor.