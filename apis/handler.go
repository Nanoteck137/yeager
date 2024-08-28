package apis

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/pyrin/api"
	"github.com/nanoteck137/yeager/core"
	"github.com/nanoteck137/yeager/types"
)

type Handler struct {
	Name        string
	Method      string
	Path        string
	DataType    any
	BodyType    types.Body
	IsMultiForm bool
	Errors      []api.ErrorType
	HandlerFunc echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

type Group interface {
	Register(handlers ...Handler)
}

func InstallHandlers(app core.App, g Group) {
	InstallMusicHandlers(app, g)
}
