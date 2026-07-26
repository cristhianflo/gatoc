## Context

The embedfixer feature currently listens to message-create events, detects URLs for a hardcoded set of social hosts, suppresses the original embed, and posts a replacement message using fixed destination domains. Today, replacement domains are hardcoded in handler functions, which requires code changes and redeploys whenever domain routing policy needs to change.

The requested behavior is to keep the supported social media list fixed to current platforms, while allowing per-platform domain customization and visibility into whether each platform is using default or custom configuration.

Current supported platforms and defaults:
- Twitter/X: source hosts `twitter.com`, `x.com`; default replacement domain `fxtwitter.com`
- Reddit: source host `reddit.com`; default replacement domain `vxreddit.com`
- Instagram: source host `instagram.com`; default replacement domain `kkinstagram.com`

Stakeholders:
- Guild moderators/admins who need to adjust replacement domains.
- End users who receive fixed embed links.
- Bot maintainers responsible for runtime reliability.

## Goals / Non-Goals

**Goals:**
- Introduce configurable replacement domains for currently supported social platforms.
- Preserve hardcoded social-platform scope (no user-created platforms).
- Support multiple source hosts mapping to one platform configuration (for example, `twitter.com` + `x.com`).
- Add command(s) to show current configuration with explicit default/custom status per platform.

**Non-Goals:**
- Supporting arbitrary new social media platforms.
- Replacing message-event architecture with slash-command-only workflows.
- Building a full admin web UI for configuration.

## Decisions

1. Define a hardcoded platform registry as source-of-truth.
   - Each entry includes: stable platform key, source-host aliases, default replacement domain, display label.
   - Rationale: guarantees fixed social scope while still enabling configurable domains.
   - Alternative considered: dynamic platform creation by admins; rejected because requirements explicitly prohibit custom social-media additions.

2. Persist only custom replacement-domain overrides.
   - Storage model: `guild_id`, `platform`, `custom_domain`, timestamps.
   - Resolution logic: if custom override exists use it; otherwise use registry default.
   - Rationale: keeps defaults centralized and simplifies fallback/rollback behavior.
   - Alternative considered: persisting full config rows for all platforms by default; rejected due to unnecessary storage and migration complexity.

3. Add an embedfixer config command group under a new slash command.
   - Proposed subcommands:
     - `show`: list all supported platforms, source hosts, active domain, and status (`default` or `custom`).
     - `set`: set custom domain for a supported platform.
     - `reset`: remove custom domain and revert to default.
   - Rationale: operational control without code edits and explicit visibility of active behavior.
   - Alternative considered: event-only commands via message syntax; rejected because slash commands provide safer validation and discoverability.

4. Use host normalization and strict domain validation.
   - Normalize host values to lowercase and strip `www.` before matching.
   - Validate `set` input is a host-like domain without path/query fragments.
   - Rationale: prevents malformed links and reduces accidental misconfiguration.

### Sequence Diagram

```text
Admin User         Discord          Bot Command Handler      Config Store
   |                  |                    |                     |
   | /embedfixer ...  |                    |                     |
   |----------------->| InteractionCreate  |                     |
   |                  |------------------->| validate input      |
   |                  |                    | read/update config   |
   |                  |                    |--------------------->|
   |                  |                    |<---------------------|
   |                  |<-------------------| command response      |

Regular User       Discord          EmbedFixer Event        Config Resolver
   |                  |                    |                     |
   | post social URL  |                    |                     |
   |----------------->| MessageCreate      |                     |
   |                  |------------------->| parse host           |
   |                  |                    | resolve platform     |
   |                  |                    | resolve active domain|
   |                  |                    |--------------------->|
   |                  |                    |<---------------------|
   |                  |                    | suppress + post fix  |
```

## Risks / Trade-offs

- [Misconfigured custom domain causes broken links] -> Mitigation: strict domain validation and clear reset-to-default command.
- [More DB/config reads on message flow] -> Mitigation: optional cache layer per guild+platform if latency becomes noticeable.
- [Permission failures for embed suppression remain possible] -> Mitigation: preserve existing fallback behavior and surface failures in logs.
- [Configuration UX ambiguity] -> Mitigation: `show` always reports both active domain and whether it is default/custom.

## Migration Plan

1. Add platform registry and config model.
2. Register migration for new model.
3. Add slash commands (`show`, `set`, `reset`).
4. Update embedfixer event flow to resolve active domain from config service with default fallback.
5. Deploy and verify config commands + event behavior in a test guild.
6. Rollback strategy: disable custom lookup and force registry defaults for all platforms.

## Open Questions

- Should config scope be per guild (recommended) or global across all guilds?
- Should non-admin users be allowed to run `show`, or restricted to management permissions only?
