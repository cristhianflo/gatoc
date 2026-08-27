# imagegen

## Purpose

`internal/imagegen` renders images with embedded text for bot responses. It is an architectural subsystem: a plain package (not a `bot.Feature`) that any feature may call to produce PNG bytes for Discord attachments.

## Responsibilities

- Render card-list images: optional title, subtitle, and note over a vertical list of flat cards.
- Encode rendered images as PNG bytes.
- Own the embedded fonts (Inter Regular/Bold, SIL OFL 1.1 — see `assets/LICENSE-INTER-OFL.txt`) and the color theme.

## Public API

- `RenderCardList(opts CardListOptions) ([]byte, error)` — the only entry point. Returns PNG bytes or an error.
- `CardListOptions{Title, Subtitle, Note, Cards, Theme}` — `Note` is optional (empty string omits it); `Theme` nil falls back to `DefaultTheme()`.
- `Card{Label, Value, Caption}` — each card draws its three text lines stacked.

## Rules

- The package is pure Go with no CGO; all assets are embedded via `go:embed` so the production image (`FROM scratch`) needs no system fonts.
- The canvas is 1024px wide; the height is computed from the option set and card count, so content is never clipped.
- `RenderCardList` SHALL return an error when `Cards` is empty.
- Font parsing happens once (`sync.Once`); font faces are created per render because they are not safe for concurrent use.
- Layout metrics are internal. Callers must not depend on exact pixel dimensions, only on "decoded PNG with width 1024".

## Failure behavior

### Scenario: no cards

- **GIVEN** `CardListOptions` with zero cards
- **WHEN** `RenderCardList` is called
- **THEN** it returns an error and no image

### Scenario: successful render

- **GIVEN** options with at least one card
- **WHEN** `RenderCardList` is called
- **THEN** the returned bytes decode as a PNG
- **AND** the decoded image is 1024px wide with a height matching the computed layout

## Boundaries

- No knowledge of Discord, features, or exchange-rate data: it draws what it is given.
- Features own the mapping from their domain data to `Card`/`CardListOptions` (see `internal/features/finance/subcommands/dollar_card.go` for the reference consumer).
