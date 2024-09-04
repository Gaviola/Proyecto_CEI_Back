package utils

import (
	"database/sql"
	"log"

	"github.com/Gaviola/Proyecto_CEI_Back.git/data"

	_ "github.com/lib/pq"
)

func connect() *sql.DB {
	connStr := "host=localhost dbname=CEI user=fgaviola password=facu1234 sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		defer db.Close()
	}
	db.Ping()

	return db

}

func DBExistUser(PassHash []byte, user string) (bool, error) {
	prueba := data.User{}
	db := connect()
	query := "SELECT * FROM users WHERE name = 'facu'"
	result := db.QueryRow(query).Scan(&prueba.ID, &prueba.Name, &prueba.Lastname, &prueba.Student_id, &prueba.Email, &prueba.Phone, &prueba.Role, &prueba.Dni, &prueba.School, &prueba.Hash, &prueba.Salt)

	if result != nil {
		return false, result
	}

	return true, nil
}
