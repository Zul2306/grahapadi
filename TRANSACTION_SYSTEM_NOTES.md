# Transaction System Implementation

## Overview

Implementasi sistem transaksi lengkap untuk inventory management dengan logika incoming/outgoing goods yang otomatis update stock di `stok_gudang` table.

## Backend Changes

### 1. Models (`backend/models/models.go`)

Ditambahkan 2 model baru:

**StockGudang** - Memetakan ke tabel `stok_gudang`

```go
type StockGudang struct {
    ID        uint      `gorm:"primaryKey;column:id" json:"id"`
    ProdukID  uint      `gorm:"index;column:produk_id" json:"produk_id"`
    GudangID  uint      `gorm:"index;column:gudang_id" json:"gudang_id"`
    Jumlah    int       `gorm:"default:0;column:jumlah" json:"jumlah"`
    CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
```

**Transaction** - Memetakan ke tabel `transaksi`

```go
type Transaction struct {
    ID                uint      `gorm:"primaryKey;column:id" json:"id"`
    ProdukID          uint      `gorm:"index;column:produk_id" json:"produk_id"`
    Tipe              string    `gorm:"type:varchar(20);column:tipe" json:"tipe"` // "in" or "out"
    Jumlah            int       `gorm:"column:jumlah" json:"jumlah"`
    PenanggungJawab   string    `gorm:"type:varchar(255);column:penanggung_jawab" json:"penanggung_jawab"`
    Tanggal           time.Time `gorm:"column:tanggal" json:"tanggal"`
    CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
    UpdatedAt         time.Time `gorm:"column:updated_at" json:"updated_at"`
}
```

### 2. Controller (`backend/controllers/stock.go`)

Refactored untuk implementasi transaction management:

#### CreateTransaction()

- Accept request dengan: produk_id, tipe (in/out), jumlah, penanggung_jawab
- Validasi produk exist
- Get atau create record di stok_gudang(gudang_id=1)
- Kalkulasi jumlah baru:
  - Jika `tipe='in'`: tambah stock
  - Jika `tipe='out'`: kurangi stock, validasi tidak negative
- Database transaction (atomic):
  - Insert transaksi record
  - Update jumlah di stok_gudang
  - Rollback jika ada error

#### GetTransactions()

- Optional filter by produk_id dan tipe
- Sort by created_at DESC

#### GetTransaction(id)

- Get single transaction by ID

#### GetStockCards() / GetStockCard()

- Refactored dari in-memory ke database-backed
- Query stok_gudang table

### 3. Routes (`backend/routes/routes.go`)

Tambah transaction endpoints di protected `/stock` group:

```go
POST   /stock/transactions       - Create transaction + update stock
GET    /stock/transactions       - List transactions (filterable)
GET    /stock/transactions/:id   - Get single transaction
```

## Frontend Changes

### Transactions.js

Complete rewrite dari template UI menjadi fungsional transaction manager:

#### Features:

1. **Toggle Incoming/Outgoing**
   - State: `transactionType` (in/out)
   - Update daftar transaksi saat switch

2. **Transaction Table**
   - Display ID, Product Name, Type, Quantity, Responsible Person, Date
   - Colorized type badge (green for incoming, orange for outgoing)
   - Search filter by product name atau responsible person
   - Delete action (UI-only, delete endpoint dapat ditambah nanti)

3. **Create Transaction Modal**
   - Form fields: Product dropdown, Quantity (number), Responsible Person
   - Validation: all fields required, quantity > 0
   - Error handling dengan display di modal
   - Submit ke `/stock/transactions` dengan payload:
     ```json
     {
       "produk_id": 1,
       "tipe": "in",
       "jumlah": 10,
       "penanggung_jawab": "John Doe"
     }
     ```

4. **Data Fetching**
   - `fetchTransactions()` - GET /stock/transactions dengan filter tipe
   - `fetchProducts()` - GET /products (untuk dropdown)
   - Auto-refresh setelah create transaction
   - Automatic retry jika pindah ke incoming/outgoing

5. **State Management**
   ```javascript
   - transactionType: 'in' or 'out'
   - transactions: array dari transaksi
   - products: dropdown data
   - loading: show loading state
   - error: display error messages
   - searchTerm: untuk filter table
   - showCreateModal: toggle modal visibility
   ```

## API Contract

### Create Transaction

```
POST /api/v1/stock/transactions
Authorization: Bearer {token}

Request:
{
  "produk_id": 1,
  "tipe": "in",
  "jumlah": 10,
  "penanggung_jawab": "Staff Name"
}

Response (200):
{
  "message": "Transaction created successfully",
  "data": {
    "transaction_id": 5,
    "product_id": 1,
    "type": "in",
    "quantity": 10,
    "new_stock": 25,
    "created_at": "2026-02-25T01:45:00Z"
  }
}

Response (400) - Insufficient Stock:
{
  "error": "Insufficient stock",
  "current_stock": 5,
  "requested": 10
}
```

### Get Transactions

```
GET /api/v1/stock/transactions?tipe=in&produk_id=1

Response:
{
  "data": [
    {
      "id": 1,
      "produk_id": 1,
      "tipe": "in",
      "jumlah": 10,
      "penanggung_jawab": "Staff",
      "tanggal": "2026-02-25T01:40:00Z",
      "created_at": "2026-02-25T01:40:00Z",
      "updated_at": "2026-02-25T01:40:00Z"
    }
  ],
  "total": 1
}
```

## Key Features

### Stock Management

- **Automatic Stock Update**: Setiap transaksi incoming/outgoing langsung update `stok_gudang.jumlah`
- **Validation**: Outgoing tidak boleh exceed current stock
- **Atomic Operations**: DB transaction ensures consistency

### User Interface

- **Real-time Filtering**: Search live filter table
- **Type Toggle**: Easy switch antara incoming/outgoing
- **Modal Form**: Clean form dalam modal untuk create
- **Status Badges**: Color-coded transaction types
- **Timestamps**: Formatted dates untuk readability

### Error Handling

- Validation errors di form (required fields, positive quantity)
- API error messages dipropagasi ke UI
- Insufficient stock error detail (current vs requested)
- Loading states dan try-again logic

## Usage Flow

1. **User Login** → Navigate ke Transactions page
2. **Select Type** → Click Incoming/Outgoing toggle
3. **View History** → Table shows transactions untuk selected type
4. **Create New** → Click "New" button
5. **Fill Form** → Select product, enter quantity & responsible person
6. **Submit** → Click Create → API updates stock + records transaction
7. **Confirm** → Modal closes, table refreshes dengan entry baru

## Testing

### Backend Test Command:

```bash
# Terminal 1 - Start Backend
cd backend && go run main.go

# Terminal 2 - Test Create Transaction
curl -X POST http://localhost:8888/api/v1/stock/transactions \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"produk_id":1,"tipe":"in","jumlah":10,"penanggung_jawab":"Tester"}'

# Test List Transactions
curl http://localhost:8888/api/v1/stock/transactions?tipe=in \
  -H "Authorization: Bearer {token}"
```

### Frontend Test:

1. Open http://localhost:3000
2. Login with credentials
3. Navigate to Transactions page
4. Try toggle between Incoming/Outgoing
5. Click "New" button
6. Fill form and submit
7. Verify table updates automatically

## Database Schema

Asumsi table sudah exist dengan struktur:

- `transaksi`: id, produk_id, tipe, jumlah, penanggung_jawab, tanggal, created_at, updated_at
- `stok_gudang`: id, produk_id, gudang_id, jumlah, created_at, updated_at
- `produk`: existing keys untuk join

## Future Enhancements

- [ ] Delete transaction endpoint (need cascade logic on stok_gudang)
- [ ] Edit transaction capability
- [ ] Batch transaction import
- [ ] Transaction history/audit trail
- [ ] Stock movement reports
- [ ] Warehouse-aware stock management (gudang_id selection)
- [ ] Barcode/QR code scanner integration
- [ ] Stock alerts untuk stok_minimal
