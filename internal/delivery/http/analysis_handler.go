package http

import (
	"fmt"
	"net/http"

	"campaign-analyzer/internal/delivery/http/presenter"
	"campaign-analyzer/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AnalysisHandler struct {
	analyzeUC *usecase.AnalyzeCampaignUsecase
	exportUC  *usecase.ExportAnalysisUsecase
}

func NewAnalysisHandler(
	analyzeUC *usecase.AnalyzeCampaignUsecase,
	exportUC *usecase.ExportAnalysisUsecase,
) *AnalysisHandler {
	return &AnalysisHandler{
		analyzeUC: analyzeUC,
		exportUC:  exportUC,
	}
}

func (h *AnalysisHandler) GetAnalyses(c *gin.Context) {
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

func (h *AnalysisHandler) GetAnalysisDetail(c *gin.Context) {
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

func (h *AnalysisHandler) AnalyzeCampaign(c *gin.Context) {
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

func (h *AnalysisHandler) ExportAnalysis(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	file, err := h.exportUC.Execute(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.ErrorResponse(err.Error()))
		return
	}

	filename := "analysis-" + id + ".pdf"

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename="+filename)

	c.Data(http.StatusOK, "application/pdf", file)
}
