package apis

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/faceair/jio"
	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	pyrinapi "github.com/nanoteck137/pyrin/api"
	"github.com/nanoteck137/yeager/core"
	"github.com/nanoteck137/yeager/tools/utils"
	"github.com/nanoteck137/yeager/types"
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
	Title  string
	Number *int
	Genres []string
	Tags   []string
}

func (PostAlbumBodyTrack) Schema() jio.Schema {
	return jio.Object().Keys(jio.K{
		"title":  jio.String().Min(1).Required(),
		"number": jio.Number(),
		"genres": jio.Array().Items(jio.String().Min(1)),
		"tags":   jio.Array().Items(jio.String().Min(1)),
	})
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
	body, err := RawBody[PostAlbumBody](data)
	if err != nil {
		return err
	}

	pretty.Println(body)

	slug := utils.Slug(body.Name)
	fmt.Printf("slug: %v\n", slug)
	artistSlug := utils.Slug(body.Artist)
	fmt.Printf("artistSlug: %v\n", artistSlug)

	// TODO(patrik): Check len(files) and len(body.tracks)
	files := form.File["files"]
	for _, f := range files {
		fmt.Printf("f.Filename: %v\n", f.Filename)
	}

	libraryDir := api.app.Config().LibraryDir

	albumDir := path.Join(libraryDir, slug)
	err = os.Mkdir(albumDir, 0755)
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
