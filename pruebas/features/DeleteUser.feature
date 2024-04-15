Feature: La API permite a un usuario eliminar un usuario

    Background:
        Given Un usuario registrado en la Base de Datos que ya está logueado
        And suministra el token JWT en la cabecera Authentication

    Scenario: El usuario elimina un user exitosamente
        When El usuario hace la peticion DELETE a la ruta "/users/{id}"
        Then La API responde con un mensaje de exito
        And La API responde con un Status Code 200

    Scenario: El usuario elimina un user y el JWT token no es valido
        And el token JWT no es valido
        When El usuario hace la peticion DELETE a la ruta "/users/{id}"
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 401
    
    Scenario: El usuario elimina un user pero en la API ocurrió un error
        When El usuario hace la peticion DELETE a la ruta "/users/{id}"
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 400