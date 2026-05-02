package repository

import "campaign-analyzer/internal/domain"

type AnalysisRepository interface {
	Save(analysis *domain.Analysis) (string, error)
}
