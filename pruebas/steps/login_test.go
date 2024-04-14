package steps

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"github.com/xeipuuv/gojsonschema"
)

// Estructura para representar la respuesta de la API
type APIResponse struct {
	StatusCode int
	Body       []byte
}

var (
	url          string
	loginUserUrl string
	requestBody  string // Esta variable contiene los datos que se envian a la API

	apiResponse *APIResponse
)

// go test -v steps/login_test.go

func setConfigs() {
	url = JsonReader("config-file.json", "API.baseUrl")
	loginUserUrl = JsonReader("config-file.json", "API.loginUserUrl")
}

// Función para enviar una solicitud POST a la API utilizando Resty
func elUsuarioHaceLaPeticionPOSTALaRuta(ruta string) error {
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(loginUserUrl)
	if err != nil {
		return err
	}

	apiResponse = &APIResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
	}
	return nil
}

// Función para verificar que la API responde con un código de estado específico
func laAPIRespondeConUnStatusCode(codigo int) error {
	// print("laAPIRespondeConUnStatusCode")
	// print("status code api: ", apiResponse.StatusCode)
	// print("status code: ", codigo)
	if apiResponse.StatusCode != codigo {
		return fmt.Errorf("el código de estado esperado es %d, pero se recibió %d", codigo, apiResponse.StatusCode)
	}
	return nil
}

// Función para verificar que la API responde con un mensaje de error específico
func laAPIRespondeConUnMensajeDeErrorIndicandoQue(mensaje string) error {
	if strings.TrimSpace(string(apiResponse.Body)) != strings.TrimSpace(mensaje) {
		return fmt.Errorf("el mensaje de error esperado es '%s', pero se recibió '%s'", mensaje, string(apiResponse.Body))
	}
	return nil
}

// Función para verificar que la API responde el JWT Token (según el json-schema)
func laAPIRespondeElTokenJWTDeAutenticacion() error {
	schemaBytes, err := ioutil.ReadFile("../schemas/jwt-schema.json")
	if err != nil {
		fmt.Println("Error al leer el JSON Schema:", err)
		return err
	}

	tokenschema := gojsonschema.NewStringLoader(string(schemaBytes))
	responseSchema := gojsonschema.NewStringLoader(string(apiResponse.Body))

	// print("tokenschema:", string(schemaBytes))
	// print("responseSchema:", string(apiResponse.Body))

	result, err := gojsonschema.Validate(responseSchema, tokenschema)
	if err != nil {
		fmt.Print("Error al validar el schema con la respuesta")
		return err
	}

	if !result.Valid() {
		fmt.Print("La respuesta no está en el esquema correcto (JWT Token Schema)")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return err
	}

	return nil
}

// Función para verificar que no se envía un dato de logueo
// TODO: TERMINAR
func noEnviaUnDatoDeLogueo(campo string) error {
	if campo == "username" {
		requestBody = JsonReader("request-body.json", "data-nouser")
	} else {
		requestBody = JsonReader("request-body.json", "data-nopass")
	}

	// Este print, se ejecuta dos veces, porque el .feature tiene dos valores en la tabla
	// print("noEnviaUnDatoDeLogueo:", requestBody)
	return nil
}

// Función para proporcionar los datos de acceso
func proporcionaLosDatosAcceso() error {
	requestBody = JsonReader("request-body.json", "data")
	return nil
}

// Esta peticion digamos que no se prueba aquí
func unUsuarioRegistradoEnLaBaseDeDatos() error {
	return nil
}

// Inicializar el contexto de la prueba
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^Un usuario registrado en la Base de Datos$`, unUsuarioRegistradoEnLaBaseDeDatos)
	ctx.Step(`^proporciona los datos acceso$`, proporcionaLosDatosAcceso)
	ctx.Step(`^El usuario hace la peticion POST a la ruta "([^"]*)"$`, elUsuarioHaceLaPeticionPOSTALaRuta)
	ctx.Step(`^La API responde con un mensaje de error indicando que "([^"]*)"$`, laAPIRespondeConUnMensajeDeErrorIndicandoQue)
	ctx.Step(`^La API responde con un Status Code (\d+)$`, laAPIRespondeConUnStatusCode)
	ctx.Step(`^La API responde el token JWT de autenticacion$`, laAPIRespondeElTokenJWTDeAutenticacion)
	ctx.Step(`^no envia un dato de logueo "([^"]*)"$`, noEnviaUnDatoDeLogueo)
}

func TestMain(m *testing.M) {
	setConfigs()

	opts := godog.Options{
		Format: "progress",
		Paths:  []string{"../features/Login.feature"}, // Se especifica que feature usa este "steptest"
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
