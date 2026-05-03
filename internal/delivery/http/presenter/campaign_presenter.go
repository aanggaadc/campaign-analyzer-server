package presenter

import (
	"campaign-analyzer/internal/domain"
)

type CampaignMetrics struct {
	CTR float64 `json:"ctr"`
	CPC float64 `json:"cpc"`
	CPA float64 `json:"cpa"`
}

type CampaignResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Platform    string          `json:"platform"`
	Clicks      int             `json:"clicks"`
	Conversions int             `json:"conversions"`
	Impressions int             `json:"impressions"`
	Cost        float64         `json:"cost"`
	DateStart   string          `json:"date_start"`
	DateEnd     string          `json:"date_end"`
	Metrics     CampaignMetrics `json:"metrics"`
}

func ToCampaignResponse(c *domain.Campaign) CampaignResponse {
	return CampaignResponse{
		ID:          c.ID,
		Name:        c.Name,
		Platform:    c.Platform,
		Clicks:      c.Clicks,
		Conversions: c.Conversions,
		Impressions: c.Impressions,
		Cost:        c.Cost,
		DateStart:   c.DateStart.Format("2006-01-02"),
		DateEnd:     c.DateEnd.Format("2006-01-02"),
		Metrics: CampaignMetrics{
			CTR: c.CTR(),
			CPC: c.CPC(),
			CPA: c.CPA(),
		},
	}
}

func ToCampaignListResponse(campaigns []*domain.Campaign) []CampaignResponse {
	var result []CampaignResponse

	for _, c := range campaigns {
		result = append(result, ToCampaignResponse(c))
	}

	return result
}
