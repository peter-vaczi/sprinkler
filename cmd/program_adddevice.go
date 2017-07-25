package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/peter.vaczi/sprinklerd/utils"
	"github.com/spf13/cobra"
)

var programAddDeviceDuration string

// programAddDeviceCmd represents the adddevice command
var programAddDeviceCmd = &cobra.Command{
	Use:   "adddevice <program> <device>",
	Short: "Add a new device to a watering program",
	Long:  `Add a new device to a watering program`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			cmd.Usage()
			os.Exit(-1)
		}

		data := make(map[string]string)
		data["device"] = args[1]
		data["duration"] = programAddDeviceDuration
		err := utils.PostRequest(
			fmt.Sprintf("%s/v1/programs/%s/devices", daemonSocket, args[0]),
			&data)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	programAddDeviceCmd.Flags().StringVarP(&programAddDeviceDuration, "duration", "d", "15m", "duration of , e.g.: 1s, 2m, 3h, 2h45m")
	programCmd.AddCommand(programAddDeviceCmd)
}
