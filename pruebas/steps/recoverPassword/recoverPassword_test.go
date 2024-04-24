package recoverpassword

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"github.com/xeipuuv/gojsonschema"
)

type APIResponse struct {
	StatusCode int
	Body       []byte
}

var (
	url                string
	recoverPasswordUrl string
	requestBody        string

	apiResponse *APIResponse
)

// go test -v steps/recoverPassword/recoverPassword_test.go

func setConfigs() {
	url = JsonReader("../config-file.json", "API.baseUrl")
	recoverPasswordUrl = JsonReader("../config-file.json", "API.recoverPasswordUrl")
}

func unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado() error {
	return nil
}

func enviaEnElRequestbodyUnJSONConElEmail() error {
	requestBody = JsonReader("request-body.json", "data")
	return nil
}

func noEnviaElCorreoElectronicoEnElRequestbody() error {
	requestBody = JsonReader("request-body.json", "data-noemail")
	return nil
}

// En una peticion GET nunca se manda un "Body", eso es un error de la API
// Aunque Postman deje usar un body en una peticion GET, no deberia, incluso
// la libreria "Resty" que se está usando para las peticiones no deja

// Para ese problema se cambió la API, ahora no es GET sino PATCH
func elUsuarioHaceLaPeticionGETALaRuta(arg1 string) error {
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Patch(recoverPasswordUrl)
	if err != nil {
		return err
	}

	apiResponse = &APIResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
	}

	return nil
}

func laAPIRespondeConUnMensajeDeError() error {
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

func laAPIRespondeElTokenJWTDeAutenticacion() error {
	schemaBytes, err := ioutil.ReadFile("../../schemas/jwt-schema.json")
	if err != nil {
		fmt.Println("Error al leer el JSON-Schema:", err)
		return err
	}

	responseSchema := gojsonschema.NewBytesLoader(apiResponse.Body)
	tokenSchema := gojsonschema.NewBytesLoader(schemaBytes)

	result, err := gojsonschema.Validate(tokenSchema, responseSchema)
	if err != nil {
		fmt.Print("Error al validar el schema con la respuesta: ", err.Error())
		return err
	}

	if !result.Valid() {
		fmt.Print("La respuesta no está en el esquema correcto (List Users Schema)")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return err
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^Un usuario registrado en la Base de Datos que ya está logueado$`, unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado)
	ctx.Step(`^envia en el request-body un JSON con el email$`, enviaEnElRequestbodyUnJSONConElEmail)
	ctx.Step(`^no envia el correo electronico en el request-body$`, noEnviaElCorreoElectronicoEnElRequestbody)
	ctx.Step(`^El usuario hace la peticion GET a la ruta "([^"]*)"$`, elUsuarioHaceLaPeticionGETALaRuta)
	ctx.Step(`^La API responde con un mensaje de error$`, laAPIRespondeConUnMensajeDeError)
	ctx.Step(`^La API responde con un Status Code (\d+)$`, laAPIRespondeConUnStatusCode)
	ctx.Step(`^La API responde el token JWT de autenticacion$`, laAPIRespondeElTokenJWTDeAutenticacion)
}

func TestMain(m *testing.M) {
	setConfigs()

	f, err := os.Create("../../reports/reports-json/report_recoverPassword.json")
	if err != nil {
		fmt.Print("error at creating report: ", err)
	}

	opts := godog.Options{
		Format: "cucumber",
		Paths:  []string{"../../features/RecoverPassword.feature"}, // Se especifica que feature usa este "steptest"
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
