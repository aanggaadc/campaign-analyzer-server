package presenter

import "campaign-analyzer/internal/domain"

type AnalysisResponse struct {
	Summary         string   `json:"summary"`
	Issues          []string `json:"issues"`
	Recommendations []string `json:"recommendations"`
	PriorityActions []string `json:"priority_actions"`
}

type AnalysisListResponse struct {
	ID              string   `json:"id"`
	UserID          string   `json:"user_id"`
	Summary         string   `json:"summary"`
	Issues          []string `json:"issues"`
	Recommendations []string `json:"recommendations"`
	PriorityActions []string `json:"priority_actions"`
	CreatedAt       string   `json:"created_at"`
}

func ToAnalysisResponse(a *domain.Analysis) AnalysisResponse {
	return AnalysisResponse{
		Summary:         a.Summary,
		Issues:          a.Issues,
		Recommendations: a.Recommendations,
		PriorityActions: a.PriorityActions,
	}
}

func ToAnalysisListResponse(campaigns []*domain.Analysis) []AnalysisListResponse {
	var result []AnalysisListResponse

	for _, c := range campaigns {
		result = append(result, AnalysisListResponse{
			ID:              c.ID,
			UserID:          c.UserID,
			Summary:         c.Summary,
			Issues:          c.Issues,
			Recommendations: c.Recommendations,
			PriorityActions: c.PriorityActions,
			CreatedAt:       c.CreatedAt.Format("2006-01-02"),
		})
	}

	return result
}
