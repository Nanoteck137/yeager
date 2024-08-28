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

var _ types.Body = (*PostAlbumBody)(nil)

type PostAlbumBody struct {
	Name   string               `json:"name"`
	Artist string               `json:"artist"`
}

func (PostAlbumBody) Schema() jio.Schema {
	return jio.Object().Keys(jio.K{
		"name":   jio.String().Min(1).Required(),
		"artist": jio.String().Min(1).Required(),
	})
}

type PostAlbum struct {
	Id string `json:"id"`
}

func (api *musicApi) HandlePostAlbum(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	data := form.Value["data"][0]
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
		vld.Required(&body.Artist).OnError(
			vld.SetField("artist", nil),
		),
	)
	if errs != nil {
		return errs
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
	fmt.Printf("files: %v\n", files)

	workDir := api.app.Config().WorkDir()

	albumDir := path.Join(workDir.AlbumsDir(), album.Id)
	err = os.MkdirAll(albumDir, 0755)
	if err != nil {
		return err
	}

	for _, f := range files {
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

		// TODO(patrik): Add created_date and maybe updated_date
		_, err = db.CreateTrack(ctx, database.CreateTrackParams{
			AlbumId: album.Id,
			ArtistId:       artist.Id,
			Title:          f.Filename,
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
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return c.JSON(200, pyrinapi.SuccessResponse(PostAlbum{
		Id: album.Id,
	}))
}

func InstallMusicHandlers(app core.App, group Group) {
	a := musicApi{app: app}

	group.Register(
		Handler{
			Name:        "ImportAlbum",
			Method:      http.MethodPost,
			Path:        "/music/album",
			DataType:    PostAlbum{},
			BodyType:    PostAlbumBody{},
			IsMultiForm: true,
			HandlerFunc: a.HandlePostAlbum,
			Middlewares: []echo.MiddlewareFunc{},
		},
	)
}
