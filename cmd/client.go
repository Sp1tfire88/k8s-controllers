package cmd

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func mustGetClientSet() *kubernetes.Clientset {
	kubeconfigPath := kubeconfig
	if kubeconfigPath == "" {
		home := os.Getenv("HOME")
		kubeconfigPath = filepath.Join(home, ".kube", "config")
	}
	flag.Set("logtostderr", "true")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load kubeconfig")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create clientset")
	}

	return clientset
}
