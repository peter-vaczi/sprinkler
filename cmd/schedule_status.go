package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/peter-vaczi/sprinkler/core"
	"github.com/peter-vaczi/sprinkler/utils"
)

var scheduleStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status",
	Long:  `Show status`,
	Run: func(cmd *cobra.Command, args []string) {

		var schs core.Schedules

		err := utils.GetRequest(daemonSocket+"/v1/schedules", &schs)
		if err != nil {
			log.Fatal(err)
		}

		printSchedules(schs)
	},
}

func printSchedules(schs core.Schedules) {
	keys := make([]string, 0, len(schs))
	for k := range schs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 5, 0, 1, ' ', 0)
	fmt.Fprintln(w, "NAME\tENABLED\tPROGRAM\tSPEC\t")

	for _, k := range keys {
		enab := "disabled"
		if schs[k].Enabled {
			enab = "enabled"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", schs[k].Name, enab, schs[k].ProgramName, schs[k].Spec)
	}

	w.Flush()
}

func init() {
	scheduleCmd.AddCommand(scheduleStatusCmd)
}
