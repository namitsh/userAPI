package database

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var ctx = context.TODO()

func dbInstance() *mongo.Client {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading the configurations")
	}
	mongoUri := os.Getenv("MONGODB_URL")
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}
	return client
}

func Connect() *mongo.Collection {
	client := dbInstance()
	collection := client.Database("UserAuth").Collection("users")
	return collection
}
