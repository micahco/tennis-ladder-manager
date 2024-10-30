package main

type League struct {
	ID   int64
	Name string
}

func (app *application) getAllLeagues() ([]League, error) {
	rows, err := app.db.Query("SELECT league_id, name FROM Leagues;")
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
