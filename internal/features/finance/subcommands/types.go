package subcommands

type DolarResponse struct {
	Source    string  `json:"fuente"`
	Name      string  `json:"nombre"`
	Buy       int     `json:"compra"`
	Sell      int     `json:"venta"`
	Average   float32 `json:"promedio"`
	UpdatedAt string  `json:"fechaActualizacion"`
}
