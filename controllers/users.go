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
)
