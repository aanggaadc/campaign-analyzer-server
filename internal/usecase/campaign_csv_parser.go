package usecase

import (
	"campaign-analyzer/internal/domain"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CampaignCSVParser struct{}

func NewCampaignCSVParser() *CampaignCSVParser {
	return &CampaignCSVParser{}
}

var allowedPlatformList = []string{
	"Facebook",
	"Instagram",
	"Google",
	"TikTok",
	"Twitter/X",
	"YouTube",
	"LinkedIn",
}

var allowedPlatformMap = map[string]bool{
	"facebook":  true,
	"instagram": true,
	"google":    true,
	"tiktok":    true,
	"twitter/x": true,
	"youtube":   true,
	"linkedin":  true,
}

func (p *CampaignCSVParser) getAllowedPlatformsString() string {
	return strings.Join(allowedPlatformList, ", ")
}

func (p *CampaignCSVParser) normalizePlatform(val string) string {
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "facebook":
		return "Facebook"
	case "instagram":
		return "Instagram"
	case "google":
		return "Google"
	case "tiktok":
		return "TikTok"
	case "twitter/x":
		return "Twitter/X"
	case "youtube":
		return "YouTube"
	case "linkedin":
		return "LinkedIn"
	default:
		return val
	}
}

func (p *CampaignCSVParser) safeGet(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func (p *CampaignCSVParser) ParseRow(userID string, row []string) (*domain.Campaign, error) {

	if len(row) < 8 {
		return nil, errors.New("invalid column count")
	}

	name := strings.TrimSpace(row[0])
	platformRaw := strings.TrimSpace(row[1])
	platformKey := strings.ToLower(platformRaw)

	if !allowedPlatformMap[platformKey] {
		return nil, fmt.Errorf(
			"invalid platform: %s (allowed: %s)",
			platformRaw,
			p.getAllowedPlatformsString(),
		)
	}

	platform := p.normalizePlatform(platformRaw)

	impressions, err := strconv.Atoi(row[2])
	if err != nil {
		return nil, fmt.Errorf("invalid impressions: %s", row[2])
	}

	clicks, err := strconv.Atoi(row[3])

	if clicks > impressions {
		return nil, fmt.Errorf("clicks (%d) cannot be greater than impressions (%d)", clicks, impressions)
	}

	if err != nil {
		return nil, fmt.Errorf("invalid clicks: %s", row[3])
	}

	conversions, err := strconv.Atoi(row[4])
	if conversions > clicks {
		return nil, fmt.Errorf("conversions (%d) cannot be greater than clicks (%d)", conversions, clicks)
	}

	if err != nil {
		return nil, fmt.Errorf("invalid conversions: %s", row[4])
	}

	cost, err := strconv.ParseFloat(row[5], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid cost: %s", row[5])
	}

	start, err := time.Parse("2006-01-02", row[6])
	if err != nil {
		return nil, errors.New("invalid start date format (YYYY-MM-DD)")
	}

	end, err := time.Parse("2006-01-02", row[7])
	if err != nil {
		return nil, errors.New("invalid end date format (YYYY-MM-DD)")
	}

	return domain.NewCampaign(
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
}
