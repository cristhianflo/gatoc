## Context

Embedfixer currently auto-processes supported social URLs whenever they appear in message content. With current behavior, users cannot intentionally bypass fixing for a message even if they want to preserve original embed behavior. The feature now also has slash-command configuration, so adding a help-oriented subcommand under the same command group is a natural extension.

This change introduces two user-control improvements:
- A message-level bypass tag: `#nofix`
- A discoverable command note: `/embedfixer nofix`

## Goals / Non-Goals

**Goals:**
- Skip embedfixer processing when message content includes `#nofix`.
- Add a subcommand that explains how to use `#nofix`.
- Keep all existing fix and config functionality unchanged when no tag is present.
- Apply bypass consistently across all currently supported social hosts.

**Non-Goals:**
- Introducing global/per-channel disable switches.
- Changing supported social platform list.
- Adding additional message directives beyond `#nofix`.

## Decisions

1. Evaluate `#nofix` early in message handler before URL parsing and host matching.
   - Rationale: cheapest and clearest control flow; avoids extra parsing/network work.
   - Alternative considered: evaluate after URL parse; rejected as unnecessary overhead and complexity.

2. Use case-insensitive tag detection (`#nofix`, `#NOFIX`) as literal token match.
   - Rationale: better UX with minimal ambiguity.
   - Alternative considered: strict lowercase-only matching; rejected as less user-friendly.

3. Add `/embedfixer nofix` as a lightweight informational subcommand.
   - Returns short usage guidance and example message.
   - Rationale: discoverability and self-documenting behavior in Discord.
   - Alternative considered: only README documentation; rejected because in-app discoverability is needed.

4. Keep command and permission model consistent with existing embedfixer command group.
   - Rationale: no surprise access changes and lower maintenance overhead.

### Sequence Diagram

```text
User             Discord            EmbedFixer Message Handler
 |                  |                         |
 | send msg + URL   |                         |
 |----------------->| MessageCreate          |
 |                  |------------------------>|
 |                  |                         | contains #nofix?
 |                  |                         |----yes----> return (skip fix)
 |                  |                         |----no-----> normal URL fix flow

User             Discord            /embedfixer nofix handler
 |                  |                         |
 | /embedfixer nofix|                         |
 |----------------->| InteractionCreate      |
 |                  |------------------------>|
 |                  |                         | respond with usage note
 |                  |<------------------------|
```

## Risks / Trade-offs

- [False positives when `#nofix` appears in unrelated text] -> Mitigation: treat as intentional explicit override; document behavior in subcommand note.
- [User confusion if bypass is not obvious] -> Mitigation: provide clear command help text with example.
- [Future directive collisions] -> Mitigation: keep directive parsing isolated so future tokens can be added predictably.

## Migration Plan

1. Add `#nofix` check in embedfixer message handler.
2. Add `/embedfixer nofix` informational subcommand.
3. Validate behavior: message with `#nofix` is skipped; message without tag follows existing fix flow.
4. Rollback: remove/disable the bypass check and subcommand to restore prior behavior.

## Open Questions

- Should `#nofix` be required as a standalone word token only, or accepted anywhere in content? (recommended: standalone token match)
