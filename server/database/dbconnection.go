package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func EstablishDBConnection() {
	databaseURL := os.Getenv("DATABASE")
	var err error
	DB, err = sql.Open("mysql", "root:12345@tcp("+databaseURL+")/apirest")
	if err != nil {
		fmt.Println("Conexion no se pudo hacer: ", err)
		return
	}

	err = DB.Ping()
	if err != nil {
		fmt.Println("No se conectó la base de datos: ", err)
		return
	}
	fmt.Println("Conexión a la Base de Datos con exito!")
}
