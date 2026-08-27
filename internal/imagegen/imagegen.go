package imagegen

// Card is a single entry in a CardList image.
type Card struct {
	Label   string // small secondary line, e.g. "DÓLAR OFICIAL (USD)"
	Value   string // big accent line, e.g. "Bs. 787,52"
	Caption string // small line, e.g. "Fuente: BCV · 26 Ago 2026, 12:00 AM"
}

// CardListOptions describes a title + optional subtitle and note + flat card
// list image.
type CardListOptions struct {
	Title    string
	Subtitle string
	Note     string // optional warning rendered under the cards
	Cards    []Card
	Theme    *Theme // nil → DefaultTheme()
}
