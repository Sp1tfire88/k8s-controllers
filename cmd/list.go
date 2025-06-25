package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployments in the specified namespace",
	Run: func(cmd *cobra.Command, args []string) {
		if err := listDeployments(); err != nil {
			log.Fatal().Err(err).Msg("Failed to list deployments")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Флаг kubeconfig и namespace уже определены как PersistentFlags в root.go
	// Здесь ничего добавлять не нужно
}

func listDeployments() error {
	log.Debug().Msgf("Using kubeconfig: %s", kubeconfig)
	log.Debug().Msgf("Using namespace: %s", namespace)

	// Подавление вывода klog в stderr
	if err := flag.Set("logtostderr", "true"); err != nil {
		log.Warn().Err(err).Msg("Failed to set flag 'logtostderr'")
	}

	clientset := mustGetClientSet()

	log.Info().Msg("Connected to cluster. Listing deployments...")

	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	if len(deployments.Items) == 0 {
		log.Warn().Msgf("No deployments found in namespace %q", namespace)
		return nil
	}

	log.Info().Msgf("Found %d deployment(s):", len(deployments.Items))
	for _, d := range deployments.Items {
		fmt.Println("📦", d.Name)
	}

	return nil
}
