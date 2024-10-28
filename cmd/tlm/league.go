package main

const InitTables = `
CREATE TABLE IF NOT EXISTS Leagues (
	league_id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR
);

CREATE TABLE IF NOT EXISTS Players (
	player_id INTEGER PRIMARY KEY AUTOINCREMENT,
	league_id INTEGER,
	name VARCHAR UNIQUE,
	FOREIGN KEY (league_id) REFERENCES Leagues(league_id)
);

CREATE TABLE IF NOT EXISTS Matches (
	match_id INTEGER PRIMARY KEY AUTOINCREMENT,
	league_id INTEGER,
	winner_id INTEGER,
	a_player_id INTEGER,
	b_player_id INTEGER,
	created_at DATE,
	FOREIGN KEY (winner_id) REFERENCES Leagues(league_id),
	FOREIGN KEY (a_player_id) REFERENCES Players(player_id),
	FOREIGN KEY (a_player_id) REFERENCES Players(player_id),
	FOREIGN KEY (b_player_id) REFERENCES Players(player_id)
);

CREATE TABLE IF NOT EXISTS Sets (
	match_id INTEGER,
	set_number INTEGER,
	a_player_games INTEGER,
	b_player_games INTEGER,
	PRIMARY KEY (match_id, set_number),
	FOREIGN KEY (match_id) REFERENCES Matches(match_id)
);
`

type League struct {
	ID   int64
	Name string
}

func (app *application) getAllLeagues() ([]League, error) {
	rows, err := app.db.Query("SELECT * FROM Leagues;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leagues []League
	for rows.Next() {
		var lg League
		err = rows.Scan(&lg.ID, &lg.Name)
		if err != nil {
			return nil, err
		}

		leagues = append(leagues, lg)
	}

	return leagues, nil
}

func leagueNames(leagues []League) []string {
	names := make([]string, len(leagues))
	for i, lg := range leagues {
		names[i] = lg.Name
	}

	return names
}

func (app *application) insertLeague(lg *League) error {
	query := "INSERT INTO Leagues (name) VALUES (?)"
	res, err := app.db.Exec(query, lg.Name)
	if err != nil {
		return err
	}

	lg.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}
