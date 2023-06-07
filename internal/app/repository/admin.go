package repository

import (
	"context"

	"elible/internal/app/models"
	"elible/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		// log.Printf("Error while inserting new admin into db, Reason: %v\n", err)
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

	// Delete any existing token associated with the same access_uuid
	_, err := TokenCollection.DeleteOne(ctx, bson.M{"accessUUID": td.AccessUUID})
	if err != nil {
		// log.Printf("Error while deleting old token from db, Reason: %v\n", err)
		// You might want to handle this error, instead of just logging
	}

	// location, _ := time.LoadLocation("Asia/Jakarta")

	// // set created_at and updated_at fields
	// td.CreatedAt = time.Now().In(location)
	// td.UpdatedAt = time.Now().In(location)

	// Insert the new token
	_, err = TokenCollection.InsertOne(ctx, td)

	if err != nil {
		return err
	}

	// Insert the new token into history collection
	_, err = HistoryCollection.InsertOne(ctx, td)

	if err != nil {
		return err
	}

	return nil
}

func (r *AdminRepository) DeleteToken(accessUUID string) error {
	TokenCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_tokens")
	ctx := context.Background()

	_, err := TokenCollection.DeleteOne(ctx, bson.M{"accessUUID": accessUUID})

	if err != nil {
		return err
	}

	return nil
}

func (r *AdminRepository) FetchToken(accessUUID string) (*models.Token, error) {
	TokenCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_tokens")
	ctx := context.Background()

	var token models.Token

	err := TokenCollection.FindOne(ctx, bson.M{"accessUUID": accessUUID}).Decode(&token)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *AdminRepository) GetAdminByToken(tokens string) (*models.Admin, error) {
	AdminCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_admins")
	TokenCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_tokens")
	ctx := context.Background()


	var token models.Token
	err := TokenCollection.FindOne(ctx, bson.M{"accessToken": tokens}).Decode(&token)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	objectId, err := primitive.ObjectIDFromHex(token.AccessUUID)
	if err != nil {
		return nil, err
	}

	var admin models.Admin
	err = AdminCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&admin)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &admin, nil
}

func (r *AdminRepository) GetTokenByValue(tokenValue string) (*models.Token, error) {
	TokenCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_tokens")
	ctx := context.Background()

	var token models.Token
	err := TokenCollection.FindOne(ctx, bson.M{"access_token": tokenValue}).Decode(&token)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}
