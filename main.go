package main

import (
	"log"
	"project/database"
	"project/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Connect to MongoDB
	database.ConnectMongo()
	defer database.Client.Disconnect(nil)

	// Initialize Fiber
	app := fiber.New()

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Routes
	app.Get("/api/todos", handlers.GetTodos)
	app.Post("/api/todos", handlers.CreateTodo)
	app.Put("/api/todos/:id", handlers.UpdateTodo)

	// Start server
	port := "5000"
	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
