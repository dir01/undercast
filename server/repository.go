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

func (r *repository) CreateTorrent(t *Torrent) error {
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

func (r *repository) getUnfinisedTorrents() ([]Torrent, error) {
	return r.queryToTorrents("SELECT * from torrents WHERE state='ADDED'", nil)
}

func (r *repository) getTorrentsList(start, count int) ([]Torrent, error) {
	args := map[string]interface{}{
		"limit":  count,
		"offset": start,
	}
	return r.queryToTorrents("SELECT * FROM torrents LIMIT :limit OFFSET :offset", args)
}

func (r *repository) queryToTorrents(query string, args interface{}) ([]Torrent, error) {
	if args == nil {
		args = struct{}{}
	}
	stmt, _ := r.db.PrepareNamed(query)
	defer stmt.Close()

	dbTorList := []dbTorrent{}
	err := stmt.Select(&dbTorList, args)
	if err != nil {
		return nil, err
	}

	result := []Torrent{}
	for _, dt := range dbTorList {
		result = append(result, dt.toEntity())
	}
	return result, nil

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
	ID             int            `db:"id"`
	State          string         `db:"state"`
	Source         string         `db:"source"`
	Name           sql.NullString `db:"name"`
	FileNames      sql.NullString `db:"filenames"`
	BytesCompleted sql.NullInt64  `db:"bytes_completed"`
	BytesMissing   sql.NullInt64  `db:"bytes_missing"`
}

func dbTorrentFromTorrent(t *Torrent) *dbTorrent {
	return &dbTorrent{
		ID:             t.ID,
		State:          string(t.State),
		Name:           sql.NullString{String: t.Name},
		Source:         t.Source,
		FileNames:      sql.NullString{String: marshalFilenames(t.FileNames)},
		BytesCompleted: sql.NullInt64{Int64: t.BytesCompleted},
		BytesMissing:   sql.NullInt64{Int64: t.BytesMissing},
	}
}

func (dt *dbTorrent) toEntity() Torrent {
	return Torrent{
		ID:             dt.ID,
		State:          state(dt.State),
		Name:           dt.Name.String,
		Source:         dt.Source,
		FileNames:      unmarshalFilenames(dt.FileNames.String),
		BytesCompleted: dt.BytesCompleted.Int64,
		BytesMissing:   dt.BytesMissing.Int64,
	}
}

func marshalFilenames(filenames []string) string {
	if f, err := json.Marshal(filenames); err != nil {
		panic(err)
	} else {
		return string(f)
	}
}

func unmarshalFilenames(fnStr string) (filenames []string) {
	if fnStr == "" {
		return nil
	}
	err := json.Unmarshal([]byte(fnStr), &filenames)
	if err != nil {
		panic(err)
	}
	return
}
