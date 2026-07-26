## Purpose

Define expected behavior for reporting multi-currency (USD and EUR) exchange-rate data in the finance command output, including resilience and transparency requirements.

## Requirements

### Requirement: Multi-currency rate response in `/dollar all`
The system SHALL include both USD and EUR exchange-rate information in the `/dollar all` command response when upstream data is available.

#### Scenario: USD and EUR are both available
- **GIVEN** the USD and EUR upstream endpoints both return valid responses
- **WHEN** a user executes `/dollar all`
- **THEN** the response SHALL include USD rate entries and EUR rate entries in the same interaction response
- **AND** each entry SHALL include identifiable currency context

### Requirement: Partial upstream failure tolerance
The system SHALL return available exchange-rate data even when one upstream currency source fails.

#### Scenario: EUR source fails and USD source succeeds
- **GIVEN** the USD upstream endpoint returns valid data
- **AND** the EUR upstream endpoint is unavailable or returns invalid data
- **WHEN** a user executes `/dollar all`
- **THEN** the response SHALL include USD rate data
- **AND** the response SHALL indicate EUR data is temporarily unavailable

#### Scenario: USD source fails and EUR source succeeds
- **GIVEN** the EUR upstream endpoint returns valid data
- **AND** the USD upstream endpoint is unavailable or returns invalid data
- **WHEN** a user executes `/dollar all`
- **THEN** the response SHALL include EUR rate data
- **AND** the response SHALL indicate USD data is temporarily unavailable

### Requirement: Source and freshness transparency
The system SHALL present source and update-time context for displayed currency data.

#### Scenario: Display metadata for each currency dataset
- **GIVEN** one or more currency datasets are available from upstream APIs
- **WHEN** a user executes `/dollar all`
- **THEN** the response SHALL include data-source labeling for displayed datasets
- **AND** the response SHALL include each dataset's update timestamp using a consistent display format
