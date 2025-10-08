package services

import (
	"Test-StructureAPI/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type bookService struct {
	collection *mongo.Collection
}

var BookService *bookService

// InitBookService initializes the book service with MongoDB collection
func InitBookService(collection *mongo.Collection) {
	BookService = &bookService{
		collection: collection,
	}
}

// GetAll retrieves all books from MongoDB
func (s *bookService) GetAll() ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []models.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	// Return empty array instead of nil
	if books == nil {
		books = []models.Book{}
	}

	return books, nil
}

// Search searches books by title or author (case-insensitive)
func (s *bookService) Search(query string) ([]models.Book, error) {
	if query == "" {
		return nil, errors.New("query is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create regex filter for case-insensitive search
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"author": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.Book
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Return empty array instead of nil
	if results == nil {
		results = []models.Book{}
	}

	return results, nil
}

// GetByID retrieves a single book by ID
func (s *bookService) GetByID(id string) (*models.Book, error) {
	// Convert string ID to ObjectID
	bookID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var book models.Book
	err = s.collection.FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	return &book, nil
}

// Create inserts a new book into MongoDB
func (s *bookService) Create(book *models.Book) (*models.Book, error) {
	// Generate ObjectID
	book.ID = primitive.NewObjectID()

	// Default status if not set
	if book.Status == "" {
		book.Status = "available"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.collection.InsertOne(ctx, book)
	if err != nil {
		return nil, err
	}

	return book, nil
}

// Update updates an existing book
func (s *bookService) Update(id string, updateData *models.Book) error {
	// Convert string ID to ObjectID
	bookID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Build update document
	update := bson.M{
		"$set": bson.M{
			"title":  updateData.Title,
			"author": updateData.Author,
			"year":   updateData.Year,
		},
	}

	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": bookID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("book not found")
	}

	return nil
}

// Delete removes a book from MongoDB
func (s *bookService) Delete(id string) error {
	// Convert string ID to ObjectID
	bookID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": bookID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("book not found")
	}

	return nil
}

// UpdateStatus toggles book status between "available" and "borrowed"
func (s *bookService) UpdateStatus(id string) (string, error) {
	// Convert string ID to ObjectID
	bookID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", errors.New("invalid ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the book first
	var book models.Book
	err = s.collection.FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("book not found")
		}
		return "", err
	}

	// Toggle status
	var newStatus string
	if book.Status == "available" {
		newStatus = "borrowed"
	} else if book.Status == "borrowed" {
		newStatus = "available"
	} else {
		return "", errors.New("invalid current status")
	}

	// Update the status
	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": bookID},
		bson.M{"$set": bson.M{"status": newStatus}},
	)
	if err != nil {
		return "", err
	}

	return newStatus, nil
}

//// GetAllWithPagination retrieves books with pagination
//func (s *bookService) GetAllWithPagination(page, limit int) ([]models.Book, int64, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	// Count total documents
//	total, err := s.collection.CountDocuments(ctx, bson.M{})
//	if err != nil {
//		return nil, 0, err
//	}
//
//	// Calculate skip
//	skip := int64((page - 1) * limit)
//
//	// Find with pagination
//	cursor, err := s.collection.Find(ctx, bson.M{}, &mongo.FindOptions{
//		Skip:  &skip,
//		Limit: int64Ptr(int64(limit)),
//	})
//	if err != nil {
//		return nil, 0, err
//	}
//	defer cursor.Close(ctx)
//
//	var books []models.Book
//	if err = cursor.All(ctx, &books); err != nil {
//		return nil, 0, err
//	}
//
//	if books == nil {
//		books = []models.Book{}
//	}
//
//	return books, total, nil
//}

// Helper function to create int64 pointer
//func int64Ptr(i int64) *int64 {
//	return &i
//}
