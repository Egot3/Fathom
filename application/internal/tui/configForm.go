package tui

import (
	"fmt"
	"net"
	"slices"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/egot3/fathom/internal/config"
)

func ConfigForm(cfg *config.Config) error {
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
			huh.NewMultiSelect[string]().
				Options(
					huh.NewOption("charm's log", "charmLog").Selected(slices.Contains(cfg.Server.Logging.Loggers, "charmLog")),
					huh.NewOption("structured log", "slog").Selected(slices.Contains(cfg.Server.Logging.Loggers, "slog")),
				).Title("Logging").
				Description(`charm is human-readable, structured is machine parsable,
					disable all if you don't read logs`).
				Value(&cfg.Server.Logging.Loggers),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	switch database {
	case "postgres":
		cfg.Database.Postgres.Used = true
		cfg.Database.Sqlite.Used = false
	case "sqlite":
		cfg.Database.Sqlite.Used = true
		cfg.Database.Postgres.Used = false
	}

	return nil
}
