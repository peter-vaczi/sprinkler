package cmd

import (
	"github.com/spf13/cobra"

	"github.com/peter.vaczi/sprinklerd/api"
	"github.com/peter.vaczi/sprinklerd/core"
)

var devs *core.Devices

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the daemon process",
	Long:  `Start the daemon process`,
	Run: func(cmd *cobra.Command, args []string) {
		runDaemon()
	},
}

func init() {
	RootCmd.AddCommand(daemonCmd)
	devs = core.New()
	devs.Add("rotor-1")
	dev, _ := devs.Add("rotor-2")
	dev.TurnOn()
}

func runDaemon() {
	mainEvents := make(chan interface{})
	api.New(mainEvents)

	for {
		select {
		case event := <-mainEvents:
			handleEvent(event)
		}
	}
}

func handleEvent(event interface{}) {

	switch event := event.(type) {
	case api.HttpStatus:
		event.ResponseChan <- api.HttpResponse{Error: nil, Body: devs}
	}
}
