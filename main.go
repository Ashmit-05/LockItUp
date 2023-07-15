package main

import (
	"context"
	"fmt"
	"log"
  "os"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
  "github.com/joho/godotenv"
)

var collection *mongo.Collection

func init()  {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error in reading .env file")
  }

  connectionString := os.Getenv("MONGODB_URI")
  dbName := os.Getenv("DB_NAME")
  colName := os.Getenv("COL_NAME")


  // connect to database
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connected successfully")

	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection instance is ready")
}

func main()  {
  fmt.Println("LockItUp")
}
