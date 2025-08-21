package tests

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	clusterv1 "github.com/cdcent/cluster-tester/cluster-operator/api/v1"
)

func TestCreateDeployment(t *testing.T) {
	serviceConfig := clusterv1.ServiceConfig{
		Image:    "test-image",
		Tag:      "latest",
		Replicas: int32Ptr(2),
		Enabled:  true,
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
	}

	deployment := createDeploymentForService(serviceConfig, "test-app", "test-namespace")

	// Test deployment metadata
	if deployment.Name != "test-app" {
		t.Errorf("Expected deployment name 'test-app', got '%s'", deployment.Name)
	}

	if deployment.Namespace != "test-namespace" {
		t.Errorf("Expected deployment namespace 'test-namespace', got '%s'", deployment.Namespace)
	}

	// Test deployment spec
	if *deployment.Spec.Replicas != 2 {
		t.Errorf("Expected 2 replicas, got %d", *deployment.Spec.Replicas)
	}

	// Test container spec
	containers := deployment.Spec.Template.Spec.Containers
	if len(containers) != 1 {
		t.Errorf("Expected 1 container, got %d", len(containers))
	}

	container := containers[0]
	if container.Name != "test-app" {
		t.Errorf("Expected container name 'test-app', got '%s'", container.Name)
	}

	expectedImage := serviceConfig.Image + ":" + serviceConfig.Tag
	if container.Image != expectedImage {
		t.Errorf("Expected container image '%s', got '%s'", expectedImage, container.Image)
	}

	// Test container port
	if len(container.Ports) != 1 {
		t.Errorf("Expected 1 container port, got %d", len(container.Ports))
	}

	if container.Ports[0].ContainerPort != 8080 {
		t.Errorf("Expected container port 8080, got %d", container.Ports[0].ContainerPort)
	}

	// Test resource requirements
	limits := container.Resources.Limits
	requests := container.Resources.Requests

	expectedCPULimit := resource.MustParse("500m")
	if !limits[corev1.ResourceCPU].Equal(expectedCPULimit) {
		t.Errorf("Expected CPU limit 500m, got %v", limits[corev1.ResourceCPU])
	}

	expectedMemoryRequest := resource.MustParse("256Mi")
	if !requests[corev1.ResourceMemory].Equal(expectedMemoryRequest) {
		t.Errorf("Expected memory request 256Mi, got %v", requests[corev1.ResourceMemory])
	}
}

func TestCreateService(t *testing.T) {
	serviceName := "test-service-app"
	port := 9090

	service := createServiceForApp(serviceName, port, "test-namespace")

	// Test service metadata
	if service.Name != "test-service-app-service" {
		t.Errorf("Expected service name 'test-service-app-service', got '%s'", service.Name)
	}

	if service.Namespace != "test-namespace" {
		t.Errorf("Expected service namespace 'test-namespace', got '%s'", service.Namespace)
	}

	// Test service spec
	if len(service.Spec.Ports) != 1 {
		t.Errorf("Expected 1 service port, got %d", len(service.Spec.Ports))
	}

	port32 := service.Spec.Ports[0]
	if port32.Port != int32(port) {
		t.Errorf("Expected service port %d, got %d", port, port32.Port)
	}

	if port32.TargetPort != intstr.FromInt(port) {
		t.Errorf("Expected target port %d, got %v", port, port32.TargetPort)
	}

	// Test selector
	if service.Spec.Selector["app"] != serviceName {
		t.Errorf("Expected selector app=%s, got %s", serviceName, service.Spec.Selector["app"])
	}
}

func TestCreateMySQLResources(t *testing.T) {
	dbConfig := clusterv1.DatabaseConfig{
		Type:        "mysql",
		StorageSize: "2Gi",
		Enabled:     true,
	}
	appName := "test-db-app"

	// Test PVC creation
	pvc := createPVCForDatabase(dbConfig, appName, "test-namespace")

	if pvc.Name != "mysql-test-db-app-pvc" {
		t.Errorf("Expected PVC name 'mysql-test-db-app-pvc', got '%s'", pvc.Name)
	}

	if pvc.Namespace != "test-namespace" {
		t.Errorf("Expected PVC namespace 'test-namespace', got '%s'", pvc.Namespace)
	}

	// Test storage size
	storageQuantity := pvc.Spec.Resources.Requests[corev1.ResourceStorage]
	expectedStorage := resource.MustParse("2Gi")
	if !storageQuantity.Equal(expectedStorage) {
		t.Errorf("Expected storage 2Gi, got %v", storageQuantity)
	}

	// Test access modes
	if len(pvc.Spec.AccessModes) != 1 || pvc.Spec.AccessModes[0] != corev1.ReadWriteOnce {
		t.Errorf("Expected access mode ReadWriteOnce, got %v", pvc.Spec.AccessModes)
	}

	// Test MySQL deployment
	deployment := createMySQLDeployment(dbConfig, appName, "test-namespace")

	if deployment.Name != "mysql-test-db-app" {
		t.Errorf("Expected MySQL deployment name 'mysql-test-db-app', got '%s'", deployment.Name)
	}

	// Test MySQL service
	service := createMySQLService(appName, "test-namespace")

	if service.Name != "mysql-test-db-app-service" {
		t.Errorf("Expected MySQL service name 'mysql-test-db-app-service', got '%s'", service.Name)
	}

	if len(service.Spec.Ports) != 1 || service.Spec.Ports[0].Port != 3306 {
		t.Errorf("Expected MySQL service port 3306, got %v", service.Spec.Ports)
	}
}

// Helper functions to test - these would normally be in the controller package
func createDeploymentForService(serviceConfig clusterv1.ServiceConfig, serviceName, namespace string) *appsv1.Deployment {
	replicas := serviceConfig.Replicas
	if replicas == nil {
		replicas = int32Ptr(1)
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": serviceName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": serviceName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": serviceName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  serviceName,
							Image: serviceConfig.Image + ":" + serviceConfig.Tag,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080, // Default port for testing
								},
							},
						},
					},
				},
			},
		},
	}

	// Add resource requirements if specified
	if serviceConfig.Resources != nil {
		resources := corev1.ResourceRequirements{
			Limits:   make(corev1.ResourceList),
			Requests: make(corev1.ResourceList),
		}

		for k, v := range serviceConfig.Resources.Limits {
			resources.Limits[corev1.ResourceName(k)] = resource.MustParse(v)
		}

		for k, v := range serviceConfig.Resources.Requests {
			resources.Requests[corev1.ResourceName(k)] = resource.MustParse(v)
		}

		deployment.Spec.Template.Spec.Containers[0].Resources = resources
	}

	return deployment
}

func createServiceForApp(serviceName string, port int, namespace string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName + "-service",
			Namespace: namespace,
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
}

func createPVCForDatabase(dbConfig clusterv1.DatabaseConfig, appName, namespace string) *corev1.PersistentVolumeClaim {
	storageQuantity := resource.MustParse(dbConfig.StorageSize)

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql-" + appName + "-pvc",
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: storageQuantity,
				},
			},
		},
	}
}

func createMySQLDeployment(dbConfig clusterv1.DatabaseConfig, appName, namespace string) *appsv1.Deployment {
	replicas := int32(1)
	mysqlName := "mysql-" + appName

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": mysqlName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": mysqlName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": mysqlName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mysql",
							Image: "mysql:8.0",
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: "rootpassword",
								},
								{
									Name:  "MYSQL_DATABASE",
									Value: appName,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3306,
								},
							},
						},
					},
				},
			},
		},
	}
}

func createMySQLService(appName, namespace string) *corev1.Service {
	mysqlName := "mysql-" + appName

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlName + "-service",
			Namespace: namespace,
			Labels: map[string]string{
				"app": mysqlName,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": mysqlName,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       3306,
					TargetPort: intstr.FromInt(3306),
				},
			},
		},
	}
}
