package main

import (
	"database/sql"
	"errors"
)

type torrent struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Magnet string `json:"magnet"`
	URL    string `json:"url"`
}

func (t *torrent) createTorrent(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO torrents(name, magnet, url) VALUES($1, $2, $3) RETURNING id",
		t.Name, t.Magnet, t.URL,
	).Scan(&t.ID)

	if err != nil {
		return err
	}

	return nil
}

func (t *torrent) deleteTorrent(db *sql.DB) error {
	if res, err := db.Exec(`DELETE FROM torrents WHERE id=$1`, t.ID); err != nil {
		return err
	} else if count, err := res.RowsAffected(); err != nil {
		return err
	} else if count == 0 {
		return errors.New("Not found")
	} else {
		return nil
	}
}

func getTorrentsList(db *sql.DB, start, count int) ([]torrent, error) {
	rows, err := db.Query("SELECT id, name, magnet, url FROM torrents LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	torrents := []torrent{}
	for rows.Next() {
		var t torrent
		if err := rows.Scan(&t.ID, &t.Name, &t.Magnet, &t.URL); err != nil {
			return nil, err
		}
		torrents = append(torrents, t)
	}
	return torrents, nil
}
