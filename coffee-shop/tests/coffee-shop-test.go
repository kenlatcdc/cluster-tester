package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Coffee struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var testCoffees = []Coffee{
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

func TestGetCoffees(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/coffees", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Espresso")
	assert.Contains(t, w.Body.String(), "Americano")
	assert.Contains(t, w.Body.String(), "Latte")
}

func TestGetCoffee(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/coffees/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Espresso")
}

func TestGetNonExistentCoffee(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/coffees/99", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Coffee not found")
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/coffees", func(c *gin.Context) {
		c.JSON(http.StatusOK, testCoffees)
	})

	r.GET("/coffees/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, coffee := range testCoffees {
			if fmt.Sprintf("%d", coffee.ID) == id {
				c.JSON(http.StatusOK, coffee)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Coffee not found"})
	})

	return r
}
