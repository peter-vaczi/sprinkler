package cmd

import (
	"log"
	"os"

	"github.com/peter-vaczi/sprinklerd/utils"
	"github.com/spf13/cobra"
)

// deviceDelCmd represents the del command
var deviceDelCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a device",
	Long:  `Delete a device`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(-1)
		}

		err := utils.DeleteRequest(daemonSocket + "/v1/devices/" + args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	deviceCmd.AddCommand(deviceDelCmd)
}
