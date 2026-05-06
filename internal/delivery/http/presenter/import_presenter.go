package presenter

type ImportResult struct {
	Imported int              `json:"imported"`
	Failed   int              `json:"failed"`
	Errors   []ImportErrorRow `json:"errors"`
}

type ImportErrorRow struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}
