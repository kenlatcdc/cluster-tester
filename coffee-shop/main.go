package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Coffee struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
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

	r.GET("/coffees", func(c *gin.Context) {
		c.JSON(http.StatusOK, coffees)
	})

	r.GET("/coffees/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, coffee := range coffees {
			if fmt.Sprintf("%d", coffee.ID) == id {
				c.JSON(http.StatusOK, coffee)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Coffee not found"})
	})

	r.POST("/coffees", func(c *gin.Context) {
		var newCoffee Coffee
		if err := c.ShouldBindJSON(&newCoffee); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		coffees = append(coffees, newCoffee)
		c.JSON(http.StatusCreated, newCoffee)
	})

	r.DELETE("/coffees/:id", func(c *gin.Context) {
		id := c.Param("id")
		for i, coffee := range coffees {
			if fmt.Sprintf("%d", coffee.ID) == id {
				coffees = append(coffees[:i], coffees[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Coffee deleted"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Coffee not found"})
	})

	r.PUT("/coffees/:id", func(c *gin.Context) {
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
	})

	r.Run(":8080")
}
