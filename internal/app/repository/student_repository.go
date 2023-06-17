// internal/app/repository/student_repository.go
package repository

import (
	"context"
	"errors"
	"log"
	"math"
	"strings"
	"time"

	"elible/internal/app/models"
	"elible/internal/config"

	"github.com/xuri/excelize/v2"
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
	schoolCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_schools")
	ctx := context.Background()

	// New code: Validate if the student email already exists
	count, err := studentCollection.CountDocuments(ctx, bson.M{"email": student.Email})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("a user with this email already exists")
	}

	// New code: Find or create the school based on the student's school name
	var school models.School
	err = schoolCollection.FindOne(ctx, bson.M{"name": student.School}).Decode(&school)
	if err == mongo.ErrNoDocuments {
		// School not found, create a new one
		school.Name = student.School
		school.CreatedAt = time.Now()
		school.UpdatedAt = time.Now()
		res, err := schoolCollection.InsertOne(ctx, school)
		if err != nil {
			return err
		}
		school.ID = res.InsertedID.(primitive.ObjectID)
	} else if err != nil {
		return err
	}

	student.SchoolID = school.ID

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
	schoolCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_schools")
	ctx := context.Background()

	// New code: Validate if the student email already exists in another document
	count, err := studentCollection.CountDocuments(ctx, bson.M{"email": student.Email, "_id": bson.M{"$ne": studentID}})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("a user with this email already exists")
	}

	// New code: Find or create the school based on the student's school name
	var school models.School
	err = schoolCollection.FindOne(ctx, bson.M{"name": student.School}).Decode(&school)
	if err == mongo.ErrNoDocuments {
		// School not found, create a new one
		school.Name = student.School
		school.CreatedAt = time.Now()
		school.UpdatedAt = time.Now()
		res, err := schoolCollection.InsertOne(ctx, school)
		if err != nil {
			return err
		}
		school.ID = res.InsertedID.(primitive.ObjectID)
	} else if err != nil {
		return err
	}
	student.SchoolID = school.ID
	// Use Jakarta's time zone
	location, _ := time.LoadLocation("Asia/Jakarta")
	// update the updated_at field
	student.UpdatedAt = time.Now().In(location)

	update := bson.M{
		"$set": student,
	}

	_, err = studentCollection.UpdateOne(ctx, bson.M{"_id": studentID}, update)

	if err != nil {
		return err
	}

	return nil
}

func (r *StudentRepository) GetAll(filter *models.StudentFilter) (*models.PagedStudents, error) {
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

	findOptions := options.Find()
	if filter.Page != nil && filter.PageSize != nil {
		skip := int64((*filter.Page - 1) * *filter.PageSize)
		limit := int64(*filter.PageSize)
		findOptions.SetSkip(skip).SetLimit(limit)
	}

	cursor, err := studentCollection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, err
	}

	var students []models.Student
	if err = cursor.All(ctx, &students); err != nil {
		return nil, err
	}

	total, err := studentCollection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(*filter.PageSize)))

	return &models.PagedStudents{
		CurrentPage:  *filter.Page,
		TotalRecords: total,
		TotalPages:   totalPages,
		Records:      students,
	}, nil
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

func (r *StudentRepository) ImportDataFromExcelStudent(filePath string) (*models.ImportResultStudent, error) {
	var schoolCreatedCount, schoolUpdatedCount, schoolFailedCount, studentCreatedCount, studentUpdatedCount, studentFailedCount int
	var schoolFailedRows, studentFailedRows []int

	studentCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_students")
	schoolCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_schools")

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	for i, row := range rows {
		if i == 0 { // Skip header row
			continue
		}

		// Schools
		schoolName := strings.ToUpper(row[3])
		schoolFailed := false

		var school models.School
		err = schoolCollection.FindOne(ctx, bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: schoolName, Options: "i"}}}).Decode(&school)
		if err == mongo.ErrNoDocuments {
			schoolCreatedCount++
			school = models.School{
				ID:          primitive.NewObjectID(),
				Name:        row[3],
				Address:     row[4],
				Province:    strings.ToUpper(row[5]),
				City:        strings.ToUpper(row[6]),
				SchoolLogo:  row[7],
				SchoolImage: row[8],
				Phone:       row[9],
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			_, err := schoolCollection.InsertOne(ctx, school)
			if err != nil {
				schoolFailedCount++
				schoolFailedRows = append(schoolFailedRows, i+1)
				schoolFailed = true
			}
		} else if err != nil {
			schoolFailedCount++
			schoolFailedRows = append(schoolFailedRows, i+1)
			schoolFailed = true
		} else {
			schoolUpdatedCount++
			updatedSchool := bson.M{
				"name":        row[3],
				"address":     row[4],
				"province":    strings.ToUpper(row[5]),
				"city":        strings.ToUpper(row[6]),
				"schoolLogo":  row[7],
				"schoolImage": row[8],
				"phone":       row[9],
				"updatedAt":   time.Now(),
			}
			_, err = schoolCollection.UpdateOne(ctx, bson.M{"_id": school.ID}, bson.M{"$set": updatedSchool})
			if err != nil {
				schoolFailedCount++
				schoolFailedRows = append(schoolFailedRows, i+1)
				continue
			}
		}

		if schoolFailed {
			continue
		}

		var student models.Student
		err = studentCollection.FindOne(ctx, bson.M{"name": row[1], "school": row[3], "phone": row[12]}).Decode(&student)
		if err == mongo.ErrNoDocuments {
			studentCreatedCount++
			student = models.Student{
				ID:               primitive.NewObjectID(),
				Name:             row[1],
				Email:            row[2],
				School:           strings.ToUpper(row[3]),
				SchoolID:         school.ID,
				Interest:         row[10],
				Gender:           strings.ToUpper(row[11]),
				Phone:            row[12],
				FinancialAbility: row[13],
				Progress:         row[14],
				Image:            row[15],
				Category:         row[16],
				Birthdate:        row[17],
				IsActive:         true,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}
			_, err = studentCollection.InsertOne(ctx, student)
			if err != nil {
				studentFailedCount++
				studentFailedRows = append(studentFailedRows, i+1)
			}
		} else if err != nil {
			studentFailedCount++
			studentFailedRows = append(studentFailedRows, i+1)
		} else {
			studentUpdatedCount++
			updatedStudent := bson.M{
				"name":             row[1],
				"email":            row[2],
				"school":           strings.ToUpper(row[3]),
				"schoolID":         school.ID,
				"interest":         row[10],
				"gender":           strings.ToUpper(row[11]),
				"phone":            row[12],
				"financialAbility": row[13],
				"progress":         row[14],
				"image":            row[15],
				"category":         row[16],
				"birthdate":        row[17],
				"updatedAt":        time.Now(),
			}
			_, err = studentCollection.UpdateOne(ctx, bson.M{"_id": student.ID}, bson.M{"$set": updatedStudent})
			if err != nil {
				studentFailedCount++
				studentFailedRows = append(studentFailedRows, i+1) // Append the row number (Excel row numbers start from 1)
				log.Printf("Failed to update student: %v", err)
				continue
			}
		}
	}
	// After processing all rows...
	result := &models.ImportResultStudent{
		SchoolStats: models.OperationStats{
			CreatedCount: schoolCreatedCount,
			UpdatedCount: schoolUpdatedCount,
			FailedCount:  schoolFailedCount,
			FailedRows:   schoolFailedRows,
		},
		StudentStats: models.OperationStats{
			CreatedCount: studentCreatedCount,
			UpdatedCount: studentUpdatedCount,
			FailedCount:  studentFailedCount,
			FailedRows:   studentFailedRows,
		},
	}

	return result, nil
}
