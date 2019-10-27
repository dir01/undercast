package server

import (
	"database/sql"
	"errors"
)

type repository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) *repository {
	r := repository{db: db}
	return &r
}

func (r *repository) createTorrent(t *torrent) error {
	err := r.db.QueryRow(
		"INSERT INTO torrents(name, magnet, url) VALUES($1, $2, $3) RETURNING id",
		t.Name, t.Magnet, t.URL,
	).Scan(&t.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) getTorrentsList(start, count int) ([]torrent, error) {
	rows, err := r.db.Query("SELECT id, name, magnet, url FROM torrents LIMIT $1 OFFSET $2", count, start)
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

func (r *repository) deleteTorrent(id int) error {
	if res, err := r.db.Exec(`DELETE FROM torrents WHERE id=$1`, id); err != nil {
		return err
	} else if count, err := res.RowsAffected(); err != nil {
		return err
	} else if count == 0 {
		return errors.New("Not found")
	} else {
		return nil
	}
}
