#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Build and test all cluster-tester services and operator.

.DESCRIPTION
    This script automates the build, test, and documentation generation 
    process for all services in the cluster-tester repository.

.PARAMETER BuildOperator
    Whether to build the Kubernetes operator (requires make)

.PARAMETER RunTests
    Whether to run unit tests for all services

.PARAMETER GenerateDocs
    Whether to generate Swagger documentation

.PARAMETER BuildImages
    Whether to build Docker images

.EXAMPLE
    .\build-all.ps1 -GenerateDocs -RunTests
    
.EXAMPLE
    .\build-all.ps1 -BuildOperator -BuildImages
#>

param(
    [switch]$BuildOperator,
    [switch]$RunTests,
    [switch]$GenerateDocs,
    [switch]$BuildImages
)

# Color output functions
function Write-Success { param($Message) Write-Host $Message -ForegroundColor Green }
function Write-Info { param($Message) Write-Host $Message -ForegroundColor Blue }
function Write-Warning { param($Message) Write-Host $Message -ForegroundColor Yellow }
function Write-Error { param($Message) Write-Host $Message -ForegroundColor Red }

# Get script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptDir

Write-Info "🚀 Starting cluster-tester build process..."

# Services to build
$Services = @(
    "coffee-shop",
    "pet-store", 
    "restaurant",
    "college-admission",
    "electronics-store",
    "electronics-store-tracing"
)

# Generate documentation if requested
if ($GenerateDocs) {
    Write-Info "📚 Generating Swagger documentation..."
    try {
        & ".\generate-docs.bat"
        Write-Success "✅ Documentation generated successfully"
    }
    catch {
        Write-Error "❌ Failed to generate documentation: $_"
        exit 1
    }
}

# Build and test services
Write-Info "🔨 Building services..."
$BuildResults = @{}

foreach ($Service in $Services) {
    Write-Info "Building $Service..."
    
    try {
        Set-Location $Service
        
        # Ensure bin directory exists
        if (!(Test-Path "bin")) {
            New-Item -ItemType Directory -Name "bin" | Out-Null
        }
        
        # Build the service
        $BinaryName = "$Service.exe"
        $BinaryPath = Join-Path "bin" $BinaryName
        
        $BuildOutput = go build -o $BinaryPath main.go 2>&1
        if ($LASTEXITCODE -ne 0) {
            throw "Build failed: $BuildOutput"
        }
        
        Write-Success "✅ Built $Service"
        $BuildResults[$Service] = @{ Status = "Success"; Binary = $BinaryPath }
        
        # Run tests if requested
        if ($RunTests) {
            Write-Info "Running tests for $Service..."
            $TestOutput = go test ./tests/... 2>&1
            if ($LASTEXITCODE -eq 0) {
                Write-Success "✅ Tests passed for $Service"
                $BuildResults[$Service].Tests = "Passed"
            } else {
                Write-Warning "⚠️ Tests failed for $Service: $TestOutput"
                $BuildResults[$Service].Tests = "Failed"
            }
        }
        
        # Build Docker image if requested
        if ($BuildImages) {
            Write-Info "Building Docker image for $Service..."
            $ImageOutput = docker build -t "$Service`:latest" . 2>&1
            if ($LASTEXITCODE -eq 0) {
                Write-Success "✅ Docker image built for $Service"
                $BuildResults[$Service].Image = "Built"
            } else {
                Write-Warning "⚠️ Docker build failed for $Service: $ImageOutput"
                $BuildResults[$Service].Image = "Failed"
            }
        }
        
    }
    catch {
        Write-Error "❌ Failed to build $Service`: $_"
        $BuildResults[$Service] = @{ Status = "Failed"; Error = $_.Exception.Message }
    }
    finally {
        Set-Location $ScriptDir
    }
}

# Build operator if requested
if ($BuildOperator) {
    Write-Info "🔧 Building Kubernetes operator..."
    try {
        Set-Location "cluster-operator"
        
        # Check if make is available
        $MakeAvailable = Get-Command make -ErrorAction SilentlyContinue
        if (!$MakeAvailable) {
            Write-Warning "⚠️ 'make' command not found. Install make or use WSL/Git Bash"
            Write-Info "Alternative: Use 'go build -o bin/manager cmd/main.go'"
        } else {
            $MakeOutput = make build 2>&1
            if ($LASTEXITCODE -eq 0) {
                Write-Success "✅ Operator built successfully"
                $BuildResults["cluster-operator"] = @{ Status = "Success" }
            } else {
                throw "Make build failed: $MakeOutput"
            }
        }
    }
    catch {
        Write-Error "❌ Failed to build operator: $_"
        $BuildResults["cluster-operator"] = @{ Status = "Failed"; Error = $_.Exception.Message }
    }
    finally {
        Set-Location $ScriptDir
    }
}

# Print summary
Write-Info "`n📊 Build Summary:"
Write-Host "=" * 50

foreach ($Service in $Services) {
    $Result = $BuildResults[$Service]
    if ($Result.Status -eq "Success") {
        Write-Success "✅ $Service - Built successfully"
        if ($Result.Tests) { Write-Host "   Tests: $($Result.Tests)" }
        if ($Result.Image) { Write-Host "   Docker: $($Result.Image)" }
    } else {
        Write-Error "❌ $Service - Build failed"
        if ($Result.Error) { Write-Host "   Error: $($Result.Error)" -ForegroundColor Red }
    }
}

if ($BuildOperator -and $BuildResults["cluster-operator"]) {
    $OpResult = $BuildResults["cluster-operator"]
    if ($OpResult.Status -eq "Success") {
        Write-Success "✅ cluster-operator - Built successfully"
    } else {
        Write-Error "❌ cluster-operator - Build failed"
        if ($OpResult.Error) { Write-Host "   Error: $($OpResult.Error)" -ForegroundColor Red }
    }
}

Write-Info "`n🎉 Build process completed!"

# Quick start instructions
Write-Info "`n🚀 Quick Start:"
Write-Host "1. Start a service:" -ForegroundColor Yellow
Write-Host "   cd coffee-shop && .\bin\coffee-shop.exe" -ForegroundColor Gray
Write-Host "2. Test health endpoint:" -ForegroundColor Yellow  
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8080/health'" -ForegroundColor Gray
Write-Host "3. View Swagger docs:" -ForegroundColor Yellow
Write-Host "   Start-Process 'http://localhost:8080/swagger/index.html'" -ForegroundColor Gray

if ($BuildOperator) {
    Write-Host "4. Deploy operator:" -ForegroundColor Yellow
    Write-Host "   cd cluster-operator && make deploy" -ForegroundColor Gray
}
