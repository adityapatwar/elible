package repository

import (
	"context"
	"elible/internal/app/models"
	"elible/internal/config"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UniversityRepository struct {
	MongoClient *mongo.Client
	cfg         *config.Config
}


func NewUniversityRepository(cfg *config.Config, mongoClient *mongo.Client) *UniversityRepository {

	return &UniversityRepository{
		cfg:         cfg,
		MongoClient: mongoClient,
	}
}


func (r *UniversityRepository) CreateUniversity(u models.University) (primitive.ObjectID, error) {
	UniversityCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_universities")
	ctx := context.Background()

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	// Check if university name already exists
	var existingUniversity models.University
	err := UniversityCollection.FindOne(ctx, bson.M{"name": u.Name}).Decode(&existingUniversity)
	if err == nil {
		return primitive.NilObjectID, errors.New("university with the same name already exists")
	}

	result, err := UniversityCollection.InsertOne(ctx, u)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *UniversityRepository) UpdateUniversity(id primitive.ObjectID, u models.University) error {
	UniversityCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_universities")
	ctx := context.Background()

	u.UpdatedAt = time.Now()

	_, err := UniversityCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": u})
	return err
}

func (r *UniversityRepository) DeleteUniversity(id primitive.ObjectID) error {
	UniversityCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_universities")
	ctx := context.Background()

	_, err := UniversityCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *UniversityRepository) GetUniversity(id primitive.ObjectID) (models.University, error) {
	UniversityCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_universities")
	ctx := context.Background()

	var u models.University
	err := UniversityCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	return u, err
}

func (r *UniversityRepository) GetUniversityByName(name string) (models.University, error) {
	UniversityCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_universities")
	ctx := context.Background()

	var u models.University
	err := UniversityCollection.FindOne(ctx, bson.M{"name": name}).Decode(&u)
	return u, err
}

func (r *UniversityRepository) GetUniversities() ([]models.University, error) {
	UniversityCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_universities")
	ctx := context.Background()

	cursor, err := UniversityCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var universities []models.University
	if err := cursor.All(ctx, &universities); err != nil {
		return nil, err
	}

	return universities, nil
}
