# Graph Report - gatoc  (2026-08-07)

## Corpus Check
- 52 files · ~133,003 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 301 nodes · 396 edges · 54 communities (36 shown, 18 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 47 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b5a21842`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- BotBuilder
- SlashCommand
- EventRouter
- BotContext
- resolveReplacementDomain
- Stack Bot Service
- fetchRates
- bot
- .EmbedFixerHandler
- Embedfixer Config Show Command
- /dollar all Command
- Compose Bot Service
- EmbedFixerFeature
- FinanceFeature
- MembersFeature
- PingFeature
- Discord Bot
- model/models.go
- opencode.json
- GitHub Actions Docker Publish Workflow
- Feature-Based Module System
- Feature-Owned GORM Models
- Cat
- deploy-secrets.sh
- EmbedFixerDomainOverride
- Supported Social Platforms
- LeaveMessage
- graphify.js
- Co-Located Living Specs
- Router Event Registration
- deploy.sh
- Embed Suppression
- dollar_status.go
- air Hot Reloading
- Main App Entrypoint
- Required Verification Checks
- github.com/bachacode/gatoc
- Name
- ParrotFeature

## God Nodes (most connected - your core abstractions)
1. `EventRouter` - 21 edges
2. `BotContext` - 18 edges
3. `BotBuilder` - 15 edges
4. `SlashCommand` - 14 edges
5. `bot` - 13 edges
6. `handleConfigSet()` - 12 edges
7. `handleConfigReset()` - 11 edges
8. `GetInteractionFailedResponse()` - 9 edges
9. `handleConfigShow()` - 9 edges
10. `main()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `NewBotBuilder()`  [INFERRED]
  cmd/bot/main.go → internal/bot/botBuilder.go
- `main()` --calls--> `LoadConfig()`  [INFERRED]
  cmd/bot/main.go → internal/config/config.go
- `main()` --calls--> `Migrate()`  [INFERRED]
  cmd/bot/main.go → internal/database/database.go
- `main()` --calls--> `New()`  [INFERRED]
  cmd/bot/main.go → internal/database/database.go
- `main()` --calls--> `NewRedis()`  [INFERRED]
  cmd/bot/main.go → internal/database/redis.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Embedfixer Show Configuration Data** — internal_features_embedfixer_spec_config_show_command, internal_features_embedfixer_spec_source_host_aliases, internal_features_embedfixer_spec_default_replacement_domain, internal_features_embedfixer_spec_custom_replacement_domain [EXTRACTED 1.00]
- **Finance Currency Response** — internal_features_finance_spec_dollar_all_command, internal_features_finance_spec_usd_exchange_rate_data, internal_features_finance_spec_eur_exchange_rate_data, internal_features_finance_spec_currency_context [EXTRACTED 1.00]
- **gatoc Runtime Deployment Stack** — docker_stack_yaml_bot_service, docker_stack_yaml_stack_db_service, docker_stack_yaml_bot_network, docker_stack_yaml_external_swarm_secrets [EXTRACTED 1.00]

## Communities (54 total, 18 thin omitted)

### Community 0 - "BotBuilder"
Cohesion: 0.10
Nodes (19): BotBuilder, Feature, main(), BotConfig, Config, DbConfig, RedisConfig, Client (+11 more)

### Community 1 - "SlashCommand"
Cohesion: 0.19
Nodes (20): ApplicationCommand, ApplicationCommandInteractionDataOption, ApplicationCommandOption, SlashCommand, SlashSubcommand, DeferReply(), EditDeferred(), GetInteractionFailedResponse() (+12 more)

### Community 2 - "EventRouter"
Cohesion: 0.26
Nodes (8): EventRouter, GuildMemberAddMiddleware, GuildMemberRemoveMiddleware, InteractionCreateMiddleware, MessageCreateMiddleware, MessageReactionAddMiddleware, ReadyMiddleware, NewEventRouter()

### Community 3 - "BotContext"
Cohesion: 0.16
Nodes (13): BotContext, GuildMemberAdd, GuildMemberRemove, Client, DB, Logger, MembersFeature, MessageCreate (+5 more)

### Community 4 - "resolveReplacementDomain"
Cohesion: 0.20
Nodes (16): platformConfig, activeDomain(), getCustomDomain(), DB, normalizeDomainHost(), resetCustomDomain(), resolveReplacementDomain(), setCustomDomain() (+8 more)

### Community 5 - "Stack Bot Service"
Cohesion: 0.11
Nodes (18): Build and Push Image Job, Deploy Job, GitHub Container Registry, Latest Docker Image Tag, Commit SHA Docker Image Tag, Docker Stack Deploy Action, Docker Swarm Deployment, Bot Overlay Network (+10 more)

### Community 6 - "fetchRates"
Cohesion: 0.21
Nodes (12): fetchRates(), formatRateFields(), Client, T, TestFetchRatesArrayResponse(), TestFetchRatesFailureResponse(), TestFetchRatesSingleObjectResponse(), TestFormatRateFields() (+4 more)

### Community 7 - "bot"
Cohesion: 0.27
Nodes (4): bot, Context, Intent, Session

### Community 8 - ".EmbedFixerHandler"
Cohesion: 0.23
Nodes (9): buildFixedEmbedMessage(), fixedURLFromDomain(), EmbedFixerFeature, MessageCreate, Session, hasNoFixTag(), T, TestEmbedFixerCommandIncludesNoFixSubcommand() (+1 more)

### Community 9 - "Embedfixer Config Show Command"
Cohesion: 0.25
Nodes (9): Embedfixer Config Set Command, Embedfixer Config Show Command, Custom Replacement Domain, Default Replacement Domain, Deterministic Platform Ordering, Replacement Domain Configuration, Source Host Aliases, Table-Like Embed Format (+1 more)

### Community 10 - "/dollar all Command"
Cohesion: 0.25
Nodes (9): Currency Context, /dollar all Command, EUR Exchange-Rate Data, finance Feature, Partial Upstream Failure Tolerance, Source and Freshness Transparency, Update Timestamp, Upstream Currency Endpoints (+1 more)

### Community 11 - "Compose Bot Service"
Cohesion: 0.25
Nodes (8): App Bind Mount, Compose Bot Service, Compose Database Volume, Compose Database Service, Development Build Target, discordbot Database, .env File, Postgres Latest Image

### Community 16 - "Discord Bot"
Cohesion: 0.50
Nodes (4): Discord Bot, discordgo, GORM, Postgres

### Community 17 - "model/models.go"
Cohesion: 0.67
Nodes (3): Model, ResponseMessage, WelcomeRole

### Community 18 - "opencode.json"
Cohesion: 0.50
Nodes (3): plugin, $schema, .opencode/plugins/graphify.js

### Community 19 - "GitHub Actions Docker Publish Workflow"
Cohesion: 0.67
Nodes (3): GitHub Actions Docker Publish Workflow, Main Branch Push Trigger, gatoc Project

### Community 20 - "Feature-Based Module System"
Cohesion: 0.67
Nodes (3): bot.Feature Interface, Feature-Based Module System, Manual Feature Registration

### Community 21 - "Feature-Owned GORM Models"
Cohesion: 0.67
Nodes (3): Database Migrate, Feature-Owned GORM Models, Members Leaf Model Package

### Community 22 - "Cat"
Cohesion: 0.67
Nodes (3): Cat, Close-up Cat Face, Cat Reaction Image

### Community 25 - "Supported Social Platforms"
Cohesion: 0.67
Nodes (3): embedfixer Feature, Hardcoded Platform Scope, Supported Social Platforms

### Community 52 - "Name"
Cohesion: 0.20
Nodes (9): Boundaries, Dependencies, Name, Purpose, Responsibilities, Rules, Scenario: Descriptive behavior, Scenarios (+1 more)

## Knowledge Gaps
- **64 isolated node(s):** `$schema`, `.opencode/plugins/graphify.js`, `deploy-secrets.sh script`, `SECRETS`, `deploy.sh script` (+59 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **18 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `BotContext` connect `BotContext` to `BotBuilder`, `SlashCommand`, `.EmbedFixerHandler`, `bot`?**
  _High betweenness centrality (0.154) - this node is a cross-community bridge._
- **Why does `SlashCommand` connect `SlashCommand` to `BotBuilder`, `BotContext`, `bot`, `EmbedFixerFeature`, `FinanceFeature`, `MembersFeature`, `PingFeature`, `ParrotFeature`?**
  _High betweenness centrality (0.089) - this node is a cross-community bridge._
- **Why does `bot` connect `bot` to `BotBuilder`, `SlashCommand`, `EventRouter`, `BotContext`?**
  _High betweenness centrality (0.074) - this node is a cross-community bridge._
- **What connects `$schema`, `.opencode/plugins/graphify.js`, `deploy-secrets.sh script` to the rest of the system?**
  _64 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `BotBuilder` be split into smaller, more focused modules?**
  _Cohesion score 0.09879032258064516 - nodes in this community are weakly interconnected._
- **Should `Stack Bot Service` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._