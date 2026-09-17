---
type: Concept
title: AI review output volume and the noise problem
description: Evidence that AI/agent review output exceeds human capacity, and the mitigation patterns tools actually use.
tags:
  - research
  - ai
  - code-review
  - noise
  - agents
status: stable
generated:
  by: omp-agent/glm-5.3
  at: "2026-09-16T21:21:05Z"
verified:
  by: omp-agent/glm-5.3
  at: "2026-09-16T21:21:05Z"
sources:
  - id: greptile-postmortem
    title: "Greptile: How to Make LLMs Shut Up (noise postmortem)"
    resource: https://www.greptile.com/blog/make-llms-shut-up
    author: human:daksh-gupta
    last_modified: "2024-12-18T00:00:00Z"
  - id: copilot-ga
    title: GitHub changelog — Copilot code review GA (2025-04-04)
    resource: https://github.blog/changelog/2025-04-04-copilot-code-review-now-generally-available/
  - id: copilot-docs
    title: GitHub docs — About GitHub Copilot code review
    resource: https://docs.github.com/en/copilot/concepts/agents/code-review
  - id: copilot-2026-changelog
    title: GitHub changelog — Auto-resolution and analysis updates in Copilot code review (2026-09-11)
    resource: https://github.blog/changelog/2026-09-11-auto-resolution-and-analysis-updates-in-copilot-code-review/
  - id: coderabbit-learnings
    title: CodeRabbit docs — Learnings (reply-driven comment suppression)
    resource: https://docs.coderabbit.ai/guides/learnings
  - id: sonar-2026
    title: Sonar State of Code Developer Survey (n=1,149, fieldwork Oct 2025)
    resource: https://www.sonarsource.com/state-of-code-developer-survey-report.pdf
  - id: arxiv-synergy
    title: Zhong et al., Human-AI Synergy in Agentic Code Review (arXiv 2603.15911)
    resource: https://arxiv.org/abs/2603.15911
  - id: arxiv-notreviews
    title: Duma et al., These Aren't the Reviews You're Looking For (arXiv 2605.02273)
    resource: https://arxiv.org/abs/2605.02273
  - id: arxiv-empirical
    title: Chowdhury et al., From Industry Claims to Empirical Reality (arXiv 2604.03196)
    resource: https://arxiv.org/abs/2604.03196
  - id: arxiv-conversation
    title: "Code Review is a Conversation: Toward Conversational AI Review Assistants (arXiv 2607.22095)"
    resource: https://arxiv.org/abs/2607.22095
  - id: cloudflare-blog
    title: Cloudflare engineering — AI code review orchestrator
    resource: https://blog.cloudflare.com/ai-code-review/
  - id: ms-devblog
    title: Microsoft devblogs — Enhancing code quality at scale with AI-powered code reviews
    resource: https://devblogs.microsoft.com/engineering-at-microsoft/enhancing-code-quality-at-scale-with-ai-powered-code-reviews/
  - id: codeant-fpr
    title: CodeAnt — AI code review false positives (vendor benchmark)
    resource: https://www.codeant.ai/blogs/ai-code-review-false-positives
  - id: anthropic-launch
    title: TechCrunch/InfoQ coverage — Anthropic launches Claude Code Review (2026-03-09)
    resource: https://techcrunch.com/2026/03/09/anthropic-launches-code-review-tool-to-check-flood-of-ai-generated-code/
  - id: hn-bubble
    title: Hacker News — "There is an AI code review bubble" (2026-02)
    resource: https://news.ycombinator.com/item?id=46766961
---

# AI review output volume and the noise problem

Empirical basis for ocre's second founding claim: AI/agent review output
exceeds what humans can manage, and the industry's mitigations point at —
but have not shipped — interactive human+AI review.

## The canonical case study: Greptile's postmortem

At launch, Greptile's bot left up to 10 comments on a 20-change PR, and
"the PR author would simply start ignoring all of them" — the defining
failure mode of AI review.[^greptile-postmortem]

Their own comment-quality analysis: **~19% good, 2% flat-out incorrect,
79% nits** ("technically true but not something the dev cared
about").[^greptile-postmortem] What failed and what worked:

* **Prompting failed**: they could not reduce nits without also reducing
  critical comments; few-shot examples made it worse (the model inferred
  superficial characteristics, not the pattern).
* **LLM-as-judge failed**: an LLM rating its own comment severities was
  "nearly random", and added an inference step to latency.
* **Learned per-team filtering worked**: nits are subjective and
  team-specific, so Greptile embeds each team's past comments that were
  addressed/upvoted or downvoted; a new comment is blocked when it is
  cosine-similar to ≥3 downvoted comments. **Address rate (share of comments
  authors actually act on, detected from subsequent commit diffs) rose from
  19% to 55+% within two weeks** — and this metric, not comment count, is
  the industry's de facto noise-quality KPI.[^greptile-postmortem]

## How GitHub constrains Copilot's review output

Copilot code review went GA 2025-04-04 after >1M developers used it in the
first month of public preview.[^copilot-ga] Its output is deliberately
constrained:

* Reviews take the non-gating "Comment" form; by default Copilot approvals
  do **not** count toward a repository's required approvals.[^copilot-docs]
* Whole file classes are excluded — dependency/lock files (package.json,
  Gemfile.lock), log files, SVGs — precisely the files that generate pure
  noise.[^copilot-docs]
* "Effort levels" trade thoroughness against volume/cost: Lite
  ($0.05–$1 in AI credits, default) vs Balanced ($0.25–$5).[^copilot-docs]
* By default it reviews a PR once (not every push).[^copilot-docs]
* Docs instruct users: "Always validate Copilot's feedback carefully.
  Supplement Copilot's feedback with a human review."[^copilot-docs]
* The interaction is one-directional: human replies on Copilot review
  comments "won't be visible to Copilot, and Copilot won't reply".
  2026 additions are feedback plumbing, not conversation: auto-resolution
  of addressed comments and resolution reasons (Addressed / Won't fix /
  Incorrect); an ensemble-of-agents-plus-tools experiment for Lite raised
  addressed comments per review by +47% high / +31% medium / +11% low
  severity at ~8% lower cost.[^copilot-2026-changelog]

## Agent reviews vs human reviews: measured differences

Three 2026 arXiv studies quantify the volume gap on GitHub data:

**Human-AI Synergy in Agentic Code Review** (278,790 inline review
conversations across 300 projects, 2022–2025, 16 identified AI review
bots):[^arxiv-synergy]

* **Verbosity**: AI reviews average **29.6 tokens per line of code reviewed
  vs 4.1 for human reviews** (~7×) — agents explain from first principles;
  humans lean on shared project knowledge. The paper explicitly flags this
  as an increased reading burden.
* **Narrowness**: >95% of agent comments are Code Improvement or Defect
  Detection. Humans additionally do Understanding (31% of comments when
  reviewing human code — and clarification questions trigger the most
  back-and-forth: 69% of such threads extend past the first exchange),
  Testing, and Knowledge Transfer (4–6%).
* **One-shot behavior**: **85–87% of AI-initiated review threads end after
  the first comment** — no follow-up. Conversations ending at AI responses
  show consistently higher PR rejection (7.1–25.8%) than those ending at
  human responses (0.9–7.8%).
* **Adoption**: suggestions from agents are adopted at **16.6% vs 56.5%**
  for human reviewers (despite agents generating 88,011 suggestions vs
  humans' 25,673); over half of unadopted agent suggestions were incorrect
  or fixed differently; adopted agent suggestions increase code
  complexity/size more.
* **Form**: agent comments carry severity labels, summary titles, lint-rule
  citations, and downstream-file lists humans don't write — a "hybrid
  artifact" mixing assessment with tool output. The paper suggests routing
  structured labels/tooling to pipelines and showing humans only the core
  assessment.

**How Humans Review AI-Generated Pull Requests** (AIDev dataset; 33,596
agent-authored PRs in ≥100-star repos):[^arxiv-notreviews]

* **61.38% of agent-authored PRs received no recorded review at all.**
* Of those reviewed, only 10.14% were human-only (58.77% agent-only,
  31.09% mixed); **71.58% of review comments on agent PRs were written by
  agents**.
* 25.92% of human comments on agent PRs were agent-steering commands
  ("@coderabbit fix lint failure") vs 1.63% on human PRs — human
  "oversight" of agent code is often automation-mediated steering, not
  evaluation.

**From Industry Claims to Empirical Reality** (3,109 commented AIDev PRs; a
direct test of the industry claim that CRAs can manage ~80% of OSS
PRs):[^arxiv-empirical]

* CRA-only-reviewed PRs merged at **45.20% vs 68.37%** for human-only
  reviewed PRs, with significantly higher abandonment (34.88% vs 21.60%).
* Of 98 closed CRA-only PRs analyzed, **60.2% fell in the 0–30% signal
  range** of comment usefulness; 12 of 13 CRAs averaged below 60% signal
  ratio.
* The paper's conclusion: CRAs should **augment, not replace** human
  reviewers.

Survey-level automation bias matches the mining data: 96% of developers
don't fully trust AI code to be functionally correct, but only 48% always
check it before committing (Sonar, n=1,149 — see
[capacity doc](/concepts/research/code-review-capacity.md) for the full
set).[^sonar-2026]

## Mitigation patterns in production tools

Every serious tool converges on volume controls; these are the baseline
expectations for any agent participant in ocre:

1. **Scoping/exclusions** — Copilot's excluded file classes and
   review-once default.[^copilot-docs]
2. **Severity/confidence gating** — Greptile's docs lead with severity
   thresholds; Ellipsis runs a multistage filter chain (dedup,
   confidence-threshold, logical-correctness/hallucination filters) and
   retains filtered-out comments with reasoning rather than silently
   dropping them.
3. **Learned per-team filters** — Greptile's embedding clusters
   (19%→55+% address rate); CodeRabbit's "learnings": replying to a
   CodeRabbit comment creates an explicit team learning that suppresses
   that comment class in future reviews — the strongest shipped
   human↔agent review interaction loop.[^coderabbit-learnings]
4. **Verification before posting** — Anthropic's Claude Code Review (launched
   2026-03-09) runs multiple specialized agents, then a verification step
   that deliberately tries to disprove each finding before posting;
   surviving findings are deduplicated and severity-ranked. Vendor-reported
   (unaudited): findings on 16%→54% of PRs, <1% incorrect, ~20 min and
   $15–25 per review.[^anthropic-launch] Cursor BugBot runs 8 parallel
   review passes with randomized diff order and flags only what multiple
   passes independently surface. Greptile's TREX executes the code during
   review.
5. **Single-coordinated-summary UX** — CodeRabbit posts one structured
   "walkthrough" summary per PR; Cloudflare's internal orchestrator runs up
   to 7 specialist reviewers per MR but posts **one** deduplicated,
   severity-judged comment via a coordinator agent, risk-tiered so light
   changes get light review (first month: 131,246 runs over 48,095 MRs
   across 5,169 repos; median 3m39s; $1.19/review), explicitly "not a
   replacement for human code review". Microsoft's internal assistant
   covers >90% of ~600K PRs/month and its feedback loop fed Copilot PR
   Reviews' GA.[^cloudflare-blog][^ms-devblog]
6. **Judge agents** — Qodo 2.0's judge deduplicates, resolves conflicts,
   and filters low-signal findings across sub-agents.

The complaint the whole field is fighting is mainstream and named — "it's
really hard to get it not to tell you 20 highly speculative reasons why the
code is problematic along with the one critical error" (Hacker News, Feb
2026); practitioners report resolving AI comments without reading them when
signal-to-noise is low.[^hn-bubble]

## The interactive frontier

Research has articulated ocre's premise directly: *Code Review is a
Conversation* (arXiv 2607.22095) argues current tools wrongly frame review
as one-shot commenting, and that assistants should ask grounded questions,
respond to developer explanations, summarize unresolved issues, and know
when to abstain.[^arxiv-conversation] The observed agent-steering commands
on agent PRs show humans already attempt this conversation in review
threads[^arxiv-notreviews] — but no shipped product implements live
multi-human + multi-agent co-review on a shared surface. The precedents are
partial: CodeRabbit's reply-learnings loop, Copilot's one-way resolution
plumbing.

## Numbers that do not exist yet

No standardized "AI comments per PR" distribution exists — tools report
address rates and precision instead; per-PR counts are anecdotal (Greptile's
"10 comments on a 20-change PR"; CodeRabbit rated "most talkative" in
comparisons). Vendor benchmarks self-measure and win (CodeAnt claims
untuned first-gen reviewers run 40–80% false positives, tuned platforms
~5–15%, and self-measures 52.2% precision[^codeant-fpr]); cross-vendor
numbers are not comparable. Anthropic's <1% incorrect claim is
vendor-reported. A tool named "Traversable" could not be identified as any
code-review product during research.

## Implications for ocre

1. The problem is quantified: agents produce ~7× the commentary per LOC,
   roughly four-fifths of it stylistic or wrong, and agent-only review
   measurably degrades outcomes (23-point merge-rate drop).
2. Scope, severity, learned filters, verification, and
   single-coordinated-summary are table stakes for any agent participating
   in ocre — not optional features.
3. Interactive, conversational review is where both the research and the
   partial product precedents point, and nothing shipped occupies it.
   ocre's live shared surface with agents as peers is the whitespace.
4. Design metric: measure address-rate-style signal quality (share of agent
   comments a human acts on), never raw comment counts.

## Source verification

Re-read against primary sources on 2026-09-16: the Greptile postmortem
(verbatim), arXiv 2603.15911 full text (all figures above confirmed
verbatim), arXiv 2605.02273 full text (table figures confirmed), arXiv
2604.03196 abstract (merge/abandonment/signal figures confirmed), the
Copilot GA changelog and docs.github.com code-review page (verbatim),
the Sonar PDF (verbatim). CodeRabbit, Cloudflare, Microsoft, Anthropic,
BugBot, Ellipsis, Qodo, CodeAnt, and HN items: research-pass reads with
vendor-bias caveats noted inline.

[^greptile-postmortem]: Greptile engineering postmortem (primary, first-party).
[^copilot-ga]: GitHub changelog, Copilot code review GA.
[^copilot-docs]: GitHub docs, About Copilot code review (primary).
[^copilot-2026-changelog]: GitHub changelog, 2026-09-11 Copilot code review updates.
[^coderabbit-learnings]: CodeRabbit learnings docs (primary vendor docs).
[^sonar-2026]: Sonar State of Code survey 2026 (vendor-conducted, n=1,149).
[^arxiv-synergy]: arXiv 2603.15911 (peer-adjacent empirical study).
[^arxiv-notreviews]: arXiv 2605.02273 (empirical study, AIDev dataset).
[^arxiv-empirical]: arXiv 2604.03196 (empirical study, AIDev dataset).
[^arxiv-conversation]: arXiv 2607.22095 (vision paper).
[^cloudflare-blog]: Cloudflare engineering blog (primary, first-party).
[^ms-devblog]: Microsoft engineering devblog (primary, first-party).
[^codeant-fpr]: CodeAnt vendor benchmark (biased source).
[^anthropic-launch]: TechCrunch launch coverage + InfoQ (vendor-reported metrics).
[^hn-bubble]: Hacker News thread (anecdotal, representative community signal).
