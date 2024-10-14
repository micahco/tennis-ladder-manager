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
