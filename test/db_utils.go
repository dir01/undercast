package main

import "database/sql"

func getDB(dbURL string) *sql.DB {
	dbURL = dbURL + "?sslmode=disable"
	if db, err := sql.Open("postgres", dbURL); err != nil {
		panic(err)
	} else {
		return db
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM torrents")
	a.DB.Exec("ALTER SEQUENCE torrents_id_seq RESTART WITH 1")
}

func dropTable() {
	a.DB.Exec("DROP TABLE torrents")
}
