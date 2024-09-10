package utils

import (
	"database/sql"
	"log"

	"github.com/Gaviola/Proyecto_CEI_Back.git/data"

	_ "github.com/lib/pq"
)

// connect
/*
Conecta a la base de datos y devuelve un puntero a la conexion.
Devuelve nil si hay un error en la conexion.
*/
func connect() *sql.DB {
	connStr := "host=localhost dbname=CEI user=fgaviola password=facu1234 sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		defer db.Close()
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db

}

// DBExistUser
/*
Busca un usuario en la base de datos segun el hash de la contraseña y el username.
Devuelve el usuario correspondiente si el usuario existe.
Devuelve un usuario vacio si el usuario no existe o si hay un error en la base de datos.
*/
func DBExistUser(PassHash []byte, user string) (data.User, error) {
	findUser := data.User{}
	db := connect()
	// TODO modificar la query para que busque por el hash y el username
	//query := "SELECT * FROM users WHERE name = 'facu'"
	query := "SELECT * FROM users WHERE name = '" + user + "' AND hash = '" + string(PassHash) + "'"
	result := db.QueryRow(query).Scan(&findUser.ID, &findUser.Name, &findUser.Lastname, &findUser.Student_id, &findUser.Email, &findUser.Phone, &findUser.Role, &findUser.Dni, &findUser.School, &findUser.Hash)

	if result != nil {
		return findUser, result
	}
	return findUser, nil
}
