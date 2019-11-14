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

func (r *repository) SaveTorrent(t *Torrent) error {
	if t.ID == 0 {
		return r.insertTorrent(t)
	}
	return r.updateTorrent(t)
}

func (r *repository) insertTorrent(t *Torrent) error {
	dt := dbTorrentFromTorrent(t)
	stmt, err := r.db.PrepareNamed(`INSERT INTO torrents(
			state, name, source, filepaths, bytes_completed, bytes_missing
		) VALUES (
			:state, :name, :source, :filepaths, :bytes_completed, :bytes_missing
		) RETURNING id`)
	if err != nil {
		return err
	}
	if err = stmt.Get(&t.ID, dt); err != nil {
		return err
	}
	return nil
}

func (r *repository) updateTorrent(t *Torrent) error {
	dt := dbTorrentFromTorrent(t)
	if _, err := r.db.NamedExec(`UPDATE torrents SET 
		state=:state,
		name=:name,
		source=:source,
		filepaths=:filepaths,
		bytes_completed=:bytes_completed,
		bytes_missing=:bytes_missing
	WHERE id=:id`, dt); err != nil {
		return err
	}
	if len(t.Episodes) == 0 {
		return nil
	}
	for _, ep := range t.Episodes {
		if ep.ID != 0 {
			continue
		}
		dEp := dbEpisodeFromEpisode(&ep, t.ID)
		stmt, err := r.db.PrepareNamed(`INSERT INTO episodes (
			torrent_id, name, filepaths
		) VALUES (
			:torrent_id, :name, :filepaths
		) RETURNING id`)
		if err != nil {
			return err
		}
		if err = stmt.Get(&ep.ID, dEp); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) GetTorrent(id int) (*Torrent, error) {
	args := map[string]interface{}{
		"id": id,
	}
	torrents, err := r.queryToTorrents("SELECT * FROM torrents WHERE id=:id", args)
	if err != nil || len(torrents) == 0 {
		return nil, err
	}
	return &torrents[0], nil
}

func (r *repository) getDownloadingTorrents() ([]Torrent, error) {
	return r.queryToTorrents("SELECT * from torrents WHERE state='DOWNLOADING'", nil)
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

	episodesMap, err := r.getEpisodesMap(dbTorList)
	if err != nil {
		return nil, err
	}

	result := []Torrent{}
	for _, dt := range dbTorList {
		t := dt.toEntity()
		t.Episodes = episodesMap[t.ID]
		result = append(result, t)
	}
	return result, nil

}

func (r *repository) getEpisodesMap(dbTorList []dbTorrent) (map[int][]Episode, error) {
	result := make(map[int][]Episode)
	if len(dbTorList) == 0 {
		return result, nil
	}

	ids := []int{}
	for _, tor := range dbTorList {
		ids = append(ids, tor.ID)
	}

	query, args, err := sqlx.In("SELECT * FROM episodes WHERE torrent_id IN (?) ORDER BY id ASC", ids)
	if err != nil {
		return nil, err
	}

	dbEpisodesList := []dbEpisode{}
	err = r.db.Select(&dbEpisodesList, r.db.Rebind(query), args...)

	if err != nil {
		return nil, err
	}
	for _, e := range dbEpisodesList {
		if _, ok := result[e.TorrentID]; !ok {
			result[e.TorrentID] = []Episode{}
		}
		result[e.TorrentID] = append(result[e.TorrentID], e.toEntity())
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
	tableCreationQueries := []string{`
	CREATE TABLE IF NOT EXISTS torrents(
		id SERIAL PRIMARY KEY,
		state VARCHAR(50),
		source TEXT NOT NULL,
		name VARCHAR(500),
		filepaths JSON,
		bytes_completed BIGINT,
		bytes_missing BIGINT
		CONSTRAINT require_source CHECK (
			(case when source is null or length(source) = 0 then FALSE else TRUE end)
		)
	)`,
		`CREATE TABLE IF NOT EXISTS episodes(
		id SERIAL PRIMARY KEY,
		torrent_id INT NOT NULL,
		name TEXT NOT NULL,
		media_url TEXT,
		filepaths JSON
	)`}
	for _, query := range tableCreationQueries {
		if _, err := r.db.Exec(query); err != nil {
			log.Fatal(err)
		}
	}
}

type dbTorrent struct {
	ID             int            `db:"id"`
	State          string         `db:"state"`
	Source         string         `db:"source"`
	Name           sql.NullString `db:"name"`
	FilePaths      sql.NullString `db:"filepaths"`
	BytesCompleted sql.NullInt64  `db:"bytes_completed"`
	BytesMissing   sql.NullInt64  `db:"bytes_missing"`
}

func dbTorrentFromTorrent(t *Torrent) *dbTorrent {
	return &dbTorrent{
		ID:    t.ID,
		State: string(t.State),
		Name: sql.NullString{
			String: t.Name,
			Valid:  true,
		},
		Source: t.Source,
		FilePaths: sql.NullString{
			String: marshalFilepaths(t.FilePaths),
			Valid:  true,
		},
		BytesCompleted: sql.NullInt64{
			Int64: t.BytesCompleted,
			Valid: true,
		},
		BytesMissing: sql.NullInt64{
			Int64: t.BytesMissing,
			Valid: true,
		},
	}
}

func (dt *dbTorrent) toEntity() Torrent {
	return Torrent{
		ID:             dt.ID,
		State:          state(dt.State),
		Name:           dt.Name.String,
		Source:         dt.Source,
		FilePaths:      unmarshalFilepaths(dt.FilePaths.String),
		BytesCompleted: dt.BytesCompleted.Int64,
		BytesMissing:   dt.BytesMissing.Int64,
	}
}

type dbEpisode struct {
	ID        int            `db:"id"`
	TorrentID int            `db:"torrent_id"`
	Name      string         `db:"name"`
	FilePaths string         `db:"filepaths"`
	MediaURL  sql.NullString `db:"media_url"`
}

func dbEpisodeFromEpisode(episode *Episode, torrentID int) *dbEpisode {
	return &dbEpisode{
		TorrentID: torrentID,
		Name:      episode.Name,
		FilePaths: marshalFilepaths(episode.FilePaths),
		MediaURL: sql.NullString{
			String: episode.MediaURL,
			Valid:  true,
		},
	}
}

func (d *dbEpisode) toEntity() Episode {
	return Episode{
		ID:        d.ID,
		Name:      d.Name,
		FilePaths: unmarshalFilepaths(d.FilePaths),
		MediaURL:  d.MediaURL.String,
	}
}

func marshalFilepaths(filepaths []string) string {
	if f, err := json.Marshal(filepaths); err != nil {
		panic(err)
	} else {
		return string(f)
	}
}

func unmarshalFilepaths(fnStr string) (filepaths []string) {
	if fnStr == "" {
		return nil
	}
	err := json.Unmarshal([]byte(fnStr), &filepaths)
	if err != nil {
		panic(err)
	}
	return
}
