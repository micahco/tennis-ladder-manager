package main

import (
	"database/sql"
	"fmt"

	"github.com/abiosoft/ishell/v2"
)

type application struct {
	db *sql.DB
	lg *League
	sh *ishell.Shell
}

func (app *application) run() error {
	app.sh.Println("Tennis Ladder Manager v0.1")

	err := app.selectLeague()
	if err != nil {
		return err
	}
	app.sh.Println("League:", app.lg.Name)

	app.registerCmds()
	app.sh.Run()

	return nil
}

func (app *application) selectLeague() error {
	leagues, err := app.getAllLeagues()
	if err != nil {
		return err
	}

	switch len(leagues) {
	case 0:
		app.sh.Println("No leagues found...")

		cont := app.promptConfirm("Would you like to create one?", true)
		if !cont {
			return fmt.Errorf("no league selected")
		}

		for app.lg.Name == "" {
			app.sh.Print("League name: ")
			app.lg.Name = app.sh.ReadLine()
		}

		err = app.insertLeague(app.lg)
		if err != nil {
			return fmt.Errorf("unable to create league: %w", err)
		}

	case 1:
		// Automatically select first league
		app.lg = &leagues[0]

	default:
		names := make([]string, len(leagues))
		for i, lg := range leagues {
			names[i] = lg.Name
		}

		choice := app.sh.MultiChoice(names, "Select a league: ")
		app.lg = &leagues[choice]
	}

	return nil
}

func (app *application) registerCmds() {
	app.sh.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "list all leagues",
		Func: app.handleList,
	})
}
