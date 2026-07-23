package cmd

import (
	"log"
	"os"

	"github.com/egot3/fathom/internal/config"
	"github.com/egot3/fathom/internal/tui"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Helps configuring fathom project",
	Long: `Starts an interactive TUI for configuring
	fathom project with pizzazis. If config is corrupted
	run "fampls config regenerate"`,
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		log.Printf("path to config: %v", PathToConfig)

		cfg := config.Config{}
		data, err := os.ReadFile(PathToConfig)
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}

		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Printf("couldn't read config data: %v", err)
			os.Exit(1)
		}

		err = tui.ConfigForm(&cfg)
		if err != nil {
			os.Exit(1)
		}

		out, err := yaml.Marshal(cfg)
		if err != nil {
			os.Exit(1)
		}

		if err := os.WriteFile(PathToConfig, out, 0644); err != nil {
			os.Exit(1)
		}
	},
}
