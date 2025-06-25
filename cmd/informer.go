package cmd

import (
	"flag"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	deploymentStore cache.Store
)

// GetDeploymentStore возвращает ссылку на informer cache
func GetDeploymentStore() cache.Store {
	return deploymentStore
}

func StartDeploymentInformerFromConfig() error {
	if err := flag.Set("logtostderr", "true"); err != nil {
		log.Warn().Err(err).Msg("Failed to set flag 'logtostderr'")
	}

	enabled := viper.GetBool("informer.enabled")
	if !enabled {
		log.Info().Msg("🔕 Informer is disabled via config")
		return nil
	}

	ns := viper.GetString("informer.namespace")
	if ns == "" {
		ns = "default"
	}

	resyncSeconds := viper.GetInt("informer.resyncPeriodSeconds")
	if resyncSeconds <= 0 {
		resyncSeconds = 60
	}
	resyncPeriod := time.Duration(resyncSeconds) * time.Second

	kubeconfigPath := viper.GetString("kubeconfig")
	inCluster := viper.GetBool("inCluster")

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

	log.Trace().Str("namespace", ns).Dur("resync", resyncPeriod).Msg("Creating informer factory")
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, resyncPeriod, informers.WithNamespace(ns))
	informer := factory.Apps().V1().Deployments().Informer()

	deploymentStore = informer.GetStore()

	_, handlerErr := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if d, ok := obj.(metav1.Object); ok {
				log.Trace().Str("deployment", d.GetName()).Msg("📦 Deployment ADDED")
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldDep, okOld := oldObj.(*appsv1.Deployment)
			newDep, okNew := newObj.(*appsv1.Deployment)
			if !okOld || !okNew {
				log.Warn().Msg("⚠️ Unable to cast Deployment object on update")
				return
			}

			oldReplicas := int32(1)
			newReplicas := int32(1)
			if oldDep.Spec.Replicas != nil {
				oldReplicas = *oldDep.Spec.Replicas
			}
			if newDep.Spec.Replicas != nil {
				newReplicas = *newDep.Spec.Replicas
			}

			if oldReplicas != newReplicas {
				log.Info().
					Str("deployment", newDep.Name).
					Int32("from", oldReplicas).
					Int32("to", newReplicas).
					Msg("🔁 Deployment scaled")
			} else {
				log.Trace().
					Str("deployment", newDep.Name).
					Msg("✏️ Deployment UPDATED (no replica change)")
			}
		},
		DeleteFunc: func(obj interface{}) {
			if d, ok := obj.(metav1.Object); ok {
				log.Trace().Str("deployment", d.GetName()).Msg("🗑️ Deployment DELETED")
			}
		},
	})
	if handlerErr != nil {
		log.Warn().Err(handlerErr).Msg("Failed to add event handler")
	}

	stop := make(chan struct{})
	// defer close(stop)

	log.Info().Msg("🚀 Starting deployment informer")
	go informer.Run(stop)

	return nil
}

var ErrNoConfigProvided = fmt.Errorf("either --kubeconfig or --in-cluster must be provided")
