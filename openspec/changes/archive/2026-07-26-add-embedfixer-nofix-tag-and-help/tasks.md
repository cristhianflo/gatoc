## 1. Nofix directive behavior

- [x] 1.1 Add `#nofix` directive detection in embedfixer message handler before URL-host processing.
- [x] 1.2 Apply case-insensitive handling for `#nofix` and ensure directive skips suppression + fixed-message posting.
- [x] 1.3 Preserve existing embedfixer behavior when directive is not present.

## 2. Nofix command guidance

- [x] 2.1 Add `/embedfixer nofix` subcommand metadata under existing embedfixer command group.
- [x] 2.2 Implement subcommand handler that returns a short explanation and at least one usage example for `#nofix`.
- [x] 2.3 Keep command permission and response conventions aligned with current embedfixer commands.

## 3. Verification and safety

- [x] 3.1 Add/adjust tests for nofix bypass behavior and nofix subcommand response expectations.
- [x] 3.2 Validate behavior manually in a guild: with `#nofix` no fix occurs, without tag normal fix flow runs.
- [x] 3.3 Run `go build -o /dev/null ./...` and resolve compile issues.
- [x] 3.4 Document rollback operation to disable nofix directive and remove nofix command if regressions occur.
