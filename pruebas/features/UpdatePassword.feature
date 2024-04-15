Feature: La API permite a un usuario actualizar su contraseña, luego de pedir la recuperación

    Background:
        Given Un usuario registrado en la Base de Datos que ya está logueado
        And envia en el request-body un JSON con la contraseña
        And suministra el token JWT de recuperación en la cabecera Authentication

    Scenario: El usuario actualiza su contraseña exitosamente
        When El usuario hace la peticion PATCH a la ruta "/password/{id}"
        Then La API responde con un mensaje de exito
        And La API responde con un Status Code 200

    Scenario: El usuario actualiza su contraseña y el JWT no es valido
        And el token JWT no es valido
        When El usuario hace la peticion PATCH a la ruta "/password/{id}"
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 401

    Scenario: El usuario actualiza su contraseña y no envia la nueva contraseña
        And no envia la nueva contraseña en el request-body
        When El usuario hace la peticion PATCH a la ruta "/password/{id}"
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 400