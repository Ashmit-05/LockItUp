package models

import (
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Passwords struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name"`
	Username string             `json:"username,omitempty"`
	URL      *url.URL           `json:"url,omitempty"`
	Notes    string             `json:"notes,omitempty"`
	Password string             `json:"password"`
}
