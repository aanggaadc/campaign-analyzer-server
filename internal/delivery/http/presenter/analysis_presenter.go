package presenter

import "campaign-analyzer/internal/domain"

type AnalysisResponse struct {
	ID              string   `json:"id"`
	UserID          string   `json:"user_id"`
	CampaignId      string   `json:"campaign_id"`
	Summary         string   `json:"summary"`
	Issues          []string `json:"issues"`
	Recommendations []string `json:"recommendations"`
	PriorityActions []string `json:"priority_actions"`
	CreatedAt       string   `json:"created_at"`
}

func ToAnalysisResponse(a *domain.Analysis) AnalysisResponse {
	return AnalysisResponse{
		ID:              a.ID,
		UserID:          a.UserID,
		CampaignId:      a.CampaignID,
		Summary:         a.Summary,
		Issues:          a.Issues,
		Recommendations: a.Recommendations,
		PriorityActions: a.PriorityActions,
		CreatedAt:       a.CreatedAt.Format("2006-01-02"),
	}
}

func ToAnalysisListResponse(campaigns []*domain.Analysis) []AnalysisResponse {
	var result []AnalysisResponse

	for _, c := range campaigns {
		result = append(result, ToAnalysisResponse(c))
	}

	return result
}
