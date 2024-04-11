Feature: La API permite a un usuario loguearse

    Background:
        Given Un usuario registrado en la Base de Datos
        And proporciona los datos acceso

    Scenario: El usuario se loguea exitosamente
        When El usuario hace la peticion POST a la ruta "/"
        Then La API responde el token JWT de autenticacion
        And La API responde con un Status Code 200

    Scenario: El usuario se loguea
        And no envia un dato de logueo "<campo>"
        When El usuario hace la peticion POST a la ruta "/"
        Then La API responde con un mensaje de error indicando que "<mensaje>"
        And La API responde con un Status Code 400
        Examples:
        | campo    | mensaje                              |
        | usuario  | El nombre de usuario es obligatorio. |
        | clave    | La clave es obligatoria.             |