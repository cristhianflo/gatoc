# Graph Report - .  (2026-08-05)

## Corpus Check
- 55 files · ~132,789 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 291 nodes · 387 edges · 52 communities (35 shown, 17 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 47 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- Core Runtime and Data
- Slash Command Handling
- Event Router Middleware
- Bot Context Features
- Embed Config Service
- Deployment Infrastructure
- Currency Rates
- Bot Builder
- Embed Message Handler
- Embedfixer Configuration
- Finance Dollar Command
- Local Docker Development
- Embedfixer Feature Slice
- Finance Feature Slice
- Members Feature Slice
- Ping Feature Slice
- External Stack
- Member Data Models
- Graphify Plugin
- Docker Publish Workflow
- Feature Module Contract
- Database Model Ownership
- Cat Asset
- Swarm Secret Deployment
- Embedfixer Domain Model
- Platform Support
- Member Message Types
- Graphify Runtime Plugin
- Co-located Specifications
- Command Event Registration
- Deployment Script
- Embed Suppression
- Dollar Status Response
- Air Hot Reloading
- Application Entrypoint
- Verification Checks
- Go Module Metadata

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

## Communities (52 total, 17 thin omitted)

### Community 0 - "Core Runtime and Data"
Cohesion: 0.11
Nodes (16): bot, main(), BotConfig, Config, DbConfig, RedisConfig, Context, Intent (+8 more)

### Community 1 - "Slash Command Handling"
Cohesion: 0.19
Nodes (20): ApplicationCommand, ApplicationCommandInteractionDataOption, ApplicationCommandOption, SlashCommand, SlashSubcommand, DeferReply(), EditDeferred(), GetInteractionFailedResponse() (+12 more)

### Community 2 - "Event Router Middleware"
Cohesion: 0.16
Nodes (10): EventRouter, GuildMemberAddMiddleware, GuildMemberRemoveMiddleware, InteractionCreateMiddleware, MessageCreateMiddleware, MessageReactionAddMiddleware, ReadyMiddleware, NewEventRouter() (+2 more)

### Community 3 - "Bot Context Features"
Cohesion: 0.15
Nodes (13): BotContext, GuildMemberAdd, GuildMemberRemove, Client, DB, Logger, MembersFeature, MessageCreate (+5 more)

### Community 4 - "Embed Config Service"
Cohesion: 0.20
Nodes (16): platformConfig, activeDomain(), getCustomDomain(), DB, normalizeDomainHost(), resetCustomDomain(), resolveReplacementDomain(), setCustomDomain() (+8 more)

### Community 5 - "Deployment Infrastructure"
Cohesion: 0.11
Nodes (18): Build and Push Image Job, Deploy Job, GitHub Container Registry, Latest Docker Image Tag, Commit SHA Docker Image Tag, Docker Stack Deploy Action, Docker Swarm Deployment, Bot Overlay Network (+10 more)

### Community 6 - "Currency Rates"
Cohesion: 0.21
Nodes (12): fetchRates(), formatRateFields(), Client, T, TestFetchRatesArrayResponse(), TestFetchRatesFailureResponse(), TestFetchRatesSingleObjectResponse(), TestFormatRateFields() (+4 more)

### Community 7 - "Bot Builder"
Cohesion: 0.20
Nodes (7): BotBuilder, Feature, Client, DB, Intent, Logger, NewBotBuilder()

### Community 8 - "Embed Message Handler"
Cohesion: 0.23
Nodes (9): buildFixedEmbedMessage(), fixedURLFromDomain(), EmbedFixerFeature, MessageCreate, Session, hasNoFixTag(), T, TestEmbedFixerCommandIncludesNoFixSubcommand() (+1 more)

### Community 9 - "Embedfixer Configuration"
Cohesion: 0.25
Nodes (9): Embedfixer Config Set Command, Embedfixer Config Show Command, Custom Replacement Domain, Default Replacement Domain, Deterministic Platform Ordering, Replacement Domain Configuration, Source Host Aliases, Table-Like Embed Format (+1 more)

### Community 10 - "Finance Dollar Command"
Cohesion: 0.25
Nodes (9): Currency Context, /dollar all Command, EUR Exchange-Rate Data, finance Feature, Partial Upstream Failure Tolerance, Source and Freshness Transparency, Update Timestamp, Upstream Currency Endpoints (+1 more)

### Community 11 - "Local Docker Development"
Cohesion: 0.25
Nodes (8): App Bind Mount, Compose Bot Service, Compose Database Volume, Compose Database Service, Development Build Target, discordbot Database, .env File, Postgres Latest Image

### Community 16 - "External Stack"
Cohesion: 0.50
Nodes (4): Discord Bot, discordgo, GORM, Postgres

### Community 17 - "Member Data Models"
Cohesion: 0.67
Nodes (3): Model, ResponseMessage, WelcomeRole

### Community 18 - "Graphify Plugin"
Cohesion: 0.50
Nodes (3): plugin, $schema, .opencode/plugins/graphify.js

### Community 19 - "Docker Publish Workflow"
Cohesion: 0.67
Nodes (3): GitHub Actions Docker Publish Workflow, Main Branch Push Trigger, gatoc Project

### Community 20 - "Feature Module Contract"
Cohesion: 0.67
Nodes (3): bot.Feature Interface, Feature-Based Module System, Manual Feature Registration

### Community 21 - "Database Model Ownership"
Cohesion: 0.67
Nodes (3): Database Migrate, Feature-Owned GORM Models, Members Leaf Model Package

### Community 22 - "Cat Asset"
Cohesion: 0.67
Nodes (3): Cat, Close-up Cat Face, Cat Reaction Image

### Community 25 - "Platform Support"
Cohesion: 0.67
Nodes (3): embedfixer Feature, Hardcoded Platform Scope, Supported Social Platforms

## Knowledge Gaps
- **57 isolated node(s):** `$schema`, `.opencode/plugins/graphify.js`, `deploy-secrets.sh script`, `SECRETS`, `deploy.sh script` (+52 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `BotContext` connect `Bot Context Features` to `Core Runtime and Data`, `Slash Command Handling`, `Embed Message Handler`?**
  _High betweenness centrality (0.165) - this node is a cross-community bridge._
- **Why does `SlashCommand` connect `Slash Command Handling` to `Core Runtime and Data`, `Event Router Middleware`, `Bot Context Features`, `Bot Builder`, `Embedfixer Feature Slice`, `Finance Feature Slice`, `Members Feature Slice`, `Ping Feature Slice`?**
  _High betweenness centrality (0.095) - this node is a cross-community bridge._
- **Why does `bot` connect `Core Runtime and Data` to `Slash Command Handling`, `Event Router Middleware`, `Bot Context Features`?**
  _High betweenness centrality (0.079) - this node is a cross-community bridge._
- **What connects `$schema`, `.opencode/plugins/graphify.js`, `deploy-secrets.sh script` to the rest of the system?**
  _57 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Core Runtime and Data` be split into smaller, more focused modules?**
  _Cohesion score 0.11396011396011396 - nodes in this community are weakly interconnected._
- **Should `Bot Context Features` be split into smaller, more focused modules?**
  _Cohesion score 0.14736842105263157 - nodes in this community are weakly interconnected._
- **Should `Deployment Infrastructure` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._