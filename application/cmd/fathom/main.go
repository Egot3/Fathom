package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/huh"
	configutils "github.com/egot3/fathom/internal/configUtils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var verbose bool

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
		// running the server

		fmt.Println("running")
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Helps configuring fathom project",
	Long: `Starts an interactive TUI for configuring
	fathom project with pizzazis`,
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := configutils.Config{}
		pathToConfig := "../../internal/config/config.yaml"
		data, _ := os.ReadFile(pathToConfig)
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Printf("couldn't read config data: %v", err)
			os.Exit(1)
		}

		database := ""
		form := huh.NewForm(

			huh.NewGroup(huh.NewNote().
				Title("Configuration").
				Description("Welcome to Configuration tool.").
				Next(true).
				NextLabel("Next"),
			),

			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Choose db").
					Options(
						huh.NewOption("Postgres(heavy)", "postgres"),
						huh.NewOption("Sqlite(light)", "sqlite"),
					).Value(&database),
			),

			huh.NewGroup(
				huh.NewInput().
					Title("Path to db").
					Value(&cfg.Database.Sqlite.Path),
			).WithHideFunc(func() bool {
				return database != "sqlite"
			}),
			huh.NewGroup(
				huh.NewInput().
					Title("user").Value(&cfg.Database.Postgres.User),
				huh.NewInput().
					Title("password").Value(&cfg.Database.Postgres.Password),
				huh.NewInput().
					Title("host").Value(&cfg.Database.Postgres.Host),
				huh.NewInput().
					Title("port").Value(&cfg.Database.Postgres.Port),
				huh.NewInput().
					Title("database name").Value(&cfg.Database.Postgres.DbName),
			).WithHideFunc(func() bool {
				return database != "postgres"
			}),
			huh.NewGroup(
				huh.NewInput().Title("server port").Value(&cfg.Server.Port).
					Validate(func(s string) error {
						port, err := strconv.Atoi(s)
						if err != nil {
							return fmt.Errorf("Port %v is not a number", s)
						}
						if port <= 1023 {
							return fmt.Errorf("Please, use another port(0-1023 are system reserved)")
						}
						ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
						if err != nil {
							return fmt.Errorf("Port is already in use/restricted")
						}
						ln.Close()
						return nil
					}),
			),
		)

		if err := form.Run(); err != nil {
			os.Exit(1)
		}

		switch database {
		case "postgres":
			cfg.Database.Postgres.Used = true
			cfg.Database.Sqlite.Used = false
		case "sqlite":
			cfg.Database.Sqlite.Used = true
			cfg.Database.Postgres.Used = false
		}

		log.Print(cfg)
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
}

func main() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}
