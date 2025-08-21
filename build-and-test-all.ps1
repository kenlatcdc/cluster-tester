# Build and Test All Applications Script
Write-Host "=== Building and Testing All Cluster-Tester Applications ===" -ForegroundColor Cyan
Write-Host

$apps = @("coffee-shop", "pet-store", "restaurant", "college-admission", "electronics-store", "electronics-store-tracing", "cluster-operator")
$results = @{}

foreach ($app in $apps) {
    Write-Host "Processing $app..." -ForegroundColor Yellow
    
    if (Test-Path $app) {
        Push-Location $app
        
        try {
            # Clean and tidy dependencies
            Write-Host "  - Running go mod tidy"
            go mod tidy 2>$null
            
            # Build the application
            Write-Host "  - Building application"
            $buildResult = go build 2>&1
            if ($LASTEXITCODE -eq 0) {
                Write-Host "  - Build: SUCCESS" -ForegroundColor Green
                $results[$app] = @{ Build = "SUCCESS" }
            } else {
                Write-Host "  - Build: FAILED" -ForegroundColor Red
                Write-Host "    Error: $buildResult" -ForegroundColor Red
                $results[$app] = @{ Build = "FAILED"; BuildError = $buildResult }
            }
            
            # Run tests if they exist
            if (Test-Path "tests" -Or (Get-ChildItem -Name "*_test.go" -ErrorAction SilentlyContinue)) {
                Write-Host "  - Running tests"
                $testResult = go test ./... 2>&1
                if ($LASTEXITCODE -eq 0) {
                    Write-Host "  - Tests: SUCCESS" -ForegroundColor Green
                    $results[$app].Test = "SUCCESS"
                } else {
                    Write-Host "  - Tests: FAILED" -ForegroundColor Red
                    Write-Host "    Error: $testResult" -ForegroundColor Red
                    $results[$app].Test = "FAILED"
                    $results[$app].TestError = $testResult
                }
            } else {
                Write-Host "  - Tests: SKIPPED (no tests found)" -ForegroundColor Yellow
                $results[$app].Test = "SKIPPED"
            }
            
        } catch {
            Write-Host "  - ERROR: $_" -ForegroundColor Red
            $results[$app] = @{ Build = "ERROR"; Error = $_.Exception.Message }
        } finally {
            Pop-Location
        }
        
        Write-Host "  - $app completed" -ForegroundColor Green
    } else {
        Write-Host "  - $app directory not found" -ForegroundColor Red
        $results[$app] = @{ Build = "NOT_FOUND" }
    }
    
    Write-Host
}

# Summary
Write-Host "=== BUILD AND TEST SUMMARY ===" -ForegroundColor Cyan
Write-Host
foreach ($app in $apps) {
    $result = $results[$app]
    $buildStatus = $result.Build
    $testStatus = if ($result.Test) { $result.Test } else { "N/A" }
    
    $buildColor = switch ($buildStatus) {
        "SUCCESS" { "Green" }
        "FAILED" { "Red" }
        "ERROR" { "Red" }
        "NOT_FOUND" { "Red" }
        default { "Yellow" }
    }
    
    $testColor = switch ($testStatus) {
        "SUCCESS" { "Green" }
        "FAILED" { "Red" }
        "SKIPPED" { "Yellow" }
        "N/A" { "Gray" }
        default { "Yellow" }
    }
    
    Write-Host "${app}:" -ForegroundColor White
    Write-Host "  Build: $buildStatus" -ForegroundColor $buildColor
    Write-Host "  Test:  $testStatus" -ForegroundColor $testColor
    Write-Host
}

# Count successes
$successfulBuilds = ($results.Values | Where-Object { $_.Build -eq "SUCCESS" }).Count
$totalApps = $apps.Count

Write-Host "=== FINAL RESULTS ===" -ForegroundColor Cyan
Write-Host "Successful builds: $successfulBuilds/$totalApps" -ForegroundColor $(if ($successfulBuilds -eq $totalApps) { "Green" } else { "Yellow" })

if ($successfulBuilds -eq $totalApps) {
    Write-Host "🎉 All applications built successfully!" -ForegroundColor Green
} else {
    Write-Host "⚠️  Some applications failed to build. Check the details above." -ForegroundColor Yellow
}

Write-Host
Write-Host "=== NEXT STEPS ===" -ForegroundColor Cyan
Write-Host "1. All applications have been updated to use native OpenAPI 3.0 specification"
Write-Host "2. Go version updated to 1.23 for all applications"
Write-Host "3. Swaggo dependencies and generated docs have been removed"
Write-Host "4. New endpoints available: /openapi.json and /docs"
Write-Host "5. Modern Swagger UI served from CDN"
Write-Host
Write-Host "You can now:"
Write-Host "- Run any application: go run main.go"
Write-Host "- View API docs at: http://localhost:8080/docs"
Write-Host "- Get OpenAPI spec at: http://localhost:8080/openapi.json"
Write-Host "- Deploy using the cluster-operator"
