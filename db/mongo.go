package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client

func Connect(uri string) {
	// v2: Connect only takes options, no context
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	// Ping still uses context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	Client = client
	log.Println("Connected to MongoDB")
}

func GetCollection(database, collection string) *mongo.Collection {
	return Client.Database(database).Collection(collection)
}
