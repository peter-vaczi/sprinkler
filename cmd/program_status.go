package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/peter.vaczi/sprinklerd/core"
	"github.com/peter.vaczi/sprinklerd/utils"
)

var programStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status",
	Long:  `Show status`,
	Run: func(cmd *cobra.Command, args []string) {

		var progs core.Programs

		err := utils.GetRequest(daemonSocket+"/v1/programs", &progs)
		if err != nil {
			log.Fatal(err)
		}

		printPrograms(progs)
	},
}

func printPrograms(progs core.Programs) {
	keys := make([]string, 0, len(progs))
	for k := range progs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 5, 0, 1, ' ', 0)
	fmt.Fprintln(w, "NAME\t")

	for _, k := range keys {
		fmt.Fprintf(w, "%s\t\n", progs[k].Name)
	}

	w.Flush()
}

func init() {
	programCmd.AddCommand(programStatusCmd)
}
