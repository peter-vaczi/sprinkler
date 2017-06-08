package cmd

import (
	"log"
	"os"

	"github.com/peter.vaczi/sprinklerd/core"
	"github.com/peter.vaczi/sprinklerd/utils"
	"github.com/spf13/cobra"
)

var addFlagOn bool

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

		dev := core.Device{Name: args[0], On: addFlagOn}
		err := utils.PostRequest(daemonSocket+"/v1/devices", &dev)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().BoolVar(&addFlagOn, "on", false, "set the device on")
}
