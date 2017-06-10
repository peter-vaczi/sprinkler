package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/spf13/cobra"

	"github.com/peter.vaczi/sprinklerd/core"
	"github.com/peter.vaczi/sprinklerd/utils"
)

var statusCmd = &cobra.Command{
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

	printLine([]string{"NAME", "PIN", "STATUS"})
	for _, k := range keys {
		printDevice(devs[k])
	}
}

func printDevice(dev *core.Device) {
	line := make([]string, 0, 2)
	line = append(line, dev.Name)
	line = append(line, fmt.Sprintf("%d", dev.Pin))
	if dev.On {
		line = append(line, "on")
	} else {
		line = append(line, "off")
	}
	printLine(line)
}

func printLine(line []string) {
	fmt.Printf("%-20s %-3s %-5s\n", line[0], line[1], line[2])
}

func init() {
	deviceCmd.AddCommand(statusCmd)
}
