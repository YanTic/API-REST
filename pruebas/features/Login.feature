Feature: La API permite a un usuario loguearse

    Background:
        Given Un usuario registrado en la Base de Datos
        And proporciona los datos acceso

    Scenario: El usuario se loguea exitosamente
        When El usuario hace la peticion POST a la ruta "/"
        Then La API responde el token JWT de autenticacion
        And La API responde con un Status Code 200

    Scenario: El usuario se loguea y no envia un dato de logueo
        And no envia un dato de logueo "<campo>"
        When El usuario hace la peticion POST a la ruta "/"
        Then La API responde con un mensaje de error indicando que "<mensaje>"
        And La API responde con un Status Code 400
        Examples:
        | campo    | mensaje                               |
        | username | Usuario y Contraseña son obligatorios |
        | password | Usuario y Contraseña son obligatorios |