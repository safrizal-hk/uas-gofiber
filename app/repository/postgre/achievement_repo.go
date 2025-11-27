package repository

import (
	"database/sql"
	"errors"

	model_postgre "github.com/safrizal-hk/uas-gofiber/app/model/postgre" 
)

type AchievementPGRepository interface {
	CreateReference(ref *model_postgre.AchievementReference) (*model_postgre.AchievementReference, error)
	FindStudentIdByUserID(userID string) (string, error) 
	UpdateStatusToSubmitted(refID string) (*model_postgre.AchievementReference, error)
	SoftDeleteReference(refID string, studentID string) (*model_postgre.AchievementReference, error) // ⚠️ KONTRAK BARU
}

type achievementPGRepositoryImpl struct {
	DB *sql.DB 
}

func NewAchievementPGRepository(db *sql.DB) AchievementPGRepository {
	return &achievementPGRepositoryImpl{DB: db}
}

// CreateReference membuat entri baru di achievement_references (PG) (FR-003)
func (r *achievementPGRepositoryImpl) CreateReference(ref *model_postgre.AchievementReference) (*model_postgre.AchievementReference, error) {
	query := `
		INSERT INTO achievement_references 
		(student_id, mongo_achievement_id, status)
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
	`
	ref.Status = model_postgre.StatusDraft 
	
	err := r.DB.QueryRow(query, ref.StudentID, ref.MongoAchievementID, ref.Status).Scan(
		&ref.ID, 
		&ref.CreatedAt, 
		&ref.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

func (r *achievementPGRepositoryImpl) FindStudentIdByUserID(userID string) (string, error) {
	var studentID string
	query := `SELECT id FROM students WHERE user_id = $1`
	err := r.DB.QueryRow(query, userID).Scan(&studentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil 
		}
		return "", err
	}
	return studentID, nil
}

// UpdateStatusToSubmitted mengubah status dari 'draft' menjadi 'submitted' (FR-004)
func (r *achievementPGRepositoryImpl) UpdateStatusToSubmitted(refID string) (*model_postgre.AchievementReference, error) {
    ref := new(model_postgre.AchievementReference)
    
    query := `
        UPDATE achievement_references 
        SET status = $1, submitted_at = NOW(), updated_at = NOW()
        WHERE id = $2 AND status = $3
        RETURNING id, student_id, mongo_achievement_id, status, submitted_at, updated_at
    `
    err := r.DB.QueryRow(query, model_postgre.StatusSubmitted, refID, model_postgre.StatusDraft).Scan(
        &ref.ID, 
        &ref.StudentID,
        &ref.MongoAchievementID,
        &ref.Status,
        &ref.SubmittedAt,
        &ref.UpdatedAt,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("prestasi tidak ditemukan atau status sudah berubah")
        }
        return nil, err
    }
    return ref, nil
}

// SoftDeleteReference menandai referensi prestasi di PG sebagai terhapus (FR-005)
func (r *achievementPGRepositoryImpl) SoftDeleteReference(refID string, studentID string) (*model_postgre.AchievementReference, error) {
    ref := new(model_postgre.AchievementReference)
    
    query := `
        UPDATE achievement_references 
        SET status = $1, updated_at = NOW() -- ⚠️ HAPUS SET deleted_at
        WHERE id = $2 AND student_id = $3 AND status = $4
        RETURNING id, student_id, mongo_achievement_id, status, created_at, updated_at, submitted_at, verified_at, verified_by, rejection_note
    `
    err := r.DB.QueryRow(query, 
        model_postgre.StatusDeleted,
        refID, 
        studentID, 
        model_postgre.StatusDraft,
    ).Scan(
        &ref.ID, 
        &ref.StudentID,
        &ref.MongoAchievementID,
        &ref.Status,
        &ref.CreatedAt,
        &ref.UpdatedAt,

        &ref.SubmittedAt,
        &ref.VerifiedAt,
        &ref.VerifiedBy,
        &ref.RejectionNote,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("prestasi tidak ditemukan, bukan milik Anda, atau status sudah disubmit/diverifikasi")
        }
        return nil, err
    }
    return ref, nil
}