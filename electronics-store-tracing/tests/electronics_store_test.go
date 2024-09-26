package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var testProducts = []Product{
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

func TestGetProducts(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Laptop")
	assert.Contains(t, w.Body.String(), "Smartphone")
	assert.Contains(t, w.Body.String(), "Tablet")
}

func TestGetProduct(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Laptop")
}

func TestGetNonExistentProduct(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/99", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Product not found")
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/products", func(c *gin.Context) {
		c.JSON(http.StatusOK, testProducts)
	})

	r.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, product := range testProducts {
			if fmt.Sprintf("%d", product.ID) == id {
				c.JSON(http.StatusOK, product)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
	})

	return r
}
