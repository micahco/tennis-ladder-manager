package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/abiosoft/ishell/v2"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	db     *sql.DB
	league *League
	shell  *ishell.Shell
}

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
	defer shell.Close()

	app := &application{
		db:     db,
		league: nil,
		shell:  shell,
	}

	if err := app.run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func (app *application) run(args []string) error {
	leagues, err := app.getAllLeagues()
	if err != nil {
		return err
	}

	// Create a league if one does not exist
	if len(leagues) == 0 {
		app.shell.Println("No leagues found...")

		cont := app.promptConfirm("Would you like to create one?", true)
		if !cont {
			return fmt.Errorf("no league selected")
		}

		var lg League
		for lg.Name == "" {
			app.shell.Print("League name: ")
			lg.Name = app.shell.ReadLine()
		}

		err = app.insertLeague(&lg)
		if err != nil {
			return fmt.Errorf("unable to create league: %w", err)
		}

		leagues = []League{lg}
	}

	// If an argument was provided, try to use that league
	n := 0
	if len(args) > 0 {
		arg := strings.TrimSpace(strings.ToLower(args[0]))
		for i, lg := range leagues {
			if strings.ToLower(lg.Name) == arg {
				n = i
				break
			}
		}
	}
	app.league = &leagues[n]
	app.setPrompt(app.league.Name)

	app.registerCmds()
	app.shell.Run()

	return nil
}

func (app *application) setPrompt(leagueName string) {
	prompt := fmt.Sprintf("(%s)> ", leagueName)
	app.shell.SetPrompt(prompt)
}

func (app *application) registerCmds() {
	app.shell.AddCmd(&ishell.Cmd{
		Name: "sl",
		Help: "select a different league",
		Func: run(app.selectLeague),
	})
	app.shell.AddCmd(&ishell.Cmd{
		Name: "player",
		Help: "player <add,remove> <username>",
		Func: run(app.player),
	})
	app.shell.AddCmd(&ishell.Cmd{
		Name: "match",
		Help: "match <add,remove>",
		Func: run(app.match),
	})
	app.shell.AddCmd(&ishell.Cmd{
		Name: "ladder",
		Help: "ladder <n>",
		Func: run(app.ladder),
	})
}
