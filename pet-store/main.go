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
    "title": "Pet Store API",
    "version": "1.0.0",
    "description": "This is a pet store service API",
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
    "/pets": {
      "get": {
        "summary": "Get all pets",
        "description": "Returns a list of all available pets",
        "operationId": "getPets",
        "tags": ["pets"],
        "responses": {
          "200": {
            "description": "List of pets",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/Pet"
                  }
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create a new pet",
        "description": "Add a new pet to the store",
        "operationId": "createPet",
        "tags": ["pets"],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/Pet"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Pet created successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Pet"
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
    "/pets/{id}": {
      "get": {
        "summary": "Get pet by ID",
        "description": "Get a specific pet by its ID",
        "operationId": "getPetByID",
        "tags": ["pets"],
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
            "description": "Pet found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Pet"
                }
              }
            }
          },
          "404": {
            "description": "Pet not found",
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
        "summary": "Update a pet",
        "description": "Update an existing pet",
        "operationId": "updatePet",
        "tags": ["pets"],
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
                "$ref": "#/components/schemas/Pet"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Pet updated successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Pet"
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
          },
          "404": {
            "description": "Pet not found",
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
      "delete": {
        "summary": "Delete a pet",
        "description": "Remove a pet from the store",
        "operationId": "deletePet",
        "tags": ["pets"],
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
            "description": "Pet deleted successfully",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "message": {
                      "type": "string"
                    }
                  }
                }
              }
            }
          },
          "404": {
            "description": "Pet not found",
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
      "Pet": {
        "type": "object",
        "required": ["id", "name", "type", "age"],
        "properties": {
          "id": {
            "type": "integer",
            "example": 1
          },
          "name": {
            "type": "string",
            "example": "Max"
          },
          "type": {
            "type": "string",
            "example": "Dog"
          },
          "age": {
            "type": "integer",
            "example": 3
          }
        }
      }
    }
  }
}`

type Pet struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Age  int    `json:"age"`
}

var pets = []Pet{
	{1, "Max", "Dog", 3},
	{2, "Bella", "Cat", 2},
	{3, "Charlie", "Dog", 4},
	{4, "Lucy", "Cat", 1},
	{5, "Buddy", "Dog", 5},
	{6, "Luna", "Cat", 3},
	{7, "Rocky", "Dog", 2},
	{8, "Molly", "Cat", 4},
	{9, "Duke", "Dog", 3},
	{10, "Daisy", "Cat", 2},
	{11, "Bear", "Dog", 1},
	{12, "Lola", "Cat", 3},
	{13, "Jack", "Dog", 5},
	{14, "Zoe", "Cat", 1},
	{15, "Toby", "Dog", 4},
}

func main() {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", healthCheck)

	// OpenAPI documentation endpoints
	r.GET("/openapi.json", getOpenAPISpec)
	r.GET("/docs", serveDocs)

	r.GET("/pets", getPets)
	r.GET("/pets/:id", getPetByID)
	r.POST("/pets", createPet)
	r.DELETE("/pets/:id", deletePet)
	r.PUT("/pets/:id", updatePet)

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
  <title>Pet Store API Documentation</title>
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
		"service": "pet-store",
		"version": "1.0.0",
	})
}

func getPets(c *gin.Context) {
	c.JSON(http.StatusOK, pets)
}

func getPetByID(c *gin.Context) {
	id := c.Param("id")
	for _, pet := range pets {
		if fmt.Sprintf("%d", pet.ID) == id {
			c.JSON(http.StatusOK, pet)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
}

func createPet(c *gin.Context) {
	var newPet Pet
	if err := c.ShouldBindJSON(&newPet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pets = append(pets, newPet)
	c.JSON(http.StatusCreated, newPet)
}

func deletePet(c *gin.Context) {
	id := c.Param("id")
	for i, pet := range pets {
		if fmt.Sprintf("%d", pet.ID) == id {
			pets = append(pets[:i], pets[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Pet deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
}

func updatePet(c *gin.Context) {
	id := c.Param("id")
	var updatedPet Pet
	if err := c.ShouldBindJSON(&updatedPet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, pet := range pets {
		if fmt.Sprintf("%d", pet.ID) == id {
			pets[i] = updatedPet
			c.JSON(http.StatusOK, updatedPet)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
}
