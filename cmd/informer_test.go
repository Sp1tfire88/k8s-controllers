package cmd_test

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
)

func TestStartDeploymentInformer_FakeClient(t *testing.T) {
	client := fake.NewSimpleClientset()
	factory := informers.NewSharedInformerFactoryWithOptions(
		client,
		0, // no resync
		informers.WithNamespace("default"),
	)

	informer := factory.Apps().V1().Deployments().Informer()
	handled := make(chan string, 1)

	informer.AddEventHandler(
		// nolint:errcheck // not needed for this callback
		newTestEventHandler(func(name string) {
			handled <- name
		}),
	)

	stop := make(chan struct{})
	defer close(stop)
	factory.Start(stop)
	factory.WaitForCacheSync(stop)

	// Create a deployment
	_, err := client.AppsV1().Deployments("default").Create(
		context.TODO(),
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-deployment",
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "test"},
				},
				Template: appsv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "test"},
					},
					Spec: appsv1.PodSpec{
						Containers: []appsv1.Container{{
							Name:  "nginx",
							Image: "nginx",
						}},
					},
				},
			},
		},
		metav1.CreateOptions{},
	)
	if err != nil {
		t.Fatalf("failed to create test deployment: %v", err)
	}

	select {
	case name := <-handled:
		log.Info().Str("event", "add").Str("deployment", name).Msg("Received deployment event")
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for deployment event")
	}
}

func newTestEventHandler(onAdd func(name string)) *testHandler {
	return &testHandler{onAdd: onAdd}
}

type testHandler struct {
	onAdd func(name string)
}

func (h *testHandler) OnAdd(obj interface{}) {
	if d, ok := obj.(metav1.Object); ok && h.onAdd != nil {
		h.onAdd(d.GetName())
	}
}

func (h *testHandler) OnUpdate(oldObj, newObj interface{}) {}

func (h *testHandler) OnDelete(obj interface{}) {}
