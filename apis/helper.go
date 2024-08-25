package apis

import (
	"errors"
	"io"

	"github.com/faceair/jio"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/yeager/tools/utils"
	"github.com/nanoteck137/yeager/types"
)

func Body[T types.Body](c echo.Context) (T, error) {
	var res T

	schema := res.Schema()

	j, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return res, err
	}

	if len(j) == 0 {
		// TODO(patrik): Fix error
		return res, errors.New("Invalid body")
	}

	data, err := jio.ValidateJSON(&j, schema)
	if err != nil {
		return res, err
	}

	err = utils.Decode(data, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
