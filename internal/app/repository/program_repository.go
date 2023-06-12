package repository

import (
	"context"
	"elible/internal/app/models"
	"elible/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudyProgramRepository struct {
	MongoClient *mongo.Client
	cfg         *config.Config
}

func NewStudyProgramRepository(cfg *config.Config, mongoClient *mongo.Client) *StudyProgramRepository {

	return &StudyProgramRepository{
		cfg:         cfg,
		MongoClient: mongoClient,
	}
}

// CreateStudyProgram creates a new study program and adds it to a specified KnowledgeBase
func (r *StudyProgramRepository) CreateStudyProgram(sp models.StudyProgram, kbYear string, kpName string) (primitive.ObjectID, error) {
	ctx := context.Background()

	sp.CreatedAt = time.Now()
	sp.UpdatedAt = time.Now()

	StudyProgramCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_study_programs")
	result, err := StudyProgramCollection.InsertOne(ctx, sp)
	if err != nil {
		return primitive.NilObjectID, err
	}

	spID := result.InsertedID.(primitive.ObjectID)

	// Adding the created study program to the specified KnowledgeProgram in a KnowledgeBase
	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")

	_, err = KnowledgeBaseCollection.UpdateOne(
		ctx,
		bson.M{"year": kbYear, "programs.name": kpName},
		bson.M{"$push": bson.M{"programs.$.study_programs": spID}},
	)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return spID, nil
}

// UpdateStudyProgram updates a study program and removes it from all KnowledgeBases it belonged to
func (r *StudyProgramRepository) UpdateStudyProgram(id primitive.ObjectID, sp models.StudyProgram) error {
	ctx := context.Background()

	sp.UpdatedAt = time.Now()

	StudyProgramCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_study_programs")
	_, err := StudyProgramCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": sp})
	if err != nil {
		return err
	}

	// Removing the updated study program from all KnowledgeBases it belonged to
	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")
	_, err = KnowledgeBaseCollection.UpdateMany(
		ctx,
		bson.M{"programs.study_programs": id},
		bson.M{"$pull": bson.M{"programs.$.study_programs": id}},
	)

	return err
}

// DeleteStudyProgram deletes a study program and removes it from all KnowledgeBases it belonged to
func (r *StudyProgramRepository) DeleteStudyProgram(id primitive.ObjectID) error {
	ctx := context.Background()

	StudyProgramCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_study_programs")
	_, err := StudyProgramCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	// Removing the deleted study program from all KnowledgeBases it belonged to
	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")
	_, err = KnowledgeBaseCollection.UpdateMany(
		ctx,
		bson.M{"programs.study_programs": id},
		bson.M{"$pull": bson.M{"programs.$[].study_programs": id}},
	)

	return err
}

// GetStudyProgram retrieves a study program by its ID
func (r *StudyProgramRepository) GetStudyProgram(id primitive.ObjectID) (models.StudyProgram, error) {
	ctx := context.Background()

	StudyProgramCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_study_programs")

	var sp models.StudyProgram
	err := StudyProgramCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&sp)
	if err != nil {
		return models.StudyProgram{}, err
	}

	return sp, nil
}

// GetStudyPrograms retrieves all study programs from a specific KnowledgeProgram in a KnowledgeBase
func (r *StudyProgramRepository) GetStudyPrograms(kbYear string, kpName string) ([]models.StudyProgram, error) {
	ctx := context.Background()

	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")

	var kb models.KnowledgeBase
	err := KnowledgeBaseCollection.FindOne(ctx, bson.M{"year": kbYear}).Decode(&kb)
	if err != nil {
		return nil, err
	}

	// Find the specified KnowledgeProgram
	var kp models.KnowledgeProgram
	for _, program := range kb.Programs {
		if program.Name == kpName {
			kp = program
			break
		}
	}

	StudyProgramCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_study_programs")

	// Find all study programs of the KnowledgeProgram
	var programs []models.StudyProgram
	for _, spID := range kp.StudyPrograms {
		var sp models.StudyProgram
		err = StudyProgramCollection.FindOne(ctx, bson.M{"_id": spID}).Decode(&sp)
		if err != nil {
			return nil, err
		}
		programs = append(programs, sp)
	}

	return programs, nil
}
