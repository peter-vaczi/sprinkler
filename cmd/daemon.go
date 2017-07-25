package cmd

import (
	"github.com/spf13/cobra"

	"github.com/peter.vaczi/sprinklerd/api"
	"github.com/peter.vaczi/sprinklerd/core"
)

var devs *core.Devices
var programs *core.Programs

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
	devs = core.NewDevices()
	programs = core.NewPrograms()
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
	case api.HttpDeviceGet:
		dev, err := devs.Get(event.Name)
		event.ResponseChan <- api.HttpResponse{Error: err, Body: dev}
	case api.HttpDeviceDel:
		err := devs.Del(event.Name)
		event.ResponseChan <- api.HttpResponse{Error: err}
	case api.HttpDeviceSet:
		err := devs.Set(event.Name, event.Device)
		event.ResponseChan <- api.HttpResponse{Error: err}
	case api.HttpProgramList:
		event.ResponseChan <- api.HttpResponse{Error: nil, Body: programs}
	case api.HttpProgramCreate:
		err := programs.Add(event.Program)
		event.ResponseChan <- api.HttpResponse{Error: err}
	case api.HttpProgramGet:
		prg, err := programs.Get(event.Name)
		event.ResponseChan <- api.HttpResponse{Error: err, Body: prg}
	case api.HttpProgramDel:
		err := programs.Del(event.Name)
		event.ResponseChan <- api.HttpResponse{Error: err}
	case api.HttpProgramAddDevice:
		prg, err := programs.Get(event.Program)
		if err != nil {
			event.ResponseChan <- api.HttpResponse{Error: err}
			return
		}
		dev, err := devs.Get(event.Device)
		if err != nil {
			event.ResponseChan <- api.HttpResponse{Error: err}
			return
		}
		err = prg.AddDevice(dev, event.Duration)
		event.ResponseChan <- api.HttpResponse{Error: err}
	case api.HttpProgramDelDevice:
		prg, err := programs.Get(event.Program)
		if err != nil {
			event.ResponseChan <- api.HttpResponse{Error: err}
			return
		}
		err = prg.DelDevice(event.Idx)
		event.ResponseChan <- api.HttpResponse{Error: err}
	}
}
