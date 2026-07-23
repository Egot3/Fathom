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

var PathToConfig string

func LoadPath() {
	PathToConfig = os.Getenv("PATH_TO_CONFIG")
}

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

func Execute() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Use to get debug-level output")

	rootCmd.AddCommand(serveCmd)

	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configRegenerateCmd)
}
