# Genesis Strategy Control Runbook

Open `/#world-strategy-control`. Advance all nine adopted-path gates and verify `adopted`; switch to the injected rollback path, advance four gates and verify `rolled_back` plus the compensated work-order commitment. Shadow week 4 must show zero business writes. The page deliberately has no direct production-policy editor.

API evidence is available at `GET /api/aese/v1/world/strategy-control`. IAOS mutations use `/api/v1/genesis/strategy/actions`, an exact release/evidence hash, expected version, idempotency key and an independent approver. Real-production targets are outside this runbook.
