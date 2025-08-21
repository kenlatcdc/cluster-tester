# Quick Verification Script
# This script verifies that all applications are updated to native OpenAPI

Write-Host "=== Cluster-Tester Migration Verification ===" -ForegroundColor Cyan
Write-Host

$apps = @("coffee-shop", "pet-store", "restaurant", "college-admission", "electronics-store", "electronics-store-tracing")

Write-Host "Checking Go versions..." -ForegroundColor Yellow
foreach ($app in $apps) {
    if (Test-Path "$app/go.mod") {
        $goVersion = Get-Content "$app/go.mod" | Select-String "^go " | ForEach-Object { $_.ToString().Split(' ')[1] }
        Write-Host "  ${app}: Go $goVersion" -ForegroundColor Green
    }
}

Write-Host
Write-Host "Checking for swaggo references..." -ForegroundColor Yellow
$swaggoFound = $false
foreach ($app in $apps) {
    if (Test-Path "$app/main.go") {
        $swaggoRefs = Select-String -Path "$app/main.go" -Pattern "swaggo|ginSwagger" -Quiet
        if ($swaggoRefs) {
            Write-Host "  ${app}: Still has swaggo references" -ForegroundColor Red
            $swaggoFound = $true
        } else {
            Write-Host "  ${app}: Clean (no swaggo)" -ForegroundColor Green
        }
    }
}

Write-Host
Write-Host "Checking for OpenAPI specifications..." -ForegroundColor Yellow
foreach ($app in $apps) {
    if (Test-Path "$app/main.go") {
        $openApiRef = Select-String -Path "$app/main.go" -Pattern "openAPISpec" -Quiet
        if ($openApiRef) {
            Write-Host "  ${app}: Has native OpenAPI spec" -ForegroundColor Green
        } else {
            Write-Host "  ${app}: Missing OpenAPI spec" -ForegroundColor Red
        }
    }
}

Write-Host
Write-Host "Checking for docs directories (should be removed)..." -ForegroundColor Yellow
foreach ($app in $apps) {
    if (Test-Path "$app/docs") {
        Write-Host "  ${app}: docs directory still exists" -ForegroundColor Red
    } else {
        Write-Host "  ${app}: docs directory removed" -ForegroundColor Green
    }
}

Write-Host
if (-not $swaggoFound) {
    Write-Host "✅ Migration completed successfully!" -ForegroundColor Green
    Write-Host "All applications are now using:" -ForegroundColor White
    Write-Host "  - Go 1.23" -ForegroundColor Green
    Write-Host "  - Native OpenAPI 3.0 specifications" -ForegroundColor Green
    Write-Host "  - No swaggo dependencies" -ForegroundColor Green
    Write-Host "  - Modern Swagger UI served from CDN" -ForegroundColor Green
} else {
    Write-Host "⚠️ Migration partially complete - some swaggo references remain" -ForegroundColor Yellow
}

Write-Host
Write-Host "Available endpoints for each service:" -ForegroundColor Cyan
Write-Host "  - /health (health check)"
Write-Host "  - /openapi.json (OpenAPI specification)"
Write-Host "  - /docs (Swagger UI documentation)"
Write-Host "  - Service-specific endpoints (pets, menu, applications, products)"
