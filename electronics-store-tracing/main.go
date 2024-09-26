package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
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

	// GET all products
	r.GET("/products", func(c *gin.Context) {
		products, err := getAllProducts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	// GET a product by ID
	r.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		product, err := getProductByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	r.Run(":8080")
}
