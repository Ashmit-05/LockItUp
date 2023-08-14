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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	passCol := os.Getenv("PASSWD_COL")
	userCol := os.Getenv("USER_COL")
	// connect to database
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connected successfully")

	passwordsCollection = client.Database(dbName).Collection(passCol)
	userCollection = client.Database(dbName).Collection(userCol)
	fmt.Println("Instance is ready")
}

func AddPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var user *models.User
	params := mux.Vars(r)
	userID := params["userId"]
	objId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Unable to fetch userId", http.StatusInternalServerError)
	}
	filter := bson.M{"_id": objId}
	err1 := userCollection.FindOne(context.Background(), filter).Decode(&user)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			http.Error(w, "No user found", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var password models.Passwords

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

	userObjId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		return
	}

	password.UserId = userObjId

	insertedId, err := passwordsCollection.InsertOne(context.Background(), password)
	if err != nil {
		http.Error(w, "Failed to insert password", http.StatusInternalServerError)
		return
	}

	passId, ok := insertedId.InsertedID.(primitive.ObjectID)
	if !ok {
		http.Error(w, "Failed to get inserted ObjectID", http.StatusInternalServerError)
		return
	}

	user.Passwords = append(user.Passwords, passId)
	// Update the user in the database
	update := bson.M{"$set": user}
	_, err = userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user":       user,
		"insertedID": insertedId.InsertedID,
	}

	// Return success response
	json.NewEncoder(w).Encode(response)
}

// get all passwords for one user
func GetAllPasswords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	userId := params["userId"]
	if userId == "" {
		http.Error(w, "Missing user id", http.StatusBadRequest)
		return
	}
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Unexpected error. Try again later", http.StatusInternalServerError)
		return
	}

	filter := bson.M{"userid": userObjId}
	cursor, err := passwordsCollection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Something went wrong! Try again later", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var results []models.Passwords
	for cursor.Next(context.Background()) {
		var password models.Passwords
		if err := cursor.Decode(&password); err != nil {
			http.Error(w, "Failed to decode password data", http.StatusInternalServerError)
			return
		}
		results = append(results, password)
	}

	response := map[string]interface{}{
		"success": results,
	}

	json.NewEncoder(w).Encode(response)
}

func GeneratePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

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

// func UpdatePassword(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Allow-Control-Allow-Methods", "POST")

// 	params := mux.Vars(r)
// 	if params["passwordId"] == "" || params["userId"] == "" {
// 		http.Error(w, "Missing essential fields", http.StatusBadRequest)
// 	}
// 	var passwordId = params["passwordId"]
// 	passIdString, err := primitive.ObjectIDFromHex(passwordId)
// 	if err != nil {
// 		http.Error(w, "Unexpected error", http.StatusInternalServerError)
// 		return
// 	}
// 	var user models.User
// 	filter := bson.M{"_id": params["userId"]}
// 	err := userCollection.FindOne(context.Background(), filter).Decode(&user)

// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			http.Error(w, "No user found", http.StatusBadRequest)
// 			return
// 		}
// 		http.Error(w, "Failed to find user", http.StatusInternalServerError)
// 		return
// 	}

// 	var reqData struct {
// 		Passwd string `json:"passwd"`
// 	}
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Failed to read request body", http.StatusBadRequest)
// 		return
// 	}
// 	err = json.Unmarshal(body, &reqData)
// 	if err != nil {
// 		http.Error(w, "Failed to parse confirm password from JSON data", http.StatusBadRequest)
// 		return
// 	}

// 	// Check if the user contains the password with the given passwdId
// 	var found bool
// 	for _, pass := range user.Passwords {

// 		if pass == passIdString {
// 			pass.Password = reqData.Passwd
// 			found = true
// 			break
// 		}
// 	}

// 	if !found {
// 		http.Error(w, "Password not found", http.StatusBadRequest)
// 		return
// 	}

// 	// Update the user in the database
// 	update := bson.M{"$set": user}
// 	_, err = userCollection.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 		http.Error(w, "Failed to update user", http.StatusInternalServerError)
// 		return
// 	}
// 	response := map[string]string{
// 		"success": "Changed the password successfully",
// 	}
// 	json.NewEncoder(w).Encode(response)
// }

func DeletePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	userID := params["userId"]
	passwordID := params["passwordId"]

	passwordObjID, err := primitive.ObjectIDFromHex(passwordID)
	if err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		return
	}
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		return
	}
	filter := bson.M{"_id": userObjID}

	var user models.User
	err1 := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			http.Error(w, "No documents found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Failed to find password document", http.StatusNotFound)
			return
		}
	}

	// Find the index of the password in the user.Passwords slice
	index := -1
	for i, p := range user.Passwords {
		if p == passwordObjID {
			index = i
			break
		}
	}

	// If the password exists in the user.Passwords slice, remove it
	if index != -1 {
		user.Passwords = append(user.Passwords[:index], user.Passwords[index+1:]...)
	}

	// Update the user document in the database
	update := bson.M{"$set": user}
	_, err = userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	passFilter := bson.M{"_id": passwordObjID}
	res := passwordsCollection.FindOneAndDelete(context.Background(), passFilter)

	response := map[string]interface{}{
		"success": "Deleted",
		"deleted": res,
	}

	json.NewEncoder(w).Encode(response)
}
