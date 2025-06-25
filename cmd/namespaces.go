package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var namespacesCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "List all Kubernetes namespaces",
	Run: func(cmd *cobra.Command, args []string) {
		clientset := mustGetClientSet()

		namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list namespaces")
		}

		if len(namespaces.Items) == 0 {
			log.Warn().Msg("No namespaces found")
			return
		}

		log.Info().Msgf("Found %d namespace(s):", len(namespaces.Items))

		// Таблица для форматированного вывода
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tSTATUS\tAGE")

		now := time.Now()
		for _, ns := range namespaces.Items {
			age := now.Sub(ns.CreationTimestamp.Time).Round(time.Second)
			fmt.Fprintf(w, "%s\t%s\t%s\n", ns.Name, ns.Status.Phase, formatDuration(age))
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(namespacesCmd)
}

// formatDuration форматирует время в стиле "3d4h" или "15m"
func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60

	switch {
	case days > 0:
		return fmt.Sprintf("%dd%dh", days, hours)
	case hours > 0:
		return fmt.Sprintf("%dh%dm", hours, mins)
	default:
		return fmt.Sprintf("%dm", mins)
	}
}
