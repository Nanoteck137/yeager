package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin/api"
)

const (
	ErrTypeRouteNotFound  api.ErrorType = "ROUTE_NOT_FOUND"
)

func RouteNotFound() *api.Error {
	return &api.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeRouteNotFound,
		Message: "Route not found",
	}
}
