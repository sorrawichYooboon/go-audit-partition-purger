package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/sorrawichyooboon/go-audit-partition-purger/internal/infrastructure/database"
	httpRouter "github.com/sorrawichyooboon/go-audit-partition-purger/internal/infrastructure/http"
	"github.com/sorrawichyooboon/go-audit-partition-purger/internal/infrastructure/http/handler"
	"github.com/sorrawichyooboon/go-audit-partition-purger/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Go Audit Partition Purger API
// @version         1.0
// @description     API for testing Audit Log insertions into partitioned DB.

// @host      localhost:8080
// @BasePath  /
func main() {
	fmt.Println("🚀 Starting Go Audit Partition Purger API...")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Proceeding with system environment variables...")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	fmt.Println("Connected to PostgreSQL successfully")

	repo := database.NewAuditRepository(db)
	auditUsecase := usecase.NewAuditUsecase(repo)
	auditHandler := handler.NewAuditHandler(auditUsecase)

	r := httpRouter.SetupRouter(auditHandler)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Swagger UI available at http://localhost:8080/swagger/index.html")

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
