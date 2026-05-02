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
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	}

	db := database.NewPostgresPool(os.Getenv("DATABASE_URL"))

	repo := database.NewCampaignRepositoryPostgres(db)
	uc := usecase.NewCampaignUsecase(repo)
	handler := http.NewCampaignHandler(uc)

	r := gin.Default()

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	r.GET("/campaigns", handler.GetCampaigns)

	r.POST("/campaigns", handler.CreateCampaign)

	r.Run(":8080")
}
