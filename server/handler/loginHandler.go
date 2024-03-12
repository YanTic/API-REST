package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	jwt "server/jwt"
	"server/payload/response"
	"server/service"
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
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tokenString))
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Error: No se pudo hacer Login"))
}

func RecoverPassword(w http.ResponseWriter, r *http.Request) {

}

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {

}

func SaludoHandler(w http.ResponseWriter, r *http.Request) {
	// Con esto se le dice al cliente que espera recibir datos en formato JSON
	w.Header().Set("Content-Type", "application/json")
	usuario := r.URL.Query().Get("nombre")

	// Se Verifique que la solicitud contenga una cabecera Authorization con un JWT valido
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Falta la cabecera Authorization")
		return
	}
	tokenString = tokenString[len("Bearer "):] // Se elimina el Bearer del token (Bearer eyJhbGciOiJ...)

	// Se verifica el token enviado
	err := jwt.VerifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Token invalido | ", err)
		return
	}

	// Se verifica si el cliente envió un usuario
	if usuario == "" {
		http.Error(w, "Solicitud no valida: El nombre es obligatorio", http.StatusNotFound)
		return
	} else {
		response := fmt.Sprintf("Hola %s", usuario)
		fmt.Fprintln(w, response)
		w.WriteHeader(http.StatusOK)
	}
}
