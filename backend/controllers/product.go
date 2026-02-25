package controllers

import (
	"errors"
	"inventory-backend/config"
	"inventory-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getDB(c *gin.Context) (*gorm.DB, bool) {
	cfgValue, ok := c.Get("config")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Config not found in context"})
		return nil, false
	}

	cfg, ok := cfgValue.(*config.Config)
	if !ok || cfg.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return nil, false
	}

	return cfg.DB, true
}

// GetProducts returns all products
func GetProducts(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	var items []models.Produk
	if err := db.Order("id ASC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"total": len(items),
	})
}

// GetProduct returns a single product by ID
func GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	db, ok := getDB(c)
	if !ok {
		return
	}

	var item models.Produk
	if err := db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// CreateProductRequest holds data for creating a product
type CreateProductRequest struct {
	KodeBarang  string  `json:"kode_barang" binding:"required"`
	NamaBarang  string  `json:"nama_barang" binding:"required"`
	JenisBarang string  `json:"jenis_barang" binding:"required"`
	Satuan      string  `json:"satuan" binding:"required"`
	StokMinimal int     `json:"stok_minimal"`
	BeratKg     float64 `json:"berat_kg"`
}

// CreateProduct creates a new product
func CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, ok := getDB(c)
	if !ok {
		return
	}

	newProduct := models.Produk{
		KodeBarang:  req.KodeBarang,
		NamaBarang:  req.NamaBarang,
		JenisBarang: req.JenisBarang,
		Satuan:      req.Satuan,
		StokMinimal: req.StokMinimal,
		BeratKg:     req.BeratKg,
	}
	if err := db.Create(&newProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data":    newProduct,
	})
}

// UpdateProductRequest holds data for updating a product
type UpdateProductRequest struct {
	KodeBarang  string  `json:"kode_barang"`
	NamaBarang  string  `json:"nama_barang"`
	JenisBarang string  `json:"jenis_barang"`
	Satuan      string  `json:"satuan"`
	StokMinimal *int    `json:"stok_minimal"`
	BeratKg     *float64 `json:"berat_kg"`
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

	db, ok := getDB(c)
	if !ok {
		return
	}

	var item models.Produk
	if err := db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	if req.KodeBarang != "" {
		item.KodeBarang = req.KodeBarang
	}
	if req.NamaBarang != "" {
		item.NamaBarang = req.NamaBarang
	}
	if req.JenisBarang != "" {
		item.JenisBarang = req.JenisBarang
	}
	if req.Satuan != "" {
		item.Satuan = req.Satuan
	}
	if req.StokMinimal != nil {
		item.StokMinimal = *req.StokMinimal
	}
	if req.BeratKg != nil {
		item.BeratKg = *req.BeratKg
	}

	if err := db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"data":    item,
	})
}

// DeleteProduct removes a product by ID
func DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	db, ok := getDB(c)
	if !ok {
		return
	}

	result := db.Delete(&models.Produk{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
