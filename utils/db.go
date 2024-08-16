package utils

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func connect() *sql.DB {
	connStr := "user=postgres..."

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		defer db.Close()
	}

	return db

}

func DBExistUser(PassHash []byte, user string) (bool, error) {
	db := connect()
	result := db.QueryRow("SELECT * FROM users WHERE pass_hash = $1 AND username = $2", PassHash, user).Scan(&PassHash, &user)

	if result != nil {
		return false, result
	}

	return true, nil
}
