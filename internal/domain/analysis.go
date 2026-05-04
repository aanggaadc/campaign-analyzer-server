package domain

import "time"

type Analysis struct {
	ID              string
	UserID          string
	CampaignID      string
	Summary         string
	Issues          []string
	Recommendations []string
	PriorityActions []string
	CreatedAt       time.Time

	CampaignName string
	CTR          float64
}
