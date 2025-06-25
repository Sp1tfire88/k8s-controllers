//go:build ignore
// +build ignore

package cmd

import (
	"context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

// testHandler implements cache.ResourceEventHandler
type testHandler struct {
	onAdd func(name string)
}

func (h *testHandler) OnAdd(obj interface{}, _ bool) {
	if d, ok := obj.(*appsv1.Deployment); ok && h.onAdd != nil {
		h.onAdd(d.Name)
	}
}

func (h *testHandler) OnUpdate(_, _ interface{}) {}
func (h *testHandler) OnDelete(_ interface{})    {}

// newTestEventHandler returns a handler that calls back on add
func newTestEventHandler(f func(name string)) cache.ResourceEventHandler {
	return &testHandler{onAdd: f}
}

func TestDeploymentInformer_AddEvent(t *testing.T) {
	client := fake.NewSimpleClientset()

	factory := informers.NewSharedInformerFactoryWithOptions(client, 0, informers.WithNamespace("default"))
	informer := factory.Apps().V1().Deployments().Informer()

	deployName := "test-deployment"
	received := make(chan string, 1)

	informer.AddEventHandler(newTestEventHandler(func(name string) {
		received <- name
	}))

	stop := make(chan struct{})
	defer close(stop)
	go factory.Start(stop)

	if !cache.WaitForCacheSync(stop, informer.HasSynced) {
		t.Fatal("Failed to sync cache")
	}

	// Create Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployName,
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": deployName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": deployName},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "nginx",
						Image: "nginx:latest",
					}},
				},
			},
		},
	}

	_, err := client.AppsV1().Deployments("default").Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create deployment: %v", err)
	}

	select {
	case name := <-received:
		if name != deployName {
			t.Errorf("Expected deployment name %q, got %q", deployName, name)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for event handler")
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}
