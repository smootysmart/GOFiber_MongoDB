package routes

import (
	"Test-StructureAPI/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Book routes - SPECIFIC ROUTES FIRST
	books := api.Group("/books")
	books.Get("/", controllers.BookController.GetAll)
	books.Get("/search", controllers.BookController.Search) // Specific route BEFORE :id
	books.Get("/:id", controllers.BookController.GetByID)   // Dynamic route AFTER specific
	books.Post("/", controllers.BookController.Create)
	books.Put("/:id", controllers.BookController.Update)
	books.Put("/:id/status", controllers.BookController.UpdateStatus)
	books.Delete("/:id", controllers.BookController.Delete)
}
