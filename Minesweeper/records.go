package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Duration struct {
	Minutes int
	Seconds int
	Raw     time.Duration
}

var dBPath = "assets/records.db"

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dBPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseDB(db *sql.DB) error {
	return db.Close()
}

func DiffStringToId(g *Game, db *sql.DB) (int, error) {
	res, err := db.Query("SELECT id FROM difficulties WHERE name = ?", g.diff)
	if err != nil {
		return 0, err
	}
	defer res.Close()

	var id int
	for res.Next() {
		err = res.Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func SaveBestTime(g *Game, name string, d Duration, db *sql.DB) error {
	diff, err := DiffStringToId(g, db)
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = db.Exec("INSERT INTO best_times (difficulty_id, player_name, time, timestamp) VALUES (?, ?, ?, ?)", diff, name, d.Raw, now)
	if err != nil {
		return err
	}

	return nil
}

func SaveStats(g *Game, db *sql.DB, win bool) error {
	diff, err := DiffStringToId(g, db)
	if err != nil {
		return err
	}

	// Getting current stats
	res, err := db.Query("SELECT wins, played FROM stats WHERE difficulty_id = ?", diff)
	if err != nil {
		return err
	}

	var wins, played int
	for res.Next() {
		err = res.Scan(&wins, &played)
		if err != nil {
			return err
		}
	}

	if win {
		wins++
	}
	played++

	// Updating new stats
	_, err = db.Exec("UPDATE stats SET wins = ?, played = ? WHERE difficulty_id = ?", wins, played, diff)
	if err != nil {
		return err
	}

	return nil

}

func QueryBestTimes(db *sql.DB, diff int) ([]bestTimes, error) {
	res, err := db.Query("SELECT player_name, time, timestamp FROM best_times WHERE difficulty_id = ? ORDER BY time ASC LIMIT 10", diff)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	bts := make([]bestTimes, 10)
	i := 0
	for res.Next() {
		var timeTemp float64
		var timeStampTemp string
		err = res.Scan(&bts[i].PlayerName, &timeTemp, &timeStampTemp)
		if err != nil {
			return nil, err
		}
		bts[i].Time.Raw = time.Duration(timeTemp)
		bts[i].Time.Minutes = int(bts[i].Time.Raw.Minutes())
		bts[i].Time.Seconds = int(bts[i].Time.Raw.Seconds()) % 60
		bts[i].Timestamp, err = time.Parse("2006-01-02 15:04:05.999999999-07:00", timeStampTemp)
		if err != nil {
			return nil, err
		}
		i++
	}

	return bts, nil
}

type bestTimes struct {
	PlayerName string
	Time       Duration
	Timestamp  time.Time
}

func QueryStats(db *sql.DB, diff int) (won, played int, err error) {
	res, err := db.Query("SELECT wins, played FROM stats WHERE difficulty_id = ?", diff)
	if err != nil {
		return 0, 0, err
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&won, &played)
		if err != nil {
			return 0, 0, err
		}
	}

	return
}
