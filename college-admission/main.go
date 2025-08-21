package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OpenAPI 3.0 specification embedded as a constant
const openAPISpec = `{
  "openapi": "3.0.3",
  "info": {
    "title": "College Admission API",
    "version": "1.0.0",
    "description": "This is a college admission service API",
    "contact": {
      "name": "API Support",
      "url": "http://www.swagger.io/support",
      "email": "support@swagger.io"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    }
  },
  "servers": [
    {
      "url": "http://localhost:8080",
      "description": "Development server"
    }
  ],
  "paths": {
    "/applications": {
      "get": {
        "summary": "Get all applications",
        "description": "Returns a list of all college applications",
        "operationId": "getApplications",
        "tags": ["applications"],
        "responses": {
          "200": {
            "description": "List of applications",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/Application"
                  }
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create a new application",
        "description": "Submit a new college application",
        "operationId": "createApplication",
        "tags": ["applications"],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/Application"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Application created successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Application"
                }
              }
            }
          },
          "400": {
            "description": "Bad request",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "error": {
                      "type": "string"
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/applications/{id}": {
      "get": {
        "summary": "Get application by ID",
        "description": "Get a specific college application by its ID",
        "operationId": "getApplicationByID",
        "tags": ["applications"],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "integer"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Application found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Application"
                }
              }
            }
          },
          "404": {
            "description": "Application not found",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "error": {
                      "type": "string"
                    }
                  }
                }
              }
            }
          }
        }
      },
      "put": {
        "summary": "Update an application",
        "description": "Update an existing college application",
        "operationId": "updateApplication",
        "tags": ["applications"],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "integer"
            }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/Application"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Application updated successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Application"
                }
              }
            }
          },
          "400": {
            "description": "Bad request"
          },
          "404": {
            "description": "Application not found"
          }
        }
      },
      "delete": {
        "summary": "Delete an application",
        "description": "Remove a college application",
        "operationId": "deleteApplication",
        "tags": ["applications"],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "integer"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Application deleted successfully"
          },
          "404": {
            "description": "Application not found"
          }
        }
      }
    },
    "/health": {
      "get": {
        "summary": "Health check",
        "description": "Check if the service is running",
        "operationId": "healthCheck",
        "tags": ["health"],
        "responses": {
          "200": {
            "description": "Service is healthy",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": {
                      "type": "string"
                    },
                    "service": {
                      "type": "string"
                    },
                    "version": {
                      "type": "string"
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Application": {
        "type": "object",
        "required": ["id", "first_name", "last_name", "age", "course"],
        "properties": {
          "id": {
            "type": "integer",
            "example": 1
          },
          "first_name": {
            "type": "string",
            "example": "John"
          },
          "last_name": {
            "type": "string",
            "example": "Doe"
          },
          "age": {
            "type": "integer",
            "example": 18
          },
          "course": {
            "type": "string",
            "example": "Computer Science"
          }
        }
      }
    }
  }
}`

type Application struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Course    string `json:"course"`
}

var applications = []Application{
	{1, "John", "Doe", 18, "Computer Science"},
	{2, "Jane", "Smith", 19, "Mechanical Engineering"},
	{3, "Bob", "Brown", 17, "Civil Engineering"},
	{4, "Alice", "Johnson", 20, "Electrical Engineering"},
	{5, "Charlie", "Davis", 21, "Business Administration"},
	{6, "David", "Wilson", 22, "Mathematics"},
	{7, "Eve", "Clark", 18, "Physics"},
	{8, "Frank", "Moore", 19, "Chemistry"},
	{9, "Grace", "Taylor", 17, "Biology"},
	{10, "Henry", "Anderson", 20, "Psychology"},
	{11, "Ivy", "Thomas", 21, "Philosophy"},
	{12, "Jack", "Jackson", 22, "Sociology"},
	{13, "Kathy", "White", 18, "History"},
	{14, "Leo", "Harris", 19, "Political Science"},
	{15, "Mia", "Martin", 17, "Art"},
}

func main() {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", healthCheck)

	// OpenAPI documentation endpoints
	r.GET("/openapi.json", getOpenAPISpec)
	r.GET("/docs", serveDocs)

	r.GET("/applications", getApplications)
	r.GET("/applications/:id", getApplicationByID)
	r.POST("/applications", createApplication)
	r.DELETE("/applications/:id", deleteApplication)
	r.PUT("/applications/:id", updateApplication)

	r.Run(":8080")
}

// getOpenAPISpec returns the OpenAPI specification as JSON
func getOpenAPISpec(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, openAPISpec)
}

// serveDocs serves the Swagger UI documentation page
func serveDocs(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
  <title>College Admission API Documentation</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui.css" />
  <style>
    html {
      box-sizing: border-box;
      overflow: -moz-scrollbars-vertical;
      overflow-y: scroll;
    }
    *, *:before, *:after {
      box-sizing: inherit;
    }
    body {
      margin:0;
      background: #fafafa;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: '/openapi.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "college-admission",
		"version": "1.0.0",
	})
}

func getApplications(c *gin.Context) {
	c.JSON(http.StatusOK, applications)
}

func getApplicationByID(c *gin.Context) {
	id := c.Param("id")
	for _, app := range applications {
		if fmt.Sprintf("%d", app.ID) == id {
			c.JSON(http.StatusOK, app)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
}

func createApplication(c *gin.Context) {
	var newApp Application
	if err := c.ShouldBindJSON(&newApp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	applications = append(applications, newApp)
	c.JSON(http.StatusCreated, newApp)
}

func deleteApplication(c *gin.Context) {
	id := c.Param("id")
	for i, app := range applications {
		if fmt.Sprintf("%d", app.ID) == id {
			applications = append(applications[:i], applications[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Application deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
}

func updateApplication(c *gin.Context) {
	id := c.Param("id")
	var updatedApp Application
	if err := c.ShouldBindJSON(&updatedApp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, app := range applications {
		if fmt.Sprintf("%d", app.ID) == id {
			applications[i] = updatedApp
			c.JSON(http.StatusOK, updatedApp)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
}
