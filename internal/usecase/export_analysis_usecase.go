package usecase

import (
	"campaign-analyzer/internal/dto"
	"campaign-analyzer/internal/repository"
)

type PDFGenerator interface {
	GenerateAnalysisPDF(data dto.AnalysisExportData) ([]byte, error)
}

type ExportAnalysisUsecase struct {
	repo         repository.AnalysisRepository
	campaignRepo repository.CampaignRepository
	pdf          PDFGenerator
}

func NewExportAnalysisUsecase(repo repository.AnalysisRepository, pdf PDFGenerator) *ExportAnalysisUsecase {
	return &ExportAnalysisUsecase{
		repo: repo,
		pdf:  pdf,
	}
}

func (u *ExportAnalysisUsecase) Execute(userID string, id string) ([]byte, error) {
	analysis, err := u.repo.FindByID(userID, id)
	campaign := analysis.Campaign

	if err != nil {
		return nil, err
	}

	data := dto.AnalysisExportData{
		CampaignName: campaign.Name,
		Platform:     campaign.Platform,

		CTR: campaign.CTR(),
		CPC: campaign.CPC(),
		CPA: campaign.CPA(),

		Summary:         analysis.Summary,
		Issues:          analysis.Issues,
		Recommendations: analysis.Recommendations,
		PriorityActions: analysis.PriorityActions,

		CreatedAt: analysis.CreatedAt,
	}

	return u.pdf.GenerateAnalysisPDF(data)
}
