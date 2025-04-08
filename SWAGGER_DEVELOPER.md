# Adding Swagger Documentation to Web Page Analyzer

This guide explains how to implement Swagger documentation for the Web Page Analyzer project.

## Prerequisites

- Go 1.21 or later
- The Web Page Analyzer project codebase

## Implementation Steps

### 1. Install Required Packages

Add the following packages to your project:

```bash
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/http-swagger
```

### 2. Update the Models

Update the models to include Swagger annotations and examples. Each field should have a proper example value that represents what the API will return.

### 3. Add Swagger Annotations to Handlers

Each handler function should include proper Swagger annotations:

- `@Summary` - A brief summary of what the endpoint does
- `@Description` - A more detailed description
- `@Tags` - Categorization for the Swagger UI
- `@Accept` - What content types are accepted
- `@Produce` - What content types are produced
- `@Param` - Parameters for the endpoint
- `@Success` - Success response with status code and type
- `@Failure` - Error responses with status codes and types
- `@Router` - The path and method for the endpoint

### 4. Create a Swagger Setup Function

Create a function to set up the Swagger UI routes in your application.

### 5. Update the Main Function

Modify the main function to initialize and use the Swagger setup.

### 6. Generate Swagger Documentation

Run the `generate-swagger.sh` script to generate the Swagger documentation from your annotations:

```bash
chmod +x generate-swagger.sh
./generate-swagger.sh
```

This will create a `docs` directory with all the necessary Swagger files.

## Accessing the Swagger UI

Once your application is running, you can access the Swagger UI at:

```
http://localhost:8080/swagger/
```

## Troubleshooting

If you encounter issues:

1. Make sure all imported packages are correct
2. Ensure the Swagger annotations are properly formatted
3. Check that the `docs` directory is being imported in your main.go file
4. Verify that the Swagger middleware is correctly set up in your router

## Benefits of Swagger Documentation

- Interactive API documentation that stays in sync with your code
- Ability to test API endpoints directly from the documentation
- Clear visibility of request and response structures
- Improved developer experience for API consumers