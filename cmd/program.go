package cmd

import "github.com/spf13/cobra"

// programCmd represents the program command
var programCmd = &cobra.Command{
	Use:   "program",
	Short: "Handle watering programs",
	Long:  `Handle watering programs`,
}

func init() {
	RootCmd.AddCommand(programCmd)
}
