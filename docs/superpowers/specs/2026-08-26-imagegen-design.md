# Design: Generic Image Generator (`internal/imagegen`)

Date: 2026-08-26
Status: Approved

## Problem

Bot responses currently use Discord embeds. The user wants image-based responses with embedded text, starting with `/dollar all`, using a generic generator that any feature can consume.

## Decisions

- **Scope:** Generic package usable by the whole bot; `/dollar all` is the first consumer.
- **Engine:** `fogleman/gg` (pure Go, no CGO — compatible with the `FROM scratch` production Dockerfile).
- **Style:** Simple flat card list; text only, no icons or embedded images.
- **Font:** Inter (Regular + Bold, SIL OFL 1.1) embedded via `go:embed` (required: prod image has no system fonts and builds with `CGO_ENABLED=0`).

## Architecture

```
internal/imagegen/
  spec.md        — subsystem spec
  imagegen.go    — public API: RenderCardList(opts) ([]byte, error)
  theme.go       — Theme struct + DefaultTheme()
  fonts.go       — go:embed Inter Regular/Bold, face factory
  cardlist.go    — card-list layout renderer
  assets/        — Inter-*.ttf + OFL license
```

### Public API

```go
type Card struct {
    Label   string // e.g. "DÓLAR OFICIAL (USD)"
    Value   string // e.g. "Bs. 787,52"
    Caption string // e.g. "Fuente: BCV · 26 Ago 2026, 12:00 AM"
}

type CardListOptions struct {
    Title    string
    Subtitle string
    Note     string // optional warning, e.g. "No disponible temporalmente: EUR"
    Cards    []Card
    Theme    *Theme // nil → DefaultTheme()
}

func RenderCardList(opts CardListOptions) ([]byte, error) // PNG bytes
```

- Canvas fixed at 1024px wide; height computed from card count (no clipping).
- Flat cards: rounded rect, thin border, three stacked text lines (label, value, caption).
- Plain package (not a `bot.Feature`); no registration in `cmd/bot/main.go`.

## Finance integration

- New `internal/features/finance/subcommands/dollar_card.go`:
  - `buildCardOptions(results) (cards, unavailable)` — one card per rate (USD then EUR), skips Bitcoin, marks currencies unavailable on fetch or timestamp-parse failure.
  - `buildCardListOptions(cards, unavailable) imagegen.CardListOptions` — title/subtitle/note assembly.
  - `cardFromRate(rate) (Card, error)` and `formatSpanishTimestamp(t)` ("26 Ago 2026, 12:00 AM", Spanish month abbreviations).
- `dollar_all.go` handler: defer → fetch (unchanged) → build options → render → send PNG via `InteractionResponseEdit` with `Files`. Existing embed path preserved as fallback.

## Error handling / degradation

| Failure | Behavior |
|---|---|
| One upstream source fails | Card list renders with `Note` naming the unavailable currency |
| All sources fail | Text error message, no image (unchanged) |
| `RenderCardList` error | Fall back to existing embed response |
| Discord send failure | Existing failure handler |

Known risk: discordgo v0.28.1 `WebhookEdit.Files` multipart behavior — verified during implementation; fallback paths cover any defect.

## Testing

- `imagegen`: in-package tests — PNG decodes with expected width/height; height math; empty-cards validation; render without optional fields.
- `subcommands`: `buildCardOptions` (ordering, Bitcoin skip, unavailable marking), `buildCardListOptions` (note), `formatSpanishTimestamp`.
- Manual: run bot, `/dollar all`, compare against mockup.

## Specs

- New `internal/imagegen/spec.md` (architectural subsystem).
- Update `internal/features/finance/spec.md` (image output + embed-fallback scenarios).
