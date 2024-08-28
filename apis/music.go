package apis

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/faceair/jio"
	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	pyrinapi "github.com/nanoteck137/pyrin/api"
	"github.com/nanoteck137/yeager/core"
	"github.com/nanoteck137/yeager/database"
	"github.com/nanoteck137/yeager/types"
	vld "github.com/tiendc/go-validator"
	"golang.org/x/net/context"
)

type Track struct {
	AlbumSlug string
	Slug      string
	Title     string
}

type Album struct {
	Slug string
	Name string

	Tracks []Track
}

type musicApi struct {
	app core.App
}

var _ types.Body = (*PostAlbumBodyTrack)(nil)

type PostAlbumBodyTrack struct {
	Title  string   `json:"title"`
	Number *int     `json:"number"`
	Genres []string `json:"genres"`
	Tags   []string `json:"tags"`
}

func (PostAlbumBodyTrack) Schema() jio.Schema {
	return jio.Object().Keys(jio.K{
		"title":  jio.String().Min(1).Required(),
		"number": jio.Number().Optional(),
		"genres": jio.Array().Items(jio.String().Min(1)),
		"tags":   jio.Array().Items(jio.String().Min(1)),
	})
}

func (d PostAlbumBodyTrack) Validate() vld.Errors {
	errs := vld.Validate(
		vld.Required(&d.Title).OnError(
			vld.SetField("title", nil),
		),
		vld.When(d.Number != nil).Then(
			vld.NumGT(d.Number, 1),
		).OnError(vld.SetField("number", nil)),
	)

	return errs
}

var _ types.Body = (*PostAlbumBody)(nil)

type PostAlbumBody struct {
	Name   string               `json:"name"`
	Artist string               `json:"artist"`
	Tracks []PostAlbumBodyTrack `json:"tracks"`
}

func (PostAlbumBody) Schema() jio.Schema {
	return jio.Object().Keys(jio.K{
		"name":   jio.String().Min(1).Required(),
		"artist": jio.String().Min(1).Required(),
		"tracks": jio.Array().Items(PostAlbumBodyTrack{}.Schema()),
	})
}

func (api *musicApi) HandlePostAlbum(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	data := form.Value["data"][0]
	fmt.Printf("data: %v\n", data)

	var body PostAlbumBody
	err = json.Unmarshal(([]byte)(data), &body)
	if err != nil {
		return err
	}

	pretty.Println(body)

	errs := vld.Validate(
		vld.Required(&body.Name).OnError(
			vld.SetField("name", nil),
		),
	)
	if errs != nil {
		return errs
	}

	for _, track := range body.Tracks {
		errs := track.Validate()
		if errs != nil {
			return errs
		}
	}

	db, tx, err := api.app.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx := context.TODO()

	artist, err := db.GetArtistByName(ctx, body.Artist)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			artist, err = db.CreateArtist(ctx, database.CreateArtistParams{
				Name:    body.Artist,
				Picture: sql.NullString{},
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	pretty.Println(artist)

	album, err := db.CreateAlbum(ctx, database.CreateAlbumParams{
		ArtistId: artist.Id,
		Name:     body.Name,
		CoverArt: sql.NullString{},
	})
	if err != nil {
		return err
	}

	pretty.Println(album)

	files := form.File["files"]

	workDir := api.app.Config().WorkDir()

	albumDir := path.Join(workDir.AlbumsDir(), album.Id)
	err = os.MkdirAll(albumDir, 0755)
	if err != nil {
		return err
	}

	for i, f := range files {
		file, err := os.OpenFile(path.Join(albumDir, f.Filename), os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}

		ff, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(file, ff)
		if err != nil {
			return err
		}

		track := body.Tracks[i]

		t, err := db.CreateTrack(ctx, database.CreateTrackParams{
			AlbumId: album.Id,
			// TODO(patrik): Wrong artist
			ArtistId:       artist.Id,
			Title:          track.Title,
			CoverArt:       sql.NullString{},
			TrackNumber:    sql.NullInt64{},
			Duration:       sql.NullInt64{},
			Year:           sql.NullInt64{},
			OriginalFile:   f.Filename,
			TranscodedFile: "",
		})
		if err != nil {
			return err
		}

		pretty.Println(t)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return c.JSON(200, pyrinapi.SuccessResponse(nil))
}

func InstallMusicHandlers(app core.App, group Group) {
	a := musicApi{app: app}

	group.Register(
		Handler{
			Name:        "ImportAlbum",
			Method:      http.MethodPost,
			Path:        "/music/album",
			DataType:    nil,
			BodyType:    nil,
			HandlerFunc: a.HandlePostAlbum,
			Middlewares: []echo.MiddlewareFunc{},
		},
	)
}
