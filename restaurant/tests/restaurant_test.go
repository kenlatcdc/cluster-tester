package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MenuItem struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var testMenuItems = []MenuItem{
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

func TestGetMenu(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/menu", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Pizza")
	assert.Contains(t, w.Body.String(), "Burger")
	assert.Contains(t, w.Body.String(), "Pasta")
}

func TestGetMenuItem(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/menu/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Pizza")
}

func TestGetNonExistentMenuItem(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/menu/99", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Item not found")
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/menu", func(c *gin.Context) {
		c.JSON(http.StatusOK, testMenuItems)
	})

	r.GET("/menu/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, item := range testMenuItems {
			if fmt.Sprintf("%d", item.ID) == id {
				c.JSON(http.StatusOK, item)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
	})

	return r
}
