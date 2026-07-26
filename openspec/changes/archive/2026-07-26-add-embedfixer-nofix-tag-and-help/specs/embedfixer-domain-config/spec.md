## MODIFIED Requirements

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
