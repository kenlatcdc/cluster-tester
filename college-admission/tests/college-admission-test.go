package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Application struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Course    string `json:"course"`
}

var testApplications = []Application{
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

func TestGetApplications(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/applications", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John")
	assert.Contains(t, w.Body.String(), "Jane")
	assert.Contains(t, w.Body.String(), "Bob")
}

func TestGetApplication(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/applications/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John")
}

func TestGetNonExistentApplication(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/applications/99", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Application not found")
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/applications", func(c *gin.Context) {
		c.JSON(http.StatusOK, testApplications)
	})

	r.GET("/applications/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, app := range testApplications {
			if fmt.Sprintf("%d", app.ID) == id {
				c.JSON(http.StatusOK, app)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
	})

	return r
}
