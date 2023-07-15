package passwordModel

import (
  "go.mongodb.org/mongo-driver/bson/primitive"
  "net/url"
)

type Passwords struct {
  ID primitive.ObjectID `json:"_id,omitempty"` `bson:"_id,omitempty"` 
  Name string `json:"name,omitempty"`
  Username string `json:"username,omitempty"`
  URL *url.URL `json:"url,omitempty"`
  Notes string `json:"notes,omitempty"`
}