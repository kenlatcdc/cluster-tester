#!/bin/bash

# Validation script for cluster-tester services
# This script builds and validates all services and the operator

set -e

echo "🔍 Validating Cluster Tester Setup..."
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        exit 1
    fi
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Check prerequisites
print_info "Checking prerequisites..."

# Check Go installation
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    print_status 0 "Go installed: $GO_VERSION"
else
    print_status 1 "Go not found"
fi

# Check Docker installation
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    print_status 0 "Docker installed: $DOCKER_VERSION"
else
    print_status 1 "Docker not found"
fi

# Check swag installation
if command -v swag &> /dev/null; then
    print_status 0 "Swag tool available"
else
    print_info "Installing swag tool..."
    go install github.com/swaggo/swag/cmd/swag@latest
    print_status $? "Swag tool installation"
fi

echo ""
print_info "Validating services..."

SERVICES=("coffee-shop" "pet-store" "restaurant" "college-admission" "electronics-store" "electronics-store-tracing")

for service in "${SERVICES[@]}"; do
    echo ""
    print_info "Validating $service..."
    
    cd "$service"
    
    # Check if go.mod exists
    if [ -f "go.mod" ]; then
        print_status 0 "$service: go.mod exists"
    else
        print_status 1 "$service: go.mod missing"
    fi
    
    # Check if main.go exists
    if [ -f "main.go" ]; then
        print_status 0 "$service: main.go exists"
    else
        print_status 1 "$service: main.go missing"
    fi
    
    # Check if catalog-info.yaml exists
    if [ -f "catalog-info.yaml" ]; then
        print_status 0 "$service: catalog-info.yaml exists"
    else
        print_status 1 "$service: catalog-info.yaml missing"
    fi
    
    # Check if docs directory exists
    if [ -d "docs" ]; then
        print_status 0 "$service: docs directory exists"
    else
        print_status 1 "$service: docs directory missing"
    fi
    
    # Clean dependencies
    print_info "$service: Cleaning Go modules..."
    go mod tidy
    print_status $? "$service: Go module cleanup"
    
    # Generate Swagger docs
    print_info "$service: Generating Swagger documentation..."
    swag init --parseDependency --parseInternal
    print_status $? "$service: Swagger doc generation"
    
    # Build service
    print_info "$service: Building service..."
    go build -o "bin/$service" main.go
    print_status $? "$service: Build successful"
    
    # Run tests if they exist
    if [ -d "tests" ]; then
        print_info "$service: Running tests..."
        go test ./tests/... -v
        print_status $? "$service: Tests passed"
    fi
    
    cd ..
done

echo ""
print_info "Validating cluster-operator..."

cd cluster-operator

# Check operator files
if [ -f "go.mod" ]; then
    print_status 0 "Operator: go.mod exists"
else
    print_status 1 "Operator: go.mod missing"
fi

if [ -f "cmd/main.go" ]; then
    print_status 0 "Operator: main.go exists"
else
    print_status 1 "Operator: main.go missing"
fi

if [ -f "Makefile" ]; then
    print_status 0 "Operator: Makefile exists"
else
    print_status 1 "Operator: Makefile missing"
fi

# Clean operator dependencies
print_info "Operator: Cleaning Go modules..."
go mod tidy
print_status $? "Operator: Go module cleanup"

# Build operator
print_info "Operator: Building operator..."
go build -o bin/manager cmd/main.go
print_status $? "Operator: Build successful"

# Check CRD files
if [ -f "config/crd/bases/cluster.cdcent.io_clustertesters.yaml" ]; then
    print_status 0 "Operator: CRD file exists"
else
    print_status 1 "Operator: CRD file missing"
fi

# Check sample files
if [ -f "config/samples/cluster_v1_clustertester.yaml" ]; then
    print_status 0 "Operator: Sample CR exists"
else
    print_status 1 "Operator: Sample CR missing"
fi

cd ..

echo ""
print_info "Validating documentation..."

# Check root files
if [ -f "catalog-info.yaml" ]; then
    print_status 0 "Root catalog-info.yaml exists"
else
    print_status 1 "Root catalog-info.yaml missing"
fi

if [ -f "INSTRUMENTATION.md" ]; then
    print_status 0 "INSTRUMENTATION.md exists"
else
    print_status 1 "INSTRUMENTATION.md missing"
fi

if [ -f "SETUP.md" ]; then
    print_status 0 "SETUP.md exists"
else
    print_status 1 "SETUP.md missing"
fi

if [ -f "generate-docs.sh" ]; then
    print_status 0 "generate-docs.sh exists"
else
    print_status 1 "generate-docs.sh missing"
fi

if [ -f "generate-docs.bat" ]; then
    print_status 0 "generate-docs.bat exists"
else
    print_status 1 "generate-docs.bat missing"
fi

echo ""
echo -e "${GREEN}🎉 All validations completed successfully!${NC}"
echo ""
echo "Next steps:"
echo "1. Review generated Swagger documentation at http://localhost:PORT/swagger/index.html"
echo "2. Import services into Backstage using the catalog-info.yaml files"
echo "3. Deploy to Kubernetes using the cluster-operator"
echo "4. Run individual services for testing and development"
echo ""
echo "For detailed instructions, see SETUP.md"
