package usecase

import (
	"campaign-analyzer/internal/domain"
	"campaign-analyzer/internal/repository"
	"time"
)

type CampaignUsecase struct {
	repo repository.CampaignRepository
}

type CampaignMetrics struct {
	CTR float64 `json:"ctr"`
	CPC float64 `json:"cpc"`
	CPA float64 `json:"cpa"`
}

func NewCampaignUsecase(repo repository.CampaignRepository) *CampaignUsecase {
	return &CampaignUsecase{
		repo: repo,
	}
}

func (uc *CampaignUsecase) GetCampaigns(
	userID string,
	page int,
	limit int,
) ([]*domain.Campaign, int, error) {

	offset := (page - 1) * limit

	campaigns, err := uc.repo.FindAll(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.repo.Count(userID)
	if err != nil {
		return nil, 0, err
	}

	return campaigns, total, nil
}

func (u *CampaignUsecase) GetCampaignByID(userID, campaignID string) (*domain.Campaign, error) {
	return u.repo.FindByID(userID, campaignID)
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

func (u *CampaignUsecase) GetCampaignMetrics(campaign *domain.Campaign) CampaignMetrics {
	return CampaignMetrics{
		CTR: campaign.CTR(),
		CPC: campaign.CPC(),
		CPA: campaign.CPA(),
	}
}
