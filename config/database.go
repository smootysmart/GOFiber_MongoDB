package config

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

var mongoDB *MongoDB

var BookCollection *mongo.Collection

func ConnectMongoDB() error {
	// Load .env file
	godotenv.Load()

	// Get MongoDB URI from environment variable
	//mongoURI := os.Getenv("MONGODB_URI")

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://root:password@localhost:27017")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// Get database and collection
	database := client.Database("library")
	collection := database.Collection("books")

	mongoDB = &MongoDB{
		Client:     client,
		Database:   database,
		Collection: collection,
	}

	BookCollection = collection

	log.Println("✅ Connected to MongoDB!")
	return nil
}

func DisconnectMongoDB() {
	if mongoDB.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoDB.Client.Disconnect(ctx); err != nil {
			log.Println("Error disconnecting MongoDB:", err)
		}
	}
}
