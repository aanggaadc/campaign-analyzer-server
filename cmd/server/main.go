package main

import (
	"log"
	"os"

	"campaign-analyzer/internal/delivery/http"
	"campaign-analyzer/internal/infrastructure/database"
	"campaign-analyzer/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"campaign-analyzer/pkg/middleware"

	"campaign-analyzer/internal/infrastructure/ai"

	"time"

	"github.com/gin-contrib/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	}

	db := database.NewPostgresPool(os.Getenv("DATABASE_URL"))

	repo := database.NewCampaignRepository(db)
	analysisRepo := database.NewAnalysisRepository(db)
	// openAIService := ai.NewOpenAIService(os.Getenv("OPENAI_API_KEY"))
	geminiService := ai.NewGeminiService(os.Getenv("GEMINI_API_KEY"))

	campaignUc := usecase.NewCampaignUsecase(repo)
	analyzeUC := usecase.NewAnalyzeCampaignUsecase(
		repo,
		analysisRepo,
		geminiService,
	)

	campaignHandler := http.NewCampaignHandler(campaignUc)
	analysisHandler := http.NewAnalysisHandler(analyzeUC)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			os.Getenv("FRONTEND_URL"),
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.GET("campaigns", campaignHandler.GetCampaigns)
	auth.POST("campaigns", campaignHandler.CreateCampaign)
	auth.GET("/campaigns/template", campaignHandler.DownloadCampaignTemplate)
	auth.POST("/campaigns/import", campaignHandler.ImportCSV)
	auth.GET("/campaigns/:id", campaignHandler.GetCampaignDetail)
	auth.GET("/campaigns/:id/analyze", analysisHandler.AnalyzeCampaign)
	auth.GET("/campaigns/:id/metrics", campaignHandler.GetCampaignMetrics)

	auth.GET("analyses", analysisHandler.GetAnalyses)
	auth.GET("analyses/:id", analysisHandler.GetAnalysisDetail)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
