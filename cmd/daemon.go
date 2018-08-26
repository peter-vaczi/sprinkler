package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/peter-vaczi/sprinkler/api"
	"github.com/peter-vaczi/sprinkler/core"
	"github.com/peter-vaczi/sprinkler/gpio"
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
	g, err := gpio.New()
	if err != nil {
		log.Fatalf("failed to initialize the gpio library: %v", err)
	}
	core.InitGpio(g)

	data := core.LoadState()
	if data != nil {
		api := api.New(daemonSocket, data)
		go api.Run()

		data.StoreState()
	}
}
