@echo off
REM Script to generate Swagger documentation for all services (Windows)
REM This script should be run from the root directory of the cluster-tester project

echo 🚀 Generating Swagger documentation for all services...

set services=coffee-shop pet-store restaurant college-admission electronics-store electronics-store-tracing

for %%s in (%services%) do (
    echo 📝 Processing %%s...
    
    if exist "%%s" (
        cd "%%s"
        
        REM Check if swag is installed
        swag version >nul 2>&1
        if errorlevel 1 (
            echo ⚠️  swag command not found. Installing...
            go install github.com/swaggo/swag/cmd/swag@latest
        )
        
        REM Generate Swagger docs
        swag init
        
        echo ✅ Generated documentation for %%s
        cd ..
    ) else (
        echo ❌ Directory %%s not found
    )
)

echo 🎉 All Swagger documentation generated successfully!
echo.
echo 📖 To view the documentation for any service:
echo    1. Start the service: cd ^<service-name^> ^&^& go run main.go
echo    2. Open browser: http://localhost:8080/swagger/
echo.
echo 🔧 To register with Backstage.io:
echo    Add this URL to your Backstage catalog:
echo    https://github.com/your-org/cluster-tester/blob/main/catalog-info.yaml

pause
