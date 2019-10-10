package main

import (
	"database/sql"
)

type episode struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Magnet string `json:"magnet"`
	URL    string `json:"url"`
}

func (e *episode) createEpisode(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO episodes(name, magnet, url) VALUES($1, $2, $3) RETURNING id",
		e.Name, e.Magnet, e.URL,
	).Scan(&e.ID)

	if err != nil {
		return err
	}

	return nil
}

func (e *episode) deleteEpisode(db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM episodes WHERE id=$1`, e.ID)
	return err
}

func getEpisodesList(db *sql.DB, start, count int) ([]episode, error) {
	rows, err := db.Query("SELECT id, name, magnet, url FROM episodes LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	episodes := []episode{}
	for rows.Next() {
		var e episode
		if err := rows.Scan(&e.ID, &e.Name, &e.Magnet, &e.URL); err != nil {
			return nil, err
		}
		episodes = append(episodes, e)
	}
	return episodes, nil
}
