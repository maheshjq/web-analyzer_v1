# Swagger Integration Plan for Web Page Analyzer

This document outlines the steps to integrate Swagger documentation into the existing Web Page Analyzer project.

## Project Structure Updates

Here's how the project structure will change after integrating Swagger:

```
web-analyzer/
├── ...
├── docs/                # Generated Swagger documentation (new)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── ...
```

## Implementation Steps

Follow these steps to integrate Swagger into the project:

### 1. Update go.mod and go.sum

Add the required Swagger dependencies by editing the `go.mod` file or running:

```bash
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/http-swagger
```

### 2. Update Model Definitions

Update the models in `internal/models/models.go` to include complete structures and Swagger annotations.

### 3. Add Swagger Annotations to Handlers

Update the handler functions in `internal/api/handlers.go` to include Swagger annotations.

### 4. Configure Swagger UI Routes

Update or create `internal/api/swagger.go` to set up the Swagger UI routes.

### 5. Update Main Function

Modify `cmd/server/main.go` to:
- Import the generated Swagger docs
- Configure the router to serve the Swagger UI

### 6. Generate Swagger Documentation

Run the provided `generate-swagger.sh` script to generate the Swagger documentation.

## Testing the Integration

After implementing these changes:

1. Build and run the application
2. Navigate to `http://localhost:8080/swagger/` in your browser
3. Verify that the Swagger UI loads and displays the API endpoints
4. Test each endpoint through the Swagger UI

## Troubleshooting Common Issues

### Issue: Swagger UI Not Loading

**Solution**: Check that the Swagger middleware is correctly set up in the router and that the path prefix is correct.

### Issue: Swagger Annotations Not Generated

**Solution**: Ensure that the `swag init` command is targeting the correct main.go file and that the annotations are properly formatted.

### Issue: Missing Package Errors

**Solution**: Run `go mod tidy` to update dependencies after adding new imports.

### Issue: Routing Conflicts

**Solution**: Ensure that the Swagger UI route prefix doesn't conflict with other API routes.

## Recommended Updates to Other Files

- **README.md**: Add information about the Swagger documentation
- **Makefile**: Add a target for generating Swagger documentation
- **Dockerfile**: Ensure the Swagger docs directory is included in the build

## Benefits of This Integration

- Self-documenting API that stays in sync with the codebase
- Interactive testing capabilities for developers
- Clear visibility of request and response structures
- Improved developer experience for API consumers