package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/egot3/fathom/internal/config"
	"github.com/egot3/fathom/internal/tui"
	"github.com/egot3/fathom/server"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var verbose bool
var pathToConfig = "../../internal/config/config.yaml"

var rootCmd = &cobra.Command{
	Use:   "fampls",
	Short: "fampls is a toolchain for fathom project",
	Long: `fampls is used for using fathom project with ease.

	It's typically used with a console to provide
	safe fathom project editing.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Using Fathom, run --help for, you guess it(help).")
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start fathom server",
	Long: `Serve starts Fathom server.
	There is really nothing more to it.`,
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		i := do.New()

		cfg := config.Config{}
		data, _ := os.ReadFile(pathToConfig)
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Printf("couldn't read config data: %v", err)
			os.Exit(1)
		}

		do.Provide(i, server.ChiServer)

		if err := http.ListenAndServe(":"+cfg.Server.Port, do.MustInvoke[chi.Router](i)); err != nil {
			log.Printf("Server execution finished: %v", err)
			os.Exit(0)
		}
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Helps configuring fathom project",
	Long: `Starts an interactive TUI for configuring
	fathom project with pizzazis. If config is corrupted
	run "fampls config regenerate"`,
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Config{}
		data, _ := os.ReadFile(pathToConfig)
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Printf("couldn't read config data: %v", err)
			os.Exit(1)
		}

		err := tui.ConfigForm(&cfg)
		if err != nil {
			os.Exit(1)
		}

		out, err := yaml.Marshal(cfg)
		if err != nil {
			os.Exit(1)
		}

		if err := os.WriteFile(pathToConfig, out, 0644); err != nil {
			os.Exit(1)
		}
	},
}

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
		cfg.Server.Logging.Logger = []string{"slog", "charmLog"}
		cfg.Server.Logging.Level = "info"

		out, err := yaml.Marshal(cfg)
		if err != nil {
			os.Exit(1)
		}

		if err := os.WriteFile(pathToConfig, out, 0644); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Use to get debug-level output")

	rootCmd.AddCommand(serveCmd)

	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configRegenerateCmd)
}

func main() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}
