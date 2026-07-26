## Context

The embedfixer config `show` command currently presents platform config in simple multi-line embed fields. While functionally complete, the output is less scannable than the structured field layout used in finance responses, where values are grouped in consistent table-like blocks.

This change focuses on output presentation only. Config storage, permissions, and command semantics remain unchanged.

## Goals / Non-Goals

**Goals:**
- Render `embedfixer show` output in a table-like style using embed fields.
- Keep all existing data points visible: platform label, source hosts, active domain, and default/custom mode.
- Match readability patterns from finance module embeds (consistent block formatting and spacing).

**Non-Goals:**
- Changing command names, permissions, or argument behavior.
- Altering domain resolution logic or persistence behavior.
- Introducing new supported platforms.

## Decisions

1. Reformat each platform row into consistent field blocks.
   - Use stable per-platform sections with predictable ordering.
   - Keep values wrapped in Discord code blocks (`fix`) where useful for alignment and readability.
   - Alternative considered: single large markdown table in embed description; rejected due to poorer mobile rendering and line-wrap risk.

2. Preserve existing data contract while changing only presentation.
   - Show command still includes: hosts, domain, mode.
   - Rationale: this is a visual change, not a behavior change.

3. Keep field ordering deterministic.
   - Follow supported-platform order from registry so users can compare results quickly.
   - Alternative considered: sort alphabetically per request; rejected to avoid mismatch with platform registry order used elsewhere.

### Sequence Diagram

```text
Admin User       Discord          /embedfixer show handler      Config lookup
   |                |                      |                         |
   | command        |                      |                         |
   |--------------->| InteractionCreate    |                         |
   |                |--------------------->| resolve config rows     |
   |                |                      |------------------------>|
   |                |                      |<------------------------|
   |                |                      | build table-like fields |
   |                |<---------------------| styled embed response   |
```

## Risks / Trade-offs

- [Overly dense formatting on mobile] -> Mitigation: keep each platform as separate field blocks with short line lengths.
- [Visual inconsistency with prior screenshots/docs] -> Mitigation: announce style change in release notes and keep content identical.
- [Discord field limits] -> Mitigation: current platform count is small and safely below embed limits.

## Migration Plan

1. Update `show` embed composition logic only.
2. Verify output in desktop and mobile Discord clients.
3. If readability regressions are reported, rollback to previous field formatting by reverting embed layout changes.

## Open Questions

- Should `mode` be color-coded with emojis/icons, or remain plain text for minimalism?
