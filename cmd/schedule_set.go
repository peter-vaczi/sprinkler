package cmd

import (
	"log"
	"os"

	"github.com/peter-vaczi/sprinkler/core"
	"github.com/peter-vaczi/sprinkler/utils"
	"github.com/spf13/cobra"
)

var scheduleSetFlagEnable bool
var scheduleSetFlagDisable bool
var scheduleSetFlagProgram string
var scheduleSetFlagSpec string

// scheduleSetCmd represents the add command
var scheduleSetCmd = &cobra.Command{
	Use:   "set <name> [flags]",
	Short: "Set parameters of a schedule",
	Long:  `Set parameters of a schedule`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 || (scheduleSetFlagEnable && scheduleSetFlagDisable) {
			cmd.Usage()
			os.Exit(-1)
		}

		var sch core.Schedule

		err := utils.GetRequest(daemonSocket+"/v1/schedules/"+args[0], &sch)
		if err != nil {
			log.Fatal(err)
		}

		if 0 < len(scheduleSetFlagProgram) {
			sch.ProgramName = scheduleSetFlagProgram
		}
		if 0 < len(scheduleSetFlagSpec) {
			sch.Spec = scheduleSetFlagSpec
		}
		if scheduleSetFlagEnable {
			sch.Enabled = true
		}
		if scheduleSetFlagDisable {
			sch.Enabled = false
		}

		err = utils.PutRequest(daemonSocket+"/v1/schedules/"+sch.Name, &sch)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	scheduleCmd.AddCommand(scheduleSetCmd)
	scheduleSetCmd.PersistentFlags().StringVar(&scheduleSetFlagProgram, "program", "", "program to start when this schedule became active")
	scheduleSetCmd.PersistentFlags().StringVar(&scheduleSetFlagSpec, "spec", "", "the scheduling specification")
	scheduleSetCmd.PersistentFlags().BoolVar(&scheduleSetFlagEnable, "enable", false, "enable the schedule")
	scheduleSetCmd.PersistentFlags().BoolVar(&scheduleSetFlagDisable, "disable", false, "disable the schedule")
}
