package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/peter-vaczi/sprinklerd/utils"
	"github.com/spf13/cobra"
)

// programDelDeviceCmd represents the deldevice command
var programDelDeviceCmd = &cobra.Command{
	Use:   "deldevice <program> <device>",
	Short: "Delete a device from a watering program",
	Long:  `Delete a device from a watering program`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			cmd.Usage()
			os.Exit(-1)
		}

		err := utils.DeleteRequest(
			fmt.Sprintf("%s/v1/programs/%s/devices/%s", daemonSocket, args[0], args[1]))
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	programCmd.AddCommand(programDelDeviceCmd)
}
