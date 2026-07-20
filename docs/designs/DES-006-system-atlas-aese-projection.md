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

Substantial AESE progress is written to `docs/progress-log.md` and a versioned `atlas-updates/*.json` declaration. CI rejects substantive changes without both records. Main-branch environments configured with Atlas endpoint secrets run `scripts/sync_system_atlas_updates.sh`; the API uses a unique `update_key`, so retries do not duplicate history. Percentages remain explicit architectural assessments supported by evidence rather than commit-count inference.

## Boundary

System Atlas tracks product construction, not HCTM simulated business facts. Scenario runtime truth remains in IAOS tenant-scoped business tables and scenario observation APIs.

## Explainable Drill-Down

- Dagre creates the initial layered layout from actual node dimensions; nodes remain draggable and the toolbar can restore automatic layout.
- Selecting a component highlights its one-hop neighbors and directional relationship labels while dimming unrelated components.
- Detail content is separated into design documents, functional entries, code locations and evidence. Registered Markdown is rendered in a modal reader through IAOS's restricted document API.
- AESE `entry_refs` resolve `#sandbox`, `#live`, `#integration` and `#atlas` to real application states. IAOS entries open the corresponding IAOS workspace route.
