---
id: DES-006
title: AESE System Atlas Projection
date: 2026-07-20
status: completed
author: Codex + User
tags: [aese, system-atlas, visualization, governance]
---

# DES-006 AESE System Atlas Projection

## Purpose

AESE needs to show both the target intelligent-enterprise simulation environment and its current completion state. The view must remain consistent with IAOS because AESE depends on IAOS identity, events, capabilities, policies, processes and AI tools.

## Design

- IAOS is the only System Atlas database and API owner; AESE does not create a competing progress database.
- AESE requests `GET /api/v1/system-atlas?view=aese` through its existing same-origin IAOS proxy and authenticated integration identity.
- The graph emphasizes virtual-enterprise model, events, scenario packages, governed ingress, role Agents, 2D sandbox, operations, experiments, economics and deferred scale/3D.
- Clicking a component reveals its target, current state, completion, design/code references and append-only history.
- Search and status filters affect the graph only; they do not alter stored completion.

## Progress Contract

Substantial AESE progress is written to `docs/progress-log.md` and registered against the relevant Atlas node with `scripts/record_system_atlas_update.sh`. Percentages are explicit architectural assessments supported by a document, test report or commit.

## Boundary

System Atlas tracks product construction, not HCTM simulated business facts. Scenario runtime truth remains in IAOS tenant-scoped business tables and scenario observation APIs.
