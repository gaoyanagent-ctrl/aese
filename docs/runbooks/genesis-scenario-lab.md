# Genesis Scenario Lab Runbook

1. Run `go run ./cmd/aese experiment validate` and `go run ./cmd/aese experiment expand`; both are offline and the latter is dry-run.
2. Execute the isolated matrix with `go run ./cmd/aese experiment run --apply`. The standard matrix contains 60 runs and zero production writes.
3. Open `http://localhost:4173/#world-experiments`. Confirm `strategy_evidence_ready = true`, baseline/lean/resilient paired comparisons, run hashes and the Simulation / Not production label.
4. Re-run `go test ./internal/experiment -count=100`. Evidence hash must remain stable.

The recommendation is a draft only. There is deliberately no “apply strategy” control. Any production change requires a separate IAOS intent and independent approval.
