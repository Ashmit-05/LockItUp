package main

import (
	"fmt"
	"log"
	"net/http"

	userRouter "github.com/Ashmit-05/LockItUp/routes"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting LockItUp server...")
	r := mux.NewRouter()

	userRouter.SetUserRoutes(r)

	log.Fatal(http.ListenAndServe(":8000", r))

}
