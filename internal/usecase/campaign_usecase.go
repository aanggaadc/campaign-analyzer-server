package usecase

import (
	"campaign-analyzer/internal/domain"
	"campaign-analyzer/internal/repository"
	"time"
)

type CampaignUsecase struct {
	repo repository.CampaignRepository
}

func NewCampaignUsecase(repo repository.CampaignRepository) *CampaignUsecase {
	return &CampaignUsecase{
		repo: repo,
	}
}

func (uc *CampaignUsecase) GetCampaigns(userID string) ([]*domain.Campaign, error) {
	return uc.repo.FindAll(userID)
}

func (uc *CampaignUsecase) CreateCampaign(
	userID string,
	name string,
	platform string,
	impressions, clicks, conversions int,
	cost float64,
	start, end time.Time,
) (*domain.Campaign, error) {

	campaign, err := domain.NewCampaign(
		"",
		userID,
		name,
		platform,
		impressions,
		clicks,
		conversions,
		cost,
		start,
		end,
	)

	if err != nil {
		return nil, err
	}

	id, err := uc.repo.Save(campaign)
	if err != nil {
		return nil, err
	}

	campaign.ID = id
	return campaign, nil
}
