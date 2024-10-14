package main

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
