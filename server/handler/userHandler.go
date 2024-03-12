package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/jwt"
	"server/payload/response"
	"server/service"

	"github.com/gorilla/mux"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if !verifyTokenHandler(w, r) {
		return
	}

	users, err := service.GetUsers()
	if err != nil {
		log.Println("ERROR: No se enviaron los usuarios: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: No se enviaron los usuarios"))
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Println("ERROR: Codificacion en JSON: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: Codificacion en JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usersJSON)
}

func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	if !verifyTokenHandler(w, r) {
		return
	}

	id := mux.Vars(r)["id"]
	user, err := service.GetUserById(id)
	if err != nil {
		log.Println("ERROR: No se envió el usuario: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: No se envió el usuario"))
		return
	}

	usersJSON, err := json.Marshal(user)
	if err != nil {
		log.Println("ERROR: Codificacion en JSON: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: Codificacion en JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usersJSON)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if !verifyTokenHandler(w, r) {
		return
	}

	var newUser response.User
	json.NewDecoder(r.Body).Decode(&newUser)
	_, err := service.CreateUser(newUser)
	if err != nil {
		log.Println("ERROR: No se creó el usuario: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: No se creó el usuario"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Usuario creado con exito!"))
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if !verifyTokenHandler(w, r) {
		return
	}

	var newUser response.User
	json.NewDecoder(r.Body).Decode(&newUser)
	userId := mux.Vars(r)["id"]

	_, err := service.UpdateUser(newUser, userId)
	if err != nil {
		log.Println("ERROR: No se actualizó el usuario: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: No se actualizó el usuario"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Usuario actualizado con exito!"))
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if !verifyTokenHandler(w, r) {
		return
	}

	userId := mux.Vars(r)["id"]

	err := service.DeleteUser(userId)
	if err != nil {
		log.Println("ERROR: No se eliminó el usuario: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: No se eliminó el usuario"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Usuario eliminado con exito!"))
}

// Funciona para verificar la cabecera Authorization y el Token que envia el user para usar el handler
func verifyTokenHandler(w http.ResponseWriter, r *http.Request) bool {
	// Se verifica que la solicitud tenga una cabecera Authorization con un JWT valido
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Falta la cabecera Authorization"))
		return false
	}
	tokenString = tokenString[len("Bearer "):] // Se elimina el Bearer del token (Bearer eyJhbGciOiJ...)

	// Se verifica el token enviado
	err := jwt.VerifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Token invalido | %v", err)))
		return false
	}

	return true
}
