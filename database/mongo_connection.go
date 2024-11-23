package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var TodoCollection *mongo.Collection

func ConnectMongo() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Check MongoDB connection
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	log.Println("Connected to MongoDB")
	Client = client
	TodoCollection = Client.Database("db_golang_todo").Collection("todo")
}