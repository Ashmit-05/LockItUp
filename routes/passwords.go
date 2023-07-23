package routes

import (
	"github.com/Ashmit-05/LockItUp/controllers"
	"github.com/Ashmit-05/LockItUp/middlewares"
	"github.com/gorilla/mux"
)

func SetPasswordRoutes(r *mux.Router) {
	r.HandleFunc("/api/passwords/add/{userId}", middlewares.Authenticate(controllers.AddPassword)).Methods("POST")
	r.HandleFunc("/api/passwords/all/{userId}", middlewares.Authenticate(controllers.GetAllPasswords)).Methods("GET")
	r.HandleFunc("/api/passwords/generate",middlewares.Authenticate(controllers.GeneratePassword)).Methods("GET")
	// r.HandleFunc("/api/passwords/update/{userId}&{passwordId}",middlewares.Authenticate(controllers.UpdatePassword)).Methods("POST")
	r.HandleFunc("/api/passwords/delete/{userId}&{passwordId}",middlewares.Authenticate(controllers.DeletePassword)).Methods("DELETE")
}
