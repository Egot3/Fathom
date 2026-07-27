package cmd

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/egot3/fathom/internal/config"
	"github.com/egot3/fathom/internal/database"
	"github.com/egot3/fathom/internal/database/repositories"
	"github.com/egot3/fathom/internal/handler"
	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/models"
	testrunner "github.com/egot3/fathom/internal/testRunner"
	"github.com/egot3/fathom/server"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start fathom server",
	Long: `Serve starts Fathom server.
	There is really nothing more to it.`,
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		i := do.New(
			do.Eager(cfg),
			do.Lazy(logging.NewLogger),
			database.DBPackage,
			repositories.RepositoryPackage,
		)

		db := do.MustInvoke[*bun.DB](i)
		if err := database.RunMigrations(context.Background(), db); err != nil {
			log.Fatalf("Fatal migration error: %v", err)
		}
		db.RegisterModel((*models.GroupsUsers)(nil))
		db.RegisterModel((*models.UserGroupsTests)(nil))
		db.RegisterModel((*models.GroupsUsers)(nil))
		db.RegisterModel((*models.TestsQuizzes)(nil))

		do.Provide(i, testrunner.NewTestRunner)

		do.Provide(i, handler.NewTestService)
		do.Provide(i, server.ChiServer)

		log.Printf("running on %v", cfg.ServerPort)
		if err := http.ListenAndServe(":"+cfg.ServerPort, do.MustInvoke[chi.Router](i)); err != nil {
			log.Printf("Server execution finished: %v", err)
			os.Exit(0)
		}
	},
}
