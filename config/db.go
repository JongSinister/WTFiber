package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var clt *mongo.Client
var DB *mongo.Database

func ConnectDB() {

	// 1) Get the MongoDB URI from the environment variables
	uri := os.Getenv("MONGO_URI")

	// 2) Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("MongoDB connection error: ", err)
	}

	// 3) Ping the MongoDB server
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Could not ping MongoDB: ", err)
	}

	clt = client
	log.Println("Connected to MongoDB")
}

func DisconnectDB() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := clt.Disconnect(ctx); err != nil {
		log.Fatal("Error disconnecting from MongoDB: ", err)
	}
	log.Println("Disconnected from MongoDB")
}

func InitDB() {
	DB = clt.Database("test")
}
