package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"server/service"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := service.GetUsers()
	if err != nil {
		log.Println("ERROR: No se enviaron los usuarios: ", err)
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Println("ERROR: Codificacion en JSON: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usersJSON)
}
