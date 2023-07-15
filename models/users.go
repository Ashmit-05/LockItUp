package userModel

import (
  "go.mongodb.org/mongo-driver/bson/primitive"
  passwordModel "github.com/Ashmit-05/LockItUp/models/passwords.go"
)

type User struct {
  ID primitive.ObjectID `json:"_id,omitempty"` `bson:"_id,omitempty"`
  Name string `json:"name"`
  Email string `json:"email"`
  PhoneNumber string `json:"phonenumber,omitempty"`
  MasterPassword string `json:"masterpassword"`
  Passwords []*passwordModel `json:"passwords,omitempty"`
}
