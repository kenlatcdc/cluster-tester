# Cluster Tester - Complete Setup Guide

This repository contains a collection of microservices and a Kubernetes operator for deploying and managing them. All services have been instrumented with OpenAPI/Swagger documentation and configured for Backstage.io API registration.

## 🏗️ Architecture Overview

The repository contains:
- **6 Microservices**: coffee-shop, pet-store, restaurant, college-admission, electronics-store, electronics-store-tracing
- **Kubernetes Operator**: cluster-operator (manages deployment of all services)
- **API Documentation**: OpenAPI/Swagger integration for all services
- **Service Discovery**: Backstage.io catalog integration

## 📋 Prerequisites

- Go 1.22+
- Docker
- Kubernetes cluster (for operator deployment)
- kubectl configured
- PowerShell 5.1+ (Windows) or PowerShell Core 7+ (cross-platform)
- (Optional) kubebuilder for operator development
- (Optional) curl or use built-in `Invoke-RestMethod` for API testing

### Windows-Specific Setup

1. **Install Go**: Download from https://golang.org/dl/
2. **Install Docker Desktop**: Download from https://docker.com/products/docker-desktop
3. **Install kubectl**: 
   ```powershell
   # Using chocolatey
   choco install kubernetes-cli
   # or using winget
   winget install Kubernetes.kubectl
   ```
4. **Install swag for documentation generation**:
   ```powershell
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

## 🚀 Quick Start

### Automated Setup (Recommended)

Use the provided PowerShell scripts for automated setup and validation:

```powershell
# Validate your environment and dependencies
.\validate-setup.ps1

# Build all services with documentation and tests
.\build-all.ps1 -GenerateDocs -RunTests -BuildImages

# For operator development
.\build-all.ps1 -BuildOperator
```

### Manual Setup

### 1. Generate API Documentation

```powershell
# Generate Swagger docs for all services
.\generate-docs.bat  # Windows (PowerShell)
# or for cross-platform:
# .\generate-docs.sh  # Linux/Mac (if using WSL)
```

### 2. Build and Test Individual Services

Each service can be built and tested independently:

```powershell
cd coffee-shop
go build -o bin/coffee-shop.exe main.go
.\bin\coffee-shop.exe  # Service runs on port 8080

# Test endpoints (using curl or Invoke-RestMethod)
curl http://localhost:8080/health
# or
Invoke-RestMethod -Uri "http://localhost:8080/health"

# Open Swagger UI in browser
Start-Process "http://localhost:8080/swagger/index.html"
```

### 3. Deploy with Kubernetes Operator

```powershell
cd cluster-operator

# Build the operator
make build

# Deploy CRDs and operator to cluster
make deploy

# Create a cluster deployment
kubectl apply -f config/samples/cluster_v1_clustertester.yaml
```

## 📊 Service Endpoints

All services expose the following standard endpoints:

| Endpoint | Description |
|----------|-------------|
| `/health` | Health check endpoint |
| `/swagger/index.html` | Swagger UI documentation |
| `/swagger/doc.json` | OpenAPI JSON spec |

### Service-Specific Endpoints

#### Coffee Shop (Port 8080)
- `GET /menu` - Get coffee menu
- `POST /order` - Place an order
- `GET /order/{id}` - Get order status

#### Pet Store (Port 8081)
- `GET /pets` - List all pets
- `POST /pets` - Add a new pet
- `GET /pets/{id}` - Get pet details

#### Restaurant (Port 8082)
- `GET /menu` - Get restaurant menu
- `POST /reservation` - Make a reservation
- `GET /reservation/{id}` - Get reservation details

#### College Admission (Port 8083)
- `POST /application` - Submit application
- `GET /application/{id}` - Get application status
- `GET /requirements` - Get admission requirements

#### Electronics Store (Port 8084)
- `GET /products` - List products
- `POST /products` - Add new product
- `GET /products/{id}` - Get product details
- `POST /order` - Place an order

#### Electronics Store Tracing (Port 8085)
- Same as Electronics Store but with distributed tracing enabled

## 🎯 Backstage.io Integration

### Service Registration

Each service includes a `catalog-info.yaml` file for Backstage registration:

```yaml
apiVersion: backstage.io/v1alpha1
kind: API
metadata:
  name: coffee-shop-api
  description: Coffee Shop API for ordering beverages
spec:
  type: openapi
  lifecycle: production
  owner: platform-team
  definition:
    $text: http://localhost:8080/swagger/doc.json
```

### Bulk Registration

Import all services at once using the root catalog file:

```
https://github.com/cdcent/cluster-tester/blob/main/catalog-info.yaml
```

## 🔧 Kubernetes Operator

The cluster-operator provides declarative management of all services:

### Custom Resource Definition

```yaml
apiVersion: cluster.cdcent.io/v1
kind: ClusterTester
metadata:
  name: my-cluster
spec:
  services:
    coffeeShop:
      enabled: true
      replicas: 2
    petStore:
      enabled: true
      replicas: 1
    # ... other services
  database:
    enabled: true
    type: mysql
  global:
    namespace: cluster-tester
    imageRegistry: "your-registry.com"
```

### Operator Commands

```powershell
# Build operator
make build

# Run locally (development)
make run

# Build and push Docker image
make docker-build docker-push IMG=your-registry/cluster-operator:tag

# Deploy to cluster
make deploy IMG=your-registry/cluster-operator:tag

# Undeploy
make undeploy

# Run tests
make test
```

## 🧪 Testing

### Unit Tests

Each service includes unit tests:

```powershell
cd coffee-shop
go test ./tests/...
```

### Integration Testing

Use the provided test files to validate service functionality:

```powershell
# Test all services using PowerShell
$services = @("coffee-shop", "pet-store", "restaurant", "college-admission", "electronics-store", "electronics-store-tracing")
foreach ($service in $services) {
    Write-Host "Testing $service..."
    Set-Location $service
    go test ./tests/...
    Set-Location ..
}
```

### Operator Testing

```powershell
cd cluster-operator
make test
```

## 🐳 Docker Support

Each service includes a Dockerfile for containerization:

```powershell
# Build service image
cd coffee-shop
docker build -t coffee-shop:latest .

# Build operator image
cd cluster-operator
make docker-build IMG=cluster-operator:latest
```

## 📁 Directory Structure

```
cluster-tester/
├── coffee-shop/           # Coffee ordering service
├── pet-store/            # Pet management service
├── restaurant/           # Restaurant reservation service
├── college-admission/    # College application service
├── electronics-store/    # Electronics e-commerce service
├── electronics-store-tracing/  # Electronics service with tracing
├── cluster-operator/     # Kubernetes operator
├── catalog-info.yaml     # Backstage bulk registration
├── INSTRUMENTATION.md    # Technical documentation
├── SETUP.md             # This file
├── generate-docs.sh     # Documentation generation (Linux/Mac)
├── generate-docs.bat    # Documentation generation (Windows)
├── build-all.ps1        # Automated build script (PowerShell)
└── validate-setup.ps1   # Setup validation script (PowerShell)
```

## 🔍 Monitoring and Observability

### Health Checks

All services expose `/health` endpoints that return:

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "service": "coffee-shop"
}
```

### Swagger Documentation

Access interactive API documentation at `/swagger/index.html` for each service.

### Tracing (Electronics Store Tracing)

The electronics-store-tracing service includes distributed tracing capabilities for monitoring request flows.

## 🚨 Troubleshooting

### Common Issues

1. **Service won't start**: Check port availability and Go version
   ```powershell
   # Check if port is in use
   Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue
   # Check Go version
   go version
   ```

2. **Swagger docs not generating**: Ensure swag is installed
   ```powershell
   go install github.com/swaggo/swag/cmd/swag@latest
   # Verify installation
   swag --version
   ```

3. **Operator deployment fails**: Verify RBAC permissions and CRD installation
   ```powershell
   kubectl get crd clustertesters.cluster.cdcent.io
   kubectl auth can-i create clustertesters --as=system:serviceaccount:cluster-tester-system:cluster-tester-controller-manager
   ```

4. **Database connection issues**: Check MySQL deployment and connection strings
   ```powershell
   kubectl get pods -n cluster-tester -l app=mysql
   kubectl logs -n cluster-tester deployment/mysql
   ```

5. **Windows file path issues**: Use PowerShell's path handling
   ```powershell
   # Use Join-Path for cross-platform compatibility
   $binaryPath = Join-Path "bin" "coffee-shop.exe"
   ```

### Debug Commands

```powershell
# Check operator logs
kubectl logs -n cluster-tester-system deployment/cluster-tester-controller-manager

# Check service pods
kubectl get pods -n cluster-tester

# Check CRD status
kubectl get clustertesters -A

# Check Windows-specific issues
Get-Process -Name "*coffee*" -ErrorAction SilentlyContinue
Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue
```

## 🤝 Contributing

1. Follow Go conventions and add tests for new features
2. Update Swagger documentation when adding new endpoints
3. Ensure Backstage catalog files are updated for API changes
4. Test operator functionality with sample CRs

## 📜 License

See LICENSE file for details.
