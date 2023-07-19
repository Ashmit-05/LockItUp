package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ashmit-05/LockItUp/routes"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("LockItUp")
	r := mux.NewRouter()

	routes.SetUserRoutes(r)
	routes.SetPasswordRoutes(r)


	log.Fatal(http.ListenAndServe(":8000", r))

}
