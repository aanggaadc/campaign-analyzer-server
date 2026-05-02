package domain

import "time"

type Analysis struct {
	ID              string
	CampaignID      string
	Summary         string
	Issues          []string
	Recommendations []string
	PriorityActions []string
	CreatedAt       time.Time
}
