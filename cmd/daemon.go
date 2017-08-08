package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/peter.vaczi/sprinklerd/api"
	"github.com/peter.vaczi/sprinklerd/core"
	"github.com/peter.vaczi/sprinklerd/gpio"
)

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
}

func runDaemon() {
	mainEvents := make(chan interface{})
	g, err := gpio.New()
	if err != nil {
		log.Fatalf("failed to initialize the gpio library: %v", err)
	}
	core.InitGpio(g)
	api.New(daemonSocket, mainEvents)

	core.LoadState()
	core.Run(mainEvents)
	core.StoreState()
}
