---
type: Reference
title: Sadowski et al. — Modern Code Review at Google (ICSE-SEIP 2018)
description: Industrial-scale data on change size, review latency, and reviewer composition at Google.
resource: https://sback.it/publications/icse2018seip.pdf
tags: [reference, research, code-review, capacity]
generated: { by: omp-agent/glm-5.3, at: 2026-09-16T21:21:05Z }
---

# Sadowski et al. — Modern Code Review: A Case Study at Google

**Resource**: [ICSE-SEIP 2018 paper PDF](https://sback.it/publications/icse2018seip.pdf)
(Sadowski, Söderberg, Ball, Katayama, Berman). Mirror:
[ACM DL](https://dl.acm.org/doi/pdf/10.1145/3183519.3183525).

9M reviewed changes at Google, Jan 2014–Jul 2016.

## What it established

* Median change: **24 modified lines**; ~90% of changes touch <10 files.
* Median initial feedback: <~1 hour for small changes, ~5 hours for very
  large — strongly non-linear in change size.
* ~70% of changes committed <24h after mailing; >75% of reviews have
  exactly one reviewer; majority of changes receive no comments other than
  approval.
* Caveat: Google engineers small changes deliberately; figures describe
  their process working, not a natural baseline.

## Verification

The primary PDF (sback.it) returns HTTP 403 from this environment; figures
cross-confirmed 2026-09-16 via the ACM DL citation and multiple independent
secondary sources quoting the paper (michaelagreiler.com,
alastairreid.github.io, UWaterloo course summaries). Used by
[code-review-capacity](/concepts/research/code-review-capacity.md); treated
as primary-quality with the 403 caveat noted.