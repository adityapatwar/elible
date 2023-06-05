// internal/app/repository/student_repository.go
package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"elible/internal/app/models"
	"elible/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudentRepository struct {
	MongoClient *mongo.Client
	cfg         *config.Config
}

func NewStudentRepository(cfg *config.Config, mongoClient *mongo.Client) *StudentRepository {
	return &StudentRepository{
		cfg:         cfg,
		MongoClient: mongoClient,
	}
}

func (r *StudentRepository) Create(student *models.Student) error {
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	ctx := context.Background()

	// Use Jakarta's time zone
	location, _ := time.LoadLocation("Asia/Jakarta")

	// set created_at and updated_at fields
	student.CreatedAt = time.Now().In(location)
	student.UpdatedAt = time.Now().In(location)

	_, err := studentCollection.InsertOne(ctx, student)
	if err != nil {
		log.Printf("Error while inserting new student into db, Reason: %v\n", err)
		return err
	}
	return nil
}

func (r *StudentRepository) GetAll() ([]*models.Student, error) {
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	ctx := context.Background()

	cursor, err := studentCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var students []*models.Student
	if err = cursor.All(ctx, &students); err != nil {
		return nil, err
	}

	return students, nil
}

func (r *StudentRepository) Delete(studentID primitive.ObjectID) error {
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	ctx := context.Background()

	_, err := studentCollection.DeleteOne(ctx, bson.M{"_id": studentID})

	if err != nil {
		return err
	}

	return nil
}

func (r *StudentRepository) Deactivate(studentID primitive.ObjectID) error {
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	ctx := context.Background()

	_, err := studentCollection.UpdateOne(ctx, bson.M{"_id": studentID}, bson.M{"$set": bson.M{"is_active": false, "updated_at": time.Now()}})

	if err != nil {
		return err
	}

	return nil
}

func (r *StudentRepository) Update(studentID primitive.ObjectID, student *models.Student) error {
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	ctx := context.Background()
	// Use Jakarta's time zone
	location, _ := time.LoadLocation("Asia/Jakarta")
	// update the updated_at field
	student.UpdatedAt = time.Now().In(location)

	update := bson.M{
		"$set": student,
	}

	_, err := studentCollection.UpdateOne(ctx, bson.M{"_id": studentID}, update)

	if err != nil {
		return err
	}

	return nil
}

func (r *StudentRepository) AddService(studentID primitive.ObjectID, service *models.TrackRecord) error {
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	ctx := context.Background()

	// Find the student by ID
	var student models.Student
	err := studentCollection.FindOne(ctx, bson.M{"_id": studentID}).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("student not found")
		}
		return err
	}

	// Check if service already exists in track_records
	for _, record := range student.TrackRecords {
		if record.ServiceName == service.ServiceName {
			return errors.New("this service already exists")
		}
	}

	// Use Jakarta's time zone
	location, _ := time.LoadLocation("Asia/Jakarta")

	// Update the service's UpdatedAt field
	service.UpdatedAt = time.Now().In(location)

	// If ServiceDate is zero value, set it to current time
	if service.ServiceDate.IsZero() {
		service.ServiceDate = time.Now().In(location)
	}

	update := bson.M{
		"$push": bson.M{
			"track_records": service,
		},
		"$set": bson.M{
			"updated_at": time.Now().In(location),
		},
	}

	_, err = studentCollection.UpdateOne(ctx, bson.M{"_id": studentID}, update)

	if err != nil {
		return err
	}

	return nil
}

func (r *StudentRepository) FindByUsername(username string) (*models.Student, error) {
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	ctx := context.Background()

	var student models.Student
	err := studentCollection.FindOne(ctx, bson.M{"name": username}).Decode(&student)
	if err != nil {
		// Handle error when the student is not found
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		// Handle other errors
		return nil, err
	}

	return &student, nil
}
