package models

import "time"

// User represents a system user
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	Email     string    `gorm:"type:varchar(100);unique" json:"email"`
	Password  string    `gorm:"type:varchar(255)" json:"-"`
	Role      string    `gorm:"type:varchar(20)" json:"role"` // admin, staff
	CreatedAt time.Time `json:"created_at"`
}

// PasswordReset represents a password reset token
type PasswordReset struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"type:varchar(100);index" json:"email"`
	Token     string    `gorm:"type:varchar(255);unique" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// Product represents an inventory item
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"type:varchar(100);unique" json:"code"`
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	Category    string    `gorm:"type:varchar(100)" json:"category"`
	Unit        string    `gorm:"type:varchar(50)" json:"unit"`
	Stock       int       `gorm:"default:0" json:"stock"`
	MinStock    int       `gorm:"default:0" json:"min_stock"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Produk represents inventory item data from table "produk"
type Produk struct {
	ID          uint    `gorm:"primaryKey;column:id" json:"id"`
	KodeBarang  string  `gorm:"type:varchar(100);column:kode_barang" json:"kode_barang"`
	NamaBarang  string  `gorm:"type:varchar(255);column:nama_barang" json:"nama_barang"`
	JenisBarang string  `gorm:"type:varchar(100);column:jenis_barang" json:"jenis_barang"`
	Satuan      string  `gorm:"type:varchar(50);column:satuan" json:"satuan"`
	StokMinimal int     `gorm:"default:0;column:stok_minimal" json:"stok_minimal"`
	BeratKg     float64 `gorm:"default:0;column:berat_kg" json:"berat_kg"`
}

func (Produk) TableName() string {
	return "produk"
}

// StockCard represents a stock movement record
type StockCard struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index" json:"product_id"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Type      string    `gorm:"type:varchar(20)" json:"type"` // in, out, opname
	Qty       int       `json:"qty"`
	Balance   int       `json:"balance"`
	Note      string    `gorm:"type:text" json:"note"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// Opname represents a stock opname session
type Opname struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProductID   uint      `gorm:"index" json:"product_id"`
	Product     *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SystemStock int       `json:"system_stock"`
	ActualStock int       `json:"actual_stock"`
	Difference  int       `json:"difference"`
	Note        string    `gorm:"type:text" json:"note"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// StockGudang represents warehouse stock levels mapped to "stok_gudang" table
type StockGudang struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	ProdukID  uint      `gorm:"index;column:produk_id" json:"produk_id"`
	GudangID  uint      `gorm:"index;column:gudang_id" json:"gudang_id"`
	Jumlah    int       `gorm:"default:0;column:jumlah" json:"jumlah"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (StockGudang) TableName() string {
	return "stok_gudang"
}

// Transaction represents inventory movement transaction mapped to "transaksi" table
type Transaction struct {
	ID              uint      `gorm:"primaryKey;column:id" json:"id"`
	ProdukID        uint      `gorm:"index;column:produk_id" json:"produk_id"`
	GudangID        uint      `gorm:"index;column:gudang_id" json:"gudang_id"`
	UserID          uint      `gorm:"index;column:user_id" json:"user_id"`
	Tipe            string    `gorm:"type:varchar(20);column:tipe" json:"tipe"` // "masuk" or "keluar"
	Jumlah          int       `gorm:"column:jumlah" json:"jumlah"`
	Tanggal         time.Time `gorm:"column:tanggal" json:"tanggal"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Transaction) TableName() string {
	return "transaksi"
}

// Gudang represents warehouse/storage location mapped to "gudang" table
type Gudang struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	Nama      string    `gorm:"type:varchar(255);column:nama" json:"nama"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Gudang) TableName() string {
	return "gudang"
}

// StockOpname represents stock opname records mapped to "stok_opname" table
type StockOpname struct {
	ID              uint      `gorm:"primaryKey;column:id" json:"id"`
	ProdukID        uint      `gorm:"index;column:produk_id" json:"produk_id"`
	StokSistem      int       `gorm:"column:stok_sistem" json:"stok_sistem"`
	StokFisik       int       `gorm:"column:stok_fisik" json:"stok_fisik"`
	Selisih         int       `gorm:"column:selisih" json:"selisih"`
	UserID          uint      `gorm:"index;column:user_id" json:"user_id"`
	Keterangan      string    `gorm:"type:text;column:keterangan" json:"keterangan"`
	SudahDisetujui  bool      `gorm:"default:false;column:sudah_disetujui" json:"sudah_disetujui"`
	Tanggal         time.Time `gorm:"column:tanggal" json:"tanggal"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (StockOpname) TableName() string {
	return "stok_opname"
}
