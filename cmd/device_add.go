package cmd

import (
	"log"
	"os"

	"github.com/peter.vaczi/sprinklerd/core"
	"github.com/peter.vaczi/sprinklerd/utils"
	"github.com/spf13/cobra"
)

var addFlagOn bool
var addFlagPin int

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <name> [flags]",
	Short: "Add a new device",
	Long:  `Add a new device`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(-1)
		}

		dev := core.Device{Name: args[0], On: addFlagOn, Pin: addFlagPin}
		err := utils.PostRequest(daemonSocket+"/v1/devices", &dev)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	deviceCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().IntVar(&addFlagPin, "pin", 0, "GPIO pin associated with this device")
	addCmd.PersistentFlags().BoolVar(&addFlagOn, "on", false, "set the device on")
}