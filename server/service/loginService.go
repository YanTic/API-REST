package service

import (
	"log"
	"net/http"
	"server/database"
)

func VerifyIdentity(username string, password string) bool {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM user WHERE username = ? AND password = ?", username, password).Scan(&count)
	if err != nil {
		log.Println("ERROR AL CONSULTAR LA BASE DE DATOS", http.StatusInternalServerError)
		return false
	}

	if count > 0 {
		return true // User exists
	}

	return false
}
