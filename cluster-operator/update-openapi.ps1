# Script to update all cluster-tester applications to use native OpenAPI instead of swaggo
param(
    [switch]$UpdateGoVersion = $true
)

$ErrorActionPreference = "Stop"

Write-Host "=== Updating Cluster-Tester Applications ===" -ForegroundColor Cyan
Write-Host

# Define the applications to update
$applications = @(
    "pet-store",
    "restaurant", 
    "college-admission",
    "electronics-store",
    "electronics-store-tracing"
)

# Coffee-shop is already updated, so we skip it

foreach ($app in $applications) {
    $appPath = "..\$app"
    
    Write-Host "Updating $app..." -ForegroundColor Green
    
    if (-not (Test-Path $appPath)) {
        Write-Warning "Application directory $appPath not found, skipping..."
        continue
    }
    
    # Update go.mod to remove swaggo dependencies and update Go version
    $goModPath = "$appPath\go.mod"
    if (Test-Path $goModPath) {
        Write-Host "  - Updating go.mod"
        
        # Read current go.mod
        $content = Get-Content $goModPath
        
        # Update Go version and remove swaggo dependencies
        $newContent = @()
        $inRequireBlock = $false
        
        foreach ($line in $content) {
            if ($UpdateGoVersion -and $line.StartsWith("go ")) {
                $newContent += "go 1.23"
            }
            elseif ($line.Trim() -eq "require (") {
                $newContent += $line
                $inRequireBlock = $true
            }
            elseif ($inRequireBlock -and $line.Trim() -eq ")") {
                $newContent += $line
                $inRequireBlock = $false
            }
            elseif ($inRequireBlock -and ($line -match "swaggo|swag")) {
                # Skip swaggo dependencies
                continue
            }
            elseif ($line.StartsWith("toolchain ")) {
                # Remove toolchain specification
                continue
            }
            else {
                $newContent += $line
            }
        }
        
        $newContent | Set-Content $goModPath
        
        # Run go mod tidy
        Push-Location $appPath
        try {
            go mod tidy
            Write-Host "  - go mod tidy completed"
        }
        catch {
            Write-Warning "  - go mod tidy failed: $_"
        }
        finally {
            Pop-Location
        }
    }
    
    # Remove docs directory if it exists
    $docsPath = "$appPath\docs"
    if (Test-Path $docsPath) {
        Write-Host "  - Removing docs directory"
        Remove-Item -Recurse -Force $docsPath
    }
    
    Write-Host "  - $app updated successfully" -ForegroundColor Green
    Write-Host
}

Write-Host "=== Update Complete ===" -ForegroundColor Cyan
Write-Host
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Update main.go files to remove swaggo imports and godoc comments"
Write-Host "2. Add native OpenAPI specification to each service"
Write-Host "3. Replace swagger routes with /openapi.json and /docs endpoints"
Write-Host "4. Test each service to ensure it builds and runs correctly"
Write-Host
Write-Host "Note: The coffee-shop application has already been updated as a reference." -ForegroundColor Green
