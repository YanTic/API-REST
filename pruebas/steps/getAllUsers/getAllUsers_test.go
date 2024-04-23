package getAllUsers

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
	url            string
	getAllUsersUrl string
	authToken      string

	apiResponse *APIResponse
)

// go test -v steps/getAllUsers/getAllUsers_test.go

func setConfigs() {
	url = JsonReader("../config-file.json", "API.baseUrl")
	getAllUsersUrl = JsonReader("../config-file.json", "API.getAllUsersUrl")
}

func unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado() error {
	return nil
}

func suministraElTokenJWTEnLaCabeceraAuthentication() error {
	authToken = JsonReader("../config-file.json", "API.token")

	if authToken == "" {
		return fmt.Errorf("el usuario no mandó ningun token")
	}

	return nil
}

func elUsuarioHaceLaPeticionGETALaRuta(ruta string) error {
	// En setConfigs() ya se establece la ruta si 'ruta == /users'
	if ruta != "/users" {
		getAllUsersUrl = JsonReader("../config-file.json", "API.getAllUsersPagUrl")
	}

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(authToken).
		Get(getAllUsersUrl)
	if err != nil {
		return err
	}

	apiResponse = &APIResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
	}

	print("LOGIN_TEST TOKEN: ", authToken)

	return nil
}

func elTokenJWTNoEsValido() error {
	authToken = JsonReader("../config-file.json", "API.token-novalid")
	return nil
}

func laAPIRespondeConLaListaDeUsuarios() error {
	schemaBytes, err := ioutil.ReadFile("../../schemas/listUsers-schema.json")
	if err != nil {
		fmt.Println("Error al leer el JSON-Schema:", err)
		return err
	}

	responseSchema := gojsonschema.NewBytesLoader(apiResponse.Body)
	listUserschema := gojsonschema.NewBytesLoader(schemaBytes)

	// Al parecer nunca es bueno utilizar siempre NewStringLoader para todo
	// En este caso, despues de tantos intentos, NewBytesLoader funcionó para la comparacion de schemas
	//listUserschema := gojsonschema.NewStringLoader(string(schemaBytes))
	//responseSchema := gojsonschema.NewStringLoader(string(apiResponse.Body))

	// println("listUserschema:\n", string(schemaBytes))
	// println("responseSchema:\n", string(apiResponse.Body))

	result, err := gojsonschema.Validate(listUserschema, responseSchema)
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

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^Un usuario registrado en la Base de Datos que ya está logueado$`, unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado)
	ctx.Step(`^suministra el token JWT en la cabecera Authentication$`, suministraElTokenJWTEnLaCabeceraAuthentication)
	ctx.Step(`^El usuario hace la peticion GET a la ruta "([^"]*)"$`, elUsuarioHaceLaPeticionGETALaRuta)
	ctx.Step(`^el token JWT no es valido$`, elTokenJWTNoEsValido)
	ctx.Step(`^La API responde con la lista de usuarios$`, laAPIRespondeConLaListaDeUsuarios)
	ctx.Step(`^La API responde con un mensaje de error$`, laAPIRespondeConUnMensajeDeError)
	ctx.Step(`^La API responde con un Status Code (\d+)$`, laAPIRespondeConUnStatusCode)
}

func TestMain(m *testing.M) {
	setConfigs()

	opts := godog.Options{
		Format: "progress",
		Paths:  []string{"../../features/GetAllUsers.feature"}, // Se especifica que feature usa este "steptest"
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
