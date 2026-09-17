---
type: Reference
title: SmartBear/Cisco code review case study (2006)
description: The largest published peer-code-review case study — source of the 200-400 LOC size cliff and inspection-rate limits.
resource: https://static1.smartbear.co/support/media/resources/cc/episode_4_thelargestcasestudyofcodereviewever.pdf
tags:
  - reference
  - research
  - code-review
  - capacity
generated:
  by: omp-agent/glm-5.3
  at: "2026-09-16T21:21:05Z"
verified:
  by: omp-agent/glm-5.3
  at: "2026-09-16T21:21:05Z"
---

# SmartBear/Cisco code review case study

**Resource**: [PDF — "Lightweight Code Review Episode 4: The Largest Case
Study of Code Review, Ever" (Jason Cohen, SmartBear,
© 2010)](https://static1.smartbear.co/support/media/resources/cc/episode_4_thelargestcasestudyofcodereviewever.pdf)

2,500 reviews, 50 developers, 3.2M LOC over 10 months at Cisco Systems
(concluded May 2006), run on CodeCollaborator.

## What it established

* Defect density (per kLOC) drops dramatically above ~200 LOC per review
  session; near zero beyond ~400. Rule: review fewer than 200–400 LOC at a
  time.
* ~400–500 LOC/hour is the maximum effective inspection rate; best detection
  below ~300 LOC/hour; >1,000 LOC/hour = "not actually looking". Expected
  yield ~15 defects/hour, higher only under 150 LOC.
* Derived time-box: <60 minutes per session (cites general psychology for
  60–90-minute concentration decline).
* "Author preparation" (author annotates the diff pre-review, ~15% of
  reviews) correlates with barely any reviewer-found defects; sampling 300
  reviews supported the self-review interpretation over the
  reviewer-bias one.

## Caveats (recorded, load-bearing)

* Vendor-funded: SmartBear sells the tool the study ran on; single company;
  never peer-reviewed; statistical treatment deferred to a book chapter
  ("Best Kept Secrets of Peer Code Review", ch. 5).
* The widely-quoted "70–90% defect discovery for 200–400 LOC in 60–90 min"
  appears only on marketing pages — **not in this PDF**; no "defects found
  in first hour" figure exists in any reachable primary.
* Independently corroborated (size effect) by Jureczko et al., IET Software
  2021 (peer-reviewed).

## Verification

PDF re-read verbatim 2026-09-16; all four conclusions confirmed; absence of
the 70–90% and first-hour figures confirmed. Used by
[code-review-capacity](/concepts/research/code-review-capacity.md).
