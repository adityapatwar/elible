// internal/app/repository/student_repository.go
package repository

import (
	"context"
	"errors"
	"time"

	"elible/internal/app/models"
	"elible/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// Check if index already exists
	cursor, err := studentCollection.Indexes().List(ctx)
	if err != nil {
		// log.Printf("Error while listing indexes, Reason: %v\n", err)
		return err
	}

	// Iterate through the returned cursor
	indexExists := false
	for cursor.Next(ctx) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			// log.Printf("Error while decoding index, Reason: %v\n", err)
			return err
		}

		// If "SearchIndex" exists, set flag to true
		if index["name"] == "SearchIndex" {
			indexExists = true
			break
		}
	}

	// Create text index if it doesn't exist
	if !indexExists {
		indexModel := mongo.IndexModel{
			Keys: bson.M{
				"$**": "text",
			},
			Options: options.Index().SetWeights(bson.M{
				"name":              1,
				"school":            1,
				"interest":          1,
				"gender":            1,
				"phone":             1,
				"financial_ability": 1,
				"progress":          1,
				"category":          1,
			}).SetName("SearchIndex"),
		}

		_, err = studentCollection.Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			// log.Printf("Error while creating text index, Reason: %v\n", err)
			return err
		}
	}

	// Use Jakarta's time zone
	location, _ := time.LoadLocation("Asia/Jakarta")

	// set created_at and updated_at fields
	student.CreatedAt = time.Now().In(location)
	student.UpdatedAt = time.Now().In(location)
	student.IsActive = true

	_, err = studentCollection.InsertOne(ctx, student)
	if err != nil {
		// log.Printf("Error while inserting new student into db, Reason: %v\n", err)
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

func (r *StudentRepository) GetAll(filter *models.StudentFilter) ([]*models.Student, error) {
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	ctx := context.Background()

	bsonFilter := make(bson.M)

	if filter != nil {
		if filter.Name != nil && *filter.Name != "" {
			// flexible match for name
			bsonFilter["name"] = bson.M{"$regex": primitive.Regex{Pattern: *filter.Name, Options: "i"}}
		}
		if filter.School != nil && *filter.School != "" {
			// flexible match for school
			bsonFilter["school"] = bson.M{"$regex": primitive.Regex{Pattern: *filter.School, Options: "i"}}
		}
		if filter.Interest != nil && *filter.Interest != "" {
			// strict match for interest
			bsonFilter["interest"] = *filter.Interest
		}
		if filter.Gender != nil && *filter.Gender != "" {
			// strict match for gender
			bsonFilter["gender"] = *filter.Gender
		}
		if filter.Phone != nil && *filter.Phone != "" {
			// flexible match for phone
			bsonFilter["phone"] = bson.M{"$regex": primitive.Regex{Pattern: *filter.Phone, Options: "i"}}
		}
		if filter.FinancialAbility != nil && *filter.FinancialAbility != "" {
			// strict match for financial ability
			bsonFilter["financial_ability"] = *filter.FinancialAbility
		}
		if filter.Progress != nil && *filter.Progress != "" {
			// strict match for progress
			bsonFilter["progress"] = *filter.Progress
		}
		if filter.Category != nil && *filter.Category != "" {
			// strict match for category
			bsonFilter["category"] = *filter.Category
		}
		if filter.IsActive != nil {
			// strict match for is_active
			bsonFilter["is_active"] = *filter.IsActive
		}
		if filter.Birthdate != nil && *filter.Birthdate != "" {
			// flexible match for birthdate
			bsonFilter["birthdate"] = bson.M{"$regex": primitive.Regex{Pattern: *filter.Birthdate, Options: "i"}}
		}
	}

	cursor, err := studentCollection.Find(ctx, bsonFilter)
	if err != nil {
		return nil, err
	}
	if cursor == nil {
		return nil, errors.New("cursor is nil")
	}

	var students []*models.Student
	if err = cursor.All(ctx, &students); err != nil {
		return nil, err
	}

	return students, nil
}

func (r *StudentRepository) GetByID(studentID primitive.ObjectID) (*models.Student, error) {
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	ctx := context.Background()

	var student models.Student
	err := studentCollection.FindOne(ctx, bson.M{"_id": studentID}).Decode(&student)
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

func (r *StudentRepository) AddService(studentID primitive.ObjectID, service *models.TrackRecord) error {
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	serviceCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_service_student")
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

	for _, existingService := range student.TrackRecords {
		if existingService.ServiceName == service.ServiceName {
			return errors.New("service already exists")
		}
	}

	// Use Jakarta's time zone
	location, _ := time.LoadLocation("Asia/Jakarta")

	service.CreatedAt = time.Now().In(location)
	service.UpdatedAt = time.Now().In(location)

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

	// Also update this service in tb_service_student collection
	serviceInStudent := models.Student{
		ID:   studentID,
		Name: student.Name,
		TrackRecords: []models.TrackRecord{
			*service,
		},
	}

	// Check if student already has a service record
	var existingService models.Student
	err = serviceCollection.FindOne(ctx, bson.M{"_id": studentID}).Decode(&existingService)

	if err == nil {
		// If record exists, update it
		update := bson.M{
			"$push": bson.M{
				"track_records": service,
			},
			"$set": bson.M{
				"updated_at": time.Now().In(location),
			},
		}

		_, err = serviceCollection.UpdateOne(ctx, bson.M{"_id": studentID}, update)

		if err != nil {
			return err
		}
	} else if err == mongo.ErrNoDocuments {
		// If no record exists, insert a new one
		_, err = serviceCollection.InsertOne(ctx, serviceInStudent)

		if err != nil {
			return err
		}
	} else {
		// Handle other errors from FindOne
		return err
	}

	r.EnsureServiceIndex()

	return nil
}

func (r *StudentRepository) UpdateService(studentID primitive.ObjectID, oldServiceName string, newService *models.TrackRecord) error {
	serviceCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_service_student")
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	ctx := context.Background()

	filter := bson.M{
		"_id":                        studentID,
		"track_records.service_name": oldServiceName,
	}

	update := bson.M{
		"$set": bson.M{
			"track_records.$": newService,
		},
	}

	_, err := serviceCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	_, err = studentCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *StudentRepository) DeleteService(studentID primitive.ObjectID, serviceName string) error {
	serviceCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_service_student")
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	ctx := context.Background()

	filter := bson.M{"_id": studentID}
	update := bson.M{
		"$pull": bson.M{
			"track_records": bson.M{"service_name": serviceName},
		},
	}

	_, err := serviceCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	_, err = studentCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *StudentRepository) EnsureServiceIndex() error {
	serviceCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_service_student")
	ctx := context.Background()

	// Check if index already exists
	cursor, err := serviceCollection.Indexes().List(ctx)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// Iterate through the returned cursor
	indexExists := false
	for cursor.Next(ctx) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			return err
		}

		// If "SearchIndex" exists, set flag to true
		if index["name"] == "SearchIndex" {
			indexExists = true
			break
		}
	}

	// Only create index if it doesn't exist
	if !indexExists {
		indexModel := mongo.IndexModel{
			Keys: bson.D{
				{Key: "name", Value: "text"},
				{Key: "track_records.service_name", Value: "text"},
				{Key: "track_records.service_date", Value: "text"},
				{Key: "track_records.status", Value: "text"},
			},
			Options: options.Index().SetName("SearchIndex"),
		}

		_, err = serviceCollection.Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *StudentRepository) FilterServices(filter models.ServiceFilter) ([]models.Student, error) {
	serviceCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_service_student")
	ctx := context.Background()

	// Build filter query
	query := bson.M{}
	if filter.Name != nil {
		query["name"] = bson.M{"$regex": filter.Name, "$options": "i"} // case insensitive search
	}
	if filter.ServiceName != nil {
		query["track_records.service_name"] = bson.M{"$regex": filter.ServiceName, "$options": "i"}
	}
	if filter.ServiceDate != nil {
		query["track_records.service_date"] = bson.M{"$regex": filter.ServiceDate, "$options": "i"}
	}
	if filter.Status != nil {
		query["track_records.status"] = bson.M{"$regex": filter.Status, "$options": "i"}
	}

	cursor, err := serviceCollection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var services []models.Student
	if err := cursor.All(ctx, &services); err != nil {
		return nil, err
	}

	return services, nil
}

func (r *StudentRepository) AddLobby(studentID primitive.ObjectID, lobby *models.Student) error {
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
	// for _, record := range student.TrackLobby {
	// 	if record.Progress == lobby.Progress {
	// 		return errors.New("this proggress already exists")
	// 	}
	// }

	// Use Jakarta's time zone
	location, _ := time.LoadLocation("Asia/Jakarta")

	lobby.CreatedAt = time.Now().In(location)
	lobby.UpdatedAt = time.Now().In(location)

	update := bson.M{
		"$push": bson.M{
			"track_lobby": lobby,
		},
		"$set": bson.M{
			"progress":   lobby.Progress,
			"updated_at": time.Now().In(location),
		},
	}

	_, err = studentCollection.UpdateOne(ctx, bson.M{"_id": studentID}, update)

	if err != nil {
		return err
	}

	return nil
}

func (r *StudentRepository) ActivateAll() error {
	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	ctx := context.Background()

	_, err := studentCollection.UpdateMany(ctx, bson.M{"is_active": bson.M{"$ne": true}}, bson.M{"$set": bson.M{"is_active": true, "updated_at": time.Now()}})

	if err != nil {
		return err
	}

	return nil
}
