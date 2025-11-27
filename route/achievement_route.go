package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/safrizal-hk/uas-gofiber/app/service" 
	"github.com/safrizal-hk/uas-gofiber/middleware"      
)

func RegisterAchievementRoutes(v1 fiber.Router, achievementService *service.AchievementService) {
	protected := v1.Group("/achievements", middleware.AuthRequired) 
	
	// POST /api/v1/achievements (Submit Prestasi)
	protected.Post("/", middleware.RBACRequired("achievement:create"), achievementService.SubmitPrestasi)
	protected.Post("/:id/submit", middleware.RBACRequired("achievement:update"), achievementService.SubmitForVerification)
	protected.Delete("/:id", middleware.RBACRequired("achievement:delete"), achievementService.DeletePrestasi)
	// GET /api/v1/achievements (List Prestasi)
	// protected.Get("/", middleware.RBACRequired("achievement:read"), achievementService.ListPrestasi)

	// TODO: Tambahkan rute CRUD, Verify, dan Reject di sini
	// protected.Put("/:id", middleware.RBACRequired("achievement:update"), achievementService.UpdatePrestasi)
}