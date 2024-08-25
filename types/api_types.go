package types

import "github.com/faceair/jio"

type Body interface {
	Schema() jio.Schema
}

