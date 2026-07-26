## Why

The embedfixer currently hardcodes replacement domains in handlers, making it difficult to adapt when a fix-domain provider degrades or policy changes are needed per server. Introducing controlled configuration now improves operational flexibility while preserving the current fixed scope of supported social platforms.

## What Changes

- Add embedfixer configuration for replacement domains per currently supported social media platform (Twitter/X, Reddit, Instagram).
- Keep supported social media hardcoded (no custom platform creation), while allowing custom replacement-domain values per supported platform.
- Support multiple source hosts per social platform (for example, Twitter includes both `twitter.com` and `x.com`) mapped to one platform configuration.
- Add slash commands to inspect and manage embedfixer domain configuration, including a command to show active config and whether each platform is using default or custom values.
- Preserve current replacement domains as defaults:
  - Twitter/X -> `fxtwitter.com`
  - Reddit -> `vxreddit.com`
  - Instagram -> `kkinstagram.com`
- Include rollback plan: disable custom-config lookup and force fallback to current defaults if runtime issues are observed after deployment.

## Capabilities

### New Capabilities
- `embedfixer-domain-config`: Configure and inspect replacement domains per supported social platform, with default/custom visibility.

### Modified Capabilities
- None.

## Impact

- Affected code:
  - `internal/features/embedfixer/handlers.go`
  - `internal/features/embedfixer/feature.go`
  - New embedfixer command handlers under `internal/features/embedfixer/`
  - `internal/database/database.go` for migration registration (if persistent model is added)
- Data model: likely new persistent config table keyed by guild and supported platform.
- User-facing behavior:
  - New config commands for viewing/updating/resetting replacement domains.
  - Existing embed fixer output now resolves replacement domain from custom config when set, otherwise defaults.
- Operational behavior: domain routing changes can be applied without code redeploy.
