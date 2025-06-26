package controller_test

import (
	"context"
	"testing"

	"github.com/Sp1tfire88/k8s-controllers/pkg/controller"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestDeploymentReconciler_Reconcile(t *testing.T) {
	// 1. Подготовим схему и фейковый клиент
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = appsv1.AddToScheme(scheme)

	// 2. Создадим фейковый Deployment
	dep := &appsv1.Deployment{}
	dep.Name = "test-deployment"
	dep.Namespace = "default"

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(dep).
		Build()

	// 3. Инициализируем контроллер
	r := &controller.DeploymentReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	// 4. Выполним Reconcile
	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-deployment",
			Namespace: "default",
		},
	}

	res, err := r.Reconcile(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if res.Requeue || res.RequeueAfter != 0 {
		t.Errorf("expected no requeue, got: %+v", res)
	}
}
