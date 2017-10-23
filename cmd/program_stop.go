package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/peter-vaczi/sprinklerd/utils"
	"github.com/spf13/cobra"
)

// programStopCmd represents the add command
var programStopCmd = &cobra.Command{
	Use:   "stop <name>",
	Short: "Stop a watering program",
	Long:  `Stop a  watering program`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(-1)
		}

		err := utils.PostRequest(
			fmt.Sprintf("%s/v1/programs/%s/stop", daemonSocket, args[0]), nil)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	programCmd.AddCommand(programStopCmd)
}
