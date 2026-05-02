package database

import (
	"campaign-analyzer/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalysisRepository struct {
	db *pgxpool.Pool
}

func NewAnalysisRepository(db *pgxpool.Pool) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

func (r *AnalysisRepository) Save(analysis *domain.Analysis) (string, error) {
	query := `
		INSERT INTO analyses (
			campaign_id,
			summary,
			issues,
			recommendations,
			priority_actions,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id string

	err := r.db.QueryRow(
		context.Background(),
		query,
		analysis.CampaignID,
		analysis.Summary,
		analysis.Issues,
		analysis.Recommendations,
		analysis.PriorityActions,
		analysis.CreatedAt,
	).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}
