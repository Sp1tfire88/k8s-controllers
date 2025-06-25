// // // informer.go
// package cmd

// import (
// 	"flag"
// 	"fmt"
// 	"time"

// 	"github.com/rs/zerolog/log"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/informers"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/rest"
// 	"k8s.io/client-go/tools/cache"
// 	"k8s.io/client-go/tools/clientcmd"
// )

// func StartDeploymentInformer(kubeconfigPath string, inCluster bool, namespace string) error {
// 	if err := flag.Set("logtostderr", "true"); err != nil {
// 		log.Warn().Err(err).Msg("Failed to set flag 'logtostderr'")
// 	}

// 	var config *rest.Config
// 	var err error

// 	switch {
// 	case inCluster:
// 		log.Trace().Msg("Using in-cluster configuration")
// 		config, err = rest.InClusterConfig()
// 	case kubeconfigPath != "":
// 		log.Trace().Str("kubeconfig", kubeconfigPath).Msg("Using external kubeconfig")
// 		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
// 	default:
// 		return ErrNoConfigProvided
// 	}

// 	if err != nil {
// 		log.Fatal().Err(err).Msg("Failed to build Kubernetes config")
// 	}

// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("Failed to create Kubernetes clientset")
// 	}

// 	log.Trace().Str("namespace", namespace).Msg("Creating informer factory")
// 	factory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Minute*10, informers.WithNamespace(namespace))

// 	informer := factory.Apps().V1().Deployments().Informer()

// 	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
// 		AddFunc: func(obj interface{}) {
// 			if d, ok := obj.(metav1.Object); ok {
// 				log.Trace().Str("deployment", d.GetName()).Msg("📦 Deployment ADDED")
// 			}
// 		},
// 		UpdateFunc: func(oldObj, newObj interface{}) {
// 			if d, ok := newObj.(metav1.Object); ok {
// 				log.Trace().Str("deployment", d.GetName()).Msg("✏️ Deployment UPDATED")
// 			}
// 		},
// 		DeleteFunc: func(obj interface{}) {
// 			if d, ok := obj.(metav1.Object); ok {
// 				log.Trace().Str("deployment", d.GetName()).Msg("🗑️ Deployment DELETED")
// 			}
// 		},
// 	})

// 	stop := make(chan struct{})
// 	defer close(stop)

// 	log.Info().Msg("🚀 Starting deployment informer")
// 	informer.Run(stop)

// 	return nil
// }

// var ErrNoConfigProvided = fmt.Errorf("either --kubeconfig or --in-cluster must be provided")
// // informer.go
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

func StartDeploymentInformer(kubeconfigPath string, inCluster bool, namespace string) error {
	if err := flag.Set("logtostderr", "true"); err != nil {
		log.Warn().Err(err).Msg("Failed to set flag 'logtostderr'")
	}

	var config *rest.Config
	var err error

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
		log.Fatal().Err(err).Msg("Failed to build Kubernetes config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Kubernetes clientset")
	}

	log.Trace().Str("namespace", namespace).Msg("Creating informer factory")
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Minute*10, informers.WithNamespace(namespace))

	informer := factory.Apps().V1().Deployments().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if d, ok := obj.(metav1.Object); ok {
				log.Trace().Str("deployment", d.GetName()).Msg("📦 Deployment ADDED")
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if d, ok := newObj.(metav1.Object); ok {
				log.Trace().Str("deployment", d.GetName()).Msg("✏️ Deployment UPDATED")
			}
		},
		DeleteFunc: func(obj interface{}) {
			if d, ok := obj.(metav1.Object); ok {
				log.Trace().Str("deployment", d.GetName()).Msg("🗑️ Deployment DELETED")
			}
		},
	})

	stop := make(chan struct{})
	defer close(stop)

	log.Info().Msg("🚀 Starting deployment informer")
	informer.Run(stop)

	return nil
}

var ErrNoConfigProvided = fmt.Errorf("either --kubeconfig or --in-cluster must be provided")
