package subcommands

type DollarSource string

const (
	OfficialDollarSource DollarSource = "oficial"
	ParallelDollarSource DollarSource = "paralelo"
)

type DolarResponse struct {
	Currency  string   `json:"moneda"`
	Source    string   `json:"fuente"`
	Name      string   `json:"nombre"`
	Buy       *float64 `json:"compra"`
	Sell      *float64 `json:"venta"`
	Average   float64  `json:"promedio"`
	UpdatedAt string   `json:"fechaActualizacion"`
}
