package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func newRepository(db *sql.DB) *repository {
	r := repository{db: sqlx.NewDb(db, "postgres")}
	r.createTables()
	return &r
}

func (r *repository) createTorrent(t *torrent) error {
	dt := dbTorrentFromTorrent(t)
	stmt, err := r.db.PrepareNamed(`INSERT INTO torrents(
			state, name, source, filenames, bytes_completed, bytes_missing
		) VALUES (
			:state, :name, :source, :filenames, :bytes_completed, :bytes_missing
		) RETURNING id`)
	if err != nil {
		return err
	}
	if err = stmt.Get(&t.ID, dt); err != nil {
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
		}

		t.FileNames = append(t.FileNames, f)
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

type dbTorrent struct {
	ID             int    `db:"id"`
	State          string `db:"state"`
	Name           string `db:"name"`
	Source         string `db:"source"`
	FileNames      string `db:"filenames"`
	BytesCompleted int64  `db:"bytes_completed"`
	BytesMissing   int64  `db:"bytes_missing"`
}

func dbTorrentFromTorrent(t *torrent) *dbTorrent {
	return &dbTorrent{
		ID:             t.ID,
		State:          string(t.State),
		Name:           t.Name,
		Source:         t.Source,
		FileNames:      marshalFilenames(t.FileNames),
		BytesCompleted: t.BytesCompleted,
		BytesMissing:   t.BytesMissing,
	}
}

func (dt *dbTorrent) toEntity() *torrent {
	return &torrent{
		ID:             dt.ID,
		State:          state(dt.State),
		Name:           dt.Name,
		Source:         dt.Source,
		FileNames:      unmarshalFilenames(dt.FileNames),
		BytesCompleted: dt.BytesCompleted,
		BytesMissing:   dt.BytesMissing,
	}
}

func marshalFilenames(filenames []string) string {
	if f, err := json.Marshal(filenames); err == nil {
		return string(f)
	} else {
		panic(err)
	}
}

func unmarshalFilenames(fnStr string) (filenames []string) {
	err := json.Unmarshal([]byte(fnStr), &filenames)
	if err != nil {
		panic(err)
	}
	return
}
