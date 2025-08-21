# Cluster Tester - Instrumented APIs

This repository contains multiple microservices that have been instrumented with OpenAPI/Swagger documentation and configured for Backstage.io registration.

## Services

All services have been enhanced with:
- **OpenAPI/Swagger documentation** available at `/swagger/` endpoint
- **Health check endpoints** at `/health`
- **Backstage.io catalog integration** via `catalog-info.yaml` files

### Available Services

1. **Coffee Shop API** (`/coffee-shop`)
   - Port: 8080
   - Endpoints: `/coffees`, `/coffees/{id}`
   - Swagger: http://localhost:8080/swagger/
   - Health: http://localhost:8080/health

2. **Pet Store API** (`/pet-store`)
   - Port: 8080
   - Endpoints: `/pets`, `/pets/{id}`
   - Swagger: http://localhost:8080/swagger/
   - Health: http://localhost:8080/health

3. **Restaurant API** (`/restaurant`)
   - Port: 8080
   - Endpoints: `/menu`, `/menu/{id}`
   - Swagger: http://localhost:8080/swagger/
   - Health: http://localhost:8080/health

4. **College Admission API** (`/college-admission`)
   - Port: 8080
   - Endpoints: `/applications`, `/applications/{id}`
   - Swagger: http://localhost:8080/swagger/
   - Health: http://localhost:8080/health

5. **Electronics Store API** (`/electronics-store`)
   - Port: 8080
   - Endpoints: `/products`, `/products/{id}`
   - Swagger: http://localhost:8080/swagger/
   - Health: http://localhost:8080/health
   - Database: MySQL

6. **Electronics Store Tracing API** (`/electronics-store-tracing`)
   - Port: 8080
   - Endpoints: `/products`, `/products/{id}`
   - Swagger: http://localhost:8080/swagger/
   - Health: http://localhost:8080/health
   - Database: MySQL
   - Features: Distributed tracing

## Getting Started

### Prerequisites

- Go 1.22+
- Docker (for MySQL databases)
- Swagger CLI (optional, for generating docs)

### Running Individual Services

1. Navigate to any service directory:
   ```bash
   cd coffee-shop
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Generate Swagger documentation (optional):
   ```bash
   swag init
   ```

4. Run the service:
   ```bash
   go run main.go
   ```

5. Access the API:
   - Service: http://localhost:8080
   - Swagger UI: http://localhost:8080/swagger/
   - Health Check: http://localhost:8080/health

### Running with Docker

Each service includes a `Dockerfile` for containerization. To build and run:

```bash
cd coffee-shop
docker build -t coffee-shop .
docker run -p 8080:8080 coffee-shop
```

### Kubernetes Deployment

Deployment files are available in each service's `deploy/` directory:

```bash
kubectl apply -f coffee-shop/deploy/
```

## Backstage.io Integration

### Registering Services

1. **Individual Service Registration:**
   Point Backstage to any service's `catalog-info.yaml`:
   ```
   https://github.com/your-org/cluster-tester/blob/main/coffee-shop/catalog-info.yaml
   ```

2. **Bulk Registration:**
   Register all services at once using the root catalog:
   ```
   https://github.com/your-org/cluster-tester/blob/main/catalog-info.yaml
   ```

### Service Catalog Features

Each service is registered with:
- **Component metadata** (name, description, owner)
- **API definitions** (OpenAPI spec)
- **Dependencies** (databases, external services)
- **Documentation links**
- **Source code links**

### API Documentation

The OpenAPI specifications include:
- Complete endpoint documentation
- Request/response schemas
- Example payloads
- Error responses
- Authentication requirements (where applicable)

## Development

### Adding New Endpoints

1. Add the endpoint handler function
2. Add Swagger annotations:
   ```go
   // @Summary Endpoint description
   // @Description Detailed description
   // @Tags tag-name
   // @Accept json
   // @Produce json
   // @Param param-name path string true "Parameter description"
   // @Success 200 {object} ResponseType
   // @Router /endpoint/{param} [get]
   ```
3. Register the route in `main()`
4. Regenerate Swagger docs: `swag init`

### Updating API Documentation

1. Modify the OpenAPI spec in `catalog-info.yaml`
2. Update Swagger annotations in Go code
3. Regenerate documentation
4. Commit changes to trigger Backstage refresh

## Testing

### Health Checks

All services expose health endpoints for monitoring:

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "service": "service-name",
  "version": "1.0.0"
}
```

### API Testing

Use the Swagger UI for interactive testing:
1. Navigate to http://localhost:8080/swagger/
2. Expand any endpoint
3. Click "Try it out"
4. Fill in parameters
5. Execute the request

### Automated Testing

Each service includes test files:
```bash
cd coffee-shop
go test ./...
```

## Contributing

1. Follow the established patterns for new services
2. Include comprehensive Swagger documentation
3. Add health check endpoints
4. Create Backstage catalog entries
5. Update this README with new services

## Support

For questions or issues:
- Check the individual service documentation
- Review the Swagger specifications
- Consult the Backstage catalog entries
- Open an issue in this repository
