package deleteuser

import (
	"encoding/json"
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
	url           string
	deleteUserUrl string
	authToken     string

	apiResponse *APIResponse
)

// go test -v steps/deleteUser/deleteUser_test.go

func setConfigs() {
	url = JsonReader("../config-file.json", "API.baseUrl")
	deleteUserUrl = JsonReader("../config-file.json", "API.deleteUserUrl")
}

// Función para obtener un jwt Token del servidor
func getToken() string {
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(JsonReader("../login/request-body.json", "data")).
		Post(JsonReader("../config-file.json", "API.loginUserUrl"))
	if err != nil {
		return ""
	}

	var tokenResp map[string]string
	if err := json.Unmarshal(resp.Body(), &tokenResp); err != nil {
		return ""
	}
	token, exists := tokenResp["token"]
	if !exists {
		return ""
	}

	// Retornar solo el token sin el {"token":
	return token
}

func unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado() error {
	return nil
}

func suministraElTokenJWTEnLaCabeceraAuthentication() error {
	authToken = getToken()

	if authToken == "" {
		return fmt.Errorf("el usuario no mandó ningun token")
	}

	return nil
}

func elUsuarioHaceLaPeticionDELETEALaRuta(arg1 string) error {
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(authToken).
		Delete(deleteUserUrl)
	if err != nil {
		return err
	}

	apiResponse = &APIResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
	}

	print(string(resp.Body()))
	return nil
}

func elTokenJWTNoEsValido() error {
	authToken = JsonReader("../config-file.json", "API.token-novalid")
	return nil
}

func laAPIRespondeConUnStatusCode(codigo int) error {
	if apiResponse.StatusCode != codigo {
		return fmt.Errorf("el código de estado esperado es %d, pero se recibió %d", codigo, apiResponse.StatusCode)
	}
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

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^Un usuario registrado en la Base de Datos que ya está logueado$`, unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado)
	ctx.Step(`^suministra el token JWT en la cabecera Authentication$`, suministraElTokenJWTEnLaCabeceraAuthentication)
	ctx.Step(`^El usuario hace la peticion DELETE a la ruta "([^"]*)"$`, elUsuarioHaceLaPeticionDELETEALaRuta)
	ctx.Step(`^el token JWT no es valido$`, elTokenJWTNoEsValido)
	ctx.Step(`^La API responde con un Status Code (\d+)$`, laAPIRespondeConUnStatusCode)
	ctx.Step(`^La API responde con un mensaje de exito$`, laAPIRespondeConUnMensajeDeExito)
	ctx.Step(`^La API responde con un mensaje de error$`, laAPIRespondeConUnMensajeDeError)
}

func TestMain(m *testing.M) {
	setConfigs()

	f, err := os.Create("../../reports/reports-json/report_deleteUser.json")
	if err != nil {
		fmt.Print("error at creating report: ", err)
	}

	opts := godog.Options{
		Format: "cucumber",
		Paths:  []string{"../../features/DeleteUser.feature"}, // Se especifica que feature usa este "steptest"
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
