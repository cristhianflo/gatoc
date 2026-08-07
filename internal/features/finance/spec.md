## Purpose

Define expected behavior for reporting multi-currency (USD and EUR) exchange-rate data in the finance command output, including resilience and transparency requirements. This is the living specification for the finance feature module.

## Responsibilities

- Report USD and EUR exchange-rate data through the `/dollar all` command.
- Preserve available data when one upstream source fails.
- Make each displayed dataset's source and freshness visible.

## Rules

When upstream data is available, `/dollar all` includes both USD and EUR exchange-rate information.

### Scenario: USD and EUR are both available
- **GIVEN** the USD and EUR upstream endpoints both return valid responses
- **WHEN** a user executes `/dollar all`
- **THEN** the response includes USD and EUR rate entries in the same interaction response
- **AND** each entry includes identifiable currency context

The command returns available exchange-rate data even when one upstream currency source fails.

### Scenario: EUR source fails and USD source succeeds
- **GIVEN** the USD upstream endpoint returns valid data
- **AND** the EUR upstream endpoint is unavailable or returns invalid data
- **WHEN** a user executes `/dollar all`
- **THEN** the response includes USD rate data
- **AND** the response indicates that EUR data is temporarily unavailable

### Scenario: USD source fails and EUR source succeeds
- **GIVEN** the EUR upstream endpoint returns valid data
- **AND** the USD upstream endpoint is unavailable or returns invalid data
- **WHEN** a user executes `/dollar all`
- **THEN** the response includes EUR rate data
- **AND** the response indicates that USD data is temporarily unavailable

Displayed currency data includes source and update-time context.

### Scenario: Display metadata for each currency dataset
- **GIVEN** one or more currency datasets are available from upstream APIs
- **WHEN** a user executes `/dollar all`
- **THEN** the response includes data-source labeling for displayed datasets
- **AND** the response includes each dataset's update timestamp using a consistent display format
