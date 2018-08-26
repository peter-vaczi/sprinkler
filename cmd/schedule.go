package cmd

import "github.com/spf13/cobra"

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Handle watering schedules",
	Long:  `Handle watering schedules`,
}

func init() {
	RootCmd.AddCommand(scheduleCmd)
}
