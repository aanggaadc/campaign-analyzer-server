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

func (r *AnalysisRepository) FindAll(
	userID string,
	limit, offset int,
) ([]*domain.Analysis, error) {

	query := `
		SELECT 
		a.id, 
		a.user_id, 
		a.campaign_id, 
		a.summary, 
		a.issues, 
		a.recommendations, 
		a.priority_actions,
		a.created_at,
		c.name,
		(c.clicks::float / NULLIF(c.impressions, 0)) AS ctr 
		FROM analyses a LEFT JOIN campaigns c ON c.id = a.campaign_id 
		WHERE a.user_id = $1
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(context.Background(), query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []*domain.Analysis

	for rows.Next() {
		var c domain.Analysis

		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.CampaignID,
			&c.Summary,
			&c.Issues,
			&c.Recommendations,
			&c.PriorityActions,
			&c.CreatedAt,
			&c.CampaignName,
			&c.CTR,
		)
		if err != nil {
			return nil, err
		}

		analyses = append(analyses, &c)
	}

	// optional tapi best practice
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return analyses, nil
}

func (r *AnalysisRepository) Count(userID string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM analyses
		WHERE user_id = $1
	`

	var total int
	err := r.db.QueryRow(context.Background(), query, userID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *AnalysisRepository) Save(analysis *domain.Analysis) (string, error) {
	query := `
		INSERT INTO analyses (
			campaign_id,
			user_id,
			summary,
			issues,
			recommendations,
			priority_actions,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id string

	err := r.db.QueryRow(
		context.Background(),
		query,
		analysis.CampaignID,
		analysis.UserID,
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

func (r *AnalysisRepository) FindByID(userID, analysisID string) (*domain.Analysis, error) {
	query := `
		SELECT 
		a.id, 
		a.user_id, 
		a.campaign_id, 
		a.summary, 
		a.issues, 
		a.recommendations, 
		a.priority_actions, 
		a.created_at, 
		c.name, 
		(c.clicks::float / NULLIF(c.impressions, 0)) AS ctr 
		FROM analyses a LEFT JOIN campaigns c ON c.id = a.campaign_id 
		WHERE a.id = $1 AND a.user_id = $2
	`

	row := r.db.QueryRow(context.Background(), query, analysisID, userID)

	var a domain.Analysis

	err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.CampaignID,
		&a.Summary,
		&a.Issues,
		&a.Recommendations,
		&a.PriorityActions,
		&a.CreatedAt,
		&a.CampaignName,
		&a.CTR,
	)

	if err != nil {
		return nil, err
	}

	return &a, nil
}
