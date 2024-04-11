package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	jwt "server/jwt"
	"server/payload/response"
	"server/service"

	"github.com/gorilla/mux"
)

type User struct {
	Username string
	Password string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Con esto se le dice al cliente que espera recibir datos en formato JSON
	w.Header().Set("Content-Type", "application/json")

	var u response.User
	json.NewDecoder(r.Body).Decode(&u)
	fmt.Printf("Valores enviados por JSON %v", u)

	if u.Username == "" || u.Password == "" {
		http.Error(w, "Usuario y Contraseña son obligatorios", http.StatusBadRequest)
	} else {

		if service.VerifyIdentity(u.Username, u.Password) {
			tokenString, err := jwt.CreateToken(u.Username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprint("Token couldn't be created", http.StatusInternalServerError)))
				return
			}

			// Se construye la respuesta en string, como se pide en el jwt-schema.json
			responseString := map[string]string{"token": tokenString}

			// Se codifica la respuesta en JSON
			// Esto se hizo porque para hacer el testing con cucumber, se necesita
			// validar el schema del JWT Token (ver jwt-schema.json en schemas/)
			responseJSON, err := json.Marshal(responseString)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Error al codificar la respuesta JSON: %v", err)))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(responseJSON)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Error: No se pudo hacer Login"))
}

func RecoverPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var u response.User
	json.NewDecoder(r.Body).Decode(&u)
	fmt.Printf("Valores enviados por JSON %v", u)

	token, err := service.RecoverPassword(u.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error: No se puede recuperar la contraseña", http.StatusInternalServerError)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func UpdateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if !verifyTokenHandler(w, r) {
		return
	}

	var u response.User
	json.NewDecoder(r.Body).Decode(&u)
	userId := mux.Vars(r)["id"]

	if u.Password == "" {
		http.Error(w, "La nueva contraseña es necesaria", http.StatusBadRequest)
	} else {

		_, err := service.UpdateUserPassword(u, userId)
		if err != nil {
			log.Println("ERROR: No se actualizó la contraseña: ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ERROR: No se actualizó la contraseña"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Contraseña actualizada con exito!"))
	}

}
