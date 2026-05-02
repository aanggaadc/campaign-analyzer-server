package database

import (
	"campaign-analyzer/internal/domain"
)

type CampaignRepositoryMemory struct {
	data []*domain.Campaign
}

func NewCampaignRepositoryMemory() *CampaignRepositoryMemory {
	return &CampaignRepositoryMemory{
		data: []*domain.Campaign{},
	}
}

func (r *CampaignRepositoryMemory) Save(campaign *domain.Campaign) error {
	// simple append (simulate DB insert)
	r.data = append(r.data, campaign)
	return nil
}
