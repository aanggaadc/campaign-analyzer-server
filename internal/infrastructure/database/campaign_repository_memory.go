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

func (r *CampaignRepositoryMemory) Save(campaign *domain.Campaign) (string, error) {
	r.data = append(r.data, campaign)
	return campaign.ID, nil
}
