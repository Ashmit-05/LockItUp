package middlewares

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/Ashmit-05/LockItUp/models"
)

// Check if a user with the given email already exists
func CheckUserExistsByEmail(email string, collection *mongo.Collection) (bool,*models.User, error) {
	filter := bson.M{"email": email}
	var user models.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false,nil, nil // User does not exist
		}
		log.Printf("Error checking user by email: %s", err)
		return true, nil, err // Error occurred, assume user exists
	}
	return true, &user, nil // User exists
}

// Check if a user with the given phone number already exists
func CheckUserExistsByPhoneNumber(phoneNumber string, collection *mongo.Collection) (bool,*models.User, error) {
	filter := bson.M{"phonenumber": phoneNumber}
	var user models.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil, nil // User does not exist
		}
		log.Printf("Error checking user by phone number: %s", err)
		return true, nil, err // Error occurred, assume user exists
	}
	return true, &user, nil // User exists
}
