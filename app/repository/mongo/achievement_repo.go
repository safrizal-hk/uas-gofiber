package repository

import (
	"context"
	"time"

	model_mongo "github.com/safrizal-hk/uas-gofiber/app/model/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementMongoRepository interface {
	Create(ctx context.Context, achievement *model_mongo.AchievementMongo) (*model_mongo.AchievementMongo, error)
	SoftDelete(ctx context.Context, id primitive.ObjectID) error
}

type achievementMongoRepositoryImpl struct {
	Collection *mongo.Collection
}

func NewAchievementMongoRepository(db *mongo.Database) AchievementMongoRepository {
	return &achievementMongoRepositoryImpl{
		Collection: db.Collection("achievements"),
	}
}

func (r *achievementMongoRepositoryImpl) Create(ctx context.Context, achievement *model_mongo.AchievementMongo) (*model_mongo.AchievementMongo, error) {
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	result, err := r.Collection.InsertOne(ctx, achievement)
	if err != nil {
		return nil, err
	}
	achievement.ID = result.InsertedID.(primitive.ObjectID)
	return achievement, nil
}

func (r *achievementMongoRepositoryImpl) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
    now := time.Now()
    update := bson.M{
        "$set": bson.M{
            "deletedAt": now,
            "updatedAt": now,
        },
    }
    _, err := r.Collection.UpdateByID(ctx, id, update)
    return err
}