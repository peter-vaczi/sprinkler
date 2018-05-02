package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/peter-vaczi/sprinkler/utils"
	"github.com/spf13/cobra"
)

// programStartCmd represents the add command
var programStartCmd = &cobra.Command{
	Use:   "start <name>",
	Short: "Start a watering program",
	Long:  `Start a  watering program`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(-1)
		}

		err := utils.PostRequest(
			fmt.Sprintf("%s/v1/programs/%s/start", daemonSocket, args[0]), nil)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	programCmd.AddCommand(programStartCmd)
}
