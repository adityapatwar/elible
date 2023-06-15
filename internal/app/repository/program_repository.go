package repository

import (
	"context"
	"elible/internal/app/models"
	"elible/internal/config"
	"math"

	"strings"
	"time"

	utils "elible/internal/app/utils"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// UpdateStudyProgram updates a study program
func (r *StudyProgramRepository) UpdateStudyProgram(id primitive.ObjectID, sp models.StudyProgram) error {
	ctx := context.Background()

	sp.UpdatedAt = time.Now()

	StudyProgramCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_study_programs")
	_, err := StudyProgramCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": sp})
	if err != nil {
		return err
	}

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
func (r *StudyProgramRepository) GetStudyProgram(id primitive.ObjectID) (models.StudyProgramWithUniversity, error) {
	ctx := context.Background()

	StudyProgramCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_study_programs")

	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$lookup": bson.M{
			"from":         "tb_universities", // adjust this to the actual university collection name
			"localField":   "program_details.university",
			"foreignField": "_id",
			"as":           "university",
		}},
		{"$unwind": "$university"},
		{"$project": bson.M{
			"study_program": "$$ROOT",
			"university":    1,
		}},
	}

	var sp models.StudyProgramWithUniversity
	cursor, err := StudyProgramCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return models.StudyProgramWithUniversity{}, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		cursor.Decode(&sp)
	}

	return sp, nil
}

// GetStudyPrograms retrieves all study programs from a specific KnowledgeProgram in a KnowledgeBase
func (r *StudyProgramRepository) GetStudyPrograms(dataFilter *models.GetStudyProgramsFilter) (*models.PagedStudyPrograms, error) {
	ctx := context.Background()

	KnowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")
	UniversityCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_universities")

	var kb models.KnowledgeBase
	err := KnowledgeBaseCollection.FindOne(ctx, bson.M{"year": dataFilter.KbYear}).Decode(&kb)
	if err != nil {
		return nil, err
	}

	// Find the specified KnowledgeProgram
	var kp models.KnowledgeProgram
	for _, program := range kb.Programs {
		if program.Name == dataFilter.KpName {
			kp = program
			break
		}
	}

	StudyProgramCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_study_programs")

	filter := bson.M{
		"_id": bson.M{
			"$in": kp.StudyPrograms,
		},
	}

	if dataFilter.SearchQuery != "" {
		filter["name"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: dataFilter.SearchQuery,
				Options: "i",
			},
		}
	}

	if dataFilter.ProgramType != "" {
		filter["program_details.program_type"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: dataFilter.ProgramType,
				Options: "i",
			},
		}
	}

	if dataFilter.Program != "" {
		filter["program_details.program"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: dataFilter.Program,
				Options: "i",
			},
		}
	}

	// Find all study programs of the KnowledgeProgram
	var programs []models.StudyProgramWithUniversity
	total, err := StudyProgramCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(dataFilter.PageSize)))

	findOptions := options.Find().SetSkip(int64((dataFilter.Page - 1) * dataFilter.PageSize)).SetLimit(int64(dataFilter.PageSize))
	cursor, err := StudyProgramCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var sp models.StudyProgram
		err = cursor.Decode(&sp)
		if err != nil {
			return nil, err
		}

		// Get the university detail
		var university models.University
		err = UniversityCollection.FindOne(ctx, bson.M{"_id": sp.ProgramDetails.University}).Decode(&university)
		if err != nil {
			return nil, err
		}

		// Add the university to the struct
		programWithUniversity := models.StudyProgramWithUniversity{StudyProgram: sp, University: university}
		programs = append(programs, programWithUniversity)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &models.PagedStudyPrograms{
		CurrentPage:  dataFilter.Page,
		TotalRecords: total,
		TotalPages:   totalPages,
		Records:      programs,
	}, nil
}

func (r *StudyProgramRepository) ImportDataFromExcel(knowledgeBaseYear, knowledgeProgramName, filePath string) (*models.ImportResult, error) {
	// Define counters
	var universityCreatedCount, universityUpdatedCount, universityFailedCount int
	var programCreatedCount, programUpdatedCount, programFailedCount int
	var universityFailedRows, programFailedRows []int
	// Load Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// Choose the collections to work with
	universityCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_universities")
	studyProgramCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_study_programs")
	knowledgeBaseCollection := r.MongoClient.Database(r.cfg.MongoDBName).Collection("tb_knowledge_bases")

	// Iterate over the rows, starting from the second row
	for i, row := range rows {
		// Skip the header row
		if i == 0 {
			continue
		}

		// Find the university
		universityFilter := bson.M{"name": row[2]}
		var existingUniversity models.University
		err = universityCollection.FindOne(ctx, universityFilter).Decode(&existingUniversity)

		universityFailed := false

		// If university doesn't exist, create it
		if err == mongo.ErrNoDocuments {
			universityCreatedCount++
			universityUpdate := bson.M{
				"$setOnInsert": bson.M{
					"_id":     primitive.NewObjectID(),
					"name":    row[2],
					"alias":   row[3],
					"address": row[4],
					"website": row[5],
					"logo":    row[6],
					"image":   row[7],
					"contact": models.Contact{
						Email: row[8],
						Phone: row[9],
						Fax:   row[10],
					},
					"socialMedia": []models.SocialMedia{
						{
							Platform: row[11],
							Link:     row[12],
						},
					},
					"createdAt": time.Now(),
					"updatedAt": time.Now(),
				},
			}

			universityOpts := options.Update().SetUpsert(true)
			_, err := universityCollection.UpdateOne(ctx, universityFilter, universityUpdate, universityOpts)
			if err != nil {
				universityFailedCount++
				universityFailedRows = append(universityFailedRows, i+1) // Append the row number (Excel row numbers start from 1)
				universityFailed = true
			}
		} else if err != nil {
			universityFailedCount++
			return nil, err
		} else {
			universityUpdatedCount++
			// If university exists, update it
			universityUpdate := bson.M{
				"$set": bson.M{
					"alias":   row[3],
					"address": row[4],
					"website": row[5],
					"logo":    row[6],
					"image":   row[7],
					"contact": models.Contact{
						Email: row[8],
						Phone: row[9],
						Fax:   row[10],
					},
					"socialMedia": []models.SocialMedia{
						{
							Platform: row[11],
							Link:     row[12],
						},
					},
					"updatedAt": time.Now(),
				},
			}

			_, err := universityCollection.UpdateOne(ctx, universityFilter, universityUpdate)
			if err != nil {
				universityFailedCount++
				universityFailedRows = append(universityFailedRows, i+1) // Append the row number (Excel row numbers start from 1)
				universityFailed = true
			}
		}
		// Do not continue if the university update/create operation failed
		if universityFailed {
			continue
		}
		universityID := existingUniversity.ID

		// Find the program
		programFilter := bson.M{
			"name":                       row[13],
			"program_details.university": universityID,
			"program_details.program":    row[14],
		}
		var existingProgram models.StudyProgram
		err = studyProgramCollection.FindOne(ctx, programFilter).Decode(&existingProgram)

		// If program doesn't exist, create it
		if err == mongo.ErrNoDocuments {
			programCreatedCount++
			programUpdate := bson.M{
				"$setOnInsert": bson.M{
					"_id":  primitive.NewObjectID(),
					"name": row[13],
					"program_details": models.Program{
						University:    universityID,
						Program:       row[14],
						ProgramType:   row[15],
						UKT:           row[16],
						SPI:           row[17],
						Capacity:      row[18],
						IsPacketC:     row[19] == "yes",
						Description:   row[20],
						Advantages:    row[21],
						Disadvantages: row[22],
						Requirements:  strings.Split(row[23], ","),
						Registration: models.RegistrationDates{
							Start: utils.ParseDate(row[24]),
							End:   utils.ParseDate(row[25]),
						},
						Exam: models.ExamDates{
							Start: utils.ParseDate(row[26]),
							End:   utils.ParseDate(row[27]),
						},
						Announcement: utils.ParseDate(row[28]),
					},
					"createdAt": time.Now(),
					"updatedAt": time.Now(),
				},
			}

			programOpts := options.Update().SetUpsert(true)
			_, err := studyProgramCollection.UpdateOne(ctx, programFilter, programUpdate, programOpts)
			if err != nil {
				programFailedCount++
				programFailedRows = append(programFailedRows, i+1)
			}
		} else if err != nil {
			programFailedCount++
			programFailedRows = append(programFailedRows, i+1)
		} else {
			programUpdatedCount++
			// If program exists, update it
			programUpdate := bson.M{
				"$set": bson.M{
					"program_details": models.Program{
						University:    universityID,
						Program:       row[14],
						ProgramType:   row[15],
						UKT:           row[16],
						SPI:           row[17],
						Capacity:      row[18],
						IsPacketC:     row[19] == "yes",
						Description:   row[20],
						Advantages:    row[21],
						Disadvantages: row[22],
						Requirements:  strings.Split(row[23], ","),
						Registration: models.RegistrationDates{
							Start: utils.ParseDate(row[24]),
							End:   utils.ParseDate(row[25]),
						},
						Exam: models.ExamDates{
							Start: utils.ParseDate(row[26]),
							End:   utils.ParseDate(row[27]),
						},
						Announcement: utils.ParseDate(row[28]),
					},
					"updatedAt": time.Now(),
				},
			}

			_, err := studyProgramCollection.UpdateOne(ctx, programFilter, programUpdate)
			if err != nil {
				programFailedCount++
				programFailedRows = append(programFailedRows, i+1)
			}
		}
		programID := existingProgram.ID

		// Check if the program already exists in the specified KnowledgeProgram
		var existingKnowledgeBase models.KnowledgeBase
		err = knowledgeBaseCollection.FindOne(ctx, bson.M{
			"year":                    knowledgeBaseYear,
			"programs.name":           knowledgeProgramName,
			"programs.study_programs": programID}).Decode(&existingKnowledgeBase)

		if err != nil {
			if err != mongo.ErrNoDocuments {
				return nil, err
			}

			// Program does not exist in the specified KnowledgeProgram, add it
			_, err = knowledgeBaseCollection.UpdateOne(
				ctx,
				bson.M{"year": knowledgeBaseYear, "programs.name": knowledgeProgramName},
				bson.M{"$push": bson.M{"programs.$.study_programs": programID}},
			)
			if err != nil {
				return nil, err
			}
		}
	}

	// After your existing code...
	result := &models.ImportResult{
		UniversityStats: models.OperationStats{
			CreatedCount: universityCreatedCount,
			UpdatedCount: universityUpdatedCount,
			FailedCount:  universityFailedCount,
			FailedRows:   universityFailedRows,
		},
		ProgramStats: models.OperationStats{
			CreatedCount: programCreatedCount,
			UpdatedCount: programUpdatedCount,
			FailedCount:  programFailedCount,
			FailedRows:   programFailedRows,
		},
	}

	return result, nil
}
