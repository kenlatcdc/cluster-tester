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
    "title": "Coffee Shop API",
    "version": "1.0.0",
    "description": "This is a coffee shop service API",
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
    "/coffees": {
      "get": {
        "summary": "Get all coffees",
        "description": "Returns a list of all available coffees",
        "operationId": "getCoffees",
        "tags": ["coffees"],
        "responses": {
          "200": {
            "description": "List of coffees",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/Coffee"
                  }
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create a new coffee",
        "description": "Add a new coffee to the menu",
        "operationId": "createCoffee",
        "tags": ["coffees"],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/Coffee"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Coffee created successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Coffee"
                }
              }
            }
          },
          "400": {
            "description": "Invalid input"
          }
        }
      }
    },
    "/coffees/{id}": {
      "get": {
        "summary": "Get coffee by ID",
        "description": "Returns a single coffee",
        "operationId": "getCoffeeById",
        "tags": ["coffees"],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "description": "Coffee ID",
            "schema": {
              "type": "integer"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Coffee details",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Coffee"
                }
              }
            }
          },
          "404": {
            "description": "Coffee not found"
          }
        }
      }
    },
    "/health": {
      "get": {
        "summary": "Health check",
        "description": "Check if the service is healthy",
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
                      "type": "string",
                      "example": "healthy"
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
      "Coffee": {
        "type": "object",
        "required": ["id", "name", "price"],
        "properties": {
          "id": {
            "type": "integer",
            "example": 1
          },
          "name": {
            "type": "string",
            "example": "Espresso"
          },
          "price": {
            "type": "number",
            "format": "float",
            "example": 2.50
          },
          "description": {
            "type": "string",
            "example": "Strong and bold coffee"
          }
        }
      }
    }
  }
}`

type Coffee struct {
	ID    int     `json:"id" example:"1"`
	Name  string  `json:"name" example:"Espresso"`
	Price float64 `json:"price" example:"2.99"`
}

var coffees = []Coffee{
	{1, "Espresso", 2.99},
	{2, "Americano", 2.49},
	{3, "Latte", 3.49},
	{4, "Cappuccino", 3.49},
	{5, "Mocha", 3.99},
	{6, "Macchiato", 3.19},
	{7, "Flat White", 3.29},
	{8, "Cold Brew", 3.49},
	{9, "Frappuccino", 4.49},
	{10, "Affogato", 3.99},
	{11, "Iced Coffee", 2.99},
	{12, "Nitro Cold Brew", 3.99},
	{13, "Cortado", 3.29},
	{14, "Red Eye", 3.99},
	{15, "Turkish Coffee", 2.99},
}

func main() {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", healthCheck)

	// OpenAPI specification endpoint
	r.GET("/openapi.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, openAPISpec)
	})

	// Swagger UI served from CDN
	r.GET("/docs", func(c *gin.Context) {
		html := `<!DOCTYPE html>
<html>
<head>
    <title>Coffee Shop API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: '/openapi.json',
            dom_id: '#swagger-ui',
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIBundle.presets.standalone
            ]
        });
    </script>
</body>
</html>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	r.GET("/coffees", getCoffees)
	r.GET("/coffees/:id", getCoffeeByID)
	r.POST("/coffees", createCoffee)
	r.DELETE("/coffees/:id", deleteCoffee)
	r.PUT("/coffees/:id", updateCoffee)

	r.Run(":8080")
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "coffee-shop",
		"version": "1.0.0",
	})
}

func getCoffees(c *gin.Context) {
	c.JSON(http.StatusOK, coffees)
}

func getCoffeeByID(c *gin.Context) {
	id := c.Param("id")
	for _, coffee := range coffees {
		if fmt.Sprintf("%d", coffee.ID) == id {
			c.JSON(http.StatusOK, coffee)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Coffee not found"})
}

func createCoffee(c *gin.Context) {
	var newCoffee Coffee
	if err := c.ShouldBindJSON(&newCoffee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coffees = append(coffees, newCoffee)
	c.JSON(http.StatusCreated, newCoffee)
}

func deleteCoffee(c *gin.Context) {
	id := c.Param("id")
	for i, coffee := range coffees {
		if fmt.Sprintf("%d", coffee.ID) == id {
			coffees = append(coffees[:i], coffees[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Coffee deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Coffee not found"})
}

func updateCoffee(c *gin.Context) {
	id := c.Param("id")
	var updatedCoffee Coffee
	if err := c.ShouldBindJSON(&updatedCoffee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, coffee := range coffees {
		if fmt.Sprintf("%d", coffee.ID) == id {
			coffees[i] = updatedCoffee
			c.JSON(http.StatusOK, updatedCoffee)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Coffee not found"})
}
