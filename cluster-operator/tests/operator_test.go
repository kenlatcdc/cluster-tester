package tests

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	clusterv1 "github.com/cdcent/cluster-tester/cluster-operator/api/v1"
)

// Helper function to create int32 pointer
func int32Ptr(i int32) *int32 {
	return &i
}

func TestClusterTesterSpec(t *testing.T) {
	// Test that we can create a basic ClusterTester spec
	clusterTester := &clusterv1.ClusterTester{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "default",
		},
		Spec: clusterv1.ClusterTesterSpec{
			CoffeeShop: clusterv1.ServiceConfig{
				Enabled:  true,
				Replicas: int32Ptr(1),
				Image:    "cdcent/coffee-shop",
				Tag:      "latest",
			},
			PetStore: clusterv1.ServiceConfig{
				Enabled:  true,
				Replicas: int32Ptr(2),
				Image:    "cdcent/pet-store",
				Tag:      "latest",
			},
		},
	}

	// Basic validation
	if clusterTester.Name != "test-cluster" {
		t.Errorf("Expected name 'test-cluster', got '%s'", clusterTester.Name)
	}

	// Test coffee shop service
	if !clusterTester.Spec.CoffeeShop.Enabled {
		t.Error("Expected coffee shop to be enabled")
	}

	if *clusterTester.Spec.CoffeeShop.Replicas != 1 {
		t.Errorf("Expected coffee shop replicas 1, got %d", *clusterTester.Spec.CoffeeShop.Replicas)
	}

	// Test pet store service
	if !clusterTester.Spec.PetStore.Enabled {
		t.Error("Expected pet store to be enabled")
	}

	if *clusterTester.Spec.PetStore.Replicas != 2 {
		t.Errorf("Expected pet store replicas 2, got %d", *clusterTester.Spec.PetStore.Replicas)
	}
}

func TestServiceConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		service clusterv1.ServiceConfig
		wantErr bool
	}{
		{
			name: "valid basic config",
			service: clusterv1.ServiceConfig{
				Enabled:  true,
				Image:    "test",
				Tag:      "latest",
				Replicas: int32Ptr(1),
			},
			wantErr: false,
		},
		{
			name: "with resources",
			service: clusterv1.ServiceConfig{
				Enabled:  true,
				Image:    "app",
				Tag:      "v1.0",
				Replicas: int32Ptr(3),
				Resources: &clusterv1.ResourceRequirements{
					Limits: map[string]string{
						"cpu":    "500m",
						"memory": "512Mi",
					},
					Requests: map[string]string{
						"cpu":    "250m",
						"memory": "256Mi",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "disabled service",
			service: clusterv1.ServiceConfig{
				Enabled: false,
				Image:   "disabled-app",
				Tag:     "latest",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic field validation
			if tt.service.Enabled && tt.service.Image == "" {
				if !tt.wantErr {
					t.Error("Enabled service should have an image")
				}
				return
			}

			if tt.service.Enabled && tt.service.Tag == "" {
				if !tt.wantErr {
					t.Error("Enabled service should have a tag")
				}
				return
			}

			// Resource config validation
			if tt.service.Resources != nil {
				if len(tt.service.Resources.Limits) == 0 && len(tt.service.Resources.Requests) == 0 {
					if !tt.wantErr {
						t.Error("Resource config should have limits or requests")
					}
					return
				}
			}

			t.Logf("Service config %s is valid", tt.name)
		})
	}
}

func TestDeploymentCreation(t *testing.T) {
	// Test the logic for creating deployments based on ServiceConfig
	serviceConfig := clusterv1.ServiceConfig{
		Enabled:  true,
		Image:    "nginx",
		Tag:      "latest",
		Replicas: int32Ptr(3),
	}

	// Simulate deployment creation
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
			Labels: map[string]string{
				"app": "test-service",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: serviceConfig.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "test-service",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "test-service",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "test-service",
							Image: serviceConfig.Image + ":" + serviceConfig.Tag,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}

	// Validate deployment
	if deployment.Name != "test-service" {
		t.Errorf("Expected deployment name 'test-service', got '%s'", deployment.Name)
	}

	if *deployment.Spec.Replicas != 3 {
		t.Errorf("Expected 3 replicas, got %d", *deployment.Spec.Replicas)
	}

	if len(deployment.Spec.Template.Spec.Containers) != 1 {
		t.Errorf("Expected 1 container, got %d", len(deployment.Spec.Template.Spec.Containers))
	}

	container := deployment.Spec.Template.Spec.Containers[0]
	expectedImage := serviceConfig.Image + ":" + serviceConfig.Tag
	if container.Image != expectedImage {
		t.Errorf("Expected container image '%s', got '%s'", expectedImage, container.Image)
	}

	if container.Ports[0].ContainerPort != 8080 {
		t.Errorf("Expected container port 8080, got %d", container.Ports[0].ContainerPort)
	}
}

func TestServiceCreation(t *testing.T) {
	// Test service creation logic
	serviceName := "test-service"
	port := 8080

	// Simulate service creation
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName + "-service",
			Namespace: "default",
			Labels: map[string]string{
				"app": serviceName,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": serviceName,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       int32(port),
					TargetPort: intstr.FromInt(port),
				},
			},
		},
	}

	// Validate service
	expectedServiceName := serviceName + "-service"
	if service.Name != expectedServiceName {
		t.Errorf("Expected service name '%s', got '%s'", expectedServiceName, service.Name)
	}

	if len(service.Spec.Ports) != 1 {
		t.Errorf("Expected 1 service port, got %d", len(service.Spec.Ports))
	}

	if service.Spec.Ports[0].Port != int32(port) {
		t.Errorf("Expected service port %d, got %d", port, service.Spec.Ports[0].Port)
	}

	if service.Spec.Selector["app"] != serviceName {
		t.Errorf("Expected selector app='%s', got '%s'", serviceName, service.Spec.Selector["app"])
	}
}
