package database

import (
	"context"
	"fmt"
	"strings"

	"campaign-analyzer/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CampaignRepositoryPostgres struct {
	db *pgxpool.Pool
}

func NewCampaignRepositoryPostgres(db *pgxpool.Pool) *CampaignRepositoryPostgres {
	return &CampaignRepositoryPostgres{db: db}
}

func (r *CampaignRepositoryPostgres) FindAll(
	userID string,
	limit, offset int,
) ([]*domain.Campaign, error) {

	query := `
		SELECT id, user_id, name, platform, impressions, clicks, conversions, cost, date_start, date_end
		FROM campaigns
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(context.Background(), query, userID, limit, offset)
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

	// optional tapi best practice
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return campaigns, nil
}

func (r *CampaignRepositoryPostgres) Count(userID string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM campaigns
		WHERE user_id = $1
	`

	var total int
	err := r.db.QueryRow(context.Background(), query, userID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *CampaignRepositoryPostgres) FindByID(userID, campaignID string) (*domain.Campaign, error) {
	query := `
		SELECT id, user_id, name, platform, impressions, clicks, conversions, cost, date_start, date_end
		FROM campaigns
		WHERE id = $1 AND user_id = $2
	`

	row := r.db.QueryRow(context.Background(), query, campaignID, userID)

	var c domain.Campaign

	err := row.Scan(
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

	return &c, nil
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

func (r *CampaignRepositoryPostgres) SaveBatch(campaigns []*domain.Campaign) error {
	if len(campaigns) == 0 {
		return nil
	}

	query := `
	INSERT INTO campaigns 
	(user_id, name, platform, impressions, clicks, conversions, cost, date_start, date_end)
	VALUES 
	`

	values := []interface{}{}
	placeholders := []string{}

	for i, c := range campaigns {
		idx := i * 10

		placeholders = append(placeholders, fmt.Sprintf(
			"($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			idx+1, idx+2, idx+3, idx+4, idx+5,
			idx+6, idx+7, idx+8, idx+9, idx+10,
		))

		values = append(values,
			c.UserID,
			c.Name,
			c.Platform,
			c.Impressions,
			c.Clicks,
			c.Conversions,
			c.Cost,
			c.DateStart,
			c.DateEnd,
		)
	}

	query += strings.Join(placeholders, ",")

	_, err := r.db.Exec(context.Background(), query, values...)
	return err
}
