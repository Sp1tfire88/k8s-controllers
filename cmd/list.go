package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployments in the default namespace",
	Run: func(cmd *cobra.Command, args []string) {
		if err := listDeployments(); err != nil {
			log.Fatal().Err(err).Msg("Failed to list deployments")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file (default: ~/.kube/config)")
}

func listDeployments() error {
	if kubeconfig == "" {
		if home := os.Getenv("HOME"); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	log.Debug().Msgf("Using kubeconfig: %s", kubeconfig)

	if err := flag.Set("logtostderr", "true"); err != nil {
		log.Warn().Err(err).Msg("Failed to set flag 'logtostderr'")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	log.Info().Msg("Connected to cluster. Listing deployments...")

	deployments, err := clientset.AppsV1().Deployments("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	if len(deployments.Items) == 0 {
		log.Warn().Msg("No deployments found in 'default' namespace")
		return nil
	}

	log.Info().Msgf("Found %d deployment(s):", len(deployments.Items))
	for _, d := range deployments.Items {
		fmt.Println("📦", d.Name)
	}

	return nil
}
