#!/bin/bash

# Script to generate Swagger documentation for all services
# This script should be run from the root directory of the cluster-tester project

set -e

echo "🚀 Generating Swagger documentation for all services..."

services=("coffee-shop" "pet-store" "restaurant" "college-admission" "electronics-store" "electronics-store-tracing")

for service in "${services[@]}"; do
    echo "📝 Processing $service..."
    
    if [ -d "$service" ]; then
        cd "$service"
        
        # Check if swag is installed
        if ! command -v swag &> /dev/null; then
            echo "⚠️  swag command not found. Installing..."
            go install github.com/swaggo/swag/cmd/swag@latest
        fi
        
        # Generate Swagger docs
        swag init
        
        echo "✅ Generated documentation for $service"
        cd ..
    else
        echo "❌ Directory $service not found"
    fi
done

echo "🎉 All Swagger documentation generated successfully!"
echo ""
echo "📖 To view the documentation for any service:"
echo "   1. Start the service: cd <service-name> && go run main.go"
echo "   2. Open browser: http://localhost:8080/swagger/"
echo ""
echo "🔧 To register with Backstage.io:"
echo "   Add this URL to your Backstage catalog:"
echo "   https://github.com/your-org/cluster-tester/blob/main/catalog-info.yaml"
