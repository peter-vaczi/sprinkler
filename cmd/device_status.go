package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/peter-vaczi/sprinklerd/core"
	"github.com/peter-vaczi/sprinklerd/utils"
)

var deviceStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status",
	Long:  `Show status`,
	Run: func(cmd *cobra.Command, args []string) {

		var devs core.Devices

		err := utils.GetRequest(daemonSocket+"/v1/devices", &devs)
		if err != nil {
			log.Fatal(err)
		}

		printDevices(devs)
	},
}

func printDevices(devs core.Devices) {
	keys := make([]string, 0, len(devs))
	for k := range devs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 5, 0, 1, ' ', 0)
	fmt.Fprintln(w, "NAME\tPIN\tSTATUS\t")

	for _, k := range keys {
		onoff := "off"
		if devs[k].On {
			onoff = "on"
		}
		fmt.Fprintf(w, "%s\t%d\t%s\t\n", devs[k].Name, devs[k].Pin, onoff)
	}

	w.Flush()
}

func init() {
	deviceCmd.AddCommand(deviceStatusCmd)
}
