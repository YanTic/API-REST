Feature: La API permite a un usuario recuperar su contraseña, recibiendo un JWT TOKEN para actualizar su contraseña

    Background:
        Given Un usuario registrado en la Base de Datos que ya está logueado
        And envia en el request-body un JSON con el email

    Scenario: El usuario pide recuperar su contraseña exitosamente
        When El usuario hace la peticion GET a la ruta "/password"
        Then La API responde el token JWT de autenticacion
        And La API responde con un Status Code 200

    Scenario: El usuario pide recuperar su contraseña y no envia el correo (o uno valido)
        And no envia el correo electronico en el request-body
        When El usuario hace la peticion GET a la ruta "/password"
        Then La API responde con un mensaje de error
        And La API responde con un Status Code 500