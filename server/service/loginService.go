package service

import (
	"log"
	"net/http"
	"server/database"
	"server/jwt"
	"server/payload/response"
)

func VerifyIdentity(username string, password string) bool {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(1) FROM user WHERE username = ? AND password = ?", username, password).Scan(&count)
	if err != nil {
		log.Println("ERROR AL CONSULTAR LA BASE DE DATOS", http.StatusInternalServerError)
		return false
	}

	if count > 0 {
		return true // User exists
	}

	return false
}

func RecoverPassword(email string) (string, error) {
	log.Println("EMAIl:", email)
	rows, err := database.DB.Query("SELECT * FROM user WHERE email = ?", email)
	if err != nil {
		log.Println("ERROR AL CONSULTAR LA BASE DE DATOS", http.StatusInternalServerError)
		return "", err
	}
	defer rows.Close()

	var user response.User
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email)
		if err != nil {
			log.Println("Error scanning row", http.StatusInternalServerError)
			return "", err
		}
	} else {
		log.Println("No se encontró el usuario", http.StatusBadRequest)
		return "", err
	}

	recoveryToken, err := jwt.CreateToken(user.Username)
	if err != nil {
		log.Println("Error al creal token", http.StatusInternalServerError)
		return "", err
	}

	return recoveryToken, nil
}

func UpdateUserPassword(user response.User, idUser string) (response.User, error) {
	rows, err := database.DB.Query("SELECT * FROM user WHERE id = ? AND email = ?", idUser, user.Email)
	if err != nil {
		log.Println("ERROR AL CONSULTAR LA BASE DE DATOS", http.StatusInternalServerError)
		return response.User{}, err
	}
	defer rows.Close()

	// Si el usuario existe, la consulta debe retornar el registro
	if rows.Next() {
		_, err := database.DB.Exec("UPDATE user SET password = ? WHERE id = ?", user.Password, idUser)
		if err != nil {
			log.Println("ERROR AL ACTUALIZAR EN LA BASE DE DATOS", http.StatusInternalServerError)
			return response.User{}, err
		}
	}

	return user, nil
}
