Feature: La API permite a un usuario actualizar los datos de un usuario

    Background:
        Given Un usuario registrado en la Base de Datos que ya está logueado
        And envia en el request-body un JSON con los datos necesarios
        And suministra el token JWT en la cabecera Authentication

    Scenario: El usuario actualiza los datos de un user exitosamente
        When El usuario hace la peticion PUT a la ruta "/users"
        Then La API responde con un mensaje de exito
        And La API responde con un Status Code 200

    Scenario: El usuario actualiza los datos de un user y el JWT token no es valido
        And el token JWT no es valido
        When El usuario hace la peticion PUT a la ruta "/users"
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 401

    Scenario: El usuario actualiza los datos de un user y la API no encuentra el usuario
        And la API no encuentra al usuario
        When El usuario hace la peticion PUT a la ruta "/users"
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 400