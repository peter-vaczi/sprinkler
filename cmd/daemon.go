package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/peter-vaczi/sprinkler/api"
	"github.com/peter-vaczi/sprinkler/core"
	"github.com/peter-vaczi/sprinkler/gpio"
)

var testMode bool

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
	daemonCmd.PersistentFlags().BoolVar(&testMode, "test-mode", false, "test mode, do no handle gpio device")
	RootCmd.AddCommand(daemonCmd)
}

func runDaemon() {
	if testMode {
		g := gpio.NewDummy()
		core.InitGpio(g)
	} else {
		g, err := gpio.New()
		if err != nil {
			log.Fatalf("failed to initialize the gpio library: %v", err)
		}
		core.InitGpio(g)
	}

	data := core.LoadState()
	if data != nil {
		api := api.New(daemonSocket, data)
		go api.Run()
		waitForSignal()
		data.Schedules.DisableAll()
		data.Programs.StopAll()
		data.StoreState()
	}
}

func waitForSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case <-sigChan:
			return
		}
	}
}
