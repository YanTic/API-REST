Feature: La API permite a un usuario la opción de registrarse

    Background:
        Given Un usuario no registrado 
        And proporciona todos los datos de registro

    Scenario: El usuario se registra exitosamente
        When El usuario hace la peticion POST a la ruta /users
        Then La API responde con un mensaje de exito
        And La API responde con un Status Code 200

    Scenario: El usuario se registra
        And no envia un dato de registro "<campo>"
        When El usuario hace la peticion GET a la ruta /users/
        Then La API responde con un mensaje de error indicando que "<mensaje>"
        And La API responde con un Status Code 400
        Examples:
        | campo    | mensaje                              |
        | usuario  | El nombre de usuario es obligatorio. |
        | clave    | La clave es obligatoria.             |
        | email    | El email es obligatorio.             |