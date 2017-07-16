package cmd

import (
	"log"
	"os"

	"github.com/peter.vaczi/sprinklerd/utils"
	"github.com/spf13/cobra"
)

// programDelCmd represents the del command
var programDelCmd = &cobra.Command{
	Use:   "del <name>",
	Short: "Delete a watering program",
	Long:  `Delete a watering program`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(-1)
		}

		err := utils.DeleteRequest(daemonSocket + "/v1/programs/" + args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	programCmd.AddCommand(programDelCmd)
}
