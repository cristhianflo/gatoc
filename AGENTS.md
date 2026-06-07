# OpenCode / AI Agent Instructions for gatoc

This repository is a Discord bot written in Go, using `discordgo` and `gorm` (Postgres).

## Architecture & Conventions

- **Side-effect Registration (CRITICAL):** Commands and Events are registered dynamically using `init()` functions.
  - `cmd/bot/main.go` uses anonymous imports (`_ "github.com/bachacode/gatoc/internal/commands"`, etc.) to trigger these.
  - To add a new command, create a `.go` file in `internal/commands/` and call `bot.RegisterCommand(...)` inside an `init()` block.
  - To add a new event, create a `.go` file in `internal/events/` and call `bot.RegisterEvent(...)` inside an `init()` block.
- **Subcommands:** Handled manually by root commands. Look at `internal/commands/dolar.go` and `internal/commands/dolar/` for the pattern. The root command defines the subcommand options and uses a `switch` statement in its `Handler` to route to the subpackage's handler.
- **Database Migrations:** No external migration files. Migrations are handled via `gorm.AutoMigrate` in `internal/database/database.go` when the bot starts up. New models must be added to the `Migrate()` function there.

## Development Workflow

- **Local Environment:** The project is containerized for development.
  1. Copy `.env.example` to `.env` and fill in required values (especially `TOKEN`, `CLIENT_ID`, `GUILD_ID`).
  2. Run `docker-compose up`.
- **Hot Reloading:** The `docker-compose` setup uses `air` (configured in `.air.toml`) to automatically recompile and restart the Go binary on file changes.
- **Testing:** There are currently no automated tests in this repository, though `github.com/stretchr/testify` is in `go.mod`. Do not attempt to run tests unless specifically instructed to add them.

## Entrypoints

- **Main App:** `cmd/bot/main.go`
- **Infrastructure / Deploy:** `docker-stack.yaml` and `deploy-secrets.sh` (used for Docker Swarm deployments).
