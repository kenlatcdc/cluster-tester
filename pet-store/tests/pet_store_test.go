package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Pet struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Age  int    `json:"age"`
}

var testPets = []Pet{
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

func TestGetPets(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/pets", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Max")
	assert.Contains(t, w.Body.String(), "Bella")
	assert.Contains(t, w.Body.String(), "Charlie")
}

func TestGetPet(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/pets/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Max")
}

func TestGetNonExistentPet(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/pets/99", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Pet not found")
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/pets", func(c *gin.Context) {
		c.JSON(http.StatusOK, testPets)
	})

	r.GET("/pets/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, pet := range testPets {
			if fmt.Sprintf("%d", pet.ID) == id {
				c.JSON(http.StatusOK, pet)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
	})

	return r
}
