// pkg/multicluster/manager.go
package multicluster

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterConfig struct {
	Name       string `mapstructure:"name"`
	Kubeconfig string `mapstructure:"kubeconfig"`
}

func StartMultiClusterInformers() error {
	if err := flag.Set("logtostderr", "true"); err != nil {
		log.Warn().Err(err).Msg("Failed to set logtostderr")
	}

	var clusters []ClusterConfig
	if err := viper.UnmarshalKey("clusters", &clusters); err != nil {
		return fmt.Errorf("failed to load clusters config: %w", err)
	}

	if len(clusters) == 0 {
		return fmt.Errorf("no clusters defined in config")
	}

	wg := sync.WaitGroup{}

	for _, cluster := range clusters {
		cfg, err := clientcmd.BuildConfigFromFlags("", cluster.Kubeconfig)
		if err != nil {
			log.Error().Err(err).Str("cluster", cluster.Name).Msg("Failed to build kubeconfig")
			continue
		}

		clientset, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			log.Error().Err(err).Str("cluster", cluster.Name).Msg("Failed to create clientset")
			continue
		}

		factory := informers.NewSharedInformerFactory(clientset, time.Minute*5)
		informer := factory.Apps().V1().Deployments().Informer()

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if d, ok := obj.(v1.Object); ok {
					log.Info().Str("cluster", cluster.Name).Str("deployment", d.GetName()).Msg("📦 Deployment ADDED")
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				if d, ok := newObj.(v1.Object); ok {
					log.Info().Str("cluster", cluster.Name).Str("deployment", d.GetName()).Msg("✏️ Deployment UPDATED")
				}
			},
			DeleteFunc: func(obj interface{}) {
				if d, ok := obj.(v1.Object); ok {
					log.Info().Str("cluster", cluster.Name).Str("deployment", d.GetName()).Msg("🗑️ Deployment DELETED")
				}
			},
		})

		wg.Add(1)
		go func(c string, inf cache.SharedIndexInformer) {
			defer wg.Done()
			log.Info().Str("cluster", c).Msg("🚀 Starting informer")
			stop := make(chan struct{})
			inf.Run(stop)
		}(cluster.Name, informer)
	}

	wg.Wait()
	return nil
}
