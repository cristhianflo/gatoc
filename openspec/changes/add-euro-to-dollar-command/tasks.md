## 1. Endpoint and data-contract preparation

- [x] 1.1 Confirm USD and EUR endpoint contracts (fields, timestamp format, source labels) and document expected mapping to internal rate structures.
- [x] 1.2 Define failure-handling rules for each upstream call (network error, non-200, decode error) and expected user-facing fallback messages.

## 2. Finance command multi-currency implementation

- [x] 2.1 Update `/dollar all` data retrieval flow to request USD and EUR endpoints concurrently with explicit timeout handling.
- [x] 2.2 Normalize USD and EUR payloads into a shared in-memory model for embed generation.
- [x] 2.3 Update embed composition to render both currencies with source and update-time context, while preserving existing formatting conventions.
- [x] 2.4 Implement partial-success behavior so available currency data is returned even when one upstream source fails.

## 3. Validation and rollout safety

- [x] 3.1 Validate behavior manually for success, USD-only, EUR-only, and total-failure scenarios using representative sample responses.
- [x] 3.2 Verify command latency and Discord interaction response behavior remain acceptable with two upstream calls.
- [x] 3.3 Prepare rollback steps to disable EUR retrieval and return to USD-only output if upstream instability is observed after deployment.
