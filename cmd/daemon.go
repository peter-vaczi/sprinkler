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
}

func runDaemon() {
	mainEvents := make(chan interface{})
	core.InitGpio()
	api.New(daemonSocket, mainEvents)

	for {
		select {
		case event := <-mainEvents:
			handleEvent(event)
		}
	}
}

func handleEvent(event interface{}) {

	switch event := event.(type) {
	case api.HttpDeviceList:
		event.ResponseChan <- api.HttpResponse{Error: nil, Body: devs}
	case api.HttpDeviceAdd:
		err := devs.Add(event.Device)
		event.ResponseChan <- api.HttpResponse{Error: err}
	case api.HttpDeviceDel:
		err := devs.Del(event.Name)
		event.ResponseChan <- api.HttpResponse{Error: err}
	}
}
