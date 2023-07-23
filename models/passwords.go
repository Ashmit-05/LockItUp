package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Passwords struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name"`
	Username string             `json:"username,omitempty"`
	URL      string             `json:"url,omitempty"`
	Notes    string             `json:"notes,omitempty"`
	Password string             `json:"password"`
	UserId   primitive.ObjectID `json:"userId"`
}
