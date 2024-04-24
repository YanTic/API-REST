package createUser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/go-faker/faker/v4"
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

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Data map[string]User

var (
	url           string
	createUserUrl string
	requestBody   string // Esta variable contiene los datos que se envian a la API
	authToken     string

	apiResponse *APIResponse
)

// go test -v steps/createUser_test.go

func setConfigs() {
	url = JsonReader("../config-file.json", "API.baseUrl")
	createUserUrl = JsonReader("../config-file.json", "API.createUserUrl")
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

// Debido a que las pruebas se ejecutan con contenedores, no se debe poner
// el token manualmente en el archivo config-file.json, por lo que ahora
// se hace una peticion POST para obtener el token.
// func suministraElTokenJWTEnLaCabeceraAuthentication() error {
// 	authToken = JsonReader("../config-file.json", "API.token")
// 	// Aqui hay un problema y es que cuando se hace el POST con resty se tiene que asignar
// 	// el token con SetAuthToken en string y no como JSON, por eso se cambia el
// 	// formato en el bloque de "request-body.json", "token"
// 	// print(authToken)

// 	if authToken == "" {
// 		return fmt.Errorf("el usuario no mandó ningun token")
// 	}

// 	return nil
// }

func suministraElTokenJWTEnLaCabeceraAuthentication() error {
	authToken = getToken()

	if authToken == "" {
		return fmt.Errorf("el usuario no mandó ningun token")
	}

	return nil
}

func elUsuarioHaceLaPeticionPOSTALaRuta(arg1 string) error {
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(authToken).
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

func elTokenJWTNoEsValido() error {
	authToken = JsonReader("../config-file.json", "API.token-novalid")
	return nil
}

func laAPIRespondeConUnMensajeDeError() error {
	if apiResponse.Body == nil || len(apiResponse.Body) == 0 {
		return fmt.Errorf("La API no mandó ningun mensaje de error")
	}

	return nil
}

func noEnviaUnDatoDeRegistro(campo string) error {
	if campo == "usuario" {
		requestBody = JsonReader("request-body.json", "data-nouser")
	} else if campo == "password" {
		requestBody = JsonReader("request-body.json", "data-nopass")
	} else {
		requestBody = JsonReader("request-body.json", "data-noemail")
	}

	return nil
}

func laAPIRespondeConUnMensajeDeErrorIndicandoQue(mensaje string) error {
	if strings.TrimSpace(string(apiResponse.Body)) != strings.TrimSpace(mensaje) {
		return fmt.Errorf("el mensaje de error esperado es '%s', pero se recibió '%s'", mensaje, string(apiResponse.Body))
	}
	return nil
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

	// SE CREAN LOS DATOS PARA LA PRUEBAS
	data := generateData()
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error al convertir los datos a JSON: ", err)
		return
	}

	if err := writeToFile(string(jsonData), "request-body.json"); err != nil {
		fmt.Println("Error al escribir en el archivo: ", err)
		return
	}

	f, err := os.Create("../../reports/reports-json/report_createUser.json")
	if err != nil {
		fmt.Print("error at creating report: ", err)
	}

	// SE EJECUTAN LAS PRUEBAS
	opts := godog.Options{
		// Format: "progress",
		Format: "cucumber",
		Paths:  []string{"../../features/CreateUser.feature"}, // Se especifica que feature usa este "steptest"
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

// Esta funcoin genera cada uno de los casos al crear un usuario, todo en un JSON
func generateData() Data {
	data := make(Data)

	data["data"] = User{
		Username: faker.Username(),
		Password: faker.Password(),
		Email:    faker.Email(),
	}

	data["data-nouser"] = User{
		Password: faker.Password(),
		Email:    faker.Email(),
	}

	data["data-nopass"] = User{
		Username: faker.Username(),
		Email:    faker.Email(),
	}

	data["data-noemail"] = User{
		Username: faker.Username(),
		Password: faker.Password(),
	}

	return data
}

func writeToFile(data, pathOfFile string) error {
	file, err := os.Create(pathOfFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}
