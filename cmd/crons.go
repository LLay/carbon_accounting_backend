package cmd

import (
	"github.com/michelaquino/golang_api_skeleton/config"
	"github.com/spf13/cobra"

	"github.com/michelaquino/golang_api_skeleton/src/energy"
)

var cronsCmd = &cobra.Command{
	Use:   "crons",
	Short: "Starts crons server",
	Long:  `Starts crons server.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		energy.FetchEIAData()
	},
}
