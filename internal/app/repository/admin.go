// internal/app/repository/admin_repository.go
package repository

import (
	"context"
	"log"
	"time"

	"elible/internal/app/models"
	"elible/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepository struct {
	MongoClient *mongo.Client
	cfg         *config.Config
}

func NewAdminRepository(cfg *config.Config, mongoClient *mongo.Client) *AdminRepository {

	return &AdminRepository{
		cfg:         cfg,
		MongoClient: mongoClient,
	}
}

func (r *AdminRepository) Create(admin *models.Admin) error {
	AdminCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_admins")
	ctx := context.Background()

	_, err := AdminCollection.InsertOne(ctx, admin)
	if err != nil {
		log.Printf("Error while inserting new admin into db, Reason: %v\n", err)
		return err
	}
	return nil
}

func (r *AdminRepository) FindByUsername(username string) (*models.Admin, error) {
	AdminCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_admins")
	ctx := context.Background()

	var admin models.Admin
	err := AdminCollection.FindOne(ctx, bson.M{"username": username}).Decode(&admin)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &admin, nil
}

func (r *AdminRepository) SaveToken(td *models.Token) error {
	TokenCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_tokens")
	HistoryCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_hitory_tokens")
	ctx := context.Background()

	at := time.Unix(td.AtExpires, 0) // converting Unix timestamp to UTC
	now := time.Now()

	// Delete any existing token associated with the same access_uuid
	_, err := TokenCollection.DeleteOne(ctx, bson.M{"access_uuid": td.AccessUUID})
	if err != nil {
		log.Printf("Error while deleting old token from db, Reason: %v\n", err)
		// You might want to handle this error, instead of just logging
	}

	// Insert the new token
	_, err = TokenCollection.InsertOne(ctx, bson.M{
		"access_token": td.AccessToken,
		"access_uuid":  td.AccessUUID,
		"at_expires":   at.UTC(),
		"created_at":   now.UTC(),
	})

	if err != nil {
		return err
	}

	// Insert the new token into history collection
	_, err = HistoryCollection.InsertOne(ctx, bson.M{
		"access_token": td.AccessToken,
		"access_uuid":  td.AccessUUID,
		"at_expires":   at.UTC(),
		"created_at":   now.UTC(),
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *AdminRepository) DeleteToken(accessUUID string) error {
	TokenCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_tokens")
	ctx := context.Background()

	_, err := TokenCollection.DeleteOne(ctx, bson.M{"access_uuid": accessUUID})

	if err != nil {
		return err
	}

	return nil
}

func (r *AdminRepository) FetchToken(accessUUID string) (string, error) {
	TokenCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("gpt_client_counter")
	ctx := context.Background()

	var token bson.M

	err := TokenCollection.FindOne(ctx, bson.M{"access_uuid": accessUUID}).Decode(&token)

	if err != nil {
		return "", err
	}

	return token["access_token"].(string), nil
}
