package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(uri, dbName string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	// Use Jakarta's time zone
	// location, _ := time.LoadLocation("Asia/Jakarta")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	client.Database(dbName)

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	// fmt.Println("Connected to MongoDB successfully at", time.Now().In(location))
	return client, nil
}
