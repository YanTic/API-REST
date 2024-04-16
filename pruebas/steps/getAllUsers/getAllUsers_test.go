package getAllUsers

import "github.com/cucumber/godog"

func unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado() error {
	return godog.ErrPending
}

func suministraElTokenJWTEnLaCabeceraAuthentication() error {
	return godog.ErrPending
}

func elUsuarioHaceLaPeticionGETALaRuta(arg1 string) error {
	return godog.ErrPending
}

func elTokenJWTNoEsValido() error {
	return godog.ErrPending
}

func laAPIRespondeConLaListaDeUsuarios() error {
	return godog.ErrPending
}

func laAPIRespondeConLaListaDeUsuariosSegnLaPaginacionDada() error {
	return godog.ErrPending
}

func laAPIRespondeConUnMensajeDeError() error {
	return godog.ErrPending
}

func laAPIRespondeConUnStatusCode(arg1 int) error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^Un usuario registrado en la Base de Datos que ya está logueado$`, unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado)
	ctx.Step(`^suministra el token JWT en la cabecera Authentication$`, suministraElTokenJWTEnLaCabeceraAuthentication)
	ctx.Step(`^El usuario hace la peticion GET a la ruta "([^"]*)"$`, elUsuarioHaceLaPeticionGETALaRuta)
	ctx.Step(`^el token JWT no es valido$`, elTokenJWTNoEsValido)
	ctx.Step(`^La API responde con la lista de usuarios$`, laAPIRespondeConLaListaDeUsuarios)
	ctx.Step(`^La API responde con la lista de usuarios, según la paginacion dada$`, laAPIRespondeConLaListaDeUsuariosSegnLaPaginacionDada)
	ctx.Step(`^La API responde con un mensaje de error$`, laAPIRespondeConUnMensajeDeError)
	ctx.Step(`^La API responde con un Status Code (\d+)$`, laAPIRespondeConUnStatusCode)
}
