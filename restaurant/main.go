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
    "title": "Restaurant API",
    "version": "1.0.0",
    "description": "This is a restaurant service API",
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
    "/menu": {
      "get": {
        "summary": "Get menu items",
        "description": "Returns a list of all menu items",
        "operationId": "getMenuItems",
        "tags": ["menu"],
        "responses": {
          "200": {
            "description": "List of menu items",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/MenuItem"
                  }
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create a new menu item",
        "description": "Add a new item to the menu",
        "operationId": "createMenuItem",
        "tags": ["menu"],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/MenuItem"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Menu item created successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/MenuItem"
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
    "/menu/{id}": {
      "get": {
        "summary": "Get menu item by ID",
        "description": "Get a specific menu item by its ID",
        "operationId": "getMenuItemByID",
        "tags": ["menu"],
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
            "description": "Menu item found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/MenuItem"
                }
              }
            }
          },
          "404": {
            "description": "Menu item not found",
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
      "MenuItem": {
        "type": "object",
        "required": ["id", "name", "description", "price", "category"],
        "properties": {
          "id": {
            "type": "integer",
            "example": 1
          },
          "name": {
            "type": "string",
            "example": "Pasta Carbonara"
          },
          "description": {
            "type": "string",
            "example": "Creamy pasta with bacon and eggs"
          },
          "price": {
            "type": "number",
            "format": "float",
            "example": 15.99
          },
          "category": {
            "type": "string",
            "example": "Main Course"
          }
        }
      }
    }
  }
}`

type MenuItem struct {
	ID    int     `json:"id" example:"1"`
	Name  string  `json:"name" example:"Pizza"`
	Price float64 `json:"price" example:"12.99"`
}

var menuItems = []MenuItem{
	{1, "Pizza", 12.99},
	{2, "Burger", 9.99},
	{3, "Pasta", 11.99},
	{4, "Salad", 7.99},
	{5, "Steak", 19.99},
	{6, "Sushi", 14.99},
	{7, "Ramen", 10.99},
	{8, "Tacos", 8.99},
	{9, "Sandwich", 6.99},
	{10, "Fries", 3.99},
	{11, "Ice Cream", 4.99},
	{12, "Soup", 5.99},
	{13, "Coffee", 2.99},
	{14, "Tea", 2.49},
	{15, "Smoothie", 4.49},
}

func main() {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", healthCheck)

	// OpenAPI documentation endpoints
	r.GET("/openapi.json", getOpenAPISpec)
	r.GET("/docs", serveDocs)

	r.GET("/menu", getMenuItems)
	r.GET("/menu/:id", getMenuItemByID)
	r.POST("/menu", createMenuItem)
	r.DELETE("/menu/:id", deleteMenuItem)
	r.PUT("/menu/:id", updateMenuItem)

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
  <title>Restaurant API Documentation</title>
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
		"service": "restaurant",
		"version": "1.0.0",
	})
}

// getMenuItems godoc
// @Summary Get all menu items
// @Description Get list of all available menu items
// @Tags menu
// @Accept json
// @Produce json
// @Success 200 {array} MenuItem
// @Router /menu [get]
func getMenuItems(c *gin.Context) {
	c.JSON(http.StatusOK, menuItems)
}

// getMenuItemByID godoc
// @Summary Get menu item by ID
// @Description Get a specific menu item by its ID
// @Tags menu
// @Accept json
// @Produce json
// @Param id path int true "Menu Item ID"
// @Success 200 {object} MenuItem
// @Failure 404 {object} map[string]string
// @Router /menu/{id} [get]
func getMenuItemByID(c *gin.Context) {
	id := c.Param("id")
	for _, item := range menuItems {
		if fmt.Sprintf("%d", item.ID) == id {
			c.JSON(http.StatusOK, item)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
}

// createMenuItem godoc
// @Summary Create a new menu item
// @Description Add a new item to the menu
// @Tags menu
// @Accept json
// @Produce json
// @Param menuItem body MenuItem true "Menu Item object"
// @Success 201 {object} MenuItem
// @Failure 400 {object} map[string]string
// @Router /menu [post]
func createMenuItem(c *gin.Context) {
	var newItem MenuItem
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	menuItems = append(menuItems, newItem)
	c.JSON(http.StatusCreated, newItem)
}

// deleteMenuItem godoc
// @Summary Delete a menu item
// @Description Remove an item from the menu
// @Tags menu
// @Accept json
// @Produce json
// @Param id path int true "Menu Item ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /menu/{id} [delete]
func deleteMenuItem(c *gin.Context) {
	id := c.Param("id")
	for i, item := range menuItems {
		if fmt.Sprintf("%d", item.ID) == id {
			menuItems = append(menuItems[:i], menuItems[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
}

// updateMenuItem godoc
// @Summary Update a menu item
// @Description Update an existing menu item
// @Tags menu
// @Accept json
// @Produce json
// @Param id path int true "Menu Item ID"
// @Param menuItem body MenuItem true "Menu Item object"
// @Success 200 {object} MenuItem
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /menu/{id} [put]
func updateMenuItem(c *gin.Context) {
	id := c.Param("id")
	var updatedItem MenuItem
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, item := range menuItems {
		if fmt.Sprintf("%d", item.ID) == id {
			menuItems[i] = updatedItem
			c.JSON(http.StatusOK, updatedItem)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
}
