package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/peter.vaczi/sprinklerd/core"
	"github.com/peter.vaczi/sprinklerd/utils"
	"github.com/spf13/cobra"
)

// programShowCmd represents the show command
var programShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show details of a watering program",
	Long:  `Show details of a watering program`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(-1)
		}

		var prg core.Program

		err := utils.GetRequest(daemonSocket+"/v1/programs/"+args[0], &prg)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Name: %s\n\n", prg.Name)
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 5, 0, 1, ' ', 0)
		fmt.Fprintln(w, "NR\tDEVICE\tDURATION\t")

		for i, e := range prg.Elements {
			fmt.Fprintf(w, "%d\t%s\t%s\t\n", i, e.Device.Name, e.Duration)
		}

		fmt.Fprintln(w)
		w.Flush()
	},
}

func init() {
	programCmd.AddCommand(programShowCmd)
}
