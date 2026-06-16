package cmd

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/egot3/fathom/internal/config"
	"github.com/egot3/fathom/internal/database"
	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/models"
	"github.com/egot3/fathom/server"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"go.yaml.in/yaml/v4"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start fathom server",
	Long: `Serve starts Fathom server.
	There is really nothing more to it.`,
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		i := do.New()

		cfg := config.Config{}
		data, _ := os.ReadFile(PathToConfig)
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Printf("couldn't read config data: %v", err)
			os.Exit(1)
		}

		do.ProvideValue(i, cfg)

		do.Provide(i, database.InitDB)
		db := do.MustInvoke[*bun.DB](i)
		if err := database.RunMigrations(context.Background(), db); err != nil {
			log.Fatalf("Fatal migration error: %v", err)
		}
		db.RegisterModel((*models.GroupsUsers)(nil))
		db.RegisterModel((*models.UserGroupsTests)(nil))
		db.RegisterModel((*models.GroupsUsers)(nil))

		do.Provide(i, logging.NewLogger)
		do.Provide(i, server.ChiServer)

		log.Printf("running on %v", cfg.Server.Port)
		if err := http.ListenAndServe(":"+cfg.Server.Port, do.MustInvoke[chi.Router](i)); err != nil {
			log.Printf("Server execution finished: %v", err)
			os.Exit(0)
		}
	},
}
