# OpenCode / AI Agent Instructions for gatoc

Discord bot in Go (`discordgo`) using GORM (Postgres).

## Architecture & Conventions

- **Feature-Based Module System:** Features are located in `internal/features/<name>`.
  - Each feature must implement the `bot.Feature` interface (`SlashCommands()`, `Models()`, `RegisterEvents()`).
  - **Manual Registration:** New features MUST be manually instantiated and added to the `features` slice in `cmd/bot/main.go`.
- **Database Models & Migrations:**
  - Each feature owns its GORM models in-package (`internal/features/<name>/...`). `members` models live in a leaf package `internal/features/members/model` to avoid an import cycle with `members/subcommands`.
  - `bot.Feature.Models()` returns the feature's models; `cmd/bot/main.go` aggregates them from all features and passes them to `database.Migrate(db, models...)`. Do NOT hardcode the model list in `internal/database/database.go`.
  - `internal/database` holds only connectivity, the variadic `Migrate` runner, and the `idx_role_user` DDL (acknowledged tech debt — would need per-feature migration hooks to relocate).
- **Slash Commands:** Defined using `bot.SlashCommand` struct.
  - Subcommands are grouped under a root command's `Options` and handled via a `switch` statement in the `Handler`.
- **Event Handling:** Features register handlers via `router.On<Event>(handler)` in their `RegisterEvents` method.

## Testing & Debugging

After ANY code change, ALL of the following MUST pass before considering the work done:

1. **Compile:** `go build -o /dev/null ./...`
2. **Vet:** `go vet ./...`
3. **Format:** `gofmt -l .` must output nothing. Run `gofmt -w .` if it lists files.
4. **Tests:** `go test ./...`

Do NOT claim a task is complete without running the checks above and confirming they pass.

## Development Workflow

- **Local Setup:**
  1. `cp .env.example .env` and fill `TOKEN`, `CLIENT_ID`, `GUILD_ID`.
  2. `docker-compose up` (uses `air` for hot-reloading).
- **Environment:** The bot runs in a containerized environment; ensure `.env` is correctly mapped.
- **Testing:** Tests live next to the code they cover (`*_test.go`). Use `github.com/stretchr/testify` for assertions.

## SDD Workflow (Spec-Driven Development)

This project uses **OpenSpec** for proposal/design/tasks scaffolding and spec deltas, with **co-located living specs** (one per feature module). OpenSpec's `openspec/specs/` management is retired; living specs live at `internal/features/<feature>/spec.md`.

The rhythm for any non-trivial change:

1. **Explore** *(optional)* — `/opsx:explore` as a thinking partner for uncertain changes. Skip for trivial fixes.
2. **Propose** — `/opsx:propose <change-name>` generates `openspec/changes/<change>/` with `proposal.md`, `design.md`, `tasks.md`, and a spec delta at `openspec/changes/<change>/specs/<feature>/spec.md` (ADDED/MODIFIED/REMOVED). **Nothing is implemented until the plan is reviewed and approved.** The delta directory MUST be named after the **feature**, not a capability.
3. **Apply** — implement `tasks.md` items one at a time. **A task may only be marked `[x]` after the 4 verification checks in "Testing & Debugging" pass** (build, vet, gofmt clean, tests green). This makes the SDD evidence-based, not trust-based.
4. **Co-located sync** *(mandatory before archive)* — apply the change's spec delta into the feature's living spec at `internal/features/<feature>/spec.md`:
   - If the feature has no living spec yet (new feature, or first change touching it), create `internal/features/<feature>/spec.md` from the delta's `ADDED` section.
   - If it exists, merge the delta by hand: apply `ADDED`/`MODIFIED` requirements and scenarios, remove `REMOVED` ones, keep one `## Purpose` section. The result is the new post-change reality of the feature.
   - This step replaces OpenSpec's auto-merge (which targeted the retired `openspec/specs/`). It is the mandatory control that keeps co-located specs from going stale.
5. **Archive** — `/opsx:archive <change> --skip-specs`. The `--skip-specs` flag is REQUIRED (OpenSpec's merge would target the non-existent `openspec/specs/`). The archived `openspec/changes/archive/<date>-<change>/specs/<feature>/spec.md` delta remains as provenance — the diff that produced this change.

### OpenSpec config
`openspec/config.yaml` injects project context (tech stack, model-ownership convention, co-located-spec rule) and per-artifact rules into every generated artifact. Keep `context:` accurate; it shapes every proposal/design/spec the AI writes.

### Don't backfill
Only write a feature's living spec when a real change touches it. Do NOT spec features that haven't been changed (e.g. `members`, `parrot`, `ping` have no spec today — they get one on their next real change). Backfilling spec for untouched code rots.

## Entrypoints & Infrastructure

- **Main App:** `cmd/bot/main.go`
- **Deployment:** `docker-stack.yaml` and `deploy-secrets.sh` for Docker Swarm.
- **CI/CD:** `.github/workflows/docker-publish.yml`.
