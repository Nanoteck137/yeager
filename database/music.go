package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/yeager/tools/utils"
)

type Artist struct {
	Id      string         `json:"id"`
	Name    string         `json:"name"`
	Picture sql.NullString `json:"picture"`
}

func artistQuery() *goqu.SelectDataset {
	return dialect.From("artists").
		Prepared(true).
		Select("artists.id", "artists.name", "artists.picture")
}

func (db *Database) GetArtistById(ctx context.Context, id string) (Artist, error) {
	ds := artistQuery().
		Where(goqu.I("artists.id").Eq(id))

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Artist{}, err
	}

	var res Artist
	err = row.Scan(&res.Id, &res.Name, &res.Picture)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Artist{}, ErrItemNotFound
		}

		return Artist{}, err
	}

	return res, nil
}

func (db *Database) GetArtistByName(ctx context.Context, name string) (Artist, error) {
	ds := artistQuery().
		Where(goqu.I("artists.name").Eq(name))

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Artist{}, err
	}

	var res Artist
	err = row.Scan(&res.Id, &res.Name, &res.Picture)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Artist{}, ErrItemNotFound
		}

		return Artist{}, err
	}

	return res, nil
}

type CreateArtistParams struct {
	Name    string
	Picture sql.NullString
}

func (db *Database) CreateArtist(ctx context.Context, params CreateArtistParams) (Artist, error) {
	slug := utils.Slug(params.Name)
	id := slug + "-" + utils.CreateShortId()

	ds := dialect.Insert("artists").
		Rows(goqu.Record{
			"id":      id,
			"name":    params.Name,
			"picture": params.Picture,
		}).
		Returning("artists.id", "artists.name", "artists.picture")

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Artist{}, err
	}

	var res Artist
	err = row.Scan(&res.Id, &res.Name, &res.Picture)
	if err != nil {
		return Artist{}, err
	}

	return res, nil
}

type Album struct {
	Id       string
	ArtistId string
	Name     string
	CoverArt sql.NullString
}

type CreateAlbumParams struct {
	ArtistId string
	Name     string
	CoverArt sql.NullString
}

func (db *Database) CreateAlbum(ctx context.Context, params CreateAlbumParams) (Album, error) {
	slug := utils.Slug(params.Name)
	id := slug + "-" + utils.CreateShortId()

	ds := dialect.Insert("albums").
		Rows(goqu.Record{
			"id":        id,
			"artist_id": params.ArtistId,
			"name":      params.Name,
			"cover_art": params.CoverArt,
		}).
		Returning("albums.id", "albums.artist_id", "albums.name", "albums.cover_art")

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Album{}, err
	}

	var res Album
	err = row.Scan(&res.Id, &res.ArtistId, &res.Name, &res.CoverArt)
	if err != nil {
		return Album{}, err
	}

	return res, nil
}

type Track struct {
	Id       string
	AlbumId  string
	ArtistId string

	Title string

	CoverArt sql.NullString

	TrackNumber sql.NullInt64
	Duration    sql.NullInt64
	Year        sql.NullInt64

	OriginalFile   string
	TranscodedFile string
}

type CreateTrackParams struct {
	AlbumId  string
	ArtistId string

	Title string

	CoverArt sql.NullString

	TrackNumber sql.NullInt64
	Duration    sql.NullInt64
	Year        sql.NullInt64

	OriginalFile   string
	TranscodedFile string
}

func (db *Database) CreateTrack(ctx context.Context, params CreateTrackParams) (Track, error) {
	slug := utils.Slug(params.Title)
	id := slug + "-" + utils.CreateShortId()

	ds := dialect.Insert("tracks").
		Rows(goqu.Record{
			"id": id,

			"album_id":  params.AlbumId,
			"artist_id": params.ArtistId,

			"title": params.Title,

			"cover_art": params.CoverArt,

			"track_number": params.TrackNumber,
			"duration":     params.Duration,
			"year":         params.Year,

			"original_file":   params.OriginalFile,
			"transcoded_file": params.TranscodedFile,
		}).
		Returning(
			"tracks.id",
			"tracks.album_id",
			"tracks.artist_id",
			"tracks.title",
			"tracks.cover_art",
			"tracks.track_number",
			"tracks.duration",
			"tracks.year",
			"tracks.original_file",
			"tracks.transcoded_file",
		)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Track{}, err
	}

	var res Track
	err = row.Scan(
		&res.Id,
		&res.AlbumId,
		&res.ArtistId,
		&res.Title,
		&res.CoverArt,
		&res.TrackNumber,
		&res.Duration,
		&res.Year,
		&res.OriginalFile,
		&res.TranscodedFile,
	)
	if err != nil {
		return Track{}, err
	}

	return res, nil
}
