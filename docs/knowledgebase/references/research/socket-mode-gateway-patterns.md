---
type: Reference
title: Slack Socket Mode and Discord Gateway (outbound bot connection patterns)
description: The two canonical proofs that a hosted platform can dispatch real-time events to a process on the user's machine with no inbound endpoint.
tags: [reference, research, agents, bridge, websocket]
generated: { by: omp-agent/glm-5.3, at: 2026-09-17T03:36:04Z }
verified: { by: omp-agent/glm-5.3, at: 2026-09-17T03:36:04Z }
---

# Slack Socket Mode and Discord Gateway

The outbound-connection patterns ocre's agent bridge composes.

## Slack Socket Mode

**Resource**:
[docs.slack.dev — Using Socket Mode](https://docs.slack.dev/apis/events-api/using-socket-mode/)

* App receives Events API + interactive payloads over an **outbound
  WebSocket** instead of a public HTTP Request URL; the wss URL is minted
  at runtime by calling `apps.connections.open` with an app-level
  (`xapp-`) token and refreshes every few hours.
* Pre-authenticated connection → no request-signature validation needed;
  each event envelope must be acked by `envelope_id`.
* Built for developers "behind a corporate firewall" — i.e., bots on
  private machines. Up to 10 simultaneous connections (anycast payload
  distribution); Socket Mode apps excluded from the public Marketplace
  (org-deploy is the distribution workaround).
* Announced at Slack Frontiers, Oct 2020 — the origin of "run your Slack
  bot on your own machine" as a supported platform pattern.

## Discord Gateway

**Resource**:
[docs.discord.com — Gateway](https://docs.discord.com/developers/events/gateway)

* Bots fetch a WSS URL, connect out, send Identify (token + intents),
  receive Ready, heartbeat to stay alive, Resume (session_id + sequence)
  to replay missed events. No inbound endpoint anywhere in the model.
* REST for writes; the Gateway is the real-time receive channel — the
  "WS in, REST out" split ocre's bridge would mirror.
* Control levers to copy: per-bot rate-limit buckets + global 50 req/s;
  IP ban at 10,000 invalid requests/10 min; verification required past
  100 servers.

## Verification

Socket Mode docs re-read verbatim 2026-09-17 (outbound WSS, runtime URL,
acks, 10-connection cap, refresh/disconnect semantics all confirmed);
Discord Gateway read in research pass from primary docs. Used by
[agent participation model](/concepts/research/agent-participation.md).