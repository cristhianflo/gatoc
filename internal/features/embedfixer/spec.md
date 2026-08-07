## Purpose

Define the behavior of the embedfixer feature: link fixing across supported social platforms, configurable replacement domains, the `show` configuration display, and the `nofix` bypass toggle. This is the living specification for the embedfixer feature module.

## Responsibilities

- Fix links from supported social platforms so Discord can display improved embeds.
- Manage replacement-domain configuration for supported platforms.
- Provide configuration display and the `#nofix` bypass behavior.

## Rules

Administrators can configure replacement domains for supported embedfixer social platforms only.

### Scenario: Set custom domain for supported platform

- **GIVEN** an administrator selects a supported platform and provides a valid domain value
- **WHEN** they execute the config set command
- **THEN** the custom replacement domain is persisted for that platform
- **AND** future fixed links for that platform use the custom replacement domain

Only the currently hardcoded social platforms are supported; users cannot create custom platforms.

### Scenario: Reject unsupported platform value

- **GIVEN** a user provides a platform not in the supported platform list
- **WHEN** they execute a config mutation command
- **THEN** the request is rejected
- **AND** the response indicates that the platform is unsupported

All configured source-host aliases for a supported platform use the same domain configuration.

### Scenario: Twitter and X hosts use shared platform config

- **GIVEN** `twitter.com` and `x.com` are source-host aliases for the same supported platform
- **WHEN** embedfixer processes either host in a message URL
- **THEN** both hosts resolve to the same platform configuration
- **AND** both fixed links use the same active replacement domain

The config display command reports the active domain mode for each supported platform.

### Scenario: Show command reports default and custom status

- **GIVEN** at least one platform has a custom domain and another uses the default domain
- **WHEN** a user executes the config show command
- **THEN** the response lists each supported platform
- **AND** each platform entry includes source-host aliases, active replacement domain, and mode (`default` or `custom`)

### Scenario: Nofix help note is available

- **GIVEN** a user needs to bypass embedfixer for a specific message
- **WHEN** they execute the `embedfixer nofix` subcommand
- **THEN** the response explains that adding `#nofix` to message content skips fix processing
- **AND** the response includes at least one example usage note

The feature falls back to default replacement domains when no custom domain exists for a platform.

### Scenario: Platform without custom override

- **GIVEN** no custom replacement domain is configured for a supported platform
- **WHEN** embedfixer processes a URL that maps to that platform
- **THEN** the fixed link uses the platform's default replacement domain

### Scenario: Message bypasses fixing with #nofix

- **GIVEN** a message contains a fixable supported social URL
- **AND** message content includes the `#nofix` tag
- **WHEN** embedfixer handles the message
- **THEN** embed suppression is skipped and no fixed-link message is posted

The `embedfixer show` command renders output in a structured, table-like embed-field layout optimized for scanning.

### Scenario: Show command returns styled layout

- **GIVEN** a user executes the `embedfixer show` command
- **WHEN** the system builds the embed response
- **THEN** the response presents platform data using consistent field-group formatting
- **AND** each platform section follows the same visual structure

Styled output preserves all current configuration information.

### Scenario: Data parity between old and new layouts

- **GIVEN** platform configuration data exists for supported platforms
- **WHEN** the styled show response is generated
- **THEN** each platform entry includes source hosts, active domain, and mode (`default` or `custom`)
- **AND** no previously shown data point is removed

The `show` output keeps a deterministic platform order to support quick visual comparison.

### Scenario: Stable row ordering across calls

- **GIVEN** a user runs `embedfixer show` multiple times with unchanged configuration
- **WHEN** each response is rendered
- **THEN** platform sections appear in the same order in every response
