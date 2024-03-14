package service

import (
	"log"
	"net/http"
	"server/database"
	"server/payload/response"
	"strconv"
)

func GetUsers(offset, pageSize int) ([]response.User, error) {
	rows, err := database.DB.Query("SELECT * FROM user ORDER BY id LIMIT " + strconv.Itoa(pageSize) + " OFFSET " + strconv.Itoa(offset))

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

func GetUserById(id string) (response.User, error) {
	rows, err := database.DB.Query("SELECT * FROM user WHERE id = " + id)
	if err != nil {
		log.Println("ERROR AL CONSULTAR LA BASE DE DATOS", http.StatusInternalServerError)
		return response.User{}, err
	}
	defer rows.Close()

	var user response.User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email)
		if err != nil {
			log.Println("Error scanning row", http.StatusInternalServerError)
			return response.User{}, err
		}
	}

	return user, nil
}

func CreateUser(user response.User) (response.User, error) {
	_, err := database.DB.Exec("INSERT INTO user (username, password, email) VALUES (?, ?, ?)",
		user.Username, user.Password, user.Email)
	if err != nil {
		log.Println("ERROR AL CONSULTAR LA BASE DE DATOS", http.StatusInternalServerError)
		return response.User{}, err
	}

	return user, nil
}

func UpdateUser(user response.User, idUser string) (response.User, error) {
	_, err := database.DB.Exec("UPDATE user SET username = ?, password = ?, email = ? WHERE id = ?",
		user.Username, user.Password, user.Email, idUser)
	if err != nil {
		log.Println("ERROR AL ACTUALIZAR EN LA BASE DE DATOS", http.StatusInternalServerError)
		return response.User{}, err
	}

	return user, nil
}

func DeleteUser(idUser string) error {
	_, err := database.DB.Exec("DELETE FROM user WHERE id = ?", idUser)
	if err != nil {
		log.Println("ERROR AL CONSULTAR LA BASE DE DATOS", http.StatusInternalServerError)
		return err
	}

	return nil
}
