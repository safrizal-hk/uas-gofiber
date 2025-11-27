package service

import (
	"context"
	"time"
	"errors"

	"github.com/gofiber/fiber/v2"
	modelMongo "github.com/safrizal-hk/uas-gofiber/app/model/mongo"
	modelPostgres "github.com/safrizal-hk/uas-gofiber/app/model/postgre"
	repoMongo "github.com/safrizal-hk/uas-gofiber/app/repository/mongo"
	repoPostgres "github.com/safrizal-hk/uas-gofiber/app/repository/postgre"
	"github.com/safrizal-hk/uas-gofiber/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService struct {
	MongoRepo repoMongo.AchievementMongoRepository
	PgRepo repoPostgres.AchievementPGRepository
}

func NewAchievementService(mongoRepo repoMongo.AchievementMongoRepository, pgRepo repoPostgres.AchievementPGRepository) *AchievementService {
	return &AchievementService{
		MongoRepo: mongoRepo,
		PgRepo: pgRepo,
	}
}

// SubmitPrestasi menangani alur FR-003
func (s *AchievementService) SubmitPrestasi(c *fiber.Ctx) error {
	profile := middleware.GetUserProfileFromContext(c)
	if profile.Role != "Mahasiswa" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Hanya Mahasiswa yang dapat submit prestasi", "code": "403"})
	}
	
	studentID, err := s.PgRepo.FindStudentIdByUserID(profile.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Kesalahan server saat mencari data mahasiswa.", "code": "500"})
	}
	if studentID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Akun tidak terhubung ke data Mahasiswa.", "code": "403"})
	}

	req := new(modelMongo.AchievementInput)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Input prestasi tidak valid", "code": "400"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoAch := modelMongo.AchievementMongo{
		StudentID: studentID, 
		AchievementType: req.AchievementType,
		Title: req.Title,
		Description: req.Description,
		Details: req.Details, 
		Attachments: req.Attachments,
		Tags: req.Tags,
		Points: req.Points,
	}
	createdMongo, err := s.MongoRepo.Create(ctx, &mongoAch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan detail ke MongoDB", "code": "500"})
	}

	pgRef := modelPostgres.AchievementReference{
		StudentID: studentID, 
		MongoAchievementID: createdMongo.ID.Hex(),
	}
	createdRef, err := s.PgRepo.CreateReference(&pgRef)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan referensi ke PostgreSQL", "code": "500"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"message": "Prestasi berhasil disimpan sebagai DRAFT",
		"id": createdRef.ID,
		"mongo_id": createdMongo.ID.Hex(),
	})
}


// SubmitForVerification menangani alur FR-004
func (s *AchievementService) SubmitForVerification(c *fiber.Ctx) error {
	profile := middleware.GetUserProfileFromContext(c)
	achievementID := c.Params("id") // Ambil ID dari URL

	if profile.Role != "Mahasiswa" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Hanya Mahasiswa yang dapat submit untuk verifikasi", "code": "403"})
    }

    updatedRef, err := s.PgRepo.UpdateStatusToSubmitted(achievementID)
    
    if err != nil {
		if errors.Is(err, errors.New("prestasi tidak ditemukan atau status sudah berubah")) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error(), "code": "404"})
		}
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error(), "code": "500"})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "message": "Prestasi berhasil disubmit untuk verifikasi",
        "new_status": updatedRef.Status,
    })
}

// DeletePrestasi menangani alur FR-005
func (s *AchievementService) DeletePrestasi(c *fiber.Ctx) error {
	profile := middleware.GetUserProfileFromContext(c)
	achievementID := c.Params("id") 

	if profile.Role != "Mahasiswa" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Hanya Mahasiswa yang dapat menghapus prestasi", "code": "403"})
	}

	studentID, err := s.PgRepo.FindStudentIdByUserID(profile.ID)
	if err != nil || studentID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Akun tidak terhubung ke data Mahasiswa.", "code": "403"})
	}

	deletedRef, err := s.PgRepo.SoftDeleteReference(achievementID, studentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "code": "400"})
	}

	mongoObjectID, err := primitive.ObjectIDFromHex(deletedRef.MongoAchievementID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengonversi Mongo ID.", "code": "500"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err = s.MongoRepo.SoftDelete(ctx, mongoObjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menghapus detail di MongoDB.", "code": "500"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"message": "Prestasi berhasil dihapus (soft deleted).",
	})
}