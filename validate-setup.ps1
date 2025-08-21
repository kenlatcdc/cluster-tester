#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Validate cluster-tester setup and health.

.DESCRIPTION
    This script validates that all services are properly configured,
    can be built, and their endpoints are accessible.

.PARAMETER Port
    Starting port number for services (default: 8080)

.PARAMETER Timeout
    Timeout in seconds for health checks (default: 30)

.EXAMPLE
    .\validate-setup.ps1
    
.EXAMPLE
    .\validate-setup.ps1 -Port 9000 -Timeout 60
#>

param(
    [int]$Port = 8080,
    [int]$Timeout = 30
)

# Color output functions
function Write-Success { param($Message) Write-Host $Message -ForegroundColor Green }
function Write-Info { param($Message) Write-Host $Message -ForegroundColor Blue }
function Write-Warning { param($Message) Write-Host $Message -ForegroundColor Yellow }
function Write-Error { param($Message) Write-Host $Message -ForegroundColor Red }

# Get script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptDir

Write-Info "🔍 Starting cluster-tester validation..."

# Check prerequisites
Write-Info "`n📋 Checking prerequisites..."

# Check Go
try {
    $GoVersion = go version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "✅ Go: $GoVersion"
    } else {
        throw "Go not found"
    }
} catch {
    Write-Error "❌ Go: Not installed or not in PATH"
    exit 1
}

# Check Docker
try {
    $DockerVersion = docker --version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "✅ Docker: $DockerVersion"
    } else {
        throw "Docker not found"
    }
} catch {
    Write-Warning "⚠️ Docker: Not installed or not running"
}

# Check kubectl
try {
    $KubectlVersion = kubectl version --client --short 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "✅ kubectl: $KubectlVersion"
    } else {
        throw "kubectl not found"
    }
} catch {
    Write-Warning "⚠️ kubectl: Not installed (needed for operator deployment)"
}

# Check swag
try {
    $SwagVersion = swag --version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "✅ swag: $SwagVersion"
    } else {
        throw "swag not found"
    }
} catch {
    Write-Warning "⚠️ swag: Not installed (run 'go install github.com/swaggo/swag/cmd/swag@latest')"
}

# Services configuration
$Services = @(
    @{ Name = "coffee-shop"; Port = $Port },
    @{ Name = "pet-store"; Port = $Port + 1 },
    @{ Name = "restaurant"; Port = $Port + 2 },
    @{ Name = "college-admission"; Port = $Port + 3 },
    @{ Name = "electronics-store"; Port = $Port + 4 },
    @{ Name = "electronics-store-tracing"; Port = $Port + 5 }
)

# Validate service structure
Write-Info "`n📁 Validating project structure..."

foreach ($Service in $Services) {
    $ServicePath = $Service.Name
    $RequiredFiles = @("main.go", "go.mod", "catalog-info.yaml")
    
    if (Test-Path $ServicePath) {
        Write-Success "✅ $($Service.Name) directory exists"
        
        foreach ($File in $RequiredFiles) {
            $FilePath = Join-Path $ServicePath $File
            if (Test-Path $FilePath) {
                Write-Success "  ✅ $File exists"
            } else {
                Write-Error "  ❌ $File missing"
            }
        }
    } else {
        Write-Error "❌ $($Service.Name) directory missing"
    }
}

# Check cluster-operator
if (Test-Path "cluster-operator") {
    Write-Success "✅ cluster-operator directory exists"
    
    $OperatorFiles = @(
        "go.mod",
        "cmd/main.go", 
        "api/v1/clustertester_types.go",
        "internal/controller/clustertester_controller.go",
        "config/crd/bases/cluster.cdcent.io_clustertesters.yaml"
    )
    
    foreach ($File in $OperatorFiles) {
        $FilePath = Join-Path "cluster-operator" $File
        if (Test-Path $FilePath) {
            Write-Success "  ✅ $File exists"
        } else {
            Write-Error "  ❌ $File missing"
        }
    }
} else {
    Write-Error "❌ cluster-operator directory missing"
}

# Build validation
Write-Info "`n🔨 Validating builds..."

$BuildJobs = @()
foreach ($Service in $Services) {
    Write-Info "Building $($Service.Name)..."
    
    $Job = Start-Job -ScriptBlock {
        param($ServicePath, $ScriptDir)
        Set-Location $ScriptDir
        Set-Location $ServicePath
        
        try {
            # Ensure bin directory exists
            if (!(Test-Path "bin")) {
                New-Item -ItemType Directory -Name "bin" | Out-Null
            }
            
            $BinaryName = "$ServicePath.exe"
            $BinaryPath = Join-Path "bin" $BinaryName
            
            $Output = go build -o $BinaryPath main.go 2>&1
            return @{
                Success = $LASTEXITCODE -eq 0
                Output = $Output
                BinaryPath = $BinaryPath
            }
        }
        catch {
            return @{
                Success = $false
                Output = $_.Exception.Message
                BinaryPath = $null
            }
        }
    } -ArgumentList $Service.Name, $ScriptDir
    
    $BuildJobs += @{ Job = $Job; Service = $Service }
}

# Wait for builds and collect results
$BuildResults = @{}
foreach ($JobInfo in $BuildJobs) {
    $Result = Receive-Job -Job $JobInfo.Job -Wait
    $BuildResults[$JobInfo.Service.Name] = $Result
    
    if ($Result.Success) {
        Write-Success "✅ $($JobInfo.Service.Name) built successfully"
    } else {
        Write-Error "❌ $($JobInfo.Service.Name) build failed: $($Result.Output)"
    }
    
    Remove-Job -Job $JobInfo.Job
}

# Runtime validation (start services and test endpoints)
Write-Info "`n🚀 Testing service endpoints..."

$RunningProcesses = @()
foreach ($Service in $Services) {
    $BuildResult = $BuildResults[$Service.Name]
    
    if ($BuildResult.Success) {
        try {
            # Check if port is available
            $PortInUse = Get-NetTCPConnection -LocalPort $Service.Port -ErrorAction SilentlyContinue
            if ($PortInUse) {
                Write-Warning "⚠️ Port $($Service.Port) is already in use, skipping $($Service.Name)"
                continue
            }
            
            Write-Info "Starting $($Service.Name) on port $($Service.Port)..."
            
            # Start the service
            $ServicePath = Join-Path $Service.Name $BuildResult.BinaryPath
            $Process = Start-Process -FilePath $ServicePath -PassThru -WindowStyle Hidden
            $RunningProcesses += @{ Process = $Process; Service = $Service }
            
            Write-Success "✅ $($Service.Name) started (PID: $($Process.Id))"
            
        } catch {
            Write-Error "❌ Failed to start $($Service.Name): $_"
        }
    }
}

if ($RunningProcesses.Count -gt 0) {
    Write-Info "`n⏳ Waiting for services to initialize..."
    Start-Sleep -Seconds 5
    
    # Test health endpoints
    foreach ($ProcessInfo in $RunningProcesses) {
        $Service = $ProcessInfo.Service
        $HealthUrl = "http://localhost:$($Service.Port)/health"
        
        try {
            Write-Info "Testing health endpoint for $($Service.Name)..."
            
            $Response = Invoke-RestMethod -Uri $HealthUrl -TimeoutSec 10 -ErrorAction Stop
            Write-Success "✅ $($Service.Name) health check passed"
            Write-Host "  Response: $($Response | ConvertTo-Json -Compress)" -ForegroundColor Gray
            
            # Test Swagger endpoint
            $SwaggerUrl = "http://localhost:$($Service.Port)/swagger/doc.json"
            try {
                $SwaggerDoc = Invoke-RestMethod -Uri $SwaggerUrl -TimeoutSec 5 -ErrorAction Stop
                Write-Success "  ✅ Swagger documentation available"
            } catch {
                Write-Warning "  ⚠️ Swagger documentation not available"
            }
            
        } catch {
            Write-Error "❌ $($Service.Name) health check failed: $_"
        }
    }
    
    # Cleanup: Stop all started processes
    Write-Info "`n🧹 Cleaning up test processes..."
    foreach ($ProcessInfo in $RunningProcesses) {
        try {
            Stop-Process -Id $ProcessInfo.Process.Id -Force -ErrorAction SilentlyContinue
            Write-Success "✅ Stopped $($ProcessInfo.Service.Name)"
        } catch {
            Write-Warning "⚠️ Could not stop $($ProcessInfo.Service.Name)"
        }
    }
}

# Generate summary report
Write-Info "`n📊 Validation Summary:"
Write-Host "=" * 50

$TotalServices = $Services.Count
$SuccessfulBuilds = ($BuildResults.Values | Where-Object { $_.Success }).Count
$SuccessfulTests = $RunningProcesses.Count

Write-Host "Services: $TotalServices" -ForegroundColor Blue
Write-Host "Successful builds: $SuccessfulBuilds/$TotalServices" -ForegroundColor $(if ($SuccessfulBuilds -eq $TotalServices) { "Green" } else { "Yellow" })
Write-Host "Successful endpoint tests: $SuccessfulTests/$SuccessfulBuilds" -ForegroundColor $(if ($SuccessfulTests -eq $SuccessfulBuilds) { "Green" } else { "Yellow" })

if ($SuccessfulBuilds -eq $TotalServices -and $SuccessfulTests -eq $SuccessfulBuilds) {
    Write-Success "`n🎉 All validations passed! Your cluster-tester setup is ready."
    Write-Info "`n🚀 Next steps:"
    Write-Host "1. Generate documentation: .\generate-docs.bat" -ForegroundColor Yellow
    Write-Host "2. Build all services: .\build-all.ps1 -RunTests" -ForegroundColor Yellow
    Write-Host "3. Deploy operator: cd cluster-operator && make deploy" -ForegroundColor Yellow
} else {
    Write-Warning "`n⚠️ Some validations failed. Check the errors above."
    Write-Info "Refer to SETUP.md for troubleshooting guidance."
}
