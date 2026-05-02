package repository

import "campaign-analyzer/internal/domain"

type CampaignRepository interface {
	FindAll(userID string) ([]*domain.Campaign, error)
	Save(campaign *domain.Campaign) error
}
