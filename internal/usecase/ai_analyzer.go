package usecase

import (
	"campaign-analyzer/internal/domain"
	"campaign-analyzer/internal/repository"
	"errors"
	"fmt"
	"time"
)

type AnalyzeInput struct {
	Name        string
	Platform    string
	Impressions int
	Clicks      int
	Conversions int
	Cost        float64

	CTR float64
	CPC float64
	CPA float64
}

type AnalyzeResult struct {
	Summary         string   `json:"summary"`
	Issues          []string `json:"issues"`
	Recommendations []string `json:"recommendations"`
	PriorityActions []string `json:"priority_actions"`
}

type AnalyzeCampaignUsecase struct {
	repo         repository.CampaignRepository
	ai           AIAnalyzer
	analysisRepo repository.AnalysisRepository
}

type AIAnalyzer interface {
	AnalyzeCampaign(input AnalyzeInput) (AnalyzeResult, error)
}

func NewAnalyzeCampaignUsecase(
	repo repository.CampaignRepository,
	analysisRepo repository.AnalysisRepository,
	ai AIAnalyzer,
) *AnalyzeCampaignUsecase {
	return &AnalyzeCampaignUsecase{
		repo:         repo,
		analysisRepo: analysisRepo,
		ai:           ai,
	}
}

func (uc *AnalyzeCampaignUsecase) GetAnalyses(
	userID string,
	page int,
	limit int,
) ([]*domain.Analysis, int, error) {

	offset := (page - 1) * limit

	analyses, err := uc.analysisRepo.FindAll(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.analysisRepo.Count(userID)
	if err != nil {
		return nil, 0, err
	}

	return analyses, total, nil
}

func (uc *AnalyzeCampaignUsecase) Execute(userID, campaignID string) (*domain.Analysis, error) {
	// 1. ambil campaign
	campaign, err := uc.repo.FindByID(userID, campaignID)
	if err != nil {
		return nil, err
	}

	if campaign == nil {
		return nil, errors.New("campaign not found")
	}

	// 2. hitung metrics
	ctr := campaign.CTR()
	cpc := campaign.CPC()
	cpa := campaign.CPA()

	// 3. build input ke AI
	input := AnalyzeInput{
		Name:        campaign.Name,
		Platform:    campaign.Platform,
		Impressions: campaign.Impressions,
		Clicks:      campaign.Clicks,
		Conversions: campaign.Conversions,
		Cost:        campaign.Cost,
		CTR:         ctr,
		CPC:         cpc,
		CPA:         cpa,
	}

	// 4. call AI
	result, err := uc.ai.AnalyzeCampaign(input)
	if err != nil {
		fmt.Println("AI ERROR:", err)
		return nil, nil
	}

	// 5. mapping ke domain
	analysis := &domain.Analysis{
		CampaignID:      campaign.ID,
		UserID:          userID,
		Summary:         result.Summary,
		Issues:          result.Issues,
		Recommendations: result.Recommendations,
		PriorityActions: result.PriorityActions,
		CreatedAt:       time.Now(),
	}

	// 6. simpan ke DB
	id, err := uc.analysisRepo.Save(analysis)
	if err != nil {
		return nil, err
	}

	analysis.ID = id

	return analysis, nil
}
