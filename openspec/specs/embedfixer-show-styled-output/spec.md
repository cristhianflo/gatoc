## Purpose

Define the expected styled output behavior for `embedfixer show` so configuration is easier to scan while preserving existing data semantics.

## Requirements

### Requirement: Table-like embed formatting for show output
The system SHALL render `embedfixer show` output in a structured, table-like embed-field layout optimized for scanning.

#### Scenario: Show command returns styled layout
- **GIVEN** a user executes the `embedfixer show` command
- **WHEN** the system builds the embed response
- **THEN** the response SHALL present platform data using consistent field-group formatting
- **AND** each platform section SHALL follow the same visual structure

### Requirement: Preserve config information in styled output
The system SHALL preserve all current config information when applying the new styling.

#### Scenario: Data parity between old and new layouts
- **GIVEN** platform configuration data exists for supported platforms
- **WHEN** the styled show response is generated
- **THEN** each platform entry SHALL include source hosts, active domain, and mode (`default` or `custom`)
- **AND** no previously shown data point SHALL be removed

### Requirement: Deterministic platform ordering
The system SHALL keep a deterministic platform order in show output to support quick visual comparison.

#### Scenario: Stable row ordering across calls
- **GIVEN** a user runs `embedfixer show` multiple times with unchanged configuration
- **WHEN** each response is rendered
- **THEN** platform sections SHALL appear in the same order in every response
