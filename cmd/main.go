// cmd/main.go
package main

import (
	elible "elible/internal/app"
	"elible/internal/app/handlers"
	"elible/internal/config"
	"elible/internal/mongodb"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	mongoClient, err := mongodb.ConnectMongoDB(cfg.MongoDBURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	deps, err := elible.InitializeDependencies(cfg, mongoClient)
	if err != nil {
		log.Fatalf("Error initializing dependencies: %v", err)
	}

	router := gin.Default()
	router.Use(corsMiddleware())

	handlers.Routes(router, cfg, deps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(router.Run(":" + port))
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}
