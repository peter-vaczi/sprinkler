package cmd

import (
	"log"
	"os"

	"github.com/peter-vaczi/sprinkler/core"
	"github.com/peter-vaczi/sprinkler/utils"
	"github.com/spf13/cobra"
)

var scheduleAddFlagEnable bool
var scheduleAddFlagProgram string
var scheduleAddFlagSpec string

// scheduleAddCmd represents the add command
var scheduleAddCmd = &cobra.Command{
	Use:   "add <name> [flags]",
	Short: "Add a new schedule",
	Long:  `Add a new schedule`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(-1)
		}

		sch := core.Schedule{Name: args[0], ProgramName: scheduleAddFlagProgram, Spec: scheduleAddFlagSpec, Enabled: scheduleAddFlagEnable}
		err := utils.PostRequest(daemonSocket+"/v1/schedules", &sch)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	scheduleCmd.AddCommand(scheduleAddCmd)
	scheduleAddCmd.PersistentFlags().StringVar(&scheduleAddFlagProgram, "program", "", "the program to start")
	scheduleAddCmd.PersistentFlags().StringVar(&scheduleAddFlagSpec, "spec", "", "the scheduling specification")
	scheduleAddCmd.PersistentFlags().BoolVar(&scheduleAddFlagEnable, "enable", false, "enable the schedule")
}
