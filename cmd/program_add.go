package cmd

import (
	"log"
	"os"

	"github.com/peter.vaczi/sprinklerd/core"
	"github.com/peter.vaczi/sprinklerd/utils"
	"github.com/spf13/cobra"
)

// programAddCmd represents the add command
var programAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new watering program",
	Long:  `Add a new watering program`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(-1)
		}

		prg := core.Program{Name: args[0]}
		err := utils.PostRequest(daemonSocket+"/v1/programs", &prg)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	programCmd.AddCommand(programAddCmd)
}
