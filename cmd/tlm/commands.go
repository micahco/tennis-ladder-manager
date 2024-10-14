package main

import (
	"github.com/abiosoft/ishell/v2"
)

func (app *application) handleList(c *ishell.Context) {
	leagues, err := app.getAllLeagues()
	if err != nil {
		c.Err(err)
		return
	}

	for _, lg := range leagues {
		app.sh.Println(lg.Name)
	}
}
