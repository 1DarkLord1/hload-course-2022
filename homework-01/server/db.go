package main

import 	"database/sql"

func createUrlStorageTable(conn *sql.DB) error {
	_, err := conn.Exec("create table if not exists url_storage ( longurl varchar primary key, tinyurl serial )")

	return err
}

func getTinyurl(conn *sql.DB, longurl string) (uint32, error) {
	stmt, err := conn.Prepare("select ( tinyurl ) from url_storage where longurl = $1")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var tinyurl uint32

	err = stmt.QueryRow(longurl).Scan(&tinyurl)

	return tinyurl, err
}

func getLongurl(conn *sql.DB, tinyurl uint32) (string, error) {
	stmt, err := conn.Prepare("select ( longurl ) from url_storage where tinyurl = $1")

	if err != nil {
		return "", err
	}

	defer stmt.Close()

	var longurl string

	err = stmt.QueryRow(tinyurl).Scan(&longurl)

	return longurl, err
}

func insertLongurl(conn *sql.DB, longurl string) error {
	stmt, err := conn.Prepare("insert into url_storage ( longurl ) values ( $1 ) on conflict do nothing")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(longurl)

	return err
}