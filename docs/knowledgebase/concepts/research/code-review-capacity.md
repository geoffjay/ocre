---
type: Concept
title: Human code review capacity and the volume problem
description: Empirical evidence on how much code review humans can effectively do, and where volume breaks review.
tags: [research, code-review, capacity, volume]
status: stable
generated: { by: omp-agent/glm-5.3, at: 2026-09-16T21:21:05Z }
verified: { by: omp-agent/glm-5.3, at: 2026-09-16T21:21:05Z }
sources:
  - id: smartbear-pdf
    title: "SmartBear: The Largest Case Study of Code Review, Ever (Cisco, 2006)"
    resource: https://static1.smartbear.co/support/media/resources/cc/episode_4_thelargestcasestudyofcodereviewever.pdf
    author: team:smartbear
  - id: smartbear-marketing
    title: SmartBear Best Practices for Code Review (marketing pages)
    resource: https://smartbear.com/learn/code-review/best-practices-for-peer-code-review/
  - id: google-eng-practices
    title: Google Code Review Developer Guide — Writing Small CLs
    resource: https://google.github.io/eng-practices/review/developer/small-cls.html
  - id: sadowski-2018
    title: "Sadowski et al., Modern Code Review: A Case Study at Google (ICSE-SEIP 2018)"
    resource: https://sback.it/publications/icse2018seip.pdf
  - id: bacchelli-2013
    title: "Bacchelli & Bird, Expectations, Outcomes, and Challenges of Modern Code Review (ICSE 2013)"
    resource: https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/ICSE202013-codereview.pdf
  - id: czerwonka-2015
    title: "Czerwonka et al., Code Reviews Do Not Find Bugs (IEEE 2015)"
    resource: https://www.microsoft.com/en-us/research/wp-content/uploads/2015/05/PID3556473.pdf
  - id: jureczko-2021
    title: "Jureczko et al., Code review effectiveness: an empirical study (IET Software 2021)"
    resource: https://digital-library.theiet.org/doi/full/10.1049/iet-sen.2020.0134
  - id: alamin-2022
    title: "Alamin et al., Do Small Code Changes Merge Faster? (MSR 2022)"
    resource: https://arxiv.org/abs/2203.05045
  - id: jetbrains-survey
    title: JetBrains State of Developer Ecosystem (cited on Qodana code-review page)
    resource: https://www.jetbrains.com/pages/qodana-use-cases/automated-code-review-tool/
  - id: codeclimate-pr-size
    title: Code Climate Velocity — PR size benchmark (June 2024)
    resource: https://docs.velocity.codeclimate.com/en/articles/2913568-pull-request-size
  - id: sonar-2026
    title: Sonar State of Code Developer Survey (n=1,149, fieldwork Oct 2025)
    resource: https://www.sonarsource.com/state-of-code-developer-survey-report.pdf
  - id: faros-2026
    title: Faros AI — AI code quality and review burden telemetry (vendor blog)
    resource: https://www.faros.ai/blog/ai-code-quality-senior-engineer-review-burden
---

# Human code review capacity and the volume problem

Empirical basis for ocre's founding claim: human review attention is
bounded, and volume — now accelerating with AI-generated code — exceeds it.

## The size cliff: effectiveness collapses beyond a few hundred LOC

The most-cited quantitative anchor is the SmartBear/Cisco case study (2,500
reviews, 50 developers, 3.2M LOC over 10 months, run on CodeCollaborator,
concluded May 2006). Defect density — defects found per kLOC, their measure
of review effectiveness — drops dramatically once a single review session
exceeds ~200 LOC and is near zero beyond ~400 LOC; the study's stated rule:
"review fewer than 200–400 LOC at a time", with 200 a good limit and 400 the
absolute maximum.[^smartbear-pdf]

The effect is independently corroborated by peer-reviewed work without the
vendor interest: Jureczko et al. (IET Software, 2021) found a
faster-than-linear decrease in code review effectiveness as changeset size
grows, replicating multiple prior reports of big-review degradation.[^jureczko-2021]

## Speed and time-box

The same study: ~400–500 LOC/hour is about the fastest anyone should
inspect (best defect detection below ~300 LOC/hour); above ~1,000 LOC/hour
"the reviewer isn't actually looking at the code at all". Expected yield is
~15 defects/hour, higher only for reviews under 150 LOC. Combining the size
limit with the rate limit gives the time-box: at most ~60 minutes per
session (the study cites general psychology for performance decline after
60–90 minutes).[^smartbear-pdf]

Two widely-quoted SmartBear numbers are **not in the published primary
source**: "200–400 LOC over 60–90 minutes yields 70–90% defect discovery",
and any "percentage of defects found in the first hour". They appear only on
marketing pages; their derivation (Chapter 5 of *Best Kept Secrets of Peer
Code Review*) is not freely accessible. Do not cite these as
established.[^smartbear-marketing]

A secondary finding worth keeping: "author preparation" — the author
annotating the diff with commentary before review (~15% of reviews in the
study) — was associated with barely any reviewer-found defects. The authors
resolved the two competing interpretations (self-review finds the defects vs
reviewers switching off) by sampling 300 prepared reviews and concluded the
optimistic one: the author finds the defects while explaining
themselves.[^smartbear-pdf]

## Google: small changes, non-linear feedback latency

Google's published guidance and internal data are the strongest industrial
corroboration:

* Guidance: ~100 lines is usually a reasonable change size; 1,000 usually
  too large. Smallness is conceptual, not a line count — a 200-line change
  in one file may be fine; the same 200 lines spread over 50 files usually
  is not.[^google-eng-practices]
* Data (9M reviewed changes, 2014–2016): median change is **24 modified
  lines**; ~90% of changes touch fewer than 10 files; median initial
  feedback is under ~1 hour for small changes vs ~5 hours for very large
  changes — strongly non-linear in size; ~70% of changes are committed
  within 24 hours; over 75% of reviews have exactly one reviewer; the
  majority of changes receive no comments other than
  approval.[^sadowski-2018]
* Caveat: Google deliberately engineers small changes; this is their
  process working, not a natural baseline.[^sadowski-2018]

## What review is actually for

Microsoft's empirical studies converge on review being primarily a
communication and knowledge-transfer activity, not a defect-finding one:

* Bacchelli & Bird (ICSE 2013; 873 developers + 165 managers surveyed, 17
  interviews, 570 CodeFlow review comments hand-classified): defect-related
  comments were only about one-eighth (~14%) of review comments, and mostly
  addressed "micro"/"superficial" concerns; the largest category was code
  improvement (readability, dead code, consistency). Yet 44% of surveyed
  developers ranked finding defects as the **#1 motivation** for review — a
  major expectations-vs-outcomes gap.[^bacchelli-2013]
* Czerwonka et al. (IEEE 2015; Windows review-comment analysis): only ~15%
  of reviewer comments indicated a possible defect, "much less a blocking
  defect"; the vast majority concerned structure, style, and minor matters.
  Interpretation caveat: this measures the *composition of comments*, not
  defect-detection rate — a small defect-comment share does not mean reviews
  don't catch the defects that matter.[^czerwonka-2015]

## Human attention budget and real-world sizes

* ~45% of developers spend one to two hours per day on code review
  (JetBrains Developer Ecosystem, n≈26k; vendor sells quality
  tooling).[^jetbrains-survey]
* Median organization's *average* PR size is ~197 LOC (Code Climate Velocity
  benchmark, June 2024) — sitting right at the effectiveness cliff — while
  well-maintained OSS median PRs are far smaller (11–28 lines added per an
  independent 10k-PR analysis by Packmind, 2024; primary URL not captured —
  flagged in the [research phase plan](/plans/research-phase.md)); Google's
  median is 24.[^codeclimate-pr-size][^sadowski-2018] The size distribution
  is heavy-tailed: a minority of large PRs consumes disproportionate review
  capacity. The same 10k-PR analysis reports discussions seldom exceeding 10
  comments, 46.9% of PRs involving exactly two people, and 31.2% involving
  only the author.
* Counter-evidence on latency: Alamin et al. (MSR 2022; 845,316 GitHub PRs
  across 100 projects in 10 languages, plus 401,790 Gerrit/Phabricator
  reviews) found **no** relationship between PR size/composition and
  time-to-merge — at OSS-project scale, latency is dominated by reviewer
  availability and first-response time. "Small → faster feedback" (inside
  Google) and "size doesn't predict merge speed" (OSS at large) are both
  true in their contexts.[^alamin-2022]

## The AI-era volume shift

* Sonar State of Code (n=1,149, fieldwork Oct 2025): **96%** of developers
  don't fully trust AI-generated code to be functionally correct, yet only
  **48%** always check AI-assisted code before committing; **38%** say
  reviewing AI code takes more effort than reviewing colleagues' code; 95%
  spend at least some effort reviewing/testing/correcting AI output (59%
  rate that effort moderate or substantial); 61% agree AI "often produces
  code that looks correct but isn't reliable". Sonar frames the result as a
  **verification bottleneck**: AI speeds code generation and moves the
  constraint to review.[^sonar-2026]
* Faros AI telemetry (vendor, 2026): ~25% of PRs now reviewed by an AI agent
  (0% in 2025); PRs ~51% larger; ~31% more PRs merged with no review at
  all. Vendor-sourced and unaudited, but directionally consistent with the
  survey evidence above.[^faros-2026]

## Implications for ocre

1. Human attention is the scarce resource: roughly 1–2 hours/day per
   engineer, effective on changes of a few hundred LOC at a few hundred
   LOC/hour. Any design that spends that attention on noise is the bug
   ocre exists to fix.
2. Review's measured value is communication, knowledge transfer, and design
   gating — not defect-finding. ocre's interactive, collaborative framing
   matches the evidence; a pure defect-finding framing does not.
3. The heavy tail of large PRs plus AI acceleration is where review breaks.
   A tool that partitions attention (who looks at what; what agents
   pre-digest; what deserves human eyes) attacks the measured bottleneck.

## Source verification

Re-read against primary sources on 2026-09-16: the SmartBear PDF (verbatim —
all four conclusions confirmed; the 70–90% and first-hour figures confirmed
absent), the Sonar report PDF (verbatim), Sadowski figures (multi-source
corroboration; the primary PDF returns 403 from this environment).
Jureczko, Alamin, Bacchelli & Bird, Czerwonka, JetBrains, Code Climate,
Packmind, and Faros: research-pass reads with the quality caveats noted
inline.

[^smartbear-pdf]: SmartBear/Cisco case study PDF (primary, vendor-funded).
[^smartbear-marketing]: SmartBear marketing pages (claims not in primary).
[^google-eng-practices]: Google eng-practices guide — small CLs.
[^sadowski-2018]: Sadowski et al. ICSE-SEIP 2018 (Google case study).
[^bacchelli-2013]: Bacchelli & Bird ICSE 2013 (Microsoft).
[^czerwonka-2015]: Czerwonka et al. IEEE 2015 (Windows).
[^jureczko-2021]: Jureczko et al. IET Software 2021.
[^alamin-2022]: Alamin et al. MSR 2022 (arXiv 2203.05045).
[^jetbrains-survey]: JetBrains Developer Ecosystem survey (vendor).
[^codeclimate-pr-size]: Code Climate Velocity PR-size benchmark (vendor telemetry).
[^sonar-2026]: Sonar State of Code survey 2026 (vendor-conducted, n=1,149).
[^faros-2026]: Faros AI vendor telemetry blog.