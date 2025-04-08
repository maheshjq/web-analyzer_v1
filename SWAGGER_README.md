# Using the Swagger UI for Web Page Analyzer

This guide explains how to use the Swagger UI to interact with the Web Page Analyzer API.

## Accessing the Swagger UI

After starting the application, navigate to:

```
http://localhost:8080/swagger/
```

You should see the Swagger UI interface displaying the available API endpoints.

## Available Endpoints

The Web Page Analyzer API has the following endpoints:

### 1. Analyze Web Page

- **Endpoint**: `/api/analyze`
- **Method**: POST
- **Description**: Analyzes a web page's structure and content

#### How to Use:

1. Click on the `/api/analyze` endpoint in the Swagger UI
2. Click the "Try it out" button
3. Enter a URL in the request body:
   ```json
   {
     "url": "https://example.com"
   }
   ```
4. Click "Execute"
5. The API will return analysis results:
   ```json
   {
     "htmlVersion": "HTML5",
     "title": "Example Domain",
     "headings": {
       "h1": 1,
       "h2": 2,
       "h3": 0,
       "h4": 0,
       "h5": 0,
       "h6": 0
     },
     "links": {
       "internal": 0,
       "external": 1,
       "inaccessible": 0
     },
     "containsLoginForm": false
   }
   ```

### 2. Health Check

- **Endpoint**: `/api/health`
- **Method**: GET
- **Description**: Returns the health status of the API

#### How to Use:

1. Click on the `/api/health` endpoint in the Swagger UI
2. Click the "Try it out" button
3. Click "Execute"
4. The API will return:
   ```json
   {
     "status": "ok"
   }
   ```

## Understanding the Responses

### Analysis Response

- **htmlVersion**: The detected HTML version of the page
- **title**: The page title
- **headings**: Count of different heading levels (h1-h6)
- **links**:
  - **internal**: Links to the same domain
  - **external**: Links to different domains
  - **inaccessible**: Links that couldn't be accessed
- **containsLoginForm**: Whether a login form was detected

### Error Response

If an error occurs, the API will return:

```json
{
  "statusCode": 400,
  "message": "Error message details"
}
```

Common status codes:
- 400: Bad request (e.g., invalid URL format)
- 502: Failed to fetch or analyze the URL

## Tips for Using the Swagger UI

1. You can expand/collapse the endpoint details by clicking on them
2. The "Schema" tab shows the data structure for requests and responses
3. After executing a request, you can see the curl command that was generated
4. You can see the full response details, including headers and status code