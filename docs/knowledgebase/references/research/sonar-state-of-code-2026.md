---
type: Reference
title: Sonar — State of Code Developer Survey (2026)
description: n=1,149 developer survey quantifying the trust-vs-verification gap and the review bottleneck for AI code.
resource: https://www.sonarsource.com/state-of-code-developer-survey-report.pdf
tags:
  - reference
  - research
  - ai
  - survey
  - code-review
generated:
  by: omp-agent/glm-5.3
  at: "2026-09-16T21:21:05Z"
verified:
  by: omp-agent/glm-5.3
  at: "2026-09-16T21:21:05Z"
---

# Sonar — State of Code Developer Survey (2026)

**Resource**:
[sonarsource.com/state-of-code-developer-survey-report.pdf](https://www.sonarsource.com/state-of-code-developer-survey-report.pdf)
Fieldwork October 2025; n=1,149 professional developers worldwide;
vendor-conducted (Sonar sells code-quality tooling) but large-sample and
widely cited.

## What it established

* **96%** of developers don't fully trust AI-generated code to be
  functionally correct, yet only **48%** always check AI-assisted code
  before committing.
* **38%** say reviewing AI code takes more effort than reviewing human
  colleagues' code (27% say less).
* 95% spend at least some effort reviewing/testing/correcting AI output;
  59% rate that effort moderate or substantial; teams spend ~24% of the
  work week validating AI output.
* 61% agree AI "often produces code that looks correct but isn't reliable".
* Framing: AI moved the bottleneck from code generation to **verification
  (review)** — ocre's core motivation in survey form.

## Verification

PDF re-read verbatim 2026-09-16; all figures above confirmed against the
question breakdowns (n=1,149 shown per chart). Used by
[code-review-capacity](/concepts/research/code-review-capacity.md) and
[ai-review-noise](/concepts/research/ai-review-noise.md).
