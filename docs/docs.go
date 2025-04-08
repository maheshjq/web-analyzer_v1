// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Web Analyzer Team"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/analyze": {
            "post": {
                "description": "Fetches and analyzes a web page by URL, returning information about its structure and content",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Analyze a web page",
                "operationId": "analyze-web-page",
                "parameters": [
                    {
                        "description": "URL to analyze",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AnalysisRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful analysis",
                        "schema": {
                            "$ref": "#/definitions/models.AnalysisResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request (invalid URL format)",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "502": {
                        "description": "Failed to fetch or analyze the URL",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/analyze": {
            "post": {
                "description": "Fetches and analyzes a web page by URL, returning information about its structure and content",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "analysis"
                ],
                "summary": "Analyze a web page",
                "parameters": [
                    {
                        "description": "URL to analyze",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.AnalysisRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful analysis",
                        "schema": {
                            "$ref": "#/definitions/api.AnalysisResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request (invalid URL format)",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "502": {
                        "description": "Failed to fetch or analyze the URL",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/health": {
            "get": {
                "description": "Returns the health status of the API",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "Service is healthy",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Returns the health status of the API",
                "produces": [
                    "application/json"
                ],
                "summary": "Health check",
                "operationId": "health-check",
                "responses": {
                    "200": {
                        "description": "Service is healthy",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.AnalysisRequest": {
            "type": "object",
            "properties": {
                "url": {
                    "type": "string",
                    "example": "https://example.com"
                }
            }
        },
        "api.AnalysisResponse": {
            "type": "object",
            "properties": {
                "containsLoginForm": {
                    "type": "boolean",
                    "example": false
                },
                "headings": {
                    "$ref": "#/definitions/api.HeadingCount"
                },
                "htmlVersion": {
                    "type": "string",
                    "example": "HTML5"
                },
                "links": {
                    "$ref": "#/definitions/api.LinkAnalysis"
                },
                "title": {
                    "type": "string",
                    "example": "Example Domain"
                }
            }
        },
        "api.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Failed to analyze URL: HTTP error 404 Not Found"
                },
                "statusCode": {
                    "type": "integer",
                    "example": 502
                }
            }
        },
        "api.HeadingCount": {
            "type": "object",
            "properties": {
                "h1": {
                    "type": "integer",
                    "example": 1
                },
                "h2": {
                    "type": "integer",
                    "example": 2
                },
                "h3": {
                    "type": "integer",
                    "example": 3
                },
                "h4": {
                    "type": "integer",
                    "example": 0
                },
                "h5": {
                    "type": "integer",
                    "example": 0
                },
                "h6": {
                    "type": "integer",
                    "example": 0
                }
            }
        },
        "api.LinkAnalysis": {
            "type": "object",
            "properties": {
                "external": {
                    "type": "integer",
                    "example": 3
                },
                "inaccessible": {
                    "type": "integer",
                    "example": 1
                },
                "internal": {
                    "type": "integer",
                    "example": 5
                }
            }
        },
        "models.AnalysisRequest": {
            "type": "object",
            "properties": {
                "url": {
                    "type": "string",
                    "example": "https://example.com"
                }
            }
        },
        "models.AnalysisResponse": {
            "type": "object",
            "properties": {
                "containsLoginForm": {
                    "type": "boolean",
                    "example": false
                },
                "headings": {
                    "$ref": "#/definitions/models.HeadingCount"
                },
                "htmlVersion": {
                    "type": "string",
                    "example": "HTML5"
                },
                "links": {
                    "$ref": "#/definitions/models.LinkAnalysis"
                },
                "title": {
                    "type": "string",
                    "example": "Example Domain"
                }
            }
        },
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Failed to analyze URL: HTTP error 404 Not Found"
                },
                "statusCode": {
                    "type": "integer",
                    "example": 502
                }
            }
        },
        "models.HeadingCount": {
            "type": "object",
            "properties": {
                "h1": {
                    "type": "integer",
                    "example": 1
                },
                "h2": {
                    "type": "integer",
                    "example": 2
                },
                "h3": {
                    "type": "integer",
                    "example": 3
                },
                "h4": {
                    "type": "integer",
                    "example": 0
                },
                "h5": {
                    "type": "integer",
                    "example": 0
                },
                "h6": {
                    "type": "integer",
                    "example": 0
                }
            }
        },
        "models.LinkAnalysis": {
            "type": "object",
            "properties": {
                "external": {
                    "type": "integer",
                    "example": 3
                },
                "inaccessible": {
                    "type": "integer",
                    "example": 1
                },
                "internal": {
                    "type": "integer",
                    "example": 5
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Web Page Analyzer API",
	Description:      "API for analyzing web pages, extracting HTML version, title, headings, links, and detecting login forms.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
