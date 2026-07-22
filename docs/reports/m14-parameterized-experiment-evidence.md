# M14 Parameterized Experiment Evidence

- Matrix: 5 scenario profiles × 4 paired seeds × 3 policies = 60 completed runs.
- Fairness: demand/supplier/equipment/quality/payment draws are policy-independent and share the same draw hash per profile/seed pair.
- Repeatability: 100 replay builds produce the same EvidenceBundle hash.
- Isolation: every run has a unique branch/run/idempotency namespace, immutable M13 parent hash and `production_writes=0`.
- Completeness: failed and cancelled manifests are present and empty; constraints, run metrics, paired deltas and Pareto flags are calculated before `strategy_evidence_ready=true`.
- Governance: IAOS revision records experiment/evidence/recommendation actions through tenant RLS, idempotency, journal and Outbox; evidence and decision require immutable references and strategy decision requires a distinct approver.
- UI/API: `/api/aese/v1/world/experiments` and `/#world-experiments` expose run-level causal hashes, comparisons, assumptions and the no-auto-apply boundary.
