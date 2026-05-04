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
	usecase   *usecase.CampaignUsecase
	analyzeUC *usecase.AnalyzeCampaignUsecase
}

func NewCampaignHandler(
	u *usecase.CampaignUsecase,
	analyzeUC *usecase.AnalyzeCampaignUsecase,
) *CampaignHandler {
	return &CampaignHandler{
		usecase:   u,
		analyzeUC: analyzeUC,
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

func (h *CampaignHandler) GetAnalyses(c *gin.Context) {
	userID := c.GetString("user_id")

	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}

	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	analyses, total, err := h.analyzeUC.GetAnalyses(userID, page, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.ErrorResponse(err.Error()))
		return
	}

	response := presenter.ToAnalysisListResponse(analyses)

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

func (h *CampaignHandler) GetAnalysisDetail(c *gin.Context) {
	userID := c.GetString("user_id")
	campaignID := c.Param("id")

	analysis, err := h.analyzeUC.GetAnalysisID(userID, campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, presenter.ErrorResponse("campaign not found"))
		return
	}

	response := presenter.ToAnalysisResponse(analysis)

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

func (h *CampaignHandler) AnalyzeCampaign(c *gin.Context) {
	userID := c.GetString("user_id")
	campaignID := c.Param("id")

	result, err := h.analyzeUC.Execute(userID, campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.ErrorResponse(err.Error()))
		return
	}

	if result == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "analysis not available",
		})
		return
	}

	response := presenter.ToAnalysisResponse(result)
	c.JSON(http.StatusOK, presenter.Success(response))

}
