package dto

import "time"

type AnalysisExportData struct {
	CampaignName string
	Platform     string

	CTR float64
	CPC float64
	CPA float64

	Summary         string
	Issues          []string
	Recommendations []string
	PriorityActions []string

	CreatedAt time.Time
}
