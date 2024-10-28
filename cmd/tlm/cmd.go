package main

import (
	"errors"

	"github.com/abiosoft/ishell/v2"
)

type cmdFunc func(c *ishell.Context)
type withError func(c *ishell.Context) error

// Simple wrapper to handle errors
func run(fn withError) cmdFunc {
	return func(c *ishell.Context) {
		if err := fn(c); err != nil {
			c.Err(err)
		}
	}
}

func (app *application) selectLeague(c *ishell.Context) error {
	leagues, err := app.getAllLeagues()
	if err != nil {
		return err
	}

	names := leagueNames(leagues)
	choice := c.MultiChoice(names, "Select a league: ")
	if choice == -1 {
		return errors.New("invalid choice")
	}

	app.league = &leagues[choice]
	app.setPrompt(app.league.Name)

	return nil
}
