package main

import (
	"github.com/alexisvisco/pocketpase-gen/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pb-gen --db-path /path/to/pocketbase.sqlite",
	Short: "Generate things around pocketbase",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		return
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&commands.FlagDBPath, "db-path", "pb_data/data.db",
		"path to pocketbase instance sqlite file")
	rootCmd.PersistentFlags().BoolVar(&commands.Verbose, "verbose", false, "enable verbose mode")
	rootCmd.AddCommand(commands.ModelsCommand)

}
