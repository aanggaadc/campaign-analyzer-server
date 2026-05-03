package presenter

import "campaign-analyzer/internal/domain"

type AnalysisResponse struct {
	Summary         string   `json:"summary"`
	Issues          []string `json:"issues"`
	Recommendations []string `json:"recommendations"`
	PriorityActions []string `json:"priority_actions"`
}

func ToAnalysisResponse(a *domain.Analysis) AnalysisResponse {
	return AnalysisResponse{
		Summary:         a.Summary,
		Issues:          a.Issues,
		Recommendations: a.Recommendations,
		PriorityActions: a.PriorityActions,
	}
}
