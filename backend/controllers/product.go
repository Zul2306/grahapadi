package controllers

import (
	"inventory-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// In-memory products store (replace with database)
var products = []models.Product{
	{
		ID: 1, Code: "PRD-001", Name: "Laptop Dell Inspiron",
		Category: "Electronics", Unit: "pcs", Stock: 10, MinStock: 3,
		Description: "Laptop for office use", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	},
	{
		ID: 2, Code: "PRD-002", Name: "Office Chair",
		Category: "Furniture", Unit: "pcs", Stock: 25, MinStock: 5,
		Description: "Ergonomic office chair", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	},
}

// GetProducts returns all products
func GetProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data":  products,
		"total": len(products),
	})
}

// GetProduct returns a single product by ID
func GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	for _, p := range products {
		if p.ID == uint(id) {
			c.JSON(http.StatusOK, gin.H{"data": p})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

// CreateProductRequest holds data for creating a product
type CreateProductRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Unit        string `json:"unit" binding:"required"`
	Stock       int    `json:"stock"`
	MinStock    int    `json:"min_stock"`
	Description string `json:"description"`
}

// CreateProduct creates a new product
func CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newProduct := models.Product{
		ID:          uint(len(products) + 1),
		Code:        req.Code,
		Name:        req.Name,
		Category:    req.Category,
		Unit:        req.Unit,
		Stock:       req.Stock,
		MinStock:    req.MinStock,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	products = append(products, newProduct)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data":    newProduct,
	})
}

// UpdateProductRequest holds data for updating a product
type UpdateProductRequest struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Unit        string `json:"unit"`
	MinStock    int    `json:"min_stock"`
	Description string `json:"description"`
}

// UpdateProduct updates an existing product
func UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, p := range products {
		if p.ID == uint(id) {
			if req.Name != "" {
				products[i].Name = req.Name
			}
			if req.Category != "" {
				products[i].Category = req.Category
			}
			if req.Unit != "" {
				products[i].Unit = req.Unit
			}
			if req.MinStock > 0 {
				products[i].MinStock = req.MinStock
			}
			if req.Description != "" {
				products[i].Description = req.Description
			}
			products[i].UpdatedAt = time.Now()

			c.JSON(http.StatusOK, gin.H{
				"message": "Product updated successfully",
				"data":    products[i],
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

// DeleteProduct removes a product by ID
func DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	for i, p := range products {
		if p.ID == uint(id) {
			products = append(products[:i], products[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}
