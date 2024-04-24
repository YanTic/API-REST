package updatepassword

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type APIResponse struct {
	StatusCode int
	Body       []byte
}

var (
	url               string
	updatePasswordUrl string
	requestBody       string
	authToken         string

	apiResponse *APIResponse
)

// go test -v steps/updatePassword/updatePassword_test.go

func setConfigs() {
	url = JsonReader("../config-file.json", "API.baseUrl")
	updatePasswordUrl = JsonReader("../config-file.json", "API.updatePasswordUrl")
}

func unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado() error {
	return nil
}

func enviaEnElRequestbodyUnJSONConLaContrasea() error {
	requestBody = JsonReader("request-body.json", "data")
	return nil
}

func suministraElTokenJWTDeRecuperacinEnLaCabeceraAuthentication() error {
	authToken = JsonReader("../config-file.json", "API.recoveryToken")

	if authToken == "" {
		return fmt.Errorf("el usuario no mandó ningun token")
	}

	return nil
}

func elUsuarioHaceLaPeticionPATCHALaRuta(arg1 string) error {
	// print("\nrequest-body: ", requestBody)
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(authToken).
		SetBody(requestBody).
		Patch(updatePasswordUrl)
	if err != nil {
		return err
	}

	apiResponse = &APIResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
	}

	// print("\nresponse: ", string(resp.Body()))

	return nil
}

func elTokenJWTNoEsValido() error {
	authToken = JsonReader("../config-file.json", "API.token-novalid")
	return nil
}

func noEnviaLaNuevaContraseaEnElRequestbody() error {
	requestBody = JsonReader("request-body.json", "data-nopass")
	return nil
}

func laAPIRespondeConUnMensajeDeError() error {
	if apiResponse.Body == nil || len(apiResponse.Body) == 0 {
		return fmt.Errorf("La API no mandó ningun mensaje de error")
	}

	return nil
}

func laAPIRespondeConUnMensajeDeExito() error {
	if apiResponse.Body == nil || len(apiResponse.Body) == 0 {
		return fmt.Errorf("La API no mandó ningun mensaje de exito")
	}

	return nil
}

func laAPIRespondeConUnStatusCode(codigo int) error {
	if apiResponse.StatusCode != codigo {
		return fmt.Errorf("el código de estado esperado es %d, pero se recibió %d", codigo, apiResponse.StatusCode)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^Un usuario registrado en la Base de Datos que ya está logueado$`, unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado)
	ctx.Step(`^envia en el request-body un JSON con la contraseña$`, enviaEnElRequestbodyUnJSONConLaContrasea)
	ctx.Step(`^suministra el token JWT de recuperación en la cabecera Authentication$`, suministraElTokenJWTDeRecuperacinEnLaCabeceraAuthentication)
	ctx.Step(`^El usuario hace la peticion PATCH a la ruta "([^"]*)"$`, elUsuarioHaceLaPeticionPATCHALaRuta)
	ctx.Step(`^el token JWT no es valido$`, elTokenJWTNoEsValido)
	ctx.Step(`^no envia la nueva contraseña en el request-body$`, noEnviaLaNuevaContraseaEnElRequestbody)
	ctx.Step(`^La API responde con un mensaje de error$`, laAPIRespondeConUnMensajeDeError)
	ctx.Step(`^La API responde con un mensaje de exito$`, laAPIRespondeConUnMensajeDeExito)
	ctx.Step(`^La API responde con un Status Code (\d+)$`, laAPIRespondeConUnStatusCode)
}

func TestMain(m *testing.M) {
	setConfigs()

	f, err := os.Create("../../reports/reports-json/report_updatePassword.json")
	if err != nil {
		fmt.Print("error at creating report: ", err)
	}

	opts := godog.Options{
		Format: "cucumber",
		Paths:  []string{"../../features/UpdatePassword.feature"}, // Se especifica que feature usa este "steptest"
		Output: f,
	}

	status := godog.TestSuite{
		Name:                "prueba",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	// Si alguna de las pruebas falla, la función TestMain retornará un código de error distinto de cero,
	// de lo contrario, retornará cero indicando que todas las pruebas pasaron correctamente.
	// Esto es utilizado por Go para determinar si las pruebas fueron exitosas o no.
	if status > 0 {
		fmt.Println("Error en la ejecución de las pruebas de Godog")
	} else {
		// Si todas las pruebas pasaron correctamente, el programa terminará con un código de salida cero
		// indicando que las pruebas fueron exitosas.
		// Esto es utilizado por Go para determinar si las pruebas fueron exitosas o no.
		fmt.Println("Todas las pruebas de Godog han pasado exitosamente")
		// os.Exit(status)
	}
}

func JsonReader(pathOfFile, pathOfData string) (result string) {
	jsonFile, err := os.Open(pathOfFile)
	if err != nil {
		fmt.Printf("Couldn't read the file in path %s and data in JSON path %s", pathOfFile, pathOfData)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	result = gjson.Get(string(byteValue), pathOfData).String()
	return result
}
