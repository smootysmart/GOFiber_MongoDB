// controllers/book_controller.go
package controllers

import (
	"fmt"

	"Test-StructureAPI/models"
	"Test-StructureAPI/services"

	"github.com/gofiber/fiber/v2"
)

type bookController struct{}

var BookController = &bookController{}

// GetAll returns all books
func (bc *bookController) GetAll(c *fiber.Ctx) error {
	books, err := services.BookService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch books",
		})
	}
	return c.JSON(books)
}

// Search searches books by title or author
func (bc *bookController) Search(c *fiber.Ctx) error {
	query := c.Query("query")

	results, err := services.BookService.Search(query)
	if err != nil {
		if err.Error() == "query is required" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Query is required",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Search failed",
		})
	}

	return c.JSON(results)
}

// GetByID returns a single book by ID
func (bc *bookController) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")

	book, err := services.BookService.GetByID(idParam)
	if err != nil {
		if err.Error() == "invalid ID format" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch book",
		})
	}

	return c.JSON(book)
}

// Create creates a new book
func (bc *bookController) Create(c *fiber.Ctx) error {
	var book models.Book

	// Parse request body
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid input",
			"details": err.Error(),
		})
	}

	// Create book via service
	createdBook, err := services.BookService.Create(&book)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to insert book",
			"details": err.Error(),
		})
	}

	fmt.Printf("📘 New Book: %+v\n", createdBook)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"inserted_id": createdBook.ID,
		"book":        createdBook,
	})
}

// Update updates an existing book
func (bc *bookController) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")

	var updateData models.Book
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	err := services.BookService.Update(idParam, &updateData)
	if err != nil {
		if err.Error() == "invalid ID format" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update book",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Book updated",
	})
}

// Delete removes a book
func (bc *bookController) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")

	err := services.BookService.Delete(idParam)
	if err != nil {
		if err.Error() == "invalid ID format" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete book",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateStatus toggles book status between available and borrowed
func (bc *bookController) UpdateStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Parse request body (optional, for future use)
	var statusReq models.StatusRequest
	if err := c.BodyParser(&statusReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	newStatus, err := services.BookService.UpdateStatus(idParam)
	if err != nil {
		if err.Error() == "invalid ID format" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		if err.Error() == "invalid current status" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid action or status",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Status updated",
		"status":  newStatus,
	})
}

// GetAllWithPagination returns books with pagination
//func (bc *bookController) GetAllWithPagination(c *fiber.Ctx) error {
//	page := c.QueryInt("page", 1)
//	limit := c.QueryInt("limit", 10)
//
//	if page < 1 {
//		page = 1
//	}
//	if limit < 1 || limit > 100 {
//		limit = 10
//	}
//
//	books, total, err := services.BookService.GetAllWithPagination(page, limit)
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Failed to fetch books",
//		})
//	}
//
//	return c.JSON(fiber.Map{
//		"data": books,
//		"meta": fiber.Map{
//			"total": total,
//			"page":  page,
//			"limit": limit,
//		},
//	})
//}
