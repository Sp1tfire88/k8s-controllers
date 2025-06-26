// cmd/multicluster_informer.go
package cmd

import (
	"errors"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var ErrNoClustersConfigured = errors.New("no clusters configured in config")

type ClusterConfig struct {
	Name       string `mapstructure:"name"`
	Kubeconfig string `mapstructure:"kubeconfig"`
	Namespace  string `mapstructure:"namespace"`
}

func StartMultiClusterInformers() error {
	if err := flag.Set("logtostderr", "true"); err != nil {
		log.Warn().Err(err).Msg("Failed to set 'logtostderr'")
	}

	var clusters []ClusterConfig
	if err := viper.UnmarshalKey("clusters", &clusters); err != nil {
		return fmt.Errorf("failed to parse cluster config: %w", err)
	}

	if len(clusters) == 0 {
		return ErrNoClustersConfigured
	}

	log.Info().Msgf("Starting informers for %d clusters", len(clusters))

	var wg sync.WaitGroup
	for _, c := range clusters {
		cluster := c // capture for goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()

			config, err := clientcmd.BuildConfigFromFlags("", cluster.Kubeconfig)
			if err != nil {
				log.Error().Err(err).Str("cluster", cluster.Name).Msg("Failed to load kubeconfig")
				return
			}

			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				log.Error().Err(err).Str("cluster", cluster.Name).Msg("Failed to create clientset")
				return
			}

			ns := cluster.Namespace
			if ns == "" {
				ns = metav1.NamespaceDefault
			}

			factory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Minute*5,
				informers.WithNamespace(ns),
			)
			informer := factory.Apps().V1().Deployments().Informer()

			_, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					log.Info().Str("cluster", cluster.Name).Msg("📦 Deployment ADDED")
				},
				UpdateFunc: func(oldObj, newObj interface{}) {
					log.Info().Str("cluster", cluster.Name).Msg("✏️ Deployment UPDATED")
				},
				DeleteFunc: func(obj interface{}) {
					log.Info().Str("cluster", cluster.Name).Msg("🗑️ Deployment DELETED")
				},
			})
			if err != nil {
				log.Error().Err(err).Str("cluster", cluster.Name).Msg("Failed to add event handler")
				return
			}

			log.Info().Str("cluster", cluster.Name).Msg("🚀 Running informer")
			stop := make(chan struct{})
			defer close(stop)
			informer.Run(stop)
		}()
	}

	wg.Wait()
	return nil
}
