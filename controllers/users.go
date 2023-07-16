package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	generateJWT "github.com/Ashmit-05/LockItUp/middlewares"
	userModel "github.com/Ashmit-05/LockItUp/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var collection *mongo.Collection

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in reading .env file")
	}

	connectionString := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	colName := os.Getenv("USER_COL")

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

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var user userModel.User

	_ = json.NewDecoder(r.Body).Decode(&user)

  fmt.Println(user)
  fmt.Println(r.Body)

	// necessary checks
	if user.Name == "" {
		log.Fatal("Please provide name")
	}
	if user.Email == "" {
		log.Fatal("Please provide email")
	}
	if user.MasterPassword == "" {
		log.Fatal("Please provide master password")
	}
	confirmPassword := r.FormValue("confirmPassword")
	if confirmPassword == "" {
		http.Error(w, "Please provide confirm password", http.StatusBadRequest)
		return
	}

	if confirmPassword != user.MasterPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.MasterPassword), 10)
	if err != nil {
		log.Fatal("Encountered an error while hashing the password")
	}
	user.MasterPassword = string(hashedPassword)
	fmt.Println("this is the hashed password : ", user.MasterPassword)

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal("Encountered an error while storing details in database")
	}
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()
	jwtToken, err := generateJWT.GenerateToken(insertedID)
	if err != nil {
		log.Fatal("Couldn't generate jwt token")
	}
	json.NewEncoder(w).Encode(jwtToken)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var user userModel.User

	_ = json.NewDecoder(r.Body).Decode(&user)

}
