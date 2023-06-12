package repository

import (
	"context"

	"elible/internal/app/models"
	"elible/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KnowledgeBaseRepository struct {
	MongoClient *mongo.Client
	cfg         *config.Config
}

func NewKnowledgeBaseRepository(cfg *config.Config, mongoClient *mongo.Client) *KnowledgeBaseRepository {

	return &KnowledgeBaseRepository{
		cfg:         cfg,
		MongoClient: mongoClient,
	}
}

func (r *KnowledgeBaseRepository) CreateKnowledgeBase(knowledgeBase *models.KnowledgeBase) error {
	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")
	ctx := context.Background()

	_, err := KnowledgeBaseCollection.InsertOne(ctx, knowledgeBase)
	if err != nil {
		return err
	}
	return nil
}

func (r *KnowledgeBaseRepository) DeleteKnowledgeBase(id primitive.ObjectID) error {
	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")
	ctx := context.Background()

	_, err := KnowledgeBaseCollection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return err
	}

	return nil
}

func (r *KnowledgeBaseRepository) UpdateKnowledgeBase(knowledgeBase *models.KnowledgeBase) error {
	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")
	ctx := context.Background()

	_, err := KnowledgeBaseCollection.UpdateOne(ctx, bson.M{"_id": knowledgeBase.ID}, bson.M{"$set": knowledgeBase})

	if err != nil {
		return err
	}

	return nil
}

func (r *KnowledgeBaseRepository) ListKnowledgeBase() ([]models.KnowledgeBase, error) {
	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")
	ctx := context.Background()

	cursor, err := KnowledgeBaseCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	var knowledgeBases []models.KnowledgeBase

	err = cursor.All(ctx, &knowledgeBases)

	if err != nil {
		return nil, err
	}

	return knowledgeBases, nil
}

func (r *KnowledgeBaseRepository) AddKnowledgeProgram(id primitive.ObjectID, program models.KnowledgeProgram) error {
	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")
	ctx := context.Background()

	_, err := KnowledgeBaseCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$push": bson.M{"programs": program}})

	if err != nil {
		return err
	}

	return nil
}

func (r *KnowledgeBaseRepository) RemoveKnowledgeProgram(id primitive.ObjectID, programName string) error {
	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")
	ctx := context.Background()

	_, err := KnowledgeBaseCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$pull": bson.M{"programs": bson.M{"name": programName}}})

	if err != nil {
		return err
	}

	return nil
}

func (r *KnowledgeBaseRepository) UpdateKnowledgeProgram(id primitive.ObjectID, oldName string, updatedProgram models.KnowledgeProgram) error {
	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")
	ctx := context.Background()

	filter := bson.M{
		"_id": id,
	}

	update := bson.M{
		"$set": bson.M{
			"programs.$[elem]": updatedProgram,
		},
	}

	arrayFilter := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.name": oldName},
		},
	})

	_, err := KnowledgeBaseCollection.UpdateOne(ctx, filter, update, arrayFilter)

	if err != nil {
		return err
	}

	return nil
}

func (r *KnowledgeBaseRepository) ListKnowledgePrograms(id primitive.ObjectID) ([]models.KnowledgeProgram, error) {
	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")
	ctx := context.Background()

	var knowledgeBase models.KnowledgeBase

	err := KnowledgeBaseCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&knowledgeBase)

	if err != nil {
		return nil, err
	}

	return knowledgeBase.Programs, nil
}
