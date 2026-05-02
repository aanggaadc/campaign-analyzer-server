package database

import (
	"context"

	"campaign-analyzer/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CampaignRepositoryPostgres struct {
	db *pgxpool.Pool
}

func NewCampaignRepositoryPostgres(db *pgxpool.Pool) *CampaignRepositoryPostgres {
	return &CampaignRepositoryPostgres{db: db}
}

func (r *CampaignRepositoryPostgres) FindAll(userID string) ([]*domain.Campaign, error) {
	query := `
		SELECT id, user_id, name, platform, impressions, clicks, conversions, cost, date_start, date_end
		FROM campaigns
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []*domain.Campaign

	for rows.Next() {
		var c domain.Campaign

		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Name,
			&c.Platform,
			&c.Impressions,
			&c.Clicks,
			&c.Conversions,
			&c.Cost,
			&c.DateStart,
			&c.DateEnd,
		)
		if err != nil {
			return nil, err
		}

		campaigns = append(campaigns, &c)
	}

	return campaigns, nil
}

func (r *CampaignRepositoryPostgres) Save(campaign *domain.Campaign) (string, error) {
	query := `
		INSERT INTO campaigns 
		(user_id, name, platform, impressions, clicks, conversions, cost, date_start, date_end)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
	`

	var id string
	err := r.db.QueryRow(context.Background(), query,
		campaign.UserID,
		campaign.Name,
		campaign.Platform,
		campaign.Impressions,
		campaign.Clicks,
		campaign.Conversions,
		campaign.Cost,
		campaign.DateStart,
		campaign.DateEnd,
	).Scan(&id)

	return id, err
}
