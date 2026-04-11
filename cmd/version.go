package cmd

import (
	"fmt"

	"github.com/jiang-zhexin/animedb/internal/app"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v%s\n", app.Appname, app.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
