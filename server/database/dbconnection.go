package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func EstablishDBConnection() {
	var err error
	DB, err = sql.Open("mysql", "root:12345@tcp(127.0.0.1:4453)/apirest")
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
