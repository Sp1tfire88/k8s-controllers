// cmd/multicluster_cmd.go
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var multiCmd = &cobra.Command{
	Use:   "multi-informers",
	Short: "Start multi-cluster informers for Deployments",
	Run: func(cmd *cobra.Command, args []string) {
		if err := StartMultiClusterInformers(); err != nil {
			log.Fatal().Err(err).Msg("Failed to start multi-cluster informers")
		}
	},
}

func init() {
	rootCmd.AddCommand(multiCmd)
}
