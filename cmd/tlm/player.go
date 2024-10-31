package main

import "strings"

type Player struct {
	ID   int64
	Name string
}

func (app *application) playerExists(username string, leagueID int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM Players WHERE name = ? AND league_id = ?)"
	err := app.db.QueryRow(query, strings.ToLower(username), leagueID).Scan(&exists)

	return exists, err
}

func (app *application) getPlayer(username string, leagueID int64) (Player, error) {
	var p Player
	query := "SELECT player_id, name FROM Players WHERE name = ? AND league_id = ?"
	err := app.db.QueryRow(query, username, leagueID).Scan(&p.ID, &p.Name)

	return p, err
}

func (app *application) insertPlayer(p *Player, leagueID int64) error {
	query := "INSERT INTO Players (league_id, name) VALUES (?, ?)"
	res, err := app.db.Exec(query, leagueID, p.Name)
	if err != nil {
		return err
	}

	p.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (app *application) deletePlayer(name string, leagueID int64) error {
	query := "DELETE FROM Players WHERE name = ? AND league_id = ?"
	_, err := app.db.Exec(query, name, leagueID)

	return err
}
