package server

import (
	"database/sql"
	"errors"
	"log"
)

type repository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) *repository {
	r := repository{db: db}
	r.createTables()
	return &r
}

func (r *repository) createTorrent(t *torrent) error {
	err := r.db.QueryRow(
		"INSERT INTO torrents(source) VALUES($1) RETURNING id",
		t.Source).Scan(&t.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) getTorrentsList(start, count int) ([]torrent, error) {
	rows, err := r.db.Query("SELECT "+
		"id, state, source, filenames, bytes_completed, bytes_missing "+
		"FROM torrents LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	torrents := []torrent{}
	for rows.Next() {
		var t torrent
		var f string
		if err := rows.Scan(&t.ID, &t.State, &t.Source, &f, &t.BytesCompleted, &t.BytesMissing); err != nil {
			return nil, err
		} else {
			t.FileNames = append(t.FileNames, f)
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

func (r *repository) createTables() {
	const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS torrents (
	id SERIAL PRIMARY KEY,
	state VARCHAR(50),
	source TEXT NOT NULL,
	name VARCHAR(500),
	filenames JSON,
	bytes_completed BIGINT,
	bytes_missing BIGINT
    CONSTRAINT require_source CHECK (
		(case when source is null or length(source) = 0 then FALSE else TRUE end)
    )

)`
	if _, err := r.db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}
