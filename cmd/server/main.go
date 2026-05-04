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

	repo := database.NewCampaignRepositoryPostgres(db)
	analysisRepo := database.NewAnalysisRepository(db)
	// openAIService := ai.NewOpenAIService(os.Getenv("OPENAI_API_KEY"))
	geminiService := ai.NewGeminiService(os.Getenv("GEMINI_API_KEY"))

	uc := usecase.NewCampaignUsecase(repo)
	analyzeUC := usecase.NewAnalyzeCampaignUsecase(
		repo,
		analysisRepo,
		geminiService,
	)

	handler := http.NewCampaignHandler(uc, analyzeUC)

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

	auth.GET("campaigns", handler.GetCampaigns)
	auth.GET("analyses", handler.GetAnalyses)
	auth.GET("/campaigns/:id", handler.GetCampaignDetail)
	auth.GET("/campaigns/:id/analyze", handler.AnalyzeCampaign)
	auth.GET("/campaigns/:id/metrics", handler.GetCampaignMetrics)
	auth.POST("campaigns", handler.CreateCampaign)

	r.Run(":8080")
}
