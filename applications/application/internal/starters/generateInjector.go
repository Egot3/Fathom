package starters

import (
	"github.com/egot3/fathom/internal/config"
	"github.com/egot3/fathom/internal/database"
	"github.com/egot3/fathom/internal/database/repositories"
	"github.com/egot3/fathom/internal/handler"
	"github.com/egot3/fathom/internal/logging"
	testrunner "github.com/egot3/fathom/internal/testRunner"
	"github.com/egot3/fathom/server"
	"github.com/samber/do/v2"
)

func GenerateInjector() do.Injector {
	i := do.New(
		config.ConfigPackage,
		logging.LogPackage,
		database.DBPackage,
		repositories.RepositoryPackage,
	)

	do.Provide(i, config.Load)

	database.RunMigrations(i)
	registerAll(i)
	initAdmin(i)

	do.Provide(i, testrunner.NewManager)
	do.Provide(i, handler.NewTestService)
	do.Provide(i, server.ChiServer)
	do.Provide(i, newHTTPStarter)

	return i
}
