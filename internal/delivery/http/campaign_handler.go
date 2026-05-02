package http

import (
	"net/http"
	"time"

	"campaign-analyzer/internal/usecase"

	"github.com/gin-gonic/gin"
)

type CreateCampaignRequest struct {
	Name        string  `json:"name"`
	Platform    string  `json:"platform"`
	Impressions int     `json:"impressions"`
	Clicks      int     `json:"clicks"`
	Conversions int     `json:"conversions"`
	Cost        float64 `json:"cost"`
	DateStart   string  `json:"date_start"`
	DateEnd     string  `json:"date_end"`
}

type CampaignHandler struct {
	usecase *usecase.CampaignUsecase
}

func NewCampaignHandler(u *usecase.CampaignUsecase) *CampaignHandler {
	return &CampaignHandler{usecase: u}
}

func (h *CampaignHandler) GetCampaigns(c *gin.Context) {
	userID := c.GetString("user_id")

	campaigns, err := h.usecase.GetCampaigns(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var req CreateCampaignRequest
	userID := c.GetString("user_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// parse date
	start, err := time.Parse("2006-01-02", req.DateStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}

	end, err := time.Parse("2006-01-02", req.DateEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}

	campaign, err := h.usecase.CreateCampaign(
		userID,
		req.Name,
		req.Platform,
		req.Impressions,
		req.Clicks,
		req.Conversions,
		req.Cost,
		start,
		end,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaign)
}
