package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var deleteName string

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Deployment",
	Run: func(cmd *cobra.Command, args []string) {
		clientset := mustGetClientSet()
		err := clientset.AppsV1().Deployments("default").Delete(context.Background(), deleteName, metav1.DeleteOptions{})
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to delete Deployment %q", deleteName)
		}
		log.Info().Msgf("🗑️ Deployment %q deleted", deleteName)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&deleteName, "name", "", "Name of the Deployment to delete")
	deleteCmd.MarkFlagRequired("name")
}
