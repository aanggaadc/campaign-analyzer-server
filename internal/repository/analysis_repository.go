package repository

import "campaign-analyzer/internal/domain"

type AnalysisRepository interface {
	FindAll(userID string, limit, offset int) ([]*domain.Analysis, error)
	Count(userID string) (int, error)
	Save(analysis *domain.Analysis) (string, error)
}
