package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/safrizal-hk/uas-gofiber/app/service" 
	repo_postgre "github.com/safrizal-hk/uas-gofiber/app/repository/postgre" // Alias untuk PostgreSQL Repo
	repo_mongo "github.com/safrizal-hk/uas-gofiber/app/repository/mongo"       // Alias untuk MongoDB Repo
	"github.com/safrizal-hk/uas-gofiber/config"
)

func RegisterAllRoutes(app *fiber.App, dbConn *config.Database) {
	
	v1 := app.Group("/api/v1") 

	// 1. Wiring Auth & Core Services (PostgreSQL)
	authRepo := repo_postgre.NewAuthRepository(dbConn.PgDB) 
	authService := service.NewAuthService(authRepo) // Auth Service ada di service/ utama
	
	// 2. Wiring Achievement Services (Hybrid DB)
	achievementPgRepo := repo_postgre.NewAchievementPGRepository(dbConn.PgDB)
	achievementMongoRepo := repo_mongo.NewAchievementMongoRepository(dbConn.MongoDB)
	achievementService := service.NewAchievementService(achievementMongoRepo, achievementPgRepo) // Service di folder postgre

	
	// ---------- Pendaftaran Grup Rute ----------
	
	// A. Rute Otentikasi (Public)
	RegisterAuthRoutes(v1, authService) 
	
	// B. Rute Prestasi (Terproteksi)
	RegisterAchievementRoutes(v1, achievementService)
	
	// ... rute lainnya
}