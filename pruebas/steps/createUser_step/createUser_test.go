package createUser_step

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

// Algunas funciones son repetidas, por lo que son declaradas en otro archivo .go
// del directorio steps. Por lo que para solucionar eso se crea para cada _test.go
// un directorio (package)
// AYUDA: https://stackoverflow.com/questions/66970531/vs-code-go-main-redeclared-in-this-block

// Estructura para representar la respuesta de la API
type APIResponse struct {
	StatusCode int
	Body       []byte
}

var (
	url           string
	createUserUrl string
	requestBody   string // Esta variable contiene los datos que se envian a la API

	apiResponse *APIResponse
)

// go test -v steps/login_test.go

func setConfigs() {
	url = JsonReader("../config-file.json", "API.baseUrl")
	createUserUrl = JsonReader("../config-file.json", "API.createUserUrl")
}

// Esta funcion es mas como por decir especificamente que se hace en el
// test, pero en testing de cucumber no se hacen consultas a la bd.
// A menos que en la API se haga.
// Ademas esto, se supone que el usuario ya está logueado, por lo que
// ya el servicio de Login debió verificar que el usuario esté en la BD
func unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado() error {
	return nil
}

func enviaEnElRequestbodyUnJSONConLosDatosNecesarios() error {
	requestBody = JsonReader("request-body.json", "data")
	return nil
}

func suministraElTokenJWTEnLaCabeceraAuthentication() error {
	return godog.ErrPending
}

func elUsuarioHaceLaPeticionPOSTALaRuta(arg1 string) error {
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(createUserUrl)

	if err != nil {
		return err
	}

	apiResponse = &APIResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
	}
	return nil
}

func laAPIRespondeConUnMensajeDeExito() error {
	return godog.ErrPending
}

func laAPIRespondeConUnStatusCode(codigo int) error {
	if apiResponse.StatusCode != codigo {
		return fmt.Errorf("el código de estado esperado es %d, pero se recibió %d", codigo, apiResponse.StatusCode)
	}
	return nil
}

func elTokenJWTNoEsValido() error {
	return godog.ErrPending
}

func laAPIRespondeConUnMensajeDeError() error {
	return godog.ErrPending
}

func noEnviaUnDatoDeRegistro(arg1 string) error {
	return godog.ErrPending
}

func laAPIRespondeConUnMensajeDeErrorIndicandoQue(arg1 string) error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^Un usuario registrado en la Base de Datos que ya está logueado$`, unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado)
	ctx.Step(`^envia en el request-body un JSON con los datos necesarios$`, enviaEnElRequestbodyUnJSONConLosDatosNecesarios)
	ctx.Step(`^suministra el token JWT en la cabecera Authentication$`, suministraElTokenJWTEnLaCabeceraAuthentication)
	ctx.Step(`^El usuario hace la peticion POST a la ruta "([^"]*)"$`, elUsuarioHaceLaPeticionPOSTALaRuta)
	ctx.Step(`^La API responde con un mensaje de exito$`, laAPIRespondeConUnMensajeDeExito)
	ctx.Step(`^La API responde con un Status Code (\d+)$`, laAPIRespondeConUnStatusCode)
	ctx.Step(`^el token JWT no es valido$`, elTokenJWTNoEsValido)
	ctx.Step(`^La API responde con un mensaje de error$`, laAPIRespondeConUnMensajeDeError)
	ctx.Step(`^no envia un dato de registro "([^"]*)"$`, noEnviaUnDatoDeRegistro)
	ctx.Step(`^La API responde con un mensaje de error indicando que "([^"]*)"$`, laAPIRespondeConUnMensajeDeErrorIndicandoQue)
}

func TestMain(m *testing.M) {
	setConfigs()

	opts := godog.Options{
		Format: "progress",
		Paths:  []string{"../../features/CreateUser.feature"}, // Se especifica que feature usa este "steptest"
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
