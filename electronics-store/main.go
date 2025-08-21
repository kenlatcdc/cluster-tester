package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// OpenAPI 3.0 specification embedded as a constant
const openAPISpec = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Electronics Store API",
    "version": "1.0.0",
    "description": "This is an electronics store service API",
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
    "/products": {
      "get": {
        "summary": "Get all products",
        "description": "Returns a list of all electronics products",
        "operationId": "getProducts",
        "tags": ["products"],
        "responses": {
          "200": {
            "description": "List of products",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/Product"
                  }
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create a new product",
        "description": "Add a new electronics product",
        "operationId": "createProduct",
        "tags": ["products"],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/Product"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Product created successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Product"
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
    "/products/{id}": {
      "get": {
        "summary": "Get product by ID",
        "description": "Get a specific product by its ID",
        "operationId": "getProductByID",
        "tags": ["products"],
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
            "description": "Product found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Product"
                }
              }
            }
          },
          "404": {
            "description": "Product not found",
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
      "Product": {
        "type": "object",
        "required": ["id", "name", "description", "price", "category", "stock"],
        "properties": {
          "id": {
            "type": "integer",
            "example": 1
          },
          "name": {
            "type": "string",
            "example": "iPhone 13"
          },
          "description": {
            "type": "string",
            "example": "Latest Apple smartphone"
          },
          "price": {
            "type": "number",
            "format": "float",
            "example": 999.99
          },
          "category": {
            "type": "string",
            "example": "Smartphones"
          },
          "stock": {
            "type": "integer",
            "example": 50
          }
        }
      }
    }
  }
}`

// @schemes http

type Product struct {
	ID    int     `json:"id" example:"1"`
	Name  string  `json:"name" example:"Laptop"`
	Price float64 `json:"price" example:"999.99"`
}

var db *sql.DB

var products = []Product{
	{1, "Laptop", 999.99},
	{2, "Smartphone", 699.99},
	{3, "Tablet", 499.99},
	{4, "Headphones", 199.99},
	{5, "Smartwatch", 299.99},
	{6, "Camera", 599.99},
	{7, "Printer", 149.99},
	{8, "Monitor", 249.99},
	{9, "Keyboard", 49.99},
	{10, "Mouse", 29.99},
	{11, "Router", 89.99},
	{12, "Speaker", 129.99},
	{13, "Microphone", 99.99},
	{14, "External Hard Drive", 79.99},
	{15, "USB Flash Drive", 19.99},
}

// Initialize database connection and create table if it doesn't exist
func initDB() {
	var err error
	dsn := "admin:password123@tcp(mysql:3306)/electronics-store"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// Check if connection is successful
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Create table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		price DECIMAL(10, 2) NOT NULL
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		panic(err)
	}
	fmt.Println("Table 'products' is ready or already exists.")

	// Insert test data into the table
	err = insertProducts()
	if err != nil {
		log.Fatalf("Error inserting products: %v", err)
	}
}

// Insert products into the database if they don't already exist
func insertProducts() error {
	for _, product := range products {
		// Check if the product already exists
		exists, err := productExistsByName(product.Name)
		if err != nil {
			return err
		}

		// Insert product if it doesn't exist
		if !exists {
			_, err := db.Exec("INSERT INTO products (name, price) VALUES (?, ?)", product.Name, product.Price)
			if err != nil {
				return err
			}
			fmt.Printf("Inserted product: %s\n", product.Name)
		}
	}
	return nil
}

// Check if a product with the same name already exists
func productExistsByName(name string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE name = ?)", name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Helper function to get all products
func getAllProducts() ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// Helper function to get a product by ID
func getProductByID(id string) (*Product, error) {
	var product Product
	err := db.QueryRow("SELECT id, name, price FROM products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func main() {
	initDB()

	r := gin.Default()

	// Health check endpoint
	r.GET("/health", healthCheck)

	// OpenAPI documentation endpoints
	r.GET("/openapi.json", getOpenAPISpec)
	r.GET("/docs", serveDocs)

	// GET all products
	r.GET("/products", getProducts)

	// GET a product by ID
	r.GET("/products/:id", getProductByIDHandler)

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
  <title>Electronics Store API Documentation</title>
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
		"service": "electronics-store",
		"version": "1.0.0",
	})
}

// getProducts godoc
// @Summary Get all products
// @Description Get list of all available products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func getProducts(c *gin.Context) {
	products, err := getAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// getProductByIDHandler godoc
// @Summary Get product by ID
// @Description Get a specific product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func getProductByIDHandler(c *gin.Context) {
	id := c.Param("id")
	product, err := getProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}
