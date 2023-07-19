package routes

import (
	"github.com/Ashmit-05/LockItUp/controllers"
	"github.com/gorilla/mux"
)

func SetUserRoutes(r *mux.Router) {
	r.HandleFunc("/api/auth/signup", controllers.SignUp).Methods("POST")
	r.HandleFunc("/api/auth/signin", controllers.SignIn).Methods("POST")
}
