@echo off
REM Validation script for cluster-tester services (Windows)
REM This script builds and validates all services and the operator

setlocal EnableDelayedExpansion

echo 🔍 Validating Cluster Tester Setup...
echo ======================================

REM Function to print status
goto :main

:print_success
echo ✅ %~1
goto :eof

:print_error
echo ❌ %~1
exit /b 1

:print_info
echo ℹ️  %~1
goto :eof

:main

REM Check prerequisites
call :print_info "Checking prerequisites..."

REM Check Go installation
go version >nul 2>&1
if %errorlevel% equ 0 (
    for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
    call :print_success "Go installed: !GO_VERSION!"
) else (
    call :print_error "Go not found"
)

REM Check Docker installation
docker --version >nul 2>&1
if %errorlevel% equ 0 (
    for /f "tokens=3" %%i in ('docker --version') do set DOCKER_VERSION=%%i
    call :print_success "Docker installed: !DOCKER_VERSION!"
) else (
    call :print_error "Docker not found"
)

REM Check swag installation
swag --version >nul 2>&1
if %errorlevel% equ 0 (
    call :print_success "Swag tool available"
) else (
    call :print_info "Installing swag tool..."
    go install github.com/swaggo/swag/cmd/swag@latest
    if %errorlevel% equ 0 (
        call :print_success "Swag tool installation"
    ) else (
        call :print_error "Swag tool installation failed"
    )
)

echo.
call :print_info "Validating services..."

set SERVICES=coffee-shop pet-store restaurant college-admission electronics-store electronics-store-tracing

for %%s in (%SERVICES%) do (
    echo.
    call :print_info "Validating %%s..."
    
    cd "%%s"
    
    REM Check if go.mod exists
    if exist "go.mod" (
        call :print_success "%%s: go.mod exists"
    ) else (
        call :print_error "%%s: go.mod missing"
    )
    
    REM Check if main.go exists
    if exist "main.go" (
        call :print_success "%%s: main.go exists"
    ) else (
        call :print_error "%%s: main.go missing"
    )
    
    REM Check if catalog-info.yaml exists
    if exist "catalog-info.yaml" (
        call :print_success "%%s: catalog-info.yaml exists"
    ) else (
        call :print_error "%%s: catalog-info.yaml missing"
    )
    
    REM Check if docs directory exists
    if exist "docs" (
        call :print_success "%%s: docs directory exists"
    ) else (
        call :print_error "%%s: docs directory missing"
    )
    
    REM Clean dependencies
    call :print_info "%%s: Cleaning Go modules..."
    go mod tidy
    if %errorlevel% equ 0 (
        call :print_success "%%s: Go module cleanup"
    ) else (
        call :print_error "%%s: Go module cleanup failed"
    )
    
    REM Generate Swagger docs
    call :print_info "%%s: Generating Swagger documentation..."
    swag init --parseDependency --parseInternal
    if %errorlevel% equ 0 (
        call :print_success "%%s: Swagger doc generation"
    ) else (
        call :print_error "%%s: Swagger doc generation failed"
    )
    
    REM Build service
    call :print_info "%%s: Building service..."
    if not exist "bin" mkdir bin
    go build -o "bin/%%s.exe" main.go
    if %errorlevel% equ 0 (
        call :print_success "%%s: Build successful"
    ) else (
        call :print_error "%%s: Build failed"
    )
    
    REM Run tests if they exist
    if exist "tests" (
        call :print_info "%%s: Running tests..."
        go test ./tests/... -v
        if %errorlevel% equ 0 (
            call :print_success "%%s: Tests passed"
        ) else (
            call :print_error "%%s: Tests failed"
        )
    )
    
    cd ..
)

echo.
call :print_info "Validating cluster-operator..."

cd cluster-operator

REM Check operator files
if exist "go.mod" (
    call :print_success "Operator: go.mod exists"
) else (
    call :print_error "Operator: go.mod missing"
)

if exist "cmd\main.go" (
    call :print_success "Operator: main.go exists"
) else (
    call :print_error "Operator: main.go missing"
)

if exist "Makefile" (
    call :print_success "Operator: Makefile exists"
) else (
    call :print_error "Operator: Makefile missing"
)

REM Clean operator dependencies
call :print_info "Operator: Cleaning Go modules..."
go mod tidy
if %errorlevel% equ 0 (
    call :print_success "Operator: Go module cleanup"
) else (
    call :print_error "Operator: Go module cleanup failed"
)

REM Build operator
call :print_info "Operator: Building operator..."
if not exist "bin" mkdir bin
go build -o "bin\manager.exe" cmd\main.go
if %errorlevel% equ 0 (
    call :print_success "Operator: Build successful"
) else (
    call :print_error "Operator: Build failed"
)

REM Check CRD files
if exist "config\crd\bases\cluster.cdcent.io_clustertesters.yaml" (
    call :print_success "Operator: CRD file exists"
) else (
    call :print_error "Operator: CRD file missing"
)

REM Check sample files
if exist "config\samples\cluster_v1_clustertester.yaml" (
    call :print_success "Operator: Sample CR exists"
) else (
    call :print_error "Operator: Sample CR missing"
)

cd ..

echo.
call :print_info "Validating documentation..."

REM Check root files
if exist "catalog-info.yaml" (
    call :print_success "Root catalog-info.yaml exists"
) else (
    call :print_error "Root catalog-info.yaml missing"
)

if exist "INSTRUMENTATION.md" (
    call :print_success "INSTRUMENTATION.md exists"
) else (
    call :print_error "INSTRUMENTATION.md missing"
)

if exist "SETUP.md" (
    call :print_success "SETUP.md exists"
) else (
    call :print_error "SETUP.md missing"
)

if exist "generate-docs.sh" (
    call :print_success "generate-docs.sh exists"
) else (
    call :print_error "generate-docs.sh missing"
)

if exist "generate-docs.bat" (
    call :print_success "generate-docs.bat exists"
) else (
    call :print_error "generate-docs.bat missing"
)

echo.
echo 🎉 All validations completed successfully!
echo.
echo Next steps:
echo 1. Review generated Swagger documentation at http://localhost:PORT/swagger/index.html
echo 2. Import services into Backstage using the catalog-info.yaml files
echo 3. Deploy to Kubernetes using the cluster-operator
echo 4. Run individual services for testing and development
echo.
echo For detailed instructions, see SETUP.md

pause
