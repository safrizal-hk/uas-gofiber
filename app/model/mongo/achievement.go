package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AchievementInput merepresentasikan data yang dikirimkan Mahasiswa melalui API
type AchievementInput struct {
	AchievementType string             `json:"achievementType" validate:"required"`
	Title           string             `json:"title" validate:"required"`
	Description     string             `json:"description"`
	Details         map[string]interface{} `json:"details"` // Field dinamis
	Attachments     []Attachment       `json:"attachments"`
	Tags            []string           `json:"tags"`
	Points          int                `json:"points"`
}

// AchievementMongo merepresentasikan dokumen yang disimpan di MongoDB
type AchievementMongo struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID       string             `bson:"studentId" json:"student_id"` // UUID dari PG
	AchievementType string             `bson:"achievementType" json:"achievement_type"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Details         primitive.M        `bson:"details" json:"details"` 
	Attachments     []Attachment       `bson:"attachments" json:"attachments"`
	Tags            []string           `bson:"tags" json:"tags"`
	Points          int                `bson:"points" json:"points"`
	CreatedAt       time.Time          `bson:"createdAt" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updated_at"`

	DeletedAt 		*time.Time 			`bson:"deletedAt,omitempty" json:"deleted_at,omitempty"`
}

// Attachment adalah struct internal untuk file pendukung
type Attachment struct {
	FileName string    `bson:"fileName" json:"file_name"`
	FileUrl  string    `bson:"fileUrl" json:"file_url"`
	FileType string    `bson:"fileType" json:"file_type"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploaded_at"`
}