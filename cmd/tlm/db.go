package main

import "database/sql"

const InitTables = `
CREATE TABLE IF NOT EXISTS Leagues (
	league_id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Players (
	player_id INTEGER PRIMARY KEY AUTOINCREMENT,
	league_id INTEGER NOT NULL,
	name TEXT NOT NULL UNIQUE COLLATE NOCASE,
	FOREIGN KEY (league_id) REFERENCES Leagues(league_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS unique_player_name_in_league ON players (name, league_id);

CREATE TABLE IF NOT EXISTS Matches (
	match_id INTEGER PRIMARY KEY AUTOINCREMENT,
	league_id INTEGER NOT NULL,
	a_player_id INTEGER NOT NULL,
	b_player_id INTEGER NOT NULL,
	winner_id INTEGER,
	date_added DATE NOT NULL,
	FOREIGN KEY (winner_id) REFERENCES Leagues(league_id),
	FOREIGN KEY (a_player_id) REFERENCES Players(player_id),
	FOREIGN KEY (b_player_id) REFERENCES Players(player_id)
);

CREATE TABLE IF NOT EXISTS Sets (
	match_id INTEGER NOT NULL,
	set_number INTEGER NOT NULL,
	a_player_games INTEGER NOT NULL,
	b_player_games INTEGER NOT NULL,
	PRIMARY KEY (match_id, set_number),
	FOREIGN KEY (match_id) REFERENCES Matches(match_id)
);
`

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(InitTables)
	if err != nil {
		return nil, err
	}

	return db, nil
}
