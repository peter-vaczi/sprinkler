package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/peter.vaczi/sprinklerd/utils"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		var ret interface{}

		err := utils.GetRequest(daemonSocket+"/v1/status", &ret)
		if err != nil {
			log.Fatal(err)
		}

		log.Print(ret)
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
