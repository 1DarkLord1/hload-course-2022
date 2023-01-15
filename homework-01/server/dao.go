package main

import (
	"database/sql"
	"fmt"
)

// TODO: Move database settings into config
const (
	host     = "localhost"
	port     = 5432
	user     = "dword"
	password = "admin"
	dbname   = "postgres"
)


type LongurlTinyurlStorage struct {
	Conn *sql.DB
}

func (storage *LongurlTinyurlStorage) init() {
	fmt.Println(sql.Drivers())

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("postgres", psqlInfo)

	storage.Conn = conn

	if err != nil {
		fmt.Println("Failed to open", err)
		panic("exit")
	}

	err = conn.Ping()

	if err != nil {
		fmt.Println("Failed to ping database", err)
		panic("exit")
	}

	err = createUrlStorageTableDB(storage.Conn)

	if err != nil {
		fmt.Println("Failed create table", err)
		panic("exit")
	}
}

func (storage *LongurlTinyurlStorage) getTinyurl(longurl string) (uint32, error) {
	return getTinyurlDB(storage.Conn, longurl)
}

func (storage *LongurlTinyurlStorage) insertLongurl(longurl string) (uint32, error) {
	return insertLongurlDB(storage.Conn, longurl)
}

func (storage *LongurlTinyurlStorage) getLongurl(tinyurl uint32) (string, error) {
	return getLongurlDB(storage.Conn, tinyurl)
}
