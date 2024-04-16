Feature: La API permite a un usuario la opción de crear usuarios

    Background:
        Given Un usuario registrado en la Base de Datos que ya está logueado
        And envia en el request-body un JSON con los datos necesarios
        And suministra el token JWT en la cabecera Authentication

    Scenario: El usuario crea un usuario exitosamente
        When El usuario hace la peticion POST a la ruta "/users" 
        Then La API responde con un mensaje de exito
        And La API responde con un Status Code 200

    Scenario: El usuario crea un usuario y el JWT token no es valido
        And el token JWT no es valido
        When El usuario hace la peticion POST a la ruta "/users" 
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 401

    Scenario: El usuario crea un usuario y no envia un dato de registro
        And no envia un dato de registro "<campo>"
        When El usuario hace la peticion POST a la ruta "/users" 
        Then La API responde con un mensaje de error indicando que "<mensaje>"
        And La API responde con un Status Code 400
        Examples:
        | campo    | mensaje                     |
        | usuario  | El username es obligatorio  |
        | password | La password es obligatoria  |
        | email    | El email es obligatorio     |