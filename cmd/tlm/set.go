package main

type Set struct {
	MatchID         int64
	SetNumber       int
	PlayerAGamesWon int
	PlayerBGamesWon int
}

func (app *application) insertSet(s *Set) error {
	query := "INSERT INTO Sets (match_id, set_number, a_player_games, b_player_games) VALUES (?, ?, ?, ?)"
	_, err := app.db.Exec(query, s.MatchID, s.SetNumber, s.PlayerAGamesWon, s.PlayerBGamesWon)
	if err != nil {
		return err
	}

	return nil
}
