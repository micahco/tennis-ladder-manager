package main

import (
	"fmt"
	"time"
)

type Match struct {
	ID      int64
	PlayerA int64
	PlayerB int64
	Winner  int64
}

func (app *application) insertMatch(m *Match, leagueID int64) error {
	query := "INSERT INTO Matches (league_id, a_player_id, b_player_id, date_added) VALUES (?, ?, ?, ?)"
	res, err := app.db.Exec(query, leagueID, m.PlayerA, m.PlayerB, time.Now())
	if err != nil {
		return err
	}

	m.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

type MatchInfo struct {
	ID      int64
	PlayerA string
	PlayerB string
	Winner  *string
}

func (app *application) getLastMatch(leagueID int64) (MatchInfo, error) {
	var m MatchInfo

	query := `SELECT 
				m.match_id,
				pa.name,
				pb.name,
				w.name
			FROM 
				Matches m
			JOIN 
				Players pa ON m.a_player_id = pa.player_id
			JOIN 
				Players pb ON m.b_player_id = pb.player_id
			LEFT JOIN 
				Players w ON m.winner_id = w.player_id
			WHERE 
				m.league_id = ?
			ORDER BY 
				m.match_id DESC
			LIMIT 1;`

	err := app.db.QueryRow(query, leagueID).Scan(&m.ID, &m.PlayerA, &m.PlayerB, &m.Winner)
	return m, err
}

func (app *application) setMatchWinner(m *Match, winnerID, leagueID int64) error {
	query := "UPDATE Matches SET winner_id = $1 WHERE match_id = $2 AND league_id = $3"
	res, err := app.db.Exec(query, winnerID, m.ID, leagueID)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("no match found with id %d", m.ID)
	}

	return nil
}

func (app *application) deleteMatch(matchID, leagueID int64) error {
	query := "DELETE FROM Matches WHERE match_id = ? AND league_id = ?"
	_, err := app.db.Exec(query, matchID, leagueID)
	return err
}
