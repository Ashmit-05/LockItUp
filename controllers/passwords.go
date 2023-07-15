package passwordController

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

  passwordModel "github.com/Ashmit-05/LockItUp/models/passwords.go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

