package repository

import (
    "context"
    models "student-performance-report/app/models/mongodb"
    "github.com/google/uuid"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type AchievementRepository interface {
    GetStudentAchievements(studentId uuid.UUID) ([]models.Achievement, error)
}

type achievementRepository struct {
    collection *mongo.Collection
}

func NewAchievementRepository(mongodb *mongo.Database) AchievementRepository {
    return &achievementRepository{
        collection: mongodb.Collection("achievements"),
    }
}

func (r *achievementRepository) GetStudentAchievements(studentId uuid.UUID) ([]models.Achievement, error) {
    ctx := context.Background()

    filter := bson.M{"studentId": studentId.String()}
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }

    var results []models.Achievement
    err = cursor.All(ctx, &results)
    return results, err
}
