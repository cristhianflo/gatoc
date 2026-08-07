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

This project uses simple, co-located Markdown specs for features and architectural subsystems. An architectural subsystem is a shared core package such as `internal/bot`, `internal/database`, or `internal/config`. A feature is a user-facing capability under `internal/features/<name>`.

Specs live beside the package they describe:

- Architectural subsystem: `internal/<package>/spec.md`
- Feature: `internal/features/<name>/spec.md`

Each package has at most one `spec.md`. Use `docs/spec-template.md` when creating a new spec. Specs should describe responsibilities, observable behavior, constraints, boundaries, and important dependencies in clear natural language. EARS keywords such as `SHALL` are not required.

For concrete behavior, edge cases, and failure behavior, prefer readable `GIVEN`, `WHEN`, and `THEN` scenarios. Use scenarios when they clarify behavior or provide a useful basis for tests; do not force every statement into a scenario.

The rhythm for a non-trivial change:

1. Identify the affected feature or architectural subsystem.
2. Read its existing `spec.md`, or create one from `docs/spec-template.md` if the package does not have one.
3. Update the spec before implementation when the change introduces or alters behavior. Review the plan before coding when the change is substantial.
4. Implement the change against the documented behavior.
5. Update the spec in the same change if implementation reality differs from the plan. Add or update scenarios for important user-visible behavior, edge cases, and failures.
6. Run all checks in **Testing & Debugging**. Do not claim completion until the code, spec, and verification results are consistent.

Formatting-only changes, behavior-preserving refactors, and test-only changes may skip spec updates. Do not backfill specs for untouched packages.

## Entrypoints & Infrastructure

- **Main App:** `cmd/bot/main.go`
- **Deployment:** `docker-stack.yaml` and `deploy-secrets.sh` for Docker Swarm.
- **CI/CD:** `.github/workflows/docker-publish.yml`.

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

When the user types `/graphify`, use the installed graphify skill or instructions before doing anything else.

Rules:
- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- Dirty graphify-out/ files are expected after hooks or incremental updates; dirty graph files are not a reason to skip graphify. Only skip graphify if the task is about stale or incorrect graph output, or the user explicitly says not to use it.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
