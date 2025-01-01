package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CheckCaptureCache() (string, string) {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Println("Failed to Open DB file")
		panic(err)
	}
	defer db.Close()

	createKeyCmd := `
		create table if not exists smash (id integer not null primary key, api);
	`
	_, err = db.Exec(createKeyCmd)
	if err != nil {
		log.Printf("%q: %s\n", err, createKeyCmd)
		panic(err)
	}

	createTournamentCmd := `
		create table if not exists tournament (id integer not null primary key, url);
	`
	_, err = db.Exec(createTournamentCmd)
	if err != nil {
		log.Printf("%q: %s\n", err, createKeyCmd)
		panic(err)
	}

	checkKeyCmd := `
		select api from smash where id = 1;
	`

	var key string
	err = db.QueryRow(checkKeyCmd).Scan(&key)
	if err != nil {
		if err == sql.ErrNoRows {
			key = ""
		} else {
			log.Printf("%q: %s\n", err, checkKeyCmd)
			panic(err)
		}
	}

	checkUrlCmd := `
		select url from tournament where id = 1;
	`

	var tournament string
	err = db.QueryRow(checkUrlCmd).Scan(&tournament)
	if err != nil {
		if err == sql.ErrNoRows {
			tournament = ""
		} else {
			log.Printf("%q: %s\n", err, checkUrlCmd)
			panic(err)
		}
	}

	return key, tournament
}

func WriteCapture(key string, url string) {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Println("Failed to Open DB file")
		panic(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println("Failed to begin insert transaction")
		panic(err)
	}
	stmt, err := tx.Prepare("insert or replace into smash (id, api) values(?,?)")
	if err != nil {
		log.Println("Failed to prepare insert transaction")
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(1, key)
	if err != nil {
		log.Println("Failed to insert key")
		panic(err)
	}

	stmt, err = tx.Prepare("insert or replace into tournament (id, url) values(?,?)")
	if err != nil {
		log.Println("Failed to prepare insert url transaction")
		panic(err)
	}

	_, err = stmt.Exec(1, url)
	if err != nil {
		log.Println("Failed to insert url")
	}
	err = tx.Commit()
	if err != nil {
		log.Println("Failed to commit insert transaction")
		panic(err)
	}
}
