package controller

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	clusterv1 "github.com/cdcent/cluster-tester/cluster-operator/api/v1"
)

func TestClusterTesterReconciler_BasicReconcile(t *testing.T) {
	// Set up the test scheme
	scheme := runtime.NewScheme()
	if err := clusterv1.AddToScheme(scheme); err != nil {
		t.Fatalf("Failed to add ClusterTester scheme: %v", err)
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("Failed to add core v1 scheme: %v", err)
	}
	if err := appsv1.AddToScheme(scheme); err != nil {
		t.Fatalf("Failed to add apps v1 scheme: %v", err)
	}

	// Create a ClusterTester resource
	clusterTester := &clusterv1.ClusterTester{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "default",
		},
		Spec: clusterv1.ClusterTesterSpec{
			Services: []clusterv1.ServiceConfig{
				{
					Name:  "test-service",
					Image: "test-image:latest",
					Port:  8080,
				},
			},
		},
	}

	// Create fake client and reconciler
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(clusterTester).
		Build()

	reconciler := &ClusterTesterReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	// Test reconciliation
	ctx := context.Background()
	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-cluster",
			Namespace: "default",
		},
	}

	// Run reconcile
	result, err := reconciler.Reconcile(ctx, req)
	if err != nil {
		t.Errorf("Reconcile failed: %v", err)
	}

	t.Logf("Reconcile result: requeue=%v, requeueAfter=%v", result.Requeue, result.RequeueAfter)

	// Verify that Deployment was created
	deployment := &appsv1.Deployment{}
	err = fakeClient.Get(ctx, types.NamespacedName{
		Name:      "test-service",
		Namespace: "default",
	}, deployment)
	if err != nil {
		t.Errorf("Expected Deployment 'test-service' to be created: %v", err)
	} else {
		t.Logf("✓ Deployment 'test-service' created successfully")

		// Verify deployment details
		if deployment.Spec.Replicas == nil || *deployment.Spec.Replicas != 1 {
			t.Errorf("Expected 1 replica, got %v", deployment.Spec.Replicas)
		}

		if len(deployment.Spec.Template.Spec.Containers) != 1 {
			t.Errorf("Expected 1 container, got %d", len(deployment.Spec.Template.Spec.Containers))
		} else {
			container := deployment.Spec.Template.Spec.Containers[0]
			if container.Image != "test-image:latest" {
				t.Errorf("Expected image 'test-image:latest', got '%s'", container.Image)
			}
		}
	}

	// Verify that Service was created
	service := &corev1.Service{}
	err = fakeClient.Get(ctx, types.NamespacedName{
		Name:      "test-service",
		Namespace: "default",
	}, service)
	if err != nil {
		t.Errorf("Expected Service 'test-service' to be created: %v", err)
	} else {
		t.Logf("✓ Service 'test-service' created successfully")

		// Verify service details
		if service.Spec.Type != corev1.ServiceTypeClusterIP {
			t.Errorf("Expected service type ClusterIP, got %s", service.Spec.Type)
		}

		if len(service.Spec.Ports) != 1 || service.Spec.Ports[0].Port != 8080 {
			t.Errorf("Expected service port 8080, got %v", service.Spec.Ports)
		}
	}
}

func TestClusterTesterReconciler_MultipleServices(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := clusterv1.AddToScheme(scheme); err != nil {
		t.Fatalf("Failed to add schemes: %v", err)
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("Failed to add schemes: %v", err)
	}
	if err := appsv1.AddToScheme(scheme); err != nil {
		t.Fatalf("Failed to add schemes: %v", err)
	}

	clusterTester := &clusterv1.ClusterTester{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "multi-service-test",
			Namespace: "default",
		},
		Spec: clusterv1.ClusterTesterSpec{
			Services: []clusterv1.ServiceConfig{
				{
					Name:  "service-a",
					Image: "image-a:latest",
					Port:  8080,
				},
				{
					Name:  "service-b",
					Image: "image-b:latest",
					Port:  8081,
				},
			},
		},
	}

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(clusterTester).
		Build()

	reconciler := &ClusterTesterReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	ctx := context.Background()
	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "multi-service-test",
			Namespace: "default",
		},
	}

	_, err := reconciler.Reconcile(ctx, req)
	if err != nil {
		t.Fatalf("Reconcile failed: %v", err)
	}

	// Verify both services were created
	serviceNames := []string{"service-a", "service-b"}
	for _, serviceName := range serviceNames {
		// Check Deployment
		deployment := &appsv1.Deployment{}
		err = fakeClient.Get(ctx, types.NamespacedName{
			Name:      serviceName,
			Namespace: "default",
		}, deployment)
		if err != nil {
			t.Errorf("Expected Deployment '%s' to be created: %v", serviceName, err)
		} else {
			t.Logf("✓ Deployment '%s' created successfully", serviceName)
		}

		// Check Service
		service := &corev1.Service{}
		err = fakeClient.Get(ctx, types.NamespacedName{
			Name:      serviceName,
			Namespace: "default",
		}, service)
		if err != nil {
			t.Errorf("Expected Service '%s' to be created: %v", serviceName, err)
		} else {
			t.Logf("✓ Service '%s' created successfully", serviceName)
		}
	}
}
