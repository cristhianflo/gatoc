## Why

The finance feature currently reports only USD-based exchange rates, which limits usefulness for users who also track EUR values. Adding EUR coverage now improves command value without introducing a new command surface.

## What Changes

- Extend the existing `/dollar all` behavior to include Euro exchange rate information alongside current USD rates.
- Retrieve and merge data from these endpoints into a single Discord response payload:
  - USD rates endpoint: `https://ve.dolarapi.com/v1/dolares` (array response)
  - EUR official endpoint: `https://ve.dolarapi.com/v1/euros/oficial` (single-object response)
- Add resilient behavior for partial upstream failures so users still receive available rate data.
- Include clear labeling for currency, source, and update time to avoid confusion between USD and EUR entries.
- Add a rollback plan: disable EUR retrieval and revert to current USD-only response if new upstream dependency is unstable in production.

## Capabilities

### New Capabilities
- `finance-multi-currency-rates`: Provide both USD and EUR exchange-rate information in the finance command response.

### Modified Capabilities
- None.

## Impact

- Affected code: `internal/features/finance/subcommands/dollar_all.go` and supporting finance subcommand types/helpers.
- External APIs:
  - `https://ve.dolarapi.com/v1/dolares` for USD rates (multiple source entries)
  - `https://ve.dolarapi.com/v1/euros/oficial` for official EUR rate (single entry)
- Payload contract notes: upstream objects expose `moneda`, `fuente`, `nombre`, `promedio`, `fechaActualizacion`, with `compra` and `venta` potentially null.
- Runtime behavior: slightly higher latency risk due to additional API call; mitigated by concurrent requests and partial-failure handling.
- User-facing output: richer embed content in `/dollar all`.
