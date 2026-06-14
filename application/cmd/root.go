/*
Copyright © 2026 ETS
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var verbose bool

var PathToConfig = "../../internal/config/config.yaml"

// rootCmd represents the base command when called without any subcommands
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fathom.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Use to get debug-level output")

	rootCmd.AddCommand(serveCmd)

	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configRegenerateCmd)
}
