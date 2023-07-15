package userController

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

  userModel "github.com/Ashmit-05/LockItUp/models/users.go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
  "golang.org/x/crypto/bcrypt"
)

func signUp(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","POST")

  var user userModel.User

	_ = json.NewDecoder(r.Body).Decode(&user)

  // necessary checks
  if user.Name == nil {
    log.Fatal("Please provide name")
  }
  if user.Email == nil {
    log.Fatal("Please provide email")
  }
  if user.MasterPassword == nil {
    log.Fatal("Please provide master password")
  }
  if r.Body.confirmPassword != user.MasterPassword {
    log.Fatal("Passwords do not match")
  }
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.MasterPassword), 10)
	if err != nil {
    log.Fatal("Encountered an error while hashing the password")
	}
  user.MasterPassword = string(hashedPassword);
  fmt.Println("this is the hashed password : ",user.MasterPassword);

  result, err := collection.InsertOne(context.Background(),user);
  if err != nil {
    log.Fatal("Encountered an error while storing details in database");
  }
  json.NewEncoder(w).Encode("success","added a new user");
}
