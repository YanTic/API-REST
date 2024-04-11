Feature: La API permite a un usuario actualizar los datos de un usuario

    Background:
        Given Un usuario registrado en la Base de Datos que ya está logueado
        And suministra el token JWT en la cabecera Authentication

    Scenario: El usuario actualiza los datos de un user
        When El usuario hace la peticion PUT a la ruta /users
        Then La API responde con un mensaje de exito
        And La API responde con un Status Code 200

    Scenario: El usuario actualiza los datos de un user
        And el token JWT no es valido
        When El usuario hace la peticion PUT a la ruta /users
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 401