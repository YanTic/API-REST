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

	router.HandleFunc("/", handler.LoginHandler).Methods("POST")
	router.HandleFunc("/password/{email}", handler.RecoverPassword).Methods("GET")
	router.HandleFunc("/password", handler.UpdateUserPassword).Methods("POST")

	usersRouter := router.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/", handler.GetUsersHandler).Methods("GET")
	usersRouter.HandleFunc("/{id}", handler.GetUserByIdHandler).Methods("GET")
	usersRouter.HandleFunc("/", handler.CreateUserHandler).Methods("POST")
	usersRouter.HandleFunc("/{id}", handler.UpdateUserHandler).Methods("PUT")
	usersRouter.HandleFunc("/{id}", handler.DeleteUserHandler).Methods("DELETE")

	servidor := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:80",
	}

	log.Fatal(servidor.ListenAndServe())
}
