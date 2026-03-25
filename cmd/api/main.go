package main

import (
	"log"
	"os"

	"github.com/example/product-api/internal/handler"
	"github.com/example/product-api/internal/repository"
	"github.com/example/product-api/internal/usecase"
	"github.com/example/product-api/pkg/database"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Product API
// @version         1.0
// @description     REST API for product management
// @host            localhost:8080
// @BasePath        /
func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Dependency injection (clean architecture wiring)
	productRepo := repository.NewProductRepository(db)
	productUC := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUC)

	r := gin.Default()

	// Swagger docs at /api-docs
	r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register product routes
	productHandler.RegisterRoutes(r)

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
