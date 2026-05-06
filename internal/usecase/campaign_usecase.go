package usecase

import (
	"campaign-analyzer/internal/domain"
	"campaign-analyzer/internal/dto"
	"campaign-analyzer/internal/repository"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"
)

type CampaignUsecase struct {
	repo   repository.CampaignRepository
	parser *CampaignCSVParser
}

type CampaignUsecaseInterface interface {
	ImportCSV(userID string, file multipart.File) (*dto.ImportResult, error)
}

type CampaignMetrics struct {
	CTR float64 `json:"ctr"`
	CPC float64 `json:"cpc"`
	CPA float64 `json:"cpa"`
}

func NewCampaignUsecase(repo repository.CampaignRepository) *CampaignUsecase {
	return &CampaignUsecase{
		repo:   repo,
		parser: NewCampaignCSVParser(),
	}
}

func (uc *CampaignUsecase) GetCampaigns(
	userID string,
	page int,
	limit int,
) ([]*domain.Campaign, int, error) {

	offset := (page - 1) * limit

	campaigns, err := uc.repo.FindAll(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.repo.Count(userID)
	if err != nil {
		return nil, 0, err
	}

	return campaigns, total, nil
}

func (u *CampaignUsecase) GetCampaignByID(userID, campaignID string) (*domain.Campaign, error) {
	return u.repo.FindByID(userID, campaignID)
}

func (uc *CampaignUsecase) CreateCampaign(
	userID string,
	name string,
	platform string,
	impressions, clicks, conversions int,
	cost float64,
	start, end time.Time,
) (*domain.Campaign, error) {

	campaign, err := domain.NewCampaign(
		"",
		userID,
		name,
		platform,
		impressions,
		clicks,
		conversions,
		cost,
		start,
		end,
	)

	if err != nil {
		return nil, err
	}

	id, err := uc.repo.Save(campaign)
	if err != nil {
		return nil, err
	}

	campaign.ID = id
	return campaign, nil
}

func (u *CampaignUsecase) GetCampaignMetrics(campaign *domain.Campaign) CampaignMetrics {
	return CampaignMetrics{
		CTR: campaign.CTR(),
		CPC: campaign.CPC(),
		CPA: campaign.CPA(),
	}
}

func (u *CampaignUsecase) ImportCSV(userID string, file multipart.File) (*dto.ImportResult, error) {
	reader := csv.NewReader(file)

	result := &dto.ImportResult{}

	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	expected := []string{
		"name", "platform", "impressions", "clicks",
		"conversions", "cost", "date_start", "date_end",
	}

	if len(header) < len(expected) {
		result.Errors = append(result.Errors, "invalid header: column count mismatch")
		return result, nil
	}

	for i := range expected {
		if strings.ToLower(strings.TrimSpace(header[i])) != expected[i] {
			result.Errors = append(result.Errors, "invalid header format")
			return result, nil
		}
	}

	rowNumber := 1
	batchSize := 100
	maxPreview := 50

	var batch []*domain.Campaign

	for {
		rowNumber++

		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			result.Failed++

			if len(result.Rows) < maxPreview {
				result.Rows = append(result.Rows, dto.ImportRowResult{
					Row:      rowNumber,
					Name:     u.parser.safeGet(row, 0),
					Platform: u.parser.safeGet(row, 1),
					Status:   "invalid",
					Error:    "failed to read row",
				})
			}

			continue
		}

		campaign, err := u.parser.ParseRow(userID, row)
		if err != nil {
			result.Failed++

			if len(result.Rows) < maxPreview {
				result.Rows = append(result.Rows, dto.ImportRowResult{
					Row:         rowNumber,
					Name:        u.parser.safeGet(row, 0),
					Platform:    u.parser.safeGet(row, 1),
					Impressions: u.parser.safeGet(row, 2),
					Clicks:      u.parser.safeGet(row, 3),
					Conversions: u.parser.safeGet(row, 4),
					Cost:        u.parser.safeGet(row, 5),
					Status:      "invalid",
					Error:       err.Error(),
				})
			}

			continue
		}

		if len(result.Rows) < maxPreview {
			result.Rows = append(result.Rows, dto.ImportRowResult{
				Row:         rowNumber,
				Name:        campaign.Name,
				Platform:    campaign.Platform,
				Impressions: u.parser.safeGet(row, 2),
				Clicks:      u.parser.safeGet(row, 3),
				Conversions: u.parser.safeGet(row, 4),
				Cost:        u.parser.safeGet(row, 5),
				Status:      "valid",
			})
		}

		batch = append(batch, campaign)

		if len(batch) >= batchSize {
			err := u.repo.SaveBatch(batch)
			if err != nil {
				result.Failed += len(batch)

				result.Errors = append(result.Errors,
					fmt.Sprintf("failed to save batch at row %d: %v", rowNumber, err),
				)

			} else {
				result.Imported += len(batch)
			}

			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		err := u.repo.SaveBatch(batch)
		if err != nil {
			result.Failed += len(batch)

			result.Errors = append(result.Errors,
				fmt.Sprintf("failed to save final batch: %v", err),
			)

		} else {
			result.Imported += len(batch)
		}
	}

	return result, nil
}
