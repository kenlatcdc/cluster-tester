/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	clusterv1 "github.com/cdcent/cluster-tester/cluster-operator/api/v1"
)

// ClusterTesterReconciler reconciles a ClusterTester object
type ClusterTesterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cluster.cdcent.io,resources=clustertesters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.cdcent.io,resources=clustertesters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cluster.cdcent.io,resources=clustertesters/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ClusterTesterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the ClusterTester instance
	var clusterTester clusterv1.ClusterTester
	if err := r.Get(ctx, req.NamespacedName, &clusterTester); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("ClusterTester resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get ClusterTester")
		return ctrl.Result{}, err
	}

	// Update status to indicate reconciliation is starting
	if clusterTester.Status.Phase == "" {
		clusterTester.Status.Phase = "Initializing"
		if err := r.Status().Update(ctx, &clusterTester); err != nil {
			logger.Error(err, "Failed to update ClusterTester status")
			return ctrl.Result{}, err
		}
	}

	// Deploy database if enabled
	if clusterTester.Spec.Database.Enabled {
		if err := r.reconcileDatabase(ctx, &clusterTester); err != nil {
			logger.Error(err, "Failed to reconcile database")
			return r.updateStatusError(ctx, &clusterTester, "DatabaseFailed", err)
		}
	}

	// Deploy services
	services := r.getServiceConfigs(&clusterTester)
	var serviceStatuses []clusterv1.ServiceStatus

	for serviceName, config := range services {
		if config.Enabled {
			status, err := r.reconcileService(ctx, &clusterTester, serviceName, config)
			if err != nil {
				logger.Error(err, "Failed to reconcile service", "service", serviceName)
				return r.updateStatusError(ctx, &clusterTester, "ServiceFailed", err)
			}
			serviceStatuses = append(serviceStatuses, status)
		}
	}

	// Update status
	clusterTester.Status.Services = serviceStatuses
	clusterTester.Status.Phase = "Ready"
	clusterTester.Status.ObservedGeneration = clusterTester.Generation

	// Set ready condition
	readyCondition := metav1.Condition{
		Type:    "Ready",
		Status:  metav1.ConditionTrue,
		Reason:  "ServicesReady",
		Message: "All services are ready",
	}
	meta.SetStatusCondition(&clusterTester.Status.Conditions, readyCondition)

	if err := r.Status().Update(ctx, &clusterTester); err != nil {
		logger.Error(err, "Failed to update ClusterTester status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

func (r *ClusterTesterReconciler) getServiceConfigs(clusterTester *clusterv1.ClusterTester) map[string]clusterv1.ServiceConfig {
	services := make(map[string]clusterv1.ServiceConfig)

	// Set defaults if not specified
	defaultReplicas := int32(1)

	// Coffee Shop
	coffeeShop := clusterTester.Spec.CoffeeShop
	if coffeeShop.Replicas == nil {
		coffeeShop.Replicas = &defaultReplicas
	}
	if coffeeShop.Image == "" {
		coffeeShop.Image = "coffee-shop"
	}
	if coffeeShop.Tag == "" {
		coffeeShop.Tag = "latest"
	}
	services["coffee-shop"] = coffeeShop

	// Pet Store
	petStore := clusterTester.Spec.PetStore
	if petStore.Replicas == nil {
		petStore.Replicas = &defaultReplicas
	}
	if petStore.Image == "" {
		petStore.Image = "pet-store"
	}
	if petStore.Tag == "" {
		petStore.Tag = "latest"
	}
	services["pet-store"] = petStore

	// Restaurant
	restaurant := clusterTester.Spec.Restaurant
	if restaurant.Replicas == nil {
		restaurant.Replicas = &defaultReplicas
	}
	if restaurant.Image == "" {
		restaurant.Image = "restaurant"
	}
	if restaurant.Tag == "" {
		restaurant.Tag = "latest"
	}
	services["restaurant"] = restaurant

	// College Admission
	collegeAdmission := clusterTester.Spec.CollegeAdmission
	if collegeAdmission.Replicas == nil {
		collegeAdmission.Replicas = &defaultReplicas
	}
	if collegeAdmission.Image == "" {
		collegeAdmission.Image = "college-admission"
	}
	if collegeAdmission.Tag == "" {
		collegeAdmission.Tag = "latest"
	}
	services["college-admission"] = collegeAdmission

	// Electronics Store
	electronicsStore := clusterTester.Spec.ElectronicsStore
	if electronicsStore.Replicas == nil {
		electronicsStore.Replicas = &defaultReplicas
	}
	if electronicsStore.Image == "" {
		electronicsStore.Image = "electronics-store"
	}
	if electronicsStore.Tag == "" {
		electronicsStore.Tag = "latest"
	}
	services["electronics-store"] = electronicsStore

	// Electronics Store Tracing
	electronicsStoreTracing := clusterTester.Spec.ElectronicsStoreTracing
	if electronicsStoreTracing.Replicas == nil {
		electronicsStoreTracing.Replicas = &defaultReplicas
	}
	if electronicsStoreTracing.Image == "" {
		electronicsStoreTracing.Image = "electronics-store-tracing"
	}
	if electronicsStoreTracing.Tag == "" {
		electronicsStoreTracing.Tag = "latest"
	}
	services["electronics-store-tracing"] = electronicsStoreTracing

	return services
}

func (r *ClusterTesterReconciler) reconcileService(ctx context.Context, clusterTester *clusterv1.ClusterTester, serviceName string, config clusterv1.ServiceConfig) (clusterv1.ServiceStatus, error) {
	logger := log.FromContext(ctx)

	namespace := clusterTester.Namespace
	if clusterTester.Spec.Global.Namespace != "" {
		namespace = clusterTester.Spec.Global.Namespace
	}

	// Create deployment
	deployment := r.createDeployment(clusterTester, serviceName, config, namespace)
	if err := controllerutil.SetControllerReference(clusterTester, deployment, r.Scheme); err != nil {
		return clusterv1.ServiceStatus{}, err
	}

	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating deployment", "deployment", deployment.Name)
		if err = r.Create(ctx, deployment); err != nil {
			return clusterv1.ServiceStatus{}, err
		}
	} else if err != nil {
		return clusterv1.ServiceStatus{}, err
	} else {
		// Update deployment if needed
		found.Spec = deployment.Spec
		if err = r.Update(ctx, found); err != nil {
			return clusterv1.ServiceStatus{}, err
		}
	}

	// Create service
	service := r.createService(clusterTester, serviceName, namespace)
	if err := controllerutil.SetControllerReference(clusterTester, service, r.Scheme); err != nil {
		return clusterv1.ServiceStatus{}, err
	}

	foundService := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating service", "service", service.Name)
		if err = r.Create(ctx, service); err != nil {
			return clusterv1.ServiceStatus{}, err
		}
	} else if err != nil {
		return clusterv1.ServiceStatus{}, err
	}

	// Get current deployment status
	if err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found); err != nil {
		return clusterv1.ServiceStatus{}, err
	}

	status := clusterv1.ServiceStatus{
		Name:          serviceName,
		Ready:         found.Status.ReadyReplicas == found.Status.Replicas && found.Status.Replicas > 0,
		Replicas:      found.Status.Replicas,
		ReadyReplicas: found.Status.ReadyReplicas,
		Endpoint:      fmt.Sprintf("%s.%s.svc.cluster.local:8080", service.Name, namespace),
	}

	return status, nil
}

func (r *ClusterTesterReconciler) createDeployment(clusterTester *clusterv1.ClusterTester, serviceName string, config clusterv1.ServiceConfig, namespace string) *appsv1.Deployment {
	labels := map[string]string{
		"app":                          serviceName,
		"app.kubernetes.io/name":       serviceName,
		"app.kubernetes.io/instance":   clusterTester.Name,
		"app.kubernetes.io/component":  "microservice",
		"app.kubernetes.io/part-of":    "cluster-tester",
		"app.kubernetes.io/managed-by": "cluster-tester-operator",
	}

	imagePullPolicy := corev1.PullIfNotPresent
	if clusterTester.Spec.Global.ImagePullPolicy != "" {
		imagePullPolicy = corev1.PullPolicy(clusterTester.Spec.Global.ImagePullPolicy)
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: config.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            serviceName,
							Image:           fmt.Sprintf("%s:%s", config.Image, config.Tag),
							ImagePullPolicy: imagePullPolicy,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health",
										Port: intstr.FromInt(8080),
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health",
										Port: intstr.FromInt(8080),
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       5,
							},
						},
					},
				},
			},
		},
	}

	// Add resource requirements if specified
	if config.Resources != nil {
		resources := corev1.ResourceRequirements{}
		if config.Resources.Limits != nil {
			resources.Limits = make(corev1.ResourceList)
			for k, v := range config.Resources.Limits {
				resources.Limits[corev1.ResourceName(k)] = *parseQuantity(v)
			}
		}
		if config.Resources.Requests != nil {
			resources.Requests = make(corev1.ResourceList)
			for k, v := range config.Resources.Requests {
				resources.Requests[corev1.ResourceName(k)] = *parseQuantity(v)
			}
		}
		deployment.Spec.Template.Spec.Containers[0].Resources = resources
	}

	// Add database environment variables for services that need them
	if serviceName == "electronics-store" || serviceName == "electronics-store-tracing" {
		deployment.Spec.Template.Spec.Containers[0].Env = []corev1.EnvVar{
			{
				Name:  "DB_HOST",
				Value: "mysql",
			},
			{
				Name:  "DB_PORT",
				Value: "3306",
			},
			{
				Name:  "DB_NAME",
				Value: "electronics-store",
			},
			{
				Name:  "DB_USER",
				Value: "admin",
			},
			{
				Name:  "DB_PASSWORD",
				Value: "password123",
			},
		}
	}

	return deployment
}

func (r *ClusterTesterReconciler) createService(clusterTester *clusterv1.ClusterTester, serviceName string, namespace string) *corev1.Service {
	labels := map[string]string{
		"app":                          serviceName,
		"app.kubernetes.io/name":       serviceName,
		"app.kubernetes.io/instance":   clusterTester.Name,
		"app.kubernetes.io/component":  "microservice",
		"app.kubernetes.io/part-of":    "cluster-tester",
		"app.kubernetes.io/managed-by": "cluster-tester-operator",
	}

	serviceType := corev1.ServiceTypeClusterIP
	if clusterTester.Spec.Global.ServiceType != "" {
		serviceType = corev1.ServiceType(clusterTester.Spec.Global.ServiceType)
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:     serviceType,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       8080,
					TargetPort: intstr.FromInt(8080),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
}

func (r *ClusterTesterReconciler) reconcileDatabase(ctx context.Context, clusterTester *clusterv1.ClusterTester) error {
	logger := log.FromContext(ctx)

	namespace := clusterTester.Namespace
	if clusterTester.Spec.Global.Namespace != "" {
		namespace = clusterTester.Spec.Global.Namespace
	}

	dbConfig := clusterTester.Spec.Database
	if dbConfig.Type == "" {
		dbConfig.Type = "mysql"
	}
	if dbConfig.Image == "" {
		dbConfig.Image = "mysql"
	}
	if dbConfig.Tag == "" {
		dbConfig.Tag = "8.0"
	}
	if dbConfig.StorageSize == "" {
		dbConfig.StorageSize = "10Gi"
	}

	// Create PVC
	pvc := r.createDatabasePVC(clusterTester, dbConfig, namespace)
	if err := controllerutil.SetControllerReference(clusterTester, pvc, r.Scheme); err != nil {
		return err
	}

	foundPVC := &corev1.PersistentVolumeClaim{}
	err := r.Get(ctx, types.NamespacedName{Name: pvc.Name, Namespace: pvc.Namespace}, foundPVC)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating PVC", "pvc", pvc.Name)
		if err = r.Create(ctx, pvc); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Create deployment
	deployment := r.createDatabaseDeployment(clusterTester, dbConfig, namespace)
	if err := controllerutil.SetControllerReference(clusterTester, deployment, r.Scheme); err != nil {
		return err
	}

	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating database deployment", "deployment", deployment.Name)
		if err = r.Create(ctx, deployment); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Create service
	service := r.createDatabaseService(clusterTester, namespace)
	if err := controllerutil.SetControllerReference(clusterTester, service, r.Scheme); err != nil {
		return err
	}

	foundService := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating database service", "service", service.Name)
		if err = r.Create(ctx, service); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func (r *ClusterTesterReconciler) createDatabasePVC(clusterTester *clusterv1.ClusterTester, dbConfig clusterv1.DatabaseConfig, namespace string) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"app":                          "mysql",
		"app.kubernetes.io/name":       "mysql",
		"app.kubernetes.io/instance":   clusterTester.Name,
		"app.kubernetes.io/component":  "database",
		"app.kubernetes.io/part-of":    "cluster-tester",
		"app.kubernetes.io/managed-by": "cluster-tester-operator",
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql-pvc",
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: *parseQuantity(dbConfig.StorageSize),
				},
			},
		},
	}

	if dbConfig.StorageClass != "" {
		pvc.Spec.StorageClassName = &dbConfig.StorageClass
	}

	return pvc
}

func (r *ClusterTesterReconciler) createDatabaseDeployment(clusterTester *clusterv1.ClusterTester, dbConfig clusterv1.DatabaseConfig, namespace string) *appsv1.Deployment {
	labels := map[string]string{
		"app":                          "mysql",
		"app.kubernetes.io/name":       "mysql",
		"app.kubernetes.io/instance":   clusterTester.Name,
		"app.kubernetes.io/component":  "database",
		"app.kubernetes.io/part-of":    "cluster-tester",
		"app.kubernetes.io/managed-by": "cluster-tester-operator",
	}

	replicas := int32(1)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql",
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mysql",
							Image: fmt.Sprintf("%s:%s", dbConfig.Image, dbConfig.Tag),
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: "rootpassword",
								},
								{
									Name:  "MYSQL_DATABASE",
									Value: "electronics-store",
								},
								{
									Name:  "MYSQL_USER",
									Value: "admin",
								},
								{
									Name:  "MYSQL_PASSWORD",
									Value: "password123",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "mysql",
									ContainerPort: 3306,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "mysql-storage",
									MountPath: "/var/lib/mysql",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "mysql-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "mysql-pvc",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *ClusterTesterReconciler) createDatabaseService(clusterTester *clusterv1.ClusterTester, namespace string) *corev1.Service {
	labels := map[string]string{
		"app":                          "mysql",
		"app.kubernetes.io/name":       "mysql",
		"app.kubernetes.io/instance":   clusterTester.Name,
		"app.kubernetes.io/component":  "database",
		"app.kubernetes.io/part-of":    "cluster-tester",
		"app.kubernetes.io/managed-by": "cluster-tester-operator",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql",
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "mysql",
					Port:       3306,
					TargetPort: intstr.FromInt(3306),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
}

func (r *ClusterTesterReconciler) updateStatusError(ctx context.Context, clusterTester *clusterv1.ClusterTester, reason string, err error) (ctrl.Result, error) {
	clusterTester.Status.Phase = "Failed"

	errorCondition := metav1.Condition{
		Type:    "Ready",
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: err.Error(),
	}
	meta.SetStatusCondition(&clusterTester.Status.Conditions, errorCondition)

	if statusErr := r.Status().Update(ctx, clusterTester); statusErr != nil {
		return ctrl.Result{}, statusErr
	}

	return ctrl.Result{RequeueAfter: time.Minute * 2}, err
}

// Helper function to parse quantity
func parseQuantity(s string) *resource.Quantity {
	q := resource.MustParse(s)
	return &q
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterTesterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1.ClusterTester{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Complete(r)
}
