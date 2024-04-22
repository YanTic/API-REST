package main

import (
	"log"
	"net/http"
	Database "server/database"
	"server/handler"

	"github.com/gorilla/mux"
)

func main() {
	// docker run --rm -p 4453:3306 database-api
	Database.EstablishDBConnection()

	router := mux.NewRouter()

	router.HandleFunc("/", handler.LoginHandler).Methods("POST")
	router.HandleFunc("/password", handler.RecoverPasswordHandler).Methods("PATCH")
	router.HandleFunc("/password/{id}", handler.UpdateUserPasswordHandler).Methods("PATCH")

	usersRouter := router.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/", handler.GetUsersHandler).Methods("GET")
	usersRouter.HandleFunc("/{id}", handler.GetUserByIdHandler).Methods("GET")
	usersRouter.HandleFunc("/", handler.CreateUserHandler).Methods("POST")
	usersRouter.HandleFunc("/{id}", handler.UpdateUserHandler).Methods("PUT")
	usersRouter.HandleFunc("/{id}", handler.DeleteUserHandler).Methods("DELETE")

	servidor := &http.Server{
		Handler: router,
		Addr:    ":80",
	}

	log.Fatal(servidor.ListenAndServe())
}
