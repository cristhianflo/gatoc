## 1. Styled embed layout design

- [x] 1.1 Review current `embedfixer show` embed field output and identify target table-like structure inspired by finance module field grouping.
- [x] 1.2 Define consistent per-platform field template (headers, value blocks, separators, and ordering) that works in desktop and mobile Discord clients.

## 2. Show command formatting update

- [x] 2.1 Update `embedfixer show` embed composition to use the new table-like field style while preserving existing data values.
- [x] 2.2 Ensure each platform entry includes source hosts, active domain, and default/custom mode in a consistent visual pattern.
- [x] 2.3 Preserve deterministic platform order in output based on supported platform registry order.

## 3. Verification and rollback safety

- [x] 3.1 Validate rendered output readability in Discord (desktop and mobile) for default-only and mixed default/custom configurations.
- [x] 3.2 Run `go build -o /dev/null ./...` and resolve any compile issues.
- [x] 3.3 Document rollback action to revert to previous show-field formatting if readability regressions are reported.
