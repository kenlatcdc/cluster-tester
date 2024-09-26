package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	r.GET("/pets", func(c *gin.Context) {
		c.JSON(http.StatusOK, pets)
	})

	r.GET("/pets/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, pet := range pets {
			if fmt.Sprintf("%d", pet.ID) == id {
				c.JSON(http.StatusOK, pet)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
	})

	r.POST("/pets", func(c *gin.Context) {
		var newPet Pet
		if err := c.ShouldBindJSON(&newPet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pets = append(pets, newPet)
		c.JSON(http.StatusCreated, newPet)
	})

	r.DELETE("/pets/:id", func(c *gin.Context) {
		id := c.Param("id")
		for i, pet := range pets {
			if fmt.Sprintf("%d", pet.ID) == id {
				pets = append(pets[:i], pets[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Pet deleted"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
	})

	r.PUT("/pets/:id", func(c *gin.Context) {
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
	})

	r.Run(":8080")
}
