package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/your_project_name/models"
	"github.com/your_project_name/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetTodos fetches all todos
func GetTodos(c *fiber.Ctx) error {
	var todos []models.Todo

	cursor, err := database.TodoCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error fetching todos: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch todos"})
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo models.Todo
		if err := cursor.Decode(&todo); err != nil {
			log.Printf("Error decoding todo: %v", err)
			continue
		}
		todos = append(todos, todo)
	}

	if err := cursor.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating over todos"})
	}

	return c.JSON(todos)
}

// CreateTodo creates a new todo
func CreateTodo(c *fiber.Ctx) error {
	todo := new(models.Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body cannot be empty"})
	}

	insertResult, err := database.TodoCollection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(todo)
}

// UpdateTodo updates a todo's completion status
func UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = database.TodoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}