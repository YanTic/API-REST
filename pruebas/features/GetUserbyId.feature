Feature: La API permite la opción de que el usuario pueda obtener un usuario solo con el ID de tal usuario

    Background:
        Given Un usuario registrado en la Base de Datos que ya está logueado
        And suministra el token JWT en la cabecera Authentication

    Scenario: El usuario pide un usuario especifico exitosamente
        When El usuario hace la peticion GET a la ruta "/users/2"
        Then La API responde con los datos del usuario
        And La API responde con un Status Code 200

    Scenario: El usuario pide un usuario especifico y el JWT token no es valido
        And el token JWT no es valido
        When El usuario hace la peticion GET a la ruta "/users/2"
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 401
    
    Scenario: El usuario pide un usuario especifico y la API no encuentra el usuario
        And la API no encuentra al usuario
        When El usuario hace la peticion GET a la ruta "/users/29999"
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 400