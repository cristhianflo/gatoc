## Why

The current embedfixer `show` command exposes useful information but lacks the visual structure used by other finance-facing outputs, making quick scanning harder in active channels. Improving presentation now will make configuration status easier to read without changing embedfixer behavior.

## What Changes

- Update embedfixer `show` command output formatting to use a table-like embed field style similar to the finance module presentation.
- Keep the same configuration data surface (platform, source hosts, active domain, default/custom state), but reorganize values into consistent visual blocks.
- Standardize spacing/separators in embed fields to improve readability across multiple platforms.
- Preserve backward compatibility for command semantics and permissions.
- Include rollback plan: revert to the previous plain field formatting if the new table-like layout causes readability or Discord render issues.

## Capabilities

### New Capabilities
- `embedfixer-show-styled-output`: Present embedfixer configuration in a structured, table-like embed format.

### Modified Capabilities
- None.

## Impact

- Affected code: `internal/features/embedfixer/config_handlers.go` and related formatting helpers in the embedfixer feature.
- User-facing behavior: command output is visually improved, while underlying config logic remains unchanged.
- Dependencies/systems: no new external integrations expected.
