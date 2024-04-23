package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("secret-key")
var emisor = "ingesis.uniquindio.edu.co"

// var SUB string # Ya no necesario para validar el token

// CADA VEZ QUE SE CREA UN TOKEN, EN jwtAuth.go la variable SUB quedá con el nuevo usuario al
// que se le creó el token, y si quiero usar un token anterior que tiene como SUB otro usuario, NO SIRVE

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iss": emisor,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// SUB = username # Ya no necesario para validar el token
	return tokenString, err
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}

	// Verifica que la firma del token sea valido y la fecha de expiración
	if !token.Valid {
		return fmt.Errorf("token invalido | Firma o Fecha de expiración (exp) no valido")
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["iss"] != emisor {
		return fmt.Errorf("ISS (emisor) incorrecto")
	}

	// TODO: ENCONTRAR UNA FORMA DE VALIDAR EL SUB DEL TOKEN (aunque no sea necesario, opcional)
	// When your authentication server receives an incoming JWT, it uses the incoming JWT's header
	// and payload segments and the shared private key to generate a signature.
	// If the signature matches, then your application knows that the incoming JWT can be trusted.
	// https://www.freecodecamp.org/news/how-to-sign-and-validate-json-web-tokens/
	// Por esta razon no se valida el SUB, porque como se obtiene el SUB (username) por la peticion
	// Aunque tambien se puede validar haciendo una consulta a la BD solo con el SUB (o username)
	// fmt.Printf("sub: %s | username: %s\n", claims["sub"], SUB)
	// if claims["sub"] != SUB {
	// 	return fmt.Errorf("SUB (usuario) incorrecto")
	// }

	return nil
}
