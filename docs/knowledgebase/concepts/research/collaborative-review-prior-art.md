---
type: Concept
title: Collaborative review prior art
description: Prior art for interactive/collaborative code review — standalone tool history, multiplayer tech, session topologies, and the confirmed gap.
tags:
  - research
  - prior-art
  - collaboration
  - multiplayer
  - crdt
status: stable
generated:
  by: omp-agent/glm-5.3
  at: "2026-09-16T21:21:05Z"
verified:
  by: omp-agent/glm-5.3
  at: "2026-09-16T21:21:05Z"
sources:
  - id: gerrit-site
    title: Gerrit Code Review (project site)
    resource: https://www.gerritcodereview.com/
  - id: phabricator-eol
    title: Phabricator repo banner — Phacility shutdown (2021-06-01)
    resource: https://github.com/phacility/phabricator
  - id: crucible-eol
    title: Atlassian — FishEye and Crucible maintenance-mode / end-of-sale notice
    resource: https://confluence.atlassian.com/fisheye/fisheye-and-crucible-are-in-basic-maintenance-mode-987143949.html
  - id: upsource-eol
    title: JetBrains — Upsource end-of-sales announcement (2022-01-31)
    resource: https://blog.jetbrains.com/upsource/2022/01/31/upsource-end-of-sales-announcement/
  - id: reviewboard-site
    title: Review Board (project site)
    resource: https://www.reviewboard.org/
  - id: github-review-docs
    title: GitHub docs — About pull request reviews
    resource: https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/about-pull-request-reviews
  - id: figma-multiplayer
    title: Figma — How Figma's multiplayer technology works (2019-10-16)
    resource: https://www.figma.com/blog/how-figmas-multiplayer-technology-works/
  - id: kleppmann-ot-crdt
    title: Kleppmann et al. — Interactive Applications for the Arma... (OT vs CRDT comparison)
    resource: https://arxiv.org/pdf/1810.02137
  - id: partykit-docs
    title: PartyKit docs — How PartyKit works (room = Durable Object, same-id routing)
    resource: https://docs.partykit.io/how-partykit-works/
  - id: phoenix-presence
    title: Phoenix Presence docs (distributed presence, CRDT-based)
    resource: https://hexdocs.pm/phoenix/presence.html
  - id: libp2p-webrtc
    title: libp2p docs — WebRTC transport (signaling, hole-punching, relays)
    resource: https://docs.libp2p.io/concepts/transports/webrtc/
  - id: inkandswitch-local-first
    title: Ink & Switch — Local-first software (Kleppmann et al., 2019)
    resource: https://www.inkandswitch.com/local-first/
  - id: matrix-crdt
    title: Matrix-CRDT (Yjs provider over Matrix federation)
    resource: https://github.com/YousefED/Matrix-CRDT
  - id: liveshare-docs
    title: Microsoft Learn — Visual Studio Live Share (maintenance mode, join-from-browser)
    resource: https://learn.microsoft.com/en-us/visualstudio/liveshare/
  - id: liveshare-security
    title: Microsoft Learn — Live Share security architecture (host/guest, relay-only, E2E SSH)
    resource: https://learn.microsoft.com/en-us/visualstudio/liveshare/reference/security
  - id: zed-collab
    title: Zed docs — Collaboration overview (channels)
    resource: https://zed.dev/docs/collaboration/overview
  - id: zed-crdt
    title: Zed blog — CRDTs (Z Sequence; DeltaDB direction)
    resource: https://zed.dev/blog/crdts
---

# Collaborative review prior art

Prior art for ocre's interactive, collaborative, human+agent review model,
across four threads: standalone review tools, multiplayer collaboration
tech, session topologies, and live review specifically.

## The standalone review category was hollowed out, not disproven

2019–2025 EOL/discontinuation history:

* **Phabricator/Differential** — ceased active maintenance 2021-06-01 when
  Phacility (the company) shut down; governance/centralization concerns
  drove the community fork **Phorge** (stable Sept 2022). A company failure,
  not a technical one.[^phabricator-eol]
* **Atlassian Crucible/FishEye** — "basic maintenance mode" since 2020; new
  sales ended 2025-05-13, support ends 2028-05-15. Atlassian explicitly
  attributes the sunset to its cloud strategy and points customers to
  Bitbucket — the most explicit vendor statement of standalone review being
  consolidated into forge-embedded review.[^crucible-eol]
* **JetBrains Upsource** — sales ended 2022-02-01, support 2023-01-31;
  review refocused into the Space platform (which itself wound down its
  git-hosting line in 2024–25 — unverified from primary).[^upsource-eol]
* Survivors occupy niches: **Gerrit** (actively maintained, 3.14.x as of
  2025; powers Android/AOSP, Chromium, Go — but for pre-commit per-commit
  patchset review, not PR-style branch review), **Review Board** (v7, 2025;
  on-prem, broad-VCS), **RhodeCode** (enterprise self-hosted; added AI
  review Dec 2025), **Reviewable** (indie GitHub-PR overlay).[^gerrit-site][^reviewboard-site]

Pattern: standalone review tools died by consolidation into forges and by
company failure — not because standalone review was technically
discredited. The niche that survives is "on-prem/enterprise or
mega-project scale"; nobody owns interactive multiplayer review.

## The gap: forge review is structurally asynchronous

GitHub PR review is batched-verdict, async, with no live presence: reviews
are submitted as Approve/Request-changes/Comment snapshots, reviewers work
on their own schedule, and there is no "who is looking at this diff right
now" — the product surface has no co-presence at
all.[^github-review-docs] Practice literature and team norms are built
around this absence. ocre's live shared surface is a genuine product gap,
not a re-skin.

## Join-via-link multiplayer: the proven UX and its topology

Figma is the canonical implementation and its writeup is effectively a
reference design for ocre's session layer:[^figma-multiplayer]

* Clients are web pages over WebSockets to a server cluster; the servers
  spin up a **separate process per document**, and every editor of that
  document connects to the same process. The URL *is* the address: "join
  via link" = route everyone holding that URL to one per-document node.
* They rejected OT as unnecessarily complex and rejected true CRDTs in
  favor of a custom **property-level last-writer-wins** system (document as
  `Map<ObjectID, Map<Property, Value>>`), with the server as central
  authority for event order, client-generated unique IDs for offline,
  re-download-and-replay on reconnect, and fractional indexing for
  ordering. Presence (colored cursors + names) rides the same connection.
  95% of edits save within 600ms.
* Consequence for ocre: review artifacts (comments, verdicts) are discrete
  objects, not long shared text — exactly the shape LWW-style CRDT-friendly
  sync handles well. Comments live separately from the multiplayer system
  (Figma keeps comments/users in Postgres).

The managed-infrastructure version is PartyKit, whose model is literally
ocre's "join the node your peers are connected to":[^partykit-docs] each
room is a Cloudflare Durable Object; the room id is an arbitrary string
(often a document/project/session id) embedded in the URL; the platform
guarantees every connection with the same id routes to the same
single-tenant room instance at the network-nearest edge datacenter. Rooms
are created on demand with near-zero startup time and can hold state
in-memory. Liveblocks sells the same shape with pre-built presence/storage/
comments primitives on Yjs.

## Sync technology: CRDT-friendly for review artifacts

* OT requires a central sequencer server and complex transformation logic;
  CRDTs (Yjs/Automerge) merge concurrent updates in any order with
  guaranteed convergence, give offline-first for free, and make conflicts
  mathematically impossible. Google Docs' OT heritage is the
  server-dependent counterpoint. Kleppmann et al.'s comparison is the
  strongest technical source.[^kleppmann-ot-crdt]
* Comment threads, presence, and cursors are append-mostly/grow-only data —
  grow-only sets and LWW registers, the easiest CRDT case. Even Figma's
  hand-rolled system is CRDT-inspired LWW rather than OT.
* **Presence and content are separable subsystems with different consistency
  needs.** Phoenix Presence tracks who's in a topic via a CRDT across
  cluster nodes with "no single point of failure... self-healing";
  Figma separates comments (durable, Postgres) from live multiplayer
  sync.[^phoenix-presence] Design cue for ocre: presence = ephemeral CRDT;
  review annotations = durable store, CRDT- or server-authoritative.

## P2P and local-first options

* True browser mesh is limited: full-mesh WebRTC connections grow O(n²),
  browsers can't accept inbound connections, and connectivity needs
  signaling + hole-punching + TURN/circuit relays (libp2p documents the
  whole dance).[^libp2p-webrtc] Mesh only makes sense for ≤4–5 peers or as
  an optimization under a relay.
* **Local-first** (Ink & Switch, Kleppmann et al. 2019) supplies the
  philosophy and tech for "server as relay, local replica authoritative":
  no spinners, multi-device, network optional, seamless collaboration,
  long-term storage, privacy by default, user ownership — with CRDTs as
  the enabling foundation.[^inkandswitch-local-first] This is the closest
  architectural match to ocre's router/gateway framing on a web platform.
* Federated precedent: Matrix-CRDT is a Yjs provider using Matrix rooms
  (federation + auth + E2EE included) as an off-the-shelf backend — any
  homeserver, same room.[^matrix-crdt]

## Live code collaboration: adjacent, but nobody reviews

* **VS/VS Code Live Share** — join-by-link, optionally from a browser with
  no install, becomes a live collaborator (shared editing, debugging,
  terminals, voice). Architecture: host/guest with P2P direct connections
  falling back to an Azure relay; all session content stays on the host;
  traffic end-to-end SSH-encrypted so the relay can't read it — a precedent
  for relay-mediation without server-owned content. But it's in maintenance
  mode "with no additional features planned", and it is an editing/pairing
  surface: no review verdicts, no comments-as-artifacts.[^liveshare-docs][^liveshare-security]
* **Zed channels** — persistent project rooms (with voice) over Zed's
  relay/collab servers; editor buffers use Zed's own CRDT (Z Sequence);
  Zed is extending the same CRDT foundation into DeltaDB, a character-level
  sync layer explicitly aimed at humans and AI agents sharing one consistent
  view of a codebase. Closest shipped product to "join a node your peers
  are on"; no dedicated review surface.[^zed-collab][^zed-crdt]
* Tuple, CodeTogether, JetBrains Code With Me, mob.sh, MobTime — pairing
  and mob-programming tooling; MobTime's "share the link to your timer to
  instantly start collaborating" is the closest small-scale match to ocre's
  link-joins-you-in UX. None do review workflows.

## Summary table

| System | Live presence | Join via link | Review verdicts | Human+AI | Status |
|---|---|---|---|---|---|
| GitHub/GitLab PRs | ✗ | ✗ (async) | ✓ | bots bolt on | dominant |
| Gerrit | ✗ | n/a | ✓ | ✗ | niche (mega-projects) |
| Figma / Google Docs | ✓ | ✓ | n/a | n/a | reference designs |
| Live Share | ✓ | ✓ | ✗ | ✗ | maintenance mode |
| Zed channels | ✓ | ✓ (rooms) | ✗ | announced direction | active |
| CodeRabbit learnings | ✗ (reply loop) | n/a | partial | one agent | active |
| **ocre (target)** | ✓ | ✓ | ✓ | ✓ | whitespace |

No shipped product was found that combines real-time multiplayer presence
with code review verdicts (approve/request-changes) and inline comments on a
shared live surface, human or agent. The whitespace claim survived targeted
searches for research prototypes and launches.

## Implications for ocre

1. The "router/node/gateway" model has a proven production shape: room id
   in the URL, same-id routing to a single-tenant node (Figma per-doc
   process; PartyKit Durable Objects). Build on that shape, not on mesh.
2. Review artifacts are CRDT-easy; live presence is a separate ephemeral
   subsystem; durable review history likely wants a server-side store (see
   open questions).
3. Live Share's relay-only, E2E-encrypted model is the closest prior art
   for a gateway that mediates without owning content — worth studying
   before deciding what the ocre node stores.
4. The graveyard (Phabricator, Crucible, Upsource) is a warning about
   standalone-review economics, not about the model: pair the standalone
   surface with forge integration from day one.

## Source verification

Re-read against primary sources on 2026-09-16: the Figma multiplayer
post, PartyKit "How PartyKit works", Live Share security/connectivity
docs (research pass), Phoenix Presence docs (research pass), libp2p WebRTC
docs (research pass), Ink & Switch essay (research pass), Atlassian
Crucible notice, JetBrains Upsource announcement, Phabricator repo banner,
GitHub PR review docs. Gerrit/Review Board/Zed/Matrix-CRDT items from
project primary sites (research pass). No figure in this doc is load-bearing
beyond the qualitative architecture claims.

[^gerrit-site]: Gerrit project site.
[^phabricator-eol]: Phabricator GitHub repo banner (primary).
[^crucible-eol]: Atlassian official maintenance/EOL notice (primary).
[^upsource-eol]: JetBrains official announcement (primary).
[^reviewboard-site]: Review Board project site.
[^github-review-docs]: GitHub docs (primary).
[^figma-multiplayer]: Figma engineering post (primary, verified by direct read).
[^kleppmann-ot-crdt]: Kleppmann et al. OT/CRDT comparison (arXiv).
[^partykit-docs]: PartyKit docs (primary, verified by direct read).
[^phoenix-presence]: Phoenix framework docs (primary).
[^libp2p-webrtc]: libp2p official docs (primary).
[^inkandswitch-local-first]: Ink & Switch essay (primary research).
[^matrix-crdt]: Matrix-CRDT GitHub repo (primary).
[^liveshare-docs]: Microsoft Learn (primary).
[^liveshare-security]: Microsoft Learn security reference (primary).
[^zed-collab]: Zed official docs.
[^zed-crdt]: Zed engineering blog.
