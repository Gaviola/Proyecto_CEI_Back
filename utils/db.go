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
func connect(isFacu bool) *sql.DB {
	var connStr string

	if isFacu {
		connStr = "host=localhost dbname=CEI user=fgaviola password=facu1234 sslmode=disable"
	} else {
		connStr = "host=localhost dbname=cei user=agus password=0811 sslmode=disable"
	}

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
	db := connect(false)
	// TODO modificar la query para que busque por el hash y el username
	//query := "SELECT * FROM users WHERE name = 'facu'"
	query := "SELECT * FROM users WHERE name = '" + user + "' AND hash = '" + string(PassHash) + "'"
	result := db.QueryRow(query).Scan(&findUser.ID, &findUser.Name, &findUser.Lastname, &findUser.Student_id, &findUser.Email, &findUser.Phone, &findUser.Role, &findUser.Dni, &findUser.School, &findUser.Hash)

	if result != nil {
		return findUser, result
	}
	return findUser, nil
}

func DBShowItemTypes() []data.ItemType {
	var itemTypes []data.ItemType

	db := connect(false)
	query := "SELECT * FROM item_type"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var itemType data.ItemType
		err := rows.Scan(&itemType.ID, &itemType.Name, &itemType.IsGeneric)
		if err != nil {
			log.Fatal(err)
		}
		itemTypes = append(itemTypes, itemType)
	}
	return itemTypes
}

func DBShowItems() []data.Item {
	var items []data.Item

	db := connect(false)
	query := "select it.id, it.name, e.code, e.price from element e join item_type it on e.type_id = it.id;"
	
	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var item data.Item
		err := rows.Scan(&item.ID, &item.ItemType, &item.Code, &item.Price)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	return items
}
