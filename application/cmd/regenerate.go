package cmd

import (
	"os"

	"github.com/egot3/fathom/internal/config"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

var configRegenerateCmd = &cobra.Command{
	Use:   "regenerate",
	Short: "Regenerates broken config",
	Long: `sets everything do default values and overwrites
		everything in this file, so be aware`,
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Config{}
		cfg.Database.Sqlite.Used = true
		cfg.Database.Sqlite.Path = "../data/fathom.db"
		cfg.Database.Postgres.Used = false

		cfg.Server.Port = "8081"
		cfg.Server.Logging.Loggers = []string{"slog", "charmLog"}
		cfg.Server.Logging.Level = "info"

		out, err := yaml.Marshal(cfg)
		if err != nil {
			os.Exit(1)
		}

		if err := os.WriteFile(PathToConfig, out, 0644); err != nil {
			os.Exit(1)
		}
	},
}
