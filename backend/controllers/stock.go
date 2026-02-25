package controllers

import (
	"inventory-backend/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateTransactionRequest holds data for creating a transaction
type CreateTransactionRequest struct {
	ProdukID uint   `json:"produk_id" binding:"required"`
	GudangID uint   `json:"gudang_id" binding:"required"`
	UserID   uint   `json:"user_id" binding:"required"`
	Tipe     string `json:"tipe" binding:"required,oneof=masuk keluar"` // masuk or keluar
	Jumlah   int    `json:"jumlah" binding:"required,gt=0"`
}

// CreateTransaction creates a new transaction and updates stock in stok_gudang
func CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, ok := getDB(c)
	if !ok {
		return
	}

	// Validate that product exists
	var produk models.Produk
	if err := db.First(&produk, req.ProdukID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Get or create stock record in stok_gudang
	var stock models.StockGudang
	result := db.Where("produk_id = ? AND gudang_id = ?", req.ProdukID, req.GudangID).First(&stock)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new stock record for this product
		stock = models.StockGudang{
			ProdukID: req.ProdukID,
			GudangID: req.GudangID,
			Jumlah:   0,
		}
		if err := db.Create(&stock).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stock record"})
			return
		}
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Calculate new stock quantity
	var newQuantity int
	if req.Tipe == "masuk" {
		newQuantity = stock.Jumlah + req.Jumlah
	} else { // keluar
		newQuantity = stock.Jumlah - req.Jumlah
		// Validate that outgoing doesn't exceed current stock
		if newQuantity < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient stock",
				"current_stock": stock.Jumlah,
				"requested": req.Jumlah,
			})
			return
		}
	}

	// Start transaction for atomic operations
	tx := db.Begin()

	// Create transaction record
	transaction := models.Transaction{
		ProdukID: req.ProdukID,
		GudangID: req.GudangID,
		UserID:   req.UserID,
		Tipe:     req.Tipe,
		Jumlah:   req.Jumlah,
		Tanggal:  time.Now(),
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Update stock quantity
	if err := tx.Model(&stock).Update("jumlah", newQuantity).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction created successfully",
		"data": gin.H{
			"transaction_id": transaction.ID,
			"product_id": transaction.ProdukID,
			"type": transaction.Tipe,
			"quantity": transaction.Jumlah,
			"new_stock": newQuantity,
			"created_at": transaction.CreatedAt,
		},
	})
}

// GetTransactions returns all transactions with optional filtering
func GetTransactions(c *gin.Context) {
	produkIDStr := c.Query("produk_id")
	tipe := c.Query("tipe") // filter by "in" or "out"

	db, ok := getDB(c)
	if !ok {
		return
	}

	var transactions []models.Transaction

	query := db

	if produkIDStr != "" {
		produkID, err := strconv.Atoi(produkIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid produk_id"})
			return
		}
		query = query.Where("produk_id = ?", produkID)
	}

	if tipe != "" {
		query = query.Where("tipe = ?", tipe)
	}

	if err := query.Order("created_at DESC").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  transactions,
		"total": len(transactions),
	})
}

// GetTransaction returns a single transaction by ID
func GetTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	db, ok := getDB(c)
	if !ok {
		return
	}

	var transaction models.Transaction

	if err := db.First(&transaction, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

// GetStockCards returns all stock card records
func GetStockCards(c *gin.Context) {
	// Optional filter by product_id
	productIDStr := c.Query("product_id")
	
	db, ok := getDB(c)
	if !ok {
		return
	}

	var stocks []models.StockGudang

	query := db

	if productIDStr != "" {
		productID, err := strconv.Atoi(productIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id"})
			return
		}
		query = query.Where("produk_id = ?", productID)
	}

	if err := query.Find(&stocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  stocks,
		"total": len(stocks),
	})
}

// GetStockCard returns a single stock record by ID
func GetStockCard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}

	db, ok := getDB(c)
	if !ok {
		return
	}

	var stock models.StockGudang

	if err := db.First(&stock, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stock record not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stock})
}

// CreateStockOpnameRequest holds data for creating a stock opname record
type CreateStockOpnameRequest struct {
	ProdukID   uint   `json:"produk_id" binding:"required"`
	StokSistem int    `json:"stok_sistem" binding:"required,gte=0"`
	StokFisik  int    `json:"stok_fisik" binding:"required,gte=0"`
	UserID     uint   `json:"user_id" binding:"required"`
	Keterangan string `json:"keterangan"`
}

// CreateStockOpname creates a new stock opname record
func CreateStockOpname(c *gin.Context) {
	var req CreateStockOpnameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, ok := getDB(c)
	if !ok {
		return
	}

	// Validate that product exists
	var produk models.Produk
	if err := db.First(&produk, req.ProdukID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Log user lookup for debugging
	log.Printf("🔍 Checking user_id: %d", req.UserID)

	// Calculate difference
	selisih := req.StokFisik - req.StokSistem

	// Create opname record
	opname := models.StockOpname{
		ProdukID:       req.ProdukID,
		StokSistem:     req.StokSistem,
		StokFisik:      req.StokFisik,
		Selisih:        selisih,
		UserID:         req.UserID,
		Keterangan:     req.Keterangan,
		SudahDisetujui: false,
		Tanggal:        time.Now(),
	}

	if err := db.Create(&opname).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    opname,
		"message": "Stock opname record created successfully",
	})
}

// GetStockOpnames returns all stock opname records
func GetStockOpnames(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	var opnames []models.StockOpname

	if err := db.Order("id ASC").Find(&opnames).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  opnames,
		"total": len(opnames),
	})
}

// GetStockOpnameByID returns a single stock opname record by ID
func GetStockOpnameByID(c *gin.Context) {
	id := c.Param("id")

	db, ok := getDB(c)
	if !ok {
		return
	}

	var opname models.StockOpname
	if err := db.First(&opname, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stock opname record not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": opname})
}

// GetGudangs returns all warehouses
func GetGudangs(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	var gudangs []models.Gudang

	if err := db.Find(&gudangs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  gudangs,
		"total": len(gudangs),
	})
}

// GetProductTotalStock returns total stock of a product across all warehouses
func GetProductTotalStock(c *gin.Context) {
	produkIDStr := c.Param("produk_id")
	log.Printf("🔍 GET /products/stock/:produk_id called with produk_id: %s", produkIDStr)

	// Convert string to uint
	produkID, err := strconv.ParseUint(produkIDStr, 10, 32)
	if err != nil {
		log.Printf("❌ Invalid produk_id format: %s", produkIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid produk_id format"})
		return
	}

	db, ok := getDB(c)
	if !ok {
		return
	}

	// Validate that product exists
	var produk models.Produk
	if err := db.First(&produk, produkID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("❌ Product not found: %d", produkID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Printf("✅ Product found: %s (%s)", produk.NamaBarang, produk.KodeBarang)

	// Debug: Check stok_gudang table contents
	var debugStocks []models.StockGudang
	db.Where("produk_id = ?", uint(produkID)).Find(&debugStocks)
	log.Printf("🔎 Found %d stock records for produk_id %d:", len(debugStocks), uint(produkID))
	for _, stock := range debugStocks {
		log.Printf("   - ID: %d, ProdukID: %d, GudangID: %d, Jumlah: %d", stock.ID, stock.ProdukID, stock.GudangID, stock.Jumlah)
	}

	// Calculate total stock across all warehouses using simpler query
	var totalStock int64 = 0
	result := db.Model(&models.StockGudang{}).
		Where("produk_id = ?", uint(produkID)).
		Select("COALESCE(SUM(jumlah), 0) as total").
		Scan(&totalStock)
	
	if result.Error != nil {
		log.Printf("❌ Error calculating stock: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	log.Printf("📦 Total stock for produk_id %s: %d", produkID, totalStock)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"produk_id":   produk.ID,
			"nama_barang": produk.NamaBarang,
			"total_stock": totalStock,
		},
	})
}
