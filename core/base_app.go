package core

import (
	"os"

	"github.com/nanoteck137/yeager/config"
	"github.com/nanoteck137/yeager/database"
	"github.com/nanoteck137/yeager/types"
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	db       *database.Database
	config   *config.Config
}

func (app *BaseApp) DB() *database.Database {
	return app.db
}

func (app *BaseApp) Config() *config.Config {
	return app.config
}

func (app *BaseApp) WorkDir() types.WorkDir {
	return app.config.WorkDir()
}

func (app *BaseApp) Bootstrap() error {
	var err error

	workDir := app.config.WorkDir()

	err = os.MkdirAll(workDir.OriginalTracksDir(), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(workDir.MobileTracksDir(), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(workDir.ImagesDir(), 0755)
	if err != nil {
		return err
	}

	app.db, err = database.Open(workDir)
	if err != nil {
		return err
	}

	err = database.RunMigrateUp(app.db)
	if err != nil {
		return err
	}

	return nil
}

func NewBaseApp(config *config.Config) *BaseApp {
	return &BaseApp{
		config: config,
	}
}
