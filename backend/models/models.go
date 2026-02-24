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
