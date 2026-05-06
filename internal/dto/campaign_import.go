package dto

type ImportResult struct {
	Imported int               `json:"imported"`
	Failed   int               `json:"failed"`
	Errors   []string          `json:"errors,omitempty"`
	Rows     []ImportRowResult `json:"rows"`
}

type ImportRowResult struct {
	Row         int    `json:"row"`
	Name        string `json:"name"`
	Platform    string `json:"platform"`
	Impressions string `json:"impressions"`
	Clicks      string `json:"clicks"`
	Conversions string `json:"conversions"`
	Cost        string `json:"cost"`
	Status      string `json:"status"` // "valid" | "invalid"
	Error       string `json:"error,omitempty"`
}
