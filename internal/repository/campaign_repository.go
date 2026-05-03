package repository

import "campaign-analyzer/internal/domain"

type CampaignRepository interface {
	FindAll(userID string, limit, offset int) ([]*domain.Campaign, error)
	Count(userID string) (int, error)
	Save(campaign *domain.Campaign) (string, error)
	FindByID(userID, campaignID string) (*domain.Campaign, error)
}
