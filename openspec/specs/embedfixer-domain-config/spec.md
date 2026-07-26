## Purpose
Define configurable replacement-domain behavior for embedfixer across supported social platforms.

## Requirements

### Requirement: Configurable replacement domains for supported platforms
The system SHALL allow configuring replacement domains for currently supported embedfixer social platforms only.

#### Scenario: Set custom domain for supported platform
- **GIVEN** an administrator selects a supported platform and provides a valid domain value
- **WHEN** they execute the config set command
- **THEN** the system SHALL persist the custom replacement domain for that platform
- **AND** future fixed links for that platform SHALL use the custom replacement domain

### Requirement: Hardcoded social-platform scope
The system SHALL support only the currently hardcoded social platforms and SHALL NOT allow custom platform creation.

#### Scenario: Reject unsupported platform value
- **GIVEN** a user provides a platform not in the supported platform list
- **WHEN** they execute a config mutation command
- **THEN** the system SHALL reject the request
- **AND** the response SHALL indicate the platform is unsupported

### Requirement: Multiple source hosts per platform
The system SHALL map all configured source-host aliases for a supported platform to the same domain configuration.

#### Scenario: Twitter and X hosts use shared platform config
- **GIVEN** `twitter.com` and `x.com` are source-host aliases for the same supported platform
- **WHEN** embedfixer processes either host in a message URL
- **THEN** the system SHALL resolve both hosts to the same platform configuration
- **AND** both fixed links SHALL use the same active replacement domain

### Requirement: Config visibility with default/custom status
The system SHALL expose a config display command that reports active domain mode per supported platform.

#### Scenario: Show command reports default and custom status
- **GIVEN** at least one platform has a custom domain and another uses the default domain
- **WHEN** a user executes the config show command
- **THEN** the response SHALL list each supported platform
- **AND** each platform entry SHALL include source-host aliases, active replacement domain, and mode (`default` or `custom`)

#### Scenario: Nofix help note is available
- **GIVEN** a user needs to bypass embedfixer for a specific message
- **WHEN** they execute the `embedfixer nofix` subcommand
- **THEN** the response SHALL explain that adding `#nofix` to message content skips fix processing
- **AND** the response SHALL include at least one example usage note

### Requirement: Default fallback behavior
The system SHALL fallback to default replacement domains when no custom domain exists for a platform.

#### Scenario: Platform without custom override
- **GIVEN** no custom replacement domain is configured for a supported platform
- **WHEN** embedfixer processes a URL that maps to that platform
- **THEN** the system SHALL generate the fixed link using the platform's default replacement domain

#### Scenario: Message bypasses fixing with #nofix
- **GIVEN** a message contains a fixable supported social URL
- **AND** message content includes the `#nofix` tag
- **WHEN** embedfixer handles the message
- **THEN** the system SHALL skip embed suppression and skip posting a fixed-link message
