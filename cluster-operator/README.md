# Cluster Tester Operator

A Kubernetes operator that deploys and manages all the cluster-tester microservices using a single Custom Resource.

## Overview

The Cluster Tester Operator simplifies the deployment and management of the cluster-tester microservices suite. Instead of managing individual deployments, services, and configurations, you can define a single `ClusterTester` custom resource that describes your desired state, and the operator will handle the rest.

## Features

- **Declarative Configuration**: Define all services in a single YAML file
- **Automatic Service Discovery**: Services are automatically configured with proper networking
- **Database Management**: Optional MySQL database deployment and configuration
- **Health Monitoring**: Built-in health checks and status reporting
- **Resource Management**: Configure CPU/memory limits and requests per service
- **Scalability**: Set replica counts for each service independently
- **Observability**: Status tracking for all deployed services

## Supported Services

The operator can deploy and manage the following services:

1. **Coffee Shop API** - Coffee menu management service
2. **Pet Store API** - Pet inventory management service
3. **Restaurant API** - Restaurant menu management service
4. **College Admission API** - Student application management service
5. **Electronics Store API** - Electronics inventory with database
6. **Electronics Store Tracing API** - Electronics inventory with distributed tracing
7. **MySQL Database** - Shared database for services that require persistence

## Installation

### Prerequisites

- Kubernetes cluster (v1.20+)
- kubectl configured to access your cluster
- Cluster admin permissions

### Deploy the Operator

1. **Install the CRDs:**
   ```bash
   kubectl apply -f config/crd/bases/cluster.cdcent.io_clustertesters.yaml
   ```

2. **Create the operator namespace:**
   ```bash
   kubectl create namespace cluster-tester-operator-system
   ```

3. **Apply RBAC configuration:**
   ```bash
   kubectl apply -f config/rbac/
   ```

4. **Deploy the operator:**
   ```bash
   kubectl apply -f config/manager/manager.yaml
   ```

### Alternative: Use Kustomize

```bash
# Install everything at once
kubectl apply -k config/default
```

## Usage

### Basic Deployment

Create a `ClusterTester` resource to deploy all services:

```yaml
apiVersion: cluster.cdcent.io/v1
kind: ClusterTester
metadata:
  name: my-cluster-tester
  namespace: default
spec:
  # Enable specific services
  coffeeShop:
    enabled: true
    replicas: 2
  
  petStore:
    enabled: true
  
  restaurant:
    enabled: true
  
  # Database for electronics services
  database:
    enabled: true
    storageSize: "10Gi"
  
  electronicsStore:
    enabled: true
  
  # Global configuration
  global:
    serviceType: NodePort
    imagePullPolicy: IfNotPresent
```

Apply the configuration:
```bash
kubectl apply -f my-cluster-tester.yaml
```

### Minimal Deployment

For development or testing, deploy only essential services:

```yaml
apiVersion: cluster.cdcent.io/v1
kind: ClusterTester
metadata:
  name: minimal-setup
spec:
  coffeeShop:
    enabled: true
  
  petStore:
    enabled: true
  
  # Disable database-dependent services
  database:
    enabled: false
```

### Production Deployment

For production environments with resource limits and multiple replicas:

```yaml
apiVersion: cluster.cdcent.io/v1
kind: ClusterTester
metadata:
  name: production-cluster-tester
  namespace: production
spec:
  coffeeShop:
    enabled: true
    replicas: 3
    image: my-registry.com/coffee-shop
    tag: v1.2.0
    resources:
      requests:
        cpu: "200m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "1Gi"
  
  petStore:
    enabled: true
    replicas: 2
    image: my-registry.com/pet-store
    tag: v1.2.0
    resources:
      requests:
        cpu: "200m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "1Gi"
  
  database:
    enabled: true
    storageSize: "100Gi"
    storageClass: "ssd"
  
  global:
    namespace: cluster-tester-prod
    serviceType: ClusterIP
    imagePullPolicy: Always
```

## Configuration Reference

### Service Configuration

Each service supports the following configuration options:

```yaml
serviceName:
  enabled: boolean          # Whether to deploy this service
  replicas: integer         # Number of replicas (default: 1)
  image: string            # Container image name
  tag: string              # Image tag (default: "latest")
  resources:               # Resource requirements
    requests:
      cpu: string          # CPU request (e.g., "100m")
      memory: string       # Memory request (e.g., "128Mi")
    limits:
      cpu: string          # CPU limit (e.g., "500m")
      memory: string       # Memory limit (e.g., "512Mi")
```

### Database Configuration

```yaml
database:
  enabled: boolean         # Whether to deploy MySQL
  type: string            # Database type (default: "mysql")
  image: string           # Database image (default: "mysql")
  tag: string             # Database tag (default: "8.0")
  storageSize: string     # Storage size (default: "10Gi")
  storageClass: string    # Storage class for PVC
```

### Global Configuration

```yaml
global:
  namespace: string        # Target namespace for deployments
  imagePullPolicy: string  # Image pull policy (IfNotPresent, Always, Never)
  serviceType: string      # Service type (ClusterIP, NodePort, LoadBalancer)
  ingressEnabled: boolean  # Whether to create ingress resources
  ingressHost: string      # Base hostname for ingress
```

## Monitoring and Status

### Check Deployment Status

```bash
# View ClusterTester resources
kubectl get clustertesters

# Get detailed status
kubectl describe clustertester my-cluster-tester

# Check operator logs
kubectl logs -n cluster-tester-operator-system deployment/cluster-tester-operator-controller-manager
```

### Service Status

The operator provides real-time status for each service:

```bash
kubectl get clustertester my-cluster-tester -o yaml
```

Status includes:
- Overall phase (Initializing, Ready, Failed)
- Individual service status (ready/not ready)
- Replica counts
- Service endpoints
- Error conditions

### Access Services

Once deployed, services are available at:

```bash
# List services
kubectl get services

# Port forward to access locally
kubectl port-forward svc/coffee-shop 8080:8080

# Access service
curl http://localhost:8080/health
curl http://localhost:8080/swagger/
```

## Troubleshooting

### Common Issues

1. **Operator not starting:**
   ```bash
   # Check operator logs
   kubectl logs -n cluster-tester-operator-system deployment/cluster-tester-operator-controller-manager
   
   # Verify RBAC permissions
   kubectl auth can-i create deployments --as=system:serviceaccount:cluster-tester-operator-system:cluster-tester-operator-controller-manager
   ```

2. **Services not deploying:**
   ```bash
   # Check ClusterTester status
   kubectl describe clustertester my-cluster-tester
   
   # Check individual deployments
   kubectl get deployments
   kubectl describe deployment coffee-shop
   ```

3. **Database connection issues:**
   ```bash
   # Check MySQL service
   kubectl get svc mysql
   kubectl logs deployment/mysql
   
   # Check PVC
   kubectl get pvc
   ```

### Debug Mode

Enable debug logging by setting environment variables in the operator deployment:

```yaml
env:
- name: LOG_LEVEL
  value: "debug"
```

## Development

### Building the Operator

```bash
# Build the binary
make build

# Build Docker image
make docker-build IMG=my-registry.com/cluster-tester-operator:v1.0.0

# Push to registry
make docker-push IMG=my-registry.com/cluster-tester-operator:v1.0.0
```

### Running Locally

```bash
# Install CRDs
make install

# Run operator locally
make run
```

### Testing

```bash
# Run unit tests
make test

# Run end-to-end tests
make test-e2e
```

## API Reference

### ClusterTester

| Field | Type | Description |
|-------|------|-------------|
| `spec.coffeeShop` | ServiceConfig | Coffee Shop service configuration |
| `spec.petStore` | ServiceConfig | Pet Store service configuration |
| `spec.restaurant` | ServiceConfig | Restaurant service configuration |
| `spec.collegeAdmission` | ServiceConfig | College Admission service configuration |
| `spec.electronicsStore` | ServiceConfig | Electronics Store service configuration |
| `spec.electronicsStoreTracing` | ServiceConfig | Electronics Store Tracing service configuration |
| `spec.database` | DatabaseConfig | Database configuration |
| `spec.global` | GlobalConfig | Global configuration options |

### ServiceConfig

| Field | Type | Description |
|-------|------|-------------|
| `enabled` | bool | Whether this service should be deployed |
| `replicas` | *int32 | Number of replicas |
| `image` | string | Container image name |
| `tag` | string | Image tag |
| `resources` | *ResourceRequirements | Resource requirements |

### DatabaseConfig

| Field | Type | Description |
|-------|------|-------------|
| `enabled` | bool | Whether to deploy the database |
| `type` | string | Database type |
| `image` | string | Database container image |
| `tag` | string | Database image tag |
| `storageSize` | string | Storage size for database |
| `storageClass` | string | Storage class |

### GlobalConfig

| Field | Type | Description |
|-------|------|-------------|
| `namespace` | string | Target namespace for deployments |
| `imagePullPolicy` | string | Image pull policy |
| `serviceType` | string | Default service type |
| `ingressEnabled` | bool | Whether to create ingress resources |
| `ingressHost` | string | Base host for ingress |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make test` and `make lint`
6. Submit a pull request

## License

Licensed under the Apache License, Version 2.0.
