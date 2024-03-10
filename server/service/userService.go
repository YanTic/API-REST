package service

import (
	"log"
	"net/http"
	"server/database"
	"server/payload/response"
)

func GetUsers() ([]response.User, error) {
	rows, err := database.DB.Query("SELECT * FROM user")
	if err != nil {
		log.Println("ERROR AL CONSULTAR LA BASE DE DATOS", http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()

	var users []response.User

	for rows.Next() {
		var user response.User

		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email)
		if err != nil {
			log.Println("Error scanning row", http.StatusInternalServerError)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
