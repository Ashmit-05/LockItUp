package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Ashmit-05/LockItUp/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-password/password"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var passwordsCollection *mongo.Collection

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in reading .env file")
	}

	connectionString := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	colName := os.Getenv("PASSWD_COL")

	// connect to database
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connected successfully")

	passwordsCollection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection instance is ready")
}

func AddPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var password models.Passwords

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &password)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// necessary checks
	if password.Name == "" {
		http.Error(w, "Please provide name for unique identification", http.StatusBadRequest)
		return
	}
	if password.Password == "" {
		http.Error(w, "Please provide password", http.StatusBadRequest)
	}

	var requestData struct {
		UserId string `json:"userId"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Failed to parse confirm password from JSON data", http.StatusBadRequest)
		return
	}
	filter := bson.M{"_id": requestData.UserId}

	var user *models.User
	err1 := userCollection.FindOne(context.Background(), filter).Decode(&user)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			http.Error(w, "No user found", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	// Add the password to the user
	user.Passwords = append(user.Passwords, &password)

	// Update the user in the database
	update := bson.M{"$set": user}
	_, err = userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	response := map[string]*models.User{
		"user": user,
	}

	// Return success response
	json.NewEncoder(w).Encode(response)
}

// get all passwords for one user
func GetAllPasswords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	var user models.User

	params := mux.Vars(r)

	filter := bson.M{"_id": params["userId"]}

	err1 := userCollection.FindOne(context.Background(), filter).Decode(&user)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			http.Error(w, "No user found", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	response := map[string][]*models.Passwords{
		"Passwords": user.Passwords,
	}

	json.NewEncoder(w).Encode(response)
}

func GeneratePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var reqBody struct {
		Length      int  `json:"length"`
		NumDigits   int  `json:"numdigits"`
		NumSymbols  int  `json:"numsymbols"`
		NoUpper     bool `json:"noupper"`
		AllowRepeat bool `json:"allowrepeat"`
	}
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		http.Error(w, "Failed to parse confirm password from JSON data", http.StatusBadRequest)
		return
	}
	passwd, err := password.Generate(
		reqBody.Length,
		reqBody.NumDigits,
		reqBody.NumSymbols,
		reqBody.NoUpper,
		reqBody.AllowRepeat,
	)
	if err != nil {
		http.Error(w, "Unable to generate password", http.StatusInternalServerError)
	}

	response := map[string]string{
		"password": passwd,
	}

	json.NewEncoder(w).Encode(response)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	params := mux.Vars(r)
	if params["passwordId"] == "" || params["userId"] == "" {
		http.Error(w,"Missing essential fields",http.StatusBadRequest)
	}
	var passwordId = params["passwordId"]
	var user models.User
	filter := bson.M{"_id": params["userId"]}
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No user found", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	var reqData struct {
		Passwd string `json:"passwd"`
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		http.Error(w, "Failed to parse confirm password from JSON data", http.StatusBadRequest)
		return
	}

	// Check if the user contains the password with the given passwdId
	var found bool
	for _, pass := range user.Passwords {
		if pass.ID.Hex() == passwordId {
			pass.Password = reqData.Passwd
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Password not found", http.StatusBadRequest)
		return
	}

	// Update the user in the database
	update := bson.M{"$set": user}
	_, err = userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	response := map[string]string {
		"success" : "Changed the password successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func DeletePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	userID := params["userId"]
	passwordID := params["passwordId"]

	// Find the user by userID
	filter := bson.M{"_id": userID}
	var user models.User
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No user found", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	// Find the password index in the user.Passwords slice
	passwordIndex := -1
	for i, password := range user.Passwords {
		if password.ID.Hex() == passwordID {
			passwordIndex = i
			break
		}
	}

	// Check if the password exists in the user.Passwords slice
	if passwordIndex == -1 {
		http.Error(w, "Password not found", http.StatusBadRequest)
		return
	}

	// Remove the password from the user.Passwords slice
	user.Passwords = append(user.Passwords[:passwordIndex], user.Passwords[passwordIndex+1:]...)

	// Update the user in the database
	update := bson.M{"$set": user}
	_, err = userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{} {
		"user" : user,
		"success" : "Deleted the user",
	}
	// Return success response
	json.NewEncoder(w).Encode(response)
}
