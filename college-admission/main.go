package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	r.GET("/applications", func(c *gin.Context) {
		c.JSON(http.StatusOK, applications)
	})

	r.GET("/applications/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, app := range applications {
			if fmt.Sprintf("%d", app.ID) == id {
				c.JSON(http.StatusOK, app)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
	})

	r.POST("/applications", func(c *gin.Context) {
		var newApp Application
		if err := c.ShouldBindJSON(&newApp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		applications = append(applications, newApp)
		c.JSON(http.StatusCreated, newApp)
	})

	r.DELETE("/applications/:id", func(c *gin.Context) {
		id := c.Param("id")
		for i, app := range applications {
			if fmt.Sprintf("%d", app.ID) == id {
				applications = append(applications[:i], applications[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Application deleted"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
	})

	r.PUT("/applications/:id", func(c *gin.Context) {
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
	})

	r.Run(":8080")
}
