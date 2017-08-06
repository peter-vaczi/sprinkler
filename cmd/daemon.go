package cmd

import (
	"github.com/spf13/cobra"

	"github.com/peter.vaczi/sprinklerd/api"
	"github.com/peter.vaczi/sprinklerd/core"
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
	core.InitGpio()
	api.New(daemonSocket, mainEvents)

	core.LoadState()
	core.Run(mainEvents)
	core.StoreState()
}
