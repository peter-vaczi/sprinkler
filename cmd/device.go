package cmd

import "github.com/spf13/cobra"

// deviceCmd represents the device command
var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Handle sprinkler devices",
	Long:  `Handle sprinkler devices`,
}

func init() {
	RootCmd.AddCommand(deviceCmd)
}
