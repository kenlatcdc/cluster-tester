package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MenuItem struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
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

	r.GET("/menu", func(c *gin.Context) {
		c.JSON(http.StatusOK, menuItems)
	})

	r.GET("/menu/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, item := range menuItems {
			if fmt.Sprintf("%d", item.ID) == id {
				c.JSON(http.StatusOK, item)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
	})

	r.POST("/menu", func(c *gin.Context) {
		var newItem MenuItem
		if err := c.ShouldBindJSON(&newItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		menuItems = append(menuItems, newItem)
		c.JSON(http.StatusCreated, newItem)
	})

	r.DELETE("/menu/:id", func(c *gin.Context) {
		id := c.Param("id")
		for i, item := range menuItems {
			if fmt.Sprintf("%d", item.ID) == id {
				menuItems = append(menuItems[:i], menuItems[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
	})

	r.PUT("/menu/:id", func(c *gin.Context) {
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
	})

	r.Run(":8080")
}
