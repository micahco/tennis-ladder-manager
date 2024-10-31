package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/abiosoft/ishell/v2"
	"github.com/fatih/color"
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

type ByName []*ishell.Cmd

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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

	// Welcome message
	app.registerCmds()
	style := color.New(color.Underline, color.FgGreen).SprintFunc()
	app.shell.Printf("\n%s - add players, record matches, and view rankings.\n", style("Tennis League Manager"))
	app.shell.Println("\nCOMMANDS:")
	cmds := app.shell.Cmds()
	sort.Sort(ByName(cmds))
	for _, cmd := range cmds {
		if s := cmd.Help; s != "" {
			app.shell.Printf("\t%s\t\t%s\n", cmd.Name, strings.ToUpper(string(s[0]))+s[1:])
		}
	}
	app.shell.Println("\nUSAGE:")
	for _, cmd := range cmds {
		if cmd.LongHelp != "" {
			app.shell.Printf("\t%s\n", cmd.LongHelp)
		}
	}
	app.shell.Println()

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

	app.shell.Run()

	return nil
}

func (app *application) setPrompt(leagueName string) {
	style := color.New(color.Bold, color.FgGreen).SprintFunc()
	prompt := fmt.Sprintf("(%s)> ", leagueName)
	app.shell.SetPrompt(style(prompt))
}

func (app *application) registerCmds() {
	app.shell.AddCmd(&ishell.Cmd{
		Name:     "select",
		Help:     "choose a different league",
		LongHelp: "select <league_name>",
		Func:     run(app.selectLeague),
	})
	app.shell.AddCmd(&ishell.Cmd{
		Name:     "player",
		Help:     "manage players",
		LongHelp: "player <add,remove> <username>",
		Func:     run(app.player),
	})
	app.shell.AddCmd(&ishell.Cmd{
		Name:     "match",
		Help:     "record match results",
		LongHelp: "match <add,remove>",
		Func:     run(app.match),
	})
	app.shell.AddCmd(&ishell.Cmd{
		Name:     "ladder",
		Help:     "view the current league rankings",
		LongHelp: "ladder <number>",
		Func:     run(app.ladder),
	})
}
