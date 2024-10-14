package main

import (
	"database/sql"
	"log"

	"github.com/abiosoft/ishell/v2"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.SetFlags(0)

	// Open database
	db, err := openDB("tlm.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Interactive shell
	shell := ishell.New()
	shell.EOF(func(c *ishell.Context) {
		c.Stop()
	})
	shell.Interrupt(func(c *ishell.Context, count int, input string) {
		c.Stop()
	})

	app := &application{
		db: db,
		lg: &League{},
		sh: shell,
	}
	defer app.sh.Close()

	if err := app.run(); err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(InitTables)
	if err != nil {
		return nil, err
	}

	return db, nil
}
