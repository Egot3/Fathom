package starters

import (
	"github.com/egot3/fathom/internal/config"
	"github.com/egot3/fathom/internal/database"
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
		server.ServerPackage,
		testrunner.ManagerPackage,
		do.Lazy(newHTTPStarter),
	)

	database.RunMigrations(i)
	registerAll(i)
	initAdmin(i)

	return i
}
