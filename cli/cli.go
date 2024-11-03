package cli

import "github.com/spf13/cobra"

func New() *cobra.Command {
	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(NewBuildCommand())

	return rootCmd
}
