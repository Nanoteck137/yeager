package routes

import (
	"github.com/nanoteck137/yeager/apis"
	"github.com/nanoteck137/yeager/core"
	"github.com/nanoteck137/yeager/types"
	"github.com/nanoteck137/pyrin/api"
)

type Route struct {
	Name       string
	Path       string
	Method     string
	ErrorTypes []api.ErrorType
	Data       any
	Body       types.Body
}

type RouteGroup struct {
	Prefix string
	Routes []Route
}

func NewRouteGroup(prefix string) *RouteGroup {
	return &RouteGroup{
		Prefix: prefix,
		Routes: []Route{},
	}
}

func (r *RouteGroup) AddRoute(name, path, method string, errorTypes []api.ErrorType, data any, body types.Body) {
	r.Routes = append(r.Routes, Route{
		Name:       name,
		Path:       path,
		Method:     method,
		ErrorTypes: errorTypes,
		Data:       data,
		Body:       body,
	})
}

func (r *RouteGroup) Register(handlers ...apis.Handler) {
	for _, h := range handlers {
		r.AddRoute(h.Name, r.Prefix+h.Path, h.Method, h.Errors, h.DataType, h.BodyType)
	}
}

func ServerRoutes(app core.App) []Route {
	g := NewRouteGroup("/api/v1")
	apis.InstallHandlers(app, g)

	return g.Routes
}
