package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/abiosoft/ishell/v2"
	"github.com/pebbe/zmq4"
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

func (app *application) ladder(c *ishell.Context) error {
	query := `SELECT 
    p.name,
    COUNT(m.winner_id) AS matches_won
FROM 
    Players p
LEFT JOIN 
    Matches m ON p.player_id = m.winner_id
GROUP BY 
    p.player_id
ORDER BY 
    matches_won DESC
`

	if len(c.Args) >= 1 {
		n, err := strconv.Atoi(c.Args[0])
		if err == nil && n > 0 {
			query += fmt.Sprintf("LIMIT %d", n)
		}
	}

	rows, err := app.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	app.shell.Println("USERNAME\tWINS")
	i := 1
	for rows.Next() {
		var name string
		var won int
		if err := rows.Scan(&name, &won); err != nil {
			return err
		}
		app.shell.Printf("%d. %s\t%d\n", i, strings.ToLower(name), won)
		i++
	}

	return nil
}

func (app *application) player(c *ishell.Context) error {
	args := c.Args
	if len(c.Args) == 0 {
		return app.addPlayer("")
	}

	args = append(args, "")
	sub := args[0]
	switch sub {
	case "add":
		return app.addPlayer(args[1])
	case "remove":
		return app.removePlayer(args[1])
	default:
		app.shell.Println("Unkown sub-command:", sub)
		return nil
	}
}

func (app *application) addPlayer(username string) error {
	p := Player{Name: username}

	// Initial instruction message
	if p.Name == "" {
		app.shell.Println("Enter player username. This should be a unique name to identify this player/team. Later, you can use this username (case insensitive) to enter match results for this player.")
	}

	for p.Name == "" {
		app.shell.Print("\nPlayer username: ")
		p.Name = app.shell.ReadLine()
	}

	exists, err := app.playerExists(p.Name, app.league.ID)
	if err != nil {
		return err
	}

	if exists {
		app.shell.Println("Cannot add player: username already exists")
		return nil
	}

	err = app.insertPlayer(&p, app.league.ID)
	if err != nil {
		return err
	}

	app.shell.Println("Successfully created player:", strings.ToLower(p.Name))

	return nil
}

func (app *application) removePlayer(username string) error {
	exists, err := app.playerExists(username, app.league.ID)
	if err != nil {
		return err
	}

	if !exists {
		app.shell.Println("Player with username does not exist.")
		return nil
	}

	app.shell.Println("\nAre you sure you want to remove this player from the database?\nDon't worry, this won't remove the matches this player has competed in.")

	cont := app.promptConfirm("Delete player?", false)
	if !cont {
		app.shell.Println("Aborted...")
		return nil
	}

	err = app.deletePlayer(username, app.league.ID)
	if err != nil {
		return err
	}

	app.shell.Println("Successfully removed player:", strings.ToLower(username))

	return nil
}

func (app *application) match(c *ishell.Context) error {
	if len(c.Args) == 0 {
		return app.newMatch()
	}

	sub := c.Args[0]
	if sub == "remove" {
		return app.removeMatch()
	}

	app.shell.Println("Unkown sub-command:", sub)
	return nil
}

func (app *application) newMatch() error {
	app.shell.Println("Please enter each player's username.\n")

	var aName, bName string
	for aName == "" {
		app.shell.Print("Player A username: ")
		aName = app.shell.ReadLine()

		exists, err := app.playerExists(aName, app.league.ID)
		if err != nil {
			return err
		}

		if !exists {
			app.shell.Println("Player with username does not exist.")
			aName = ""
		}
	}

	for bName == "" {
		app.shell.Print("Player B username: ")
		bName = app.shell.ReadLine()

		exists, err := app.playerExists(bName, app.league.ID)
		if err != nil {
			return err
		}

		if !exists {
			app.shell.Println("Player with username does not exist.")
			bName = ""
		}
	}

	pa, err := app.getPlayer(aName, app.league.ID)
	if err != nil {
		return err
	}

	pb, err := app.getPlayer(bName, app.league.ID)
	if err != nil {
		return err
	}

	m := Match{
		PlayerA: pa.ID,
		PlayerB: pb.ID,
	}

	err = app.insertMatch(&m, app.league.ID)
	if err != nil {
		return err
	}

	sets, err := app.enterSetResults(&pa, &pb, m.ID)
	if err != nil {
		return err
	}

	app.shell.Println("\nMatch Results:")

	aWins := 0
	bWins := 0

	app.shell.Printf("\n%s\t\t%s\n", pa.Name, pb.Name)
	for _, s := range sets {
		if s.PlayerAGamesWon > s.PlayerBGamesWon {
			aWins++
		} else {
			bWins++
		}
		app.shell.Printf("%d\t\t%d\n", s.PlayerAGamesWon, s.PlayerBGamesWon)
	}

	var winner *Player
	if aWins > bWins {
		winner = &pa
	} else if bWins > aWins {
		winner = &pb
	}

	app.shell.Println()
	if winner != nil {
		err = app.setMatchWinner(&m, winner.ID, app.league.ID)
		if err != nil {
			return err
		}
		app.shell.Println(winner.Name, "won!")
	} else {
		app.shell.Println("Draw...")
	}

	return nil
}

func (app *application) enterSetResults(pa, pb *Player, matchID int64) ([]*Set, error) {
	sets := make([]*Set, 0)
	for i := 1; true; i++ {
		if i == 1 {
			app.shell.Println("\nEnter set 1 results.")
		} else {
			msg := fmt.Sprintf("\nEnter set %d results?", i)
			cont := app.promptConfirm(msg, true)
			if !cont {
				return sets, nil
			}
		}

		var err error
		s := Set{
			MatchID:         matchID,
			SetNumber:       i,
			PlayerAGamesWon: -1,
			PlayerBGamesWon: -1,
		}
		validated := false
		for !validated {
			for s.PlayerAGamesWon < 0 {
				app.shell.Printf("%s games won: ", pa.Name)
				in := app.shell.ReadLine()
				s.PlayerAGamesWon, err = strconv.Atoi(in)
				if err != nil {
					s.PlayerAGamesWon = -1
				}
			}
			for s.PlayerBGamesWon < 0 {
				app.shell.Printf("%s games won: ", pb.Name)
				in := app.shell.ReadLine()
				s.PlayerBGamesWon, err = strconv.Atoi(in)
				if err != nil {
					s.PlayerBGamesWon = -1
				}
			}
			validated, err = app.validate_set(s.PlayerAGamesWon, s.PlayerBGamesWon)
			if err != nil {
				return nil, err
			}
			if !validated {
				app.shell.Println("\nInvlaid set games. Please enter a valid set.")
				s.PlayerAGamesWon = -1
				s.PlayerBGamesWon = -1
			}
		}

		err = app.insertSet(&s)
		if err != nil {
			return nil, err
		}

		sets = append(sets, &s)
	}

	return sets, nil
}

func (app *application) removeMatch() error {
	m, err := app.getLastMatch(app.league.ID)
	if err != nil {
		return err
	}

	app.shell.Printf("Last match:\t%s vs %s\n", strings.ToUpper(m.PlayerA), strings.ToUpper(m.PlayerB))

	app.shell.Println("\nAre you sure you want to remove this match from the database?\nWARNING: Once deleted the stats of the match cannot be recovered.")

	cont := app.promptConfirm("Delete match?", false)
	if !cont {
		app.shell.Println("Aborted...")
		return nil
	}

	err = app.deleteMatch(m.ID, app.league.ID)
	if err != nil {
		return err
	}

	app.shell.Println("Successfully removed match.")

	return nil
}

func (app *application) artwork(c *ishell.Context) error {
	color := ""
	if len(c.Args) == 1 {
		color = c.Args[0]
	}

	context, err := zmq4.NewContext()
	if err != nil {
		return err
	}
	defer context.Term()

	socket, err := context.NewSocket(zmq4.REQ)
	if err != nil {
		return err
	}
	defer socket.Close()
	err = socket.Connect("tcp://localhost:9875")
	if err != nil {
		return err
	}
	b, err := json.Marshal(color)
	if err != nil {
		return err
	}
	// Send the request to the server
	_, err = socket.SendBytes(b, 0)
	if err != nil {
		return err
	}
	// Wait for a response
	res, err := socket.Recv(0)
	if err != nil {
		return err
	}

	app.shell.Println(res)

	return nil
}

func (app *application) stats(c *ishell.Context) error {
	if len(c.Args) != 1 {
		return errors.New("missing username")
	}

	query := `SELECT
    COALESCE(wins.count, 0) AS wins,
    COALESCE(losses.count, 0) AS losses
FROM 
    Players p
LEFT JOIN (
    SELECT 
        winner_id AS player_id, 
        COUNT(*) AS count
    FROM 
        Matches
    WHERE 
        winner_id IS NOT NULL
    GROUP BY 
        winner_id
) AS wins
ON 
    p.player_id = wins.player_id
LEFT JOIN (
    SELECT 
        CASE 
            WHEN a_player_id = winner_id THEN b_player_id
            WHEN b_player_id = winner_id THEN a_player_id
        END AS player_id,
        COUNT(*) AS count
    FROM 
        Matches
    WHERE 
        winner_id IS NOT NULL
    GROUP BY 
        player_id
) AS losses
ON 
    p.player_id = losses.player_id
WHERE 
    p.league_id = ?
    AND LOWER(p.name) = LOWER(?);`

	username := c.Args[0]
	var wins, losses int
	err := app.db.QueryRow(query, app.league.ID, username).Scan(&wins, &losses)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("player with username does not exist")
		}
		return err
	}

	context, err := zmq4.NewContext()
	if err != nil {
		return err
	}
	defer context.Term()

	socket, err := context.NewSocket(zmq4.REQ)
	if err != nil {
		return err
	}
	defer socket.Close()
	err = socket.Connect("tcp://localhost:9877")
	if err != nil {
		return err
	}
	b, err := json.Marshal([]int{wins, losses})
	if err != nil {
		return err
	}
	// Send the request to the server
	_, err = socket.SendBytes(b, 0)
	if err != nil {
		return err
	}

	// Wait for a response
	res, err := socket.Recv(0)
	if err != nil {
		return err
	}

	app.shell.Println(res)

	return nil
}
