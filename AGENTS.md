# OpenCode / AI Agent Instructions for gatoc

Discord bot in Go (`discordgo`) using GORM (Postgres).

## Architecture & Conventions

- **Feature-Based Module System:** Features are located in `internal/features/<name>`.
  - Each feature must implement the `bot.Feature` interface (`SlashCommands()`, `Models()`, `RegisterEvents()`).
  - **Manual Registration:** New features MUST be manually instantiated and added to the `features` slice in `cmd/bot/main.go`.
- **Database Models & Migrations:** 
  - Centralized in `internal/database/database.go`.
  - Migrations are handled via `gorm.AutoMigrate` in the `Migrate()` function. Add new models there.
- **Slash Commands:** Defined using `bot.SlashCommand` struct.
  - Subcommands are grouped under a root command's `Options` and handled via a `switch` statement in the `Handler`.
- **Event Handling:** Features register handlers via `router.On<Event>(handler)` in their `RegisterEvents` method.

## Development Workflow

- **Local Setup:**
  1. `cp .env.example .env` and fill `TOKEN`, `CLIENT_ID`, `GUILD_ID`.
  2. `docker-compose up` (uses `air` for hot-reloading).
- **Environment:** The bot runs in a containerized environment; ensure `.env` is correctly mapped.
- **Testing:** No automated tests currently. `github.com/stretchr/testify` is available if needed.

## Entrypoints & Infrastructure

- **Main App:** `cmd/bot/main.go`
- **Deployment:** `docker-stack.yaml` and `deploy-secrets.sh` for Docker Swarm.
- **CI/CD:** `.github/workflows/docker-publish.yml`.
