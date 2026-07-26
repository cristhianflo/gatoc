## Why

Users currently have no simple way to opt out of automatic embed fixing for a specific message, which can be inconvenient when sharing links without transformation. Adding a `#nofix` escape hatch and a discoverable help subcommand improves control and reduces confusion.

## What Changes

- Add support for `#nofix` tag in message content so embedfixer skips fix processing for that message.
- Add a new `/embedfixer nofix` subcommand that explains how and when to use the `#nofix` tag.
- Keep existing fix behavior unchanged when `#nofix` is not present.
- Ensure the `#nofix` rule applies across currently supported social media hosts.
- Include rollback plan: remove the new subcommand and disable tag check to restore previous always-fix behavior.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `embedfixer-domain-config`: Extend embedfixer runtime and command UX to support per-message fix bypass via `#nofix` and add guidance subcommand.

## Impact

- Affected code:
  - `internal/features/embedfixer/handlers.go`
  - `internal/features/embedfixer/commands.go`
  - `internal/features/embedfixer/config_handlers.go` (or new command handler file for note response)
- User-facing behavior:
  - New command for explaining nofix usage.
  - Message-level opt-out using `#nofix`.
- APIs/dependencies: no new external dependencies expected.
