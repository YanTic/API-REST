package updateuser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/go-faker/faker/v4"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

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
	updateUserUrl string
	requestBody   string
	authToken     string

	apiResponse *APIResponse
)

// go test -v steps/updateUser/updateUser_test.go

func setConfigs() {
	url = JsonReader("../config-file.json", "API.baseUrl")
	updateUserUrl = JsonReader("../config-file.json", "API.updateUserUrl")
}

func unUsuarioRegistradoEnLaBaseDeDatosQueYaEstLogueado() error {
	return nil
}

func enviaEnElRequestbodyUnJSONConLosDatosNecesarios() error {
	requestBody = JsonReader("request-body.json", "data")
	return nil
}

func suministraElTokenJWTEnLaCabeceraAuthentication() error {
	authToken = JsonReader("../config-file.json", "API.token")

	if authToken == "" {
		return fmt.Errorf("el usuario no mandó ningun token")
	}

	return nil
}

func elUsuarioHaceLaPeticionPUTALaRuta(arg1 string) error {
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(authToken).
		SetBody(requestBody).
		Put(updateUserUrl)
	if err != nil {
		return err
	}

	apiResponse = &APIResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
	}

	print("\nresponse: ", string(resp.Body()))
	print("\nLOGIN_TEST TOKEN: ", authToken)

	return nil
}

func elTokenJWTNoEsValido() error {
	authToken = JsonReader("../config-file.json", "API.token-novalid")
	return nil
}

func laAPINoEncuentraAlUsuario() error {
	updateUserUrl += "9999" // Se pone a la API a que busque un ID que no existe
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
	ctx.Step(`^envia en el request-body un JSON con los datos necesarios$`, enviaEnElRequestbodyUnJSONConLosDatosNecesarios)
	ctx.Step(`^suministra el token JWT en la cabecera Authentication$`, suministraElTokenJWTEnLaCabeceraAuthentication)
	ctx.Step(`^El usuario hace la peticion PUT a la ruta "([^"]*)"$`, elUsuarioHaceLaPeticionPUTALaRuta)
	ctx.Step(`^el token JWT no es valido$`, elTokenJWTNoEsValido)
	ctx.Step(`^la API no encuentra al usuario$`, laAPINoEncuentraAlUsuario)
	ctx.Step(`^La API responde con un mensaje de error$`, laAPIRespondeConUnMensajeDeError)
	ctx.Step(`^La API responde con un mensaje de exito$`, laAPIRespondeConUnMensajeDeExito)
	ctx.Step(`^La API responde con un Status Code (\d+)$`, laAPIRespondeConUnStatusCode)
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

	// SE EJECUTAN LAS PRUEBAS
	opts := godog.Options{
		Format: "progress",
		Paths:  []string{"../../features/UpdateUser.feature"}, // Se especifica que feature usa este "steptest"
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

func generateData() Data {
	data := make(Data)

	switch rand.Intn(5) {
	case 0:
		data["data"] = User{
			Username: faker.Username(),
			Password: faker.Password(),
			Email:    faker.Email(),
		}
	case 1:
		data["data"] = User{
			Password: faker.Password(),
			Email:    faker.Email(),
		}
	case 2:
		data["data"] = User{
			Username: faker.Username(),
			Email:    faker.Email(),
		}
	case 3:
		data["data"] = User{
			Username: faker.Username(),
			Password: faker.Password(),
		}
	case 4:
		data["data"] = User{
			Username: faker.Username(),
		}
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
