---
type: Concept
title: Desktop vs web platform tradeoffs
description: Evidence base for the desktop-vs-web platform decision for ocre — assembled, not decided.
tags: [research, platform, desktop, web, electron, tauri]
status: stable
generated: { by: omp-agent/glm-5.3, at: 2026-09-16T21:21:05Z }
verified: { by: omp-agent/glm-5.3, at: 2026-09-16T21:21:05Z }
sources:
  - id: figma-web
    title: Figma — Building a professional design tool on the web
    resource: https://www.figma.com/blog/building-a-professional-design-tool-on-the-web/
  - id: figma-multiplayer
    title: Figma — How multiplayer technology works
    resource: https://www.figma.com/blog/how-figmas-multiplayer-technology-works/
  - id: gh-repo-limits
    title: GitHub docs — Repository limits (diff size caps)
    resource: https://docs.github.com/en/repositories/creating-and-managing-repositories/repository-limits
  - id: gh-virtualized
    title: GitHub blog — The uphill climb of making diff lines performant (2026-04-03)
    resource: https://github.blog/engineering/architecture-optimization/the-uphill-climb-of-making-diff-lines-performant/
  - id: gh-desktop-electron
    title: GitHub blog — How four native developers wrote an Electron app / Everyone's moving to GitHub Desktop
    resource: https://github.blog/2018-05-02-everyones-moving-to-github-desktop/
  - id: slack-electron
    title: Slack engineering — Building hybrid applications with Electron
    resource: https://slack.engineering/building-hybrid-applications-with-electron/
  - id: viral-loops
    title: Appcues — Viral SaaS product loops (Google Docs case)
    resource: https://www.appcues.com/blog/viral-saas-product
  - id: zed-crdts
    title: Zed blog — CRDTs (native collaboration)
    resource: https://zed.dev/blog/crdts
  - id: liveshare-security
    title: Microsoft Learn — Live Share security (relay-only, E2E)
    resource: https://learn.microsoft.com/en-us/visualstudio/liveshare/reference/security
  - id: yjs-awareness
    title: Yjs docs — Awareness protocol (presence/cursors)
    resource: https://docs.yjs.dev/api/about-awareness
  - id: hopp-tauri
    title: Hoppscotch — Tauri vs Electron migration report (first-party)
    resource: https://www.gethopp.app/blog/tauri-vs-electron
  - id: inkandswitch
    title: Ink & Switch — Local-first software
    resource: https://www.inkandswitch.com/local-first/
  - id: liveblocks-pricing
    title: Liveblocks pricing (managed multiplayer MAU tiers)
    resource: https://liveblocks.io/pricing
  - id: pwa-ios
    title: MagicBell — PWA limitations on iOS/Safari
    resource: https://www.magicbell.com/blog/pwa-ios-limitations-safari-support-complete-guide
---

# Desktop vs web platform tradeoffs

Evidence base for ocre's platform decision. **Assembled, not decided** —
no decision doc exists yet. The core loop being evaluated: share a link to
a review; a coworker clicks it and instantly becomes a live contributor.

## 1. The link-share loop favors web unambiguously

Every product with ocre's exact adoption shape — Figma, Google Docs, Miro,
CodeSandbox — is browser-first. The documented growth mechanism for
collaborative tools is zero-install link sharing: "a viral loop is formed
any time a prospect is able to interact with your product before they've
adopted it"; Google Docs' shared-doc experience is the canonical case, and
Figma's "jump into a design file and start working in seconds" is the
design-tool case.[^viral-loops] No counterexample of an install-first
collaborative tool with a click-and-you're-in contributor path was found.
Slack built its desktop app *after* web growth carried adoption, as an
Electron wrapper around the web app.[^slack-electron] Caveat: no public
k-factor or cohort data isolates link-sharing as the *cause* of adoption —
the mechanism is consistently described, but quantitative claims are
practitioner lore, not measured.

## 2. Large-diff rendering: a real web constraint, solved with engineering

* GitHub — the highest-traffic web review UI — enforces hard limits:
  ~20,000 loadable diff lines / 1 MB raw per PR; 20,000 lines / 500 KB per
  file with only ~400 lines / 20 KB auto-loaded; 300 files per diff
  (rising to ~1,000–3,000 with the 2026 experimental virtualized "Files
  Changed" mode).[^gh-repo-limits] Their engineering post describes window
  virtualization for 10k+ line PRs — i.e., even the best-funded web review
  surface treats large diffs as an active performance problem solved in
  code, not by the platform.[^gh-virtualized]
* Figma handles document scale with a custom WebGL→WebGPU rendering engine
  (C++ to WASM) rather than the DOM — proof the browser can do
  professional-grade, large-canvas work, at the cost of building a custom
  engine.[^figma-web]
* Native git clients (Sublime Merge) handle huge repos with less
  engineering effort (vendor marketing + third-party roundups;
  medium-confidence) — but desktop git clients historically *couldn't do
  hosted review workflows* (Sublime Merge lacks PR integration; Tower
  bolted it on). Desktop got the local-git performance; the collaborative
  review surface gravitated to web.

## 3. Desktop's strengths map to the reviewer's environment, not to the collaboration

Local git operations without server round-trips, offline access,
terminal/editor integration, OS notifications — all real, none of them
about reviewing *together*. The strongest native collaboration precedents
are about pair-editing a live workspace, not reviewing shared artifacts:

* **Zed** built collaboration into its native core with CRDTs (rejecting
  OT), Protobuf RPC, Lamport timestamps, vector clocks — years of bespoke
  engineering — and is extending the same foundation into DeltaDB for
  human+agent shared codebase views.[^zed-crdts]
* **VS Live Share**: host/guest with P2P direct + Azure relay fallback;
  all content stays on the host; relay is E2E-SSH-blind — the strongest
  "gateway that mediates without owning content" prior art. It was
  ultimately sunset as a standalone effort, itself a data point about
  sustaining that model as a feature rather than a product.[^liveshare-security]

## 4. The web multiplayer stack is commodity; desktop multiplayer is bespoke

* Yjs ships an Awareness CRDT (presence/cursors) with interchangeable
  y-websocket (client/server, cross-tab, PubSub-scaled) and y-webrtc (P2P,
  signaling-only; mesh grows quadratically, unsuitable for hundreds of
  collaborators) providers — cursor/comment/presence sync in a browser is a
  weeks-scale integration.[^yjs-awareness]
* Figma's own multiplayer is simpler than full CRDT/OT (property-level LWW
  over WebSockets; 95% of edits saved within 600ms) and sufficed because
  concurrent edits on discrete objects rarely conflict.[^figma-multiplayer]
* Managed cost ladder: Liveblocks is free to ~1k MAU then roughly
  $299+/mo at 10k MAU (vendor pricing page; figures conflict across
  sources — treat as indicative); self-hosting Hocuspocus at 1k MAU is
  estimated ~$20–40/mo in raw infra (third-party estimate). The durable
  finding is the two-order-of-magnitude managed-vs-self-hosted gap and the
  MAU-metering risk when "occasional reviewers" dominate the user
  base.[^liveblocks-pricing]

## 5. Hybrid is the dominant endgame pattern

* GitHub Desktop replaced two native codebases with one Electron app
  ("twice the work, twice the bugs"); within six months of 1.0 the Electron
  version had more users than both classic native versions combined — a
  web-tech app doing local git operations via Node.[^gh-desktop-electron]
* Slack chose Electron explicitly for cross-platform parity + native
  capabilities from one codebase; the web app carried growth
  first.[^slack-electron]
* Tauri offers the same wrap-later shape at ~20–50× smaller footprint
  (Hoppscotch migration: 165 MB → 8 MB bundle, ~70% memory reduction —
  first-party), at the cost of per-platform WebView inconsistency
  (WebKit/WebView2/WebKitGTK).[^hopp-tauri]
* PWA gets web reach but install/notification friction is high and
  platform-constrained (iOS: 4+ taps via Share menu, push only after Home
  Screen install, EU DMA removed standalone-PWA support) — irrelevant to
  the link-share loop itself, but rules out "PWA as our desktop story"
  complacency.[^pwa-ios]

## 6. Local-first web: the pattern matching ocre's stated topology

Ink & Switch's local-first architecture (CRDTs, network optional, server
as relay, local replica authoritative) is the closest match to ocre's
router/gateway framing while keeping the instant-link loop — Zed's DeltaDB
shows the same convergence arriving from the native side.[^inkandswitch]

## Strongest tradeoffs (evidence-based, not a recommendation)

1. **Adoption loop**: web-only strength. Zero-install link → instant live
   contributor. No observed counterexample.
2. **Large-diff performance**: native-leaning strength, but web closes the
   gap with deliberate engineering (virtualization, GPU rendering) — GitHub
   and Figma prove it at scale.
3. **Environment integration** (git, terminal, notifications, offline):
   desktop strength, irrelevant to collaboration itself.
4. **Multiplayer infrastructure**: web-only commodity (Yjs + managed
   rooms); native requires bespoke work (Zed's years of CRDT engineering).
5. **Migration asymmetry**: web-first → wrap in Electron/Tauri later is
   proven (Desktop, Slack); desktop-first → add zero-install link sharing
   has no observed precedent.
6. **Cost/control**: web real-time infra has MAU-metered managed pricing
   or self-host ops; local-first CRDT architecture is the escape hatch that
   also matches ocre's topology.

## Open questions (decided by the platform decision, not by this doc)

* How much review state can live client-side? A review (comments, verdicts,
  presence) is far smaller than a repo — but durable audit history and
  cross-PR references may inherently want a server-side store.
* Enterprise posture: Live Share-style E2E relay vs GitHub-style
  server-held state — what do paying teams actually require (data
  residency, audit trails)?
* Zed's DeltaDB is pre-release; embeddable infrastructure or competitive
  feature?

## Source verification

Re-read on 2026-09-16: Figma multiplayer post (primary, direct read —
earlier verification), PartyKit docs (direct read), GitHub repository-limits
docs + Desktop Electron posts + Copilot changelog (direct reads during
research/verification passes). Slack, Zed, Live Share, Yjs, Hoppscotch,
Liveblocks, PWA-limit, Ink & Switch, Appcues: research-pass reads with
vendor-bias caveats noted inline. Pricing and adoption-quantification
claims are explicitly indicative, not load-bearing.

[^figma-web]: Figma engineering blog (primary).
[^figma-multiplayer]: Figma engineering blog (primary).
[^gh-repo-limits]: GitHub official docs (primary).
[^gh-virtualized]: GitHub engineering blog (primary).
[^gh-desktop-electron]: GitHub blog (primary).
[^slack-electron]: Slack engineering blog (primary).
[^viral-loops]: Appcues practitioner analysis (secondary; consistent across sources).
[^zed-crdts]: Zed engineering blog (primary, vendor bias).
[^liveshare-security]: Microsoft Learn (primary).
[^yjs-awareness]: Yjs official docs (primary, open source).
[^hopp-tauri]: Hoppscotch first-party migration report.
[^inkandswitch]: Ink & Switch research essay (primary research).
[^liveblocks-pricing]: Liveblocks vendor pricing (indicative figures).
[^pwa-ios]: MagicBell guide + MDN (platform constraints).