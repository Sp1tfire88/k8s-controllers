package cmd

import (
	"flag"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var ErrNoConfigProvided = fmt.Errorf("either --kubeconfig or --in-cluster must be provided")

// StartDeploymentInformer запускает SharedInformer для Deployments
func StartDeploymentInformer(kubeconfigPath string, inCluster bool, namespace string) error {
	// Подавление логов client-go
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
	store := informer.GetStore()

	handlerFuncs := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if d, ok := obj.(metav1.Object); ok {
				log.Trace().Str("deployment", d.GetName()).Msg("📦 Deployment ADDED")
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldDep, okOld := oldObj.(*appsv1.Deployment)
			newDep, okNew := newObj.(*appsv1.Deployment)
			if okOld && okNew {
				oldReplicas := int32(0)
				newReplicas := int32(0)
				if oldDep.Spec.Replicas != nil {
					oldReplicas = *oldDep.Spec.Replicas
				}
				if newDep.Spec.Replicas != nil {
					newReplicas = *newDep.Spec.Replicas
				}

				if oldReplicas != newReplicas {
					log.Info().
						Str("deployment", newDep.Name).
						Int32("old", oldReplicas).
						Int32("new", newReplicas).
						Msg("🔁 Replicas count changed")
				} else {
					log.Trace().
						Str("deployment", newDep.Name).
						Msg("✏️ Deployment UPDATED (no replica change)")
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			var name string

			switch t := obj.(type) {
			case cache.DeletedFinalStateUnknown:
				if d, ok := t.Obj.(metav1.Object); ok {
					name = d.GetName()
				}
			case metav1.Object:
				name = t.GetName()
			default:
				log.Warn().Msg("Unknown type for deleted object")
				return
			}

			log.Trace().Str("deployment", name).Msg("🗑️ Deployment DELETED")

			key := fmt.Sprintf("%s/%s", namespace, name)
			_, exists, err := store.GetByKey(key)
			if err != nil {
				log.Error().Err(err).Str("deployment", name).Msg("Failed to retrieve from cache")
			} else if exists {
				log.Warn().Str("deployment", name).Msg("⚠️ Deployment still in cache (possibly stale)")
			} else {
				log.Info().Str("deployment", name).Msg("✅ Confirmed deletion from cache")
			}

		},
	}

	if _, err := informer.AddEventHandler(handlerFuncs); err != nil {
		log.Fatal().Err(err).Msg("Failed to register event handler")
	}

	stop := make(chan struct{})
	defer close(stop)

	log.Info().Msg("🚀 Starting deployment informer")
	informer.Run(stop)

	return nil
}
