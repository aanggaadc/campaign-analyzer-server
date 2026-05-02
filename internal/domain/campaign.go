package domain

import (
	"errors"
	"time"
)

type Campaign struct {
	ID          string
	UserID      string
	Name        string
	Platform    string
	Impressions int
	Clicks      int
	Conversions int
	Cost        float64
	DateStart   time.Time
	DateEnd     time.Time
}

func NewCampaign(
	ID, userID, name, platform string,
	impressions, clicks, conversions int,
	cost float64,
	start, end time.Time,
) (*Campaign, error) {

	campaign := &Campaign{
		ID:          ID,
		UserID:      userID,
		Name:        name,
		Platform:    platform,
		Impressions: impressions,
		Clicks:      clicks,
		Conversions: conversions,
		Cost:        cost,
		DateStart:   start,
		DateEnd:     end,
	}

	if err := campaign.Validate(); err != nil {
		return nil, err
	}

	return campaign, nil
}

func (c *Campaign) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}

	if c.Platform == "" {
		return errors.New("platform is required")
	}

	if c.Impressions < 0 || c.Clicks < 0 || c.Conversions < 0 {
		return errors.New("metrics cannot be negative")
	}

	if c.DateEnd.Before(c.DateStart) {
		return errors.New("end date cannot be before start date")
	}

	return nil
}

func (c *Campaign) CTR() float64 {
	if c.Impressions == 0 {
		return 0
	}
	return float64(c.Clicks) / float64(c.Impressions)
}

func (c *Campaign) CPC() float64 {
	if c.Clicks == 0 {
		return 0
	}
	return c.Cost / float64(c.Clicks)
}

func (c *Campaign) CPA() float64 {
	if c.Conversions == 0 {
		return 0
	}
	return c.Cost / float64(c.Conversions)
}
