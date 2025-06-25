package cmd

import (
	"flag"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var ErrNoConfigProvided = fmt.Errorf("either --kubeconfig or --in-cluster must be provided")

// StartDeploymentInformer запускает SharedInformer для деплойментов в указанном namespace
func StartDeploymentInformer(kubeconfigPath string, inCluster bool, namespace string) error {
	if err := flag.Set("logtostderr", "true"); err != nil {
		log.Warn().Err(err).Msg("Failed to set flag 'logtostderr'")
	}

	var (
		config *rest.Config
		err    error
	)

	switch {
	case inCluster:
		log.Trace().Msg("Using in-cluster configuration")
		config, err = rest.InClusterConfig()
	case kubeconfigPath != "":
		log.Trace().Str("kubeconfig", kubeconfigPath).Msg("Using external kubeconfig")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	default:
		return ErrNoConfigProvided
	}

	if err != nil {
		return fmt.Errorf("failed to build Kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes clientset: %w", err)
	}

	log.Trace().Str("namespace", namespace).Msg("Creating informer factory")
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 10*time.Minute, informers.WithNamespace(namespace))

	informer := factory.Apps().V1().Deployments().Informer()

	handlerFuncs := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if d, ok := obj.(metav1.Object); ok {
				log.Trace().Str("deployment", d.GetName()).Msg("📦 Deployment ADDED")
			}
		},
		UpdateFunc: func(_, newObj interface{}) {
			if d, ok := newObj.(metav1.Object); ok {
				log.Trace().Str("deployment", d.GetName()).Msg("✏️ Deployment UPDATED")
			}
		},
		DeleteFunc: func(obj interface{}) {
			if d, ok := obj.(metav1.Object); ok {
				log.Trace().Str("deployment", d.GetName()).Msg("🗑️ Deployment DELETED")
			}
		},
	}

	informer.AddEventHandler(handlerFuncs)

	stop := make(chan struct{})
	defer close(stop)

	log.Info().Msg("🚀 Starting deployment informer")
	informer.Run(stop)

	return nil
}
