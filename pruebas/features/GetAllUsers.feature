Feature: La API permite la opción de que el usuario puede obtener una lista de todos los usuarios

    Background:
        Given Un usuario registrado en la Base de Datos que ya está logueado
        And suministra el token JWT en la cabecera Authentication

    Scenario: El usuario pide todos los usuarios exitosamente
        When El usuario hace la peticion GET a la ruta "/users"
        Then La API responde con la lista de usuarios 
        And La API responde con un Status Code 200

    Scenario: El usuario pide todos los usuarios con paginacion exitosamente
        When El usuario hace la peticion GET a la ruta "/users?offset=1&pagesize=2"
        Then La API responde con la lista de usuarios
        And La API responde con un Status Code 200

    Scenario: El usuario pide todos los usuarios y el JWT token no es valido
        And el token JWT no es valido
        When El usuario hace la peticion GET a la ruta "/users"
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 401
    
    Scenario: El usuario pide todos los usuarios con paginacion y el JWT token no es valido
        And el token JWT no es valido
        When El usuario hace la peticion GET a la ruta "/users?offset=1&pagesize=2"
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 401