package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string                 `json:"name"`
	Email          string                 `json:"email"`
	PhoneNumber    string                 `json:"phonenumber,omitempty"`
	MasterPassword string                 `json:"masterpassword"`
	Passwords      []*Passwords `json:"passwords,omitempty"`
}
