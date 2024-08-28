package core

import (
	"github.com/nanoteck137/yeager/config"
	"github.com/nanoteck137/yeager/database"
)

// Inspiration from Pocketbase: https://github.com/pocketbase/pocketbase
// File: https://github.com/pocketbase/pocketbase/blob/master/core/app.go
type App interface {
	DB() *database.Database
	Config() *config.Config

	Bootstrap() error
}
