package cmd

import (
	"log"
	"os"

	"github.com/peter-vaczi/sprinkler/utils"
	"github.com/spf13/cobra"
)

// scheduleDelCmd represents the del command
var scheduleDelCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a schedule",
	Long:  `Delete a schedule`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(-1)
		}

		err := utils.DeleteRequest(daemonSocket + "/v1/schedules/" + args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	scheduleCmd.AddCommand(scheduleDelCmd)
}
