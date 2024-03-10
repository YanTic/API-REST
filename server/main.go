package main

import (
	"log"
	"net/http"
	Database "server/database"
	"server/handler"

	"github.com/gorilla/mux"
)

func main() {
	Database.EstablishDBConnection()

	router := mux.NewRouter()

	usersRouter := router.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/", handler.GetUsersHandler).Methods("GET")

	servidor := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:80",
	}

	log.Fatal(servidor.ListenAndServe())
}
