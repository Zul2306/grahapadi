package controllers

import (
	"inventory-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// In-memory stock cards store (replace with database)
var stockCards = []models.StockCard{}

// GetStockCards returns all stock card records
func GetStockCards(c *gin.Context) {
	// Optional filter by product_id
	productIDStr := c.Query("product_id")
	if productIDStr != "" {
		productID, err := strconv.Atoi(productIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id"})
			return
		}
		var filtered []models.StockCard
		for _, s := range stockCards {
			if s.ProductID == uint(productID) {
				filtered = append(filtered, s)
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": filtered, "total": len(filtered)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  stockCards,
		"total": len(stockCards),
	})
}

// GetStockCard returns a single stock card record by ID
func GetStockCard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock card ID"})
		return
	}

	for _, s := range stockCards {
		if s.ID == uint(id) {
			c.JSON(http.StatusOK, gin.H{"data": s})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Stock card not found"})
}

// OpnameRequest holds data for a stock opname entry
type OpnameRequest struct {
	ProductID   uint   `json:"product_id" binding:"required"`
	ActualStock int    `json:"actual_stock" binding:"required"`
	Note        string `json:"note"`
}

// CreateOpname creates a new opname record and adjusts stock
func CreateOpname(c *gin.Context) {
	var req OpnameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the product
	var productIdx int = -1
	for i, p := range products {
		if p.ID == req.ProductID {
			productIdx = i
			break
		}
	}
	if productIdx == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	systemStock := products[productIdx].Stock
	difference := req.ActualStock - systemStock

	// Create opname record
	opname := models.Opname{
		ID:          uint(len(stockCards) + 1),
		ProductID:   req.ProductID,
		SystemStock: systemStock,
		ActualStock: req.ActualStock,
		Difference:  difference,
		Note:        req.Note,
		CreatedAt:   time.Now(),
	}

	// Create stock card entry for opname adjustment
	newCard := models.StockCard{
		ID:        uint(len(stockCards) + 1),
		ProductID: req.ProductID,
		Type:      "opname",
		Qty:       difference,
		Balance:   req.ActualStock,
		Note:      "Opname adjustment: " + req.Note,
		CreatedAt: time.Now(),
	}
	stockCards = append(stockCards, newCard)

	// Update product stock
	products[productIdx].Stock = req.ActualStock
	products[productIdx].UpdatedAt = time.Now()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Opname recorded successfully",
		"data":    opname,
	})
}
