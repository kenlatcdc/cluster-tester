#!/bin/bash

echo "=== Building All Cluster-Tester Applications ==="
echo

apps=("coffee-shop" "pet-store" "restaurant" "college-admission" "electronics-store" "electronics-store-tracing" "cluster-operator")

for app in "${apps[@]}"; do
    echo "Building $app..."
    if [ -d "$app" ]; then
        cd "$app" || exit 1
        
        echo "  - Running go mod tidy"
        go mod tidy
        
        echo "  - Building application"
        if go build; then
            echo "  - Build: SUCCESS"
        else
            echo "  - Build: FAILED"
        fi
        
        echo "  - Running tests"
        if go test ./... 2>/dev/null; then
            echo "  - Tests: SUCCESS"
        else
            echo "  - Tests: FAILED or SKIPPED"
        fi
        
        cd .. || exit 1
        echo "  - $app completed"
    else
        echo "  - $app directory not found"
    fi
    echo
done

echo "=== Build process completed ==="
