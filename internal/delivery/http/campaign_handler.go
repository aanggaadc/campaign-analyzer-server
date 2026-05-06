package http

import (
	"fmt"
	"net/http"
	"time"

	"campaign-analyzer/internal/delivery/http/presenter"
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

func NewCampaignHandler(
	u *usecase.CampaignUsecase,
) *CampaignHandler {
	return &CampaignHandler{
		usecase: u,
	}
}

func (h *CampaignHandler) GetCampaigns(c *gin.Context) {
	userID := c.GetString("user_id")

	// default values
	page := 1
	limit := 10

	// query params
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}

	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	campaigns, total, err := h.usecase.GetCampaigns(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.ErrorResponse(err.Error()))
		return
	}

	response := presenter.ToCampaignListResponse(campaigns)

	meta := presenter.PaginationMeta(page, limit, total)

	c.JSON(http.StatusOK, presenter.SuccessWithMeta(response, meta))
}

func (h *CampaignHandler) GetCampaignDetail(c *gin.Context) {
	userID := c.GetString("user_id")
	campaignID := c.Param("id")

	campaign, err := h.usecase.GetCampaignByID(userID, campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, presenter.ErrorResponse("campaign not found"))
		return
	}

	response := presenter.ToCampaignResponse(campaign)

	c.JSON(http.StatusOK, presenter.Success(response))
}

func (h *CampaignHandler) GetCampaignMetrics(c *gin.Context) {
	userID := c.GetString("user_id")
	campaignID := c.Param("id")

	campaign, err := h.usecase.GetCampaignByID(userID, campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, presenter.ErrorResponse("campaign not found"))
		return
	}

	response := presenter.CampaignMetrics{
		CTR: campaign.CTR(),
		CPC: campaign.CPC(),
		CPA: campaign.CPA(),
	}

	c.JSON(http.StatusOK, presenter.Success(response))
}

func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var req CreateCampaignRequest
	userID := c.GetString("user_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, presenter.ErrorResponse(err.Error()))
		return
	}

	// parse date
	start, err := time.Parse("2006-01-02", req.DateStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, presenter.ErrorResponse("invalid start date"))
		return
	}

	end, err := time.Parse("2006-01-02", req.DateEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, presenter.ErrorResponse("invalid end date"))
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
		c.JSON(http.StatusBadRequest, presenter.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, presenter.Success(presenter.ToCampaignResponse(campaign)))
}

func (h *CampaignHandler) DownloadCampaignTemplate(c *gin.Context) {
	csvContent := `name,platform,impressions,clicks,conversions,cost,date_start,date_end`

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=campaign_template.csv")
	c.Header("Content-Type", "text/csv")

	c.String(http.StatusOK, csvContent)
}

func (h *CampaignHandler) ImportCSV(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, presenter.ErrorResponse("file is required"))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, presenter.ErrorResponse("cannot open file"))
		return
	}
	defer file.Close()

	userID := c.GetString("user_id")

	result, err := h.usecase.ImportCSV(userID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, presenter.Success(result))
}
