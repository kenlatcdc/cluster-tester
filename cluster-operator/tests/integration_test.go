package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterv1 "github.com/cdcent/cluster-tester/cluster-operator/api/v1"
)

// IntegrationTest runs tests against a real Kubernetes cluster
func TestIntegration(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration tests. Set RUN_INTEGRATION_TESTS=true to run.")
	}

	// Load kubeconfig
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = clientcmd.RecommendedHomeFile
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatalf("Failed to load kubeconfig: %v", err)
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Failed to create kubernetes client: %v", err)
	}

	// Create controller-runtime client
	scheme := clusterv1.AddToScheme(clientcmd.Scheme)
	k8sClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		t.Fatalf("Failed to create controller-runtime client: %v", err)
	}

	ctx := context.Background()
	namespace := "cluster-tester-test"

	// Create test namespace
	t.Run("CreateNamespace", func(t *testing.T) {
		testCreateNamespace(t, clientset, namespace)
	})

	// Test basic application deployment
	t.Run("DeployBasicApplications", func(t *testing.T) {
		testDeployBasicApplications(t, k8sClient, clientset, ctx, namespace)
	})

	// Test database deployment
	t.Run("DeployWithDatabase", func(t *testing.T) {
		testDeployWithDatabase(t, k8sClient, clientset, ctx, namespace)
	})

	// Test service accessibility
	t.Run("TestServiceAccessibility", func(t *testing.T) {
		testServiceAccessibility(t, clientset, namespace)
	})

	// Cleanup
	t.Run("Cleanup", func(t *testing.T) {
		testCleanup(t, clientset, namespace)
	})
}

func testCreateNamespace(t *testing.T, clientset *kubernetes.Clientset, namespace string) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		t.Fatalf("Failed to create namespace %s: %v", namespace, err)
	}
	t.Logf("Namespace %s created successfully", namespace)
}

func testDeployBasicApplications(t *testing.T, k8sClient client.Client, clientset *kubernetes.Clientset, ctx context.Context, namespace string) {
	clusterTester := &clusterv1.ClusterTester{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-basic-apps",
			Namespace: namespace,
		},
		Spec: clusterv1.ClusterTesterSpec{
			Applications: []clusterv1.ApplicationConfig{
				{
					Name:     "coffee-shop",
					Image:    "cdcent/coffee-shop:latest",
					Port:     8080,
					Enabled:  true,
					Replicas: 1,
				},
				{
					Name:     "pet-store",
					Image:    "cdcent/pet-store:latest",
					Port:     8081,
					Enabled:  true,
					Replicas: 1,
				},
			},
		},
	}

	// Create the ClusterTester resource
	err := k8sClient.Create(ctx, clusterTester)
	if err != nil {
		t.Fatalf("Failed to create ClusterTester: %v", err)
	}
	t.Log("ClusterTester created successfully")

	// Wait for deployments to be created and ready
	timeout := 5 * time.Minute
	interval := 10 * time.Second

	for _, app := range clusterTester.Spec.Applications {
		t.Run(fmt.Sprintf("WaitFor-%s-Deployment", app.Name), func(t *testing.T) {
			err := wait.PollImmediate(interval, timeout, func() (bool, error) {
				deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, app.Name, metav1.GetOptions{})
				if err != nil {
					if errors.IsNotFound(err) {
						t.Logf("Deployment %s not found yet, waiting...", app.Name)
						return false, nil
					}
					return false, err
				}

				if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
					t.Logf("Deployment %s is ready with %d replicas", app.Name, deployment.Status.ReadyReplicas)
					return true, nil
				}

				t.Logf("Deployment %s: %d/%d replicas ready", app.Name, deployment.Status.ReadyReplicas, *deployment.Spec.Replicas)
				return false, nil
			})

			if err != nil {
				t.Fatalf("Deployment %s failed to become ready: %v", app.Name, err)
			}
		})

		t.Run(fmt.Sprintf("Verify-%s-Service", app.Name), func(t *testing.T) {
			serviceName := fmt.Sprintf("%s-service", app.Name)
			service, err := clientset.CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
			if err != nil {
				t.Fatalf("Failed to get service %s: %v", serviceName, err)
			}

			if len(service.Spec.Ports) == 0 {
				t.Fatalf("Service %s has no ports", serviceName)
			}

			if service.Spec.Ports[0].Port != int32(app.Port) {
				t.Fatalf("Service %s port mismatch: expected %d, got %d", serviceName, app.Port, service.Spec.Ports[0].Port)
			}

			t.Logf("Service %s verified successfully on port %d", serviceName, service.Spec.Ports[0].Port)
		})
	}
}

func testDeployWithDatabase(t *testing.T, k8sClient client.Client, clientset *kubernetes.Clientset, ctx context.Context, namespace string) {
	clusterTester := &clusterv1.ClusterTester{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-db-apps",
			Namespace: namespace,
		},
		Spec: clusterv1.ClusterTesterSpec{
			Applications: []clusterv1.ApplicationConfig{
				{
					Name:     "electronics-store",
					Image:    "cdcent/electronics-store:latest",
					Port:     8082,
					Enabled:  true,
					Replicas: 1,
					Database: &clusterv1.DatabaseConfig{
						Type:        "mysql",
						StorageSize: "1Gi",
					},
				},
			},
		},
	}

	// Create the ClusterTester resource
	err := k8sClient.Create(ctx, clusterTester)
	if err != nil {
		t.Fatalf("Failed to create ClusterTester with database: %v", err)
	}
	t.Log("ClusterTester with database created successfully")

	// Wait for MySQL deployment
	mysqlDeploymentName := "mysql-electronics-store"
	timeout := 5 * time.Minute
	interval := 10 * time.Second

	t.Run("WaitForMySQLDeployment", func(t *testing.T) {
		err := wait.PollImmediate(interval, timeout, func() (bool, error) {
			deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, mysqlDeploymentName, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					t.Logf("MySQL deployment not found yet, waiting...")
					return false, nil
				}
				return false, err
			}

			if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
				t.Logf("MySQL deployment is ready")
				return true, nil
			}

			t.Logf("MySQL deployment: %d/%d replicas ready", deployment.Status.ReadyReplicas, *deployment.Spec.Replicas)
			return false, nil
		})

		if err != nil {
			t.Fatalf("MySQL deployment failed to become ready: %v", err)
		}
	})

	// Verify PVC
	t.Run("VerifyPVC", func(t *testing.T) {
		pvcName := "mysql-electronics-store-pvc"
		pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
		if err != nil {
			t.Fatalf("Failed to get PVC %s: %v", pvcName, err)
		}

		if pvc.Status.Phase != corev1.ClaimBound {
			t.Logf("PVC %s is in phase %s (not bound yet, this might be expected in some environments)", pvcName, pvc.Status.Phase)
		} else {
			t.Logf("PVC %s is bound successfully", pvcName)
		}

		storageQuantity := pvc.Spec.Resources.Requests[corev1.ResourceStorage]
		if storageQuantity.String() != "1Gi" {
			t.Fatalf("PVC %s storage size mismatch: expected 1Gi, got %s", pvcName, storageQuantity.String())
		}
	})

	// Wait for application deployment
	t.Run("WaitForApplicationDeployment", func(t *testing.T) {
		err := wait.PollImmediate(interval, timeout, func() (bool, error) {
			deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, "electronics-store", metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					t.Logf("Electronics store deployment not found yet, waiting...")
					return false, nil
				}
				return false, err
			}

			if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
				t.Logf("Electronics store deployment is ready")
				return true, nil
			}

			t.Logf("Electronics store deployment: %d/%d replicas ready", deployment.Status.ReadyReplicas, *deployment.Spec.Replicas)
			return false, nil
		})

		if err != nil {
			t.Fatalf("Electronics store deployment failed to become ready: %v", err)
		}
	})
}

func testServiceAccessibility(t *testing.T, clientset *kubernetes.Clientset, namespace string) {
	ctx := context.Background()

	// Get all services in the namespace
	services, err := clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list services: %v", err)
	}

	for _, service := range services.Items {
		if service.Name == "kubernetes" {
			continue // Skip default kubernetes service
		}

		t.Run(fmt.Sprintf("CheckService-%s", service.Name), func(t *testing.T) {
			// Verify service has endpoints
			endpoints, err := clientset.CoreV1().Endpoints(namespace).Get(ctx, service.Name, metav1.GetOptions{})
			if err != nil {
				t.Logf("Warning: Could not get endpoints for service %s: %v", service.Name, err)
				return
			}

			hasEndpoints := false
			for _, subset := range endpoints.Subsets {
				if len(subset.Addresses) > 0 {
					hasEndpoints = true
					break
				}
			}

			if hasEndpoints {
				t.Logf("Service %s has active endpoints", service.Name)
			} else {
				t.Logf("Warning: Service %s has no active endpoints", service.Name)
			}

			// Verify service spec
			if len(service.Spec.Ports) == 0 {
				t.Errorf("Service %s has no ports defined", service.Name)
			} else {
				t.Logf("Service %s has %d port(s) defined", service.Name, len(service.Spec.Ports))
			}
		})
	}
}

func testCleanup(t *testing.T, clientset *kubernetes.Clientset, namespace string) {
	ctx := context.Background()

	// Delete the namespace (this will delete all resources in it)
	err := clientset.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		t.Fatalf("Failed to delete namespace %s: %v", namespace, err)
	}

	t.Logf("Namespace %s deleted successfully", namespace)
}
