package repository

import (
    "context"
	"time"
    models "student-performance-report/app/models/mongodb"
    "github.com/google/uuid"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementRepository interface {
    GetStudentAchievements(studentId uuid.UUID) ([]models.Achievement, error)
    InsertOne(ctx context.Context, achievement models.Achievement) (string, error)
    FindAllDetails(ctx context.Context, mongoIDs []string) ([]models.Achievement, error)
	FindOne(ctx context.Context, mongoID string) (*models.Achievement, error)
	DeleteAchievement(ctx context.Context, mongoID string) error
	UpdateOne(ctx context.Context, mongoID string, data models.Achievement) error
	AddAttachment(ctx context.Context, mongoID string, attachment models.Attachment) error
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

func (r *achievementRepository) InsertOne(ctx context.Context, achievement models.Achievement) (string, error) {
	collection := r.collection
	result, err := collection.InsertOne(ctx, achievement)
	if err != nil {
		return "", err
	}
	// Mengembalikan ID object mongo sebagai string
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *achievementRepository) FindAllDetails(ctx context.Context, mongoIDs []string) ([]models.Achievement, error) {
    // Convert string IDs to ObjectIDs
	var objectIDs []primitive.ObjectID
	for _, id := range mongoIDs {
		oid, _ := primitive.ObjectIDFromHex(id)
		objectIDs = append(objectIDs, oid)
	}

	collection := r.collection
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []models.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}
	return achievements, nil
}

func (r *achievementRepository) FindOne(ctx context.Context, mongoID string) (*models.Achievement, error) {
    oid, err := primitive.ObjectIDFromHex(mongoID)
    if err != nil {
        return nil, err
    }

    var result models.Achievement
    filter := bson.M{"_id": oid}
    
    err = r.collection.FindOne(ctx, filter).Decode(&result)
    if err != nil {
        return nil, err
    }
    
    return &result, nil
}

func (r *achievementRepository) DeleteAchievement(ctx context.Context, mongoID string) error {
    oid, err := primitive.ObjectIDFromHex(mongoID)
    if err != nil {
        return err
    }

    filter := bson.M{"_id": oid}
    _, err = r.collection.DeleteOne(ctx, filter)
    return err
}

func (r *achievementRepository) UpdateOne(ctx context.Context, mongoID string, data models.Achievement) error {
    oid, err := primitive.ObjectIDFromHex(mongoID)
    if err != nil { return err }

    // Kita update field-field tertentu
    update := bson.M{
        "$set": bson.M{
            "title":           data.Title,
            "description":     data.Description,
            "achievementType": data.AchievementType,
            "details":         data.Details,
            "tags":            data.Tags,
            "points":          data.Points,
            "updatedAt":       time.Now(),
        },
    }

    _, err = r.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
    return err
}

func (r *achievementRepository) AddAttachment(ctx context.Context, mongoID string, attachment models.Attachment) error {
    oid, err := primitive.ObjectIDFromHex(mongoID)
    if err != nil { return err }

    // Gunakan operator $push untuk menambah item ke array attachments
    update := bson.M{
        "$push": bson.M{"attachments": attachment},
        "$set":  bson.M{"updatedAt": time.Now()},
    }

    _, err = r.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
    return err
}

