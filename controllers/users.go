package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Ashmit-05/LockItUp/middlewares"
	userModel "github.com/Ashmit-05/LockItUp/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

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

	userCollection = client.Database(dbName).Collection(colName)
	fmt.Println("userCollection instance is ready")
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var user userModel.User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// necessary checks
	if user.Name == "" {
		http.Error(w, "Please provide name", http.StatusBadRequest)
		return
	}
	if user.Email == "" {
		http.Error(w, "Please provide email", http.StatusBadRequest)
		return
	}
	if user.MasterPassword == "" {
		http.Error(w, "Please provide master password", http.StatusBadRequest)
		return
	}

	existingEmail, _, err1 := middlewares.CheckUserExistsByEmail(user.Email, userCollection)
	existingPhone, _, err2 := middlewares.CheckUserExistsByPhoneNumber(user.PhoneNumber, userCollection)

	if err1 != nil || err2 != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		return
	}

	if existingEmail || existingPhone {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	var requestData struct {
		ConfirmPassword string `json:"confirmpassword"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Failed to parse confirm password from JSON data", http.StatusBadRequest)
		return
	}

	confirmPassword := requestData.ConfirmPassword
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
		http.Error(w, "unable to hash password", http.StatusInternalServerError)
	}
	user.MasterPassword = string(hashedPassword)
	fmt.Println("this is the hashed password : ", user.MasterPassword)

	result, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Encountered an error while trying to store details in database", http.StatusInternalServerError)
	}
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()
	jwtToken, err := middlewares.GenerateToken(insertedID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to generate jwt token", http.StatusInternalServerError)
	}
	// Create a JSON response
	response := map[string]string{
		"token": jwtToken,
	}
	json.NewEncoder(w).Encode(response)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var user userModel.User

	_ = json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" && user.PhoneNumber == "" {
		http.Error(w, "Please provide email or phone number", http.StatusBadRequest)
		return
	}
	if user.MasterPassword == "" {
		http.Error(w, "Please provide master password", http.StatusBadRequest)
		return
	}
	if user.Email != "" {
		exists, existingUser, err := middlewares.CheckUserExistsByEmail(user.Email, userCollection)
		if err != nil {
			http.Error(w, "Unable to fetch user details, please try later", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "User not found in the database", http.StatusConflict)
			return
		}
		correctPassword := bcrypt.CompareHashAndPassword([]byte(existingUser.MasterPassword), []byte(user.MasterPassword))
		if correctPassword != nil {
			http.Error(w, "Incorrect Password", http.StatusBadRequest)
		}
		jwtToken, err := middlewares.GenerateToken(existingUser.ID.Hex())
		if err != nil {
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
		}
		// Create a JSON response
		response := map[string]string{
			"token": jwtToken,
		}
		json.NewEncoder(w).Encode(response)
		return
	} else if user.PhoneNumber != "" {
		exists, existingUser, err := middlewares.CheckUserExistsByPhoneNumber(user.PhoneNumber, userCollection)
		if err != nil {
			http.Error(w, "Unable to fetch user details, please try later", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "User not found in the database", http.StatusConflict)
			return
		}
		correctPassword := bcrypt.CompareHashAndPassword([]byte(existingUser.MasterPassword), []byte(user.MasterPassword))
		if correctPassword != nil {
			http.Error(w, "Incorrect Password", http.StatusBadRequest)
		}
		jwtToken, err := middlewares.GenerateToken(existingUser.ID.Hex())
		if err != nil {
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
		}
		// Create a JSON response
		response := map[string]string{
			"token": jwtToken,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

}
