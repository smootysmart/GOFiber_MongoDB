// main.go
package main

import (
	"Test-StructureAPI/services"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"Test-StructureAPI/config"
	"Test-StructureAPI/middleware"
	"Test-StructureAPI/routes"
)

func main() {
	// Connect to MongoDB
	if err := config.ConnectMongoDB(); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer config.DisconnectMongoDB()

	services.InitBookService(config.BookCollection)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.CustomErrorHandler,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Setup routes
	routes.SetupRoutes(app)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
			"time":     time.Now(),
		})
	})

	// Start server
	log.Println("🚀 Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}
