package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/rs/cors"
)

// Todo represents a todo item
type Todo struct {
	ID        primitive.ObjectID    `json:"_id, omitempty" bson:"_id,omitempty"` // Pastikan ID sesuai dengan tipe data di MongoDB
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

var collection *mongo.Collection

func main() {
	// Log startup message
	fmt.Println("Starting application...")

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error disconnecting MongoDB: %v", err)
		}
	}()

	// Check MongoDB connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
	fmt.Println("Connected to MongoDB Atlas")

	// Set collection
	collection = client.Database("db_golang_todo").Collection("todo")

	// Initialize Fiber
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin,Content-Type, Accept",
	}))

	// Routes
	app.Get("/api/todos", getTodos)
	app.Post("/api/todos",createTodo)

	// Get port from environment or default to 5000
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Start server
	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port)) // Tambahkan ":" di depan port
}

// getTodos fetches all todos from the database
func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	// Find all todos in the collection
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error fetching todos: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch todos",
		})
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor and decode each document
	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			log.Printf("Error decoding todo: %v", err)
			continue
		}
		todos = append(todos, todo)
	}

	// Check for errors during cursor iteration
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error iterating over todos",
		})
	}

	// Return todos as JSON response
	return c.JSON(todos)
}


func createTodo(c *fiber.Ctx)error{
	todo := new(Todo)
	// {id:0, completed:false,body:""}

	if err := c.BodyParser(todo); err != nil{
		return err;
	}

	if todo.Body == ""{
		return c.Status(400).JSON(fiber.Map{"error" : "Todo body cannot be empty"})
	}

	insertResult, err := collection.InsertOne(context.Background(),todo)

	if err != nil{
		return  err
	}
	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
	
}


func updateTodo(c *fiber.Ctx)error{
	id := c.Params("id")
	objectID,err := primitive.ObjectIDFromHex(id)

	if err != nil{
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id" : objectID}
	update := bson.M{"$set": bson.M{"completed":true}}

	_,err = collection.UpdateOne(context.Background(),filter,update)

	if err != nil{
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})


}