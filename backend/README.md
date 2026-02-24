# Inventory Backend API

Backend REST API untuk sistem inventory menggunakan **Go** + **Gin Framework**.

## Struktur Folder

```
backend/
├── main.go                 # Entry point aplikasi
├── config/
│   └── config.go           # Konfigurasi aplikasi (env vars)
├── routes/
│   └── routes.go           # Definisi semua route API
├── middleware/
│   ├── cors.go             # CORS middleware
│   └── auth.go             # JWT Auth middleware
├── controllers/
│   ├── auth.go             # Handler login & logout
│   ├── user.go             # Handler CRUD user
│   ├── product.go          # Handler CRUD produk
│   └── stock.go            # Handler kartu stok & opname
├── models/
│   └── models.go           # Struct data model
├── .env.example            # Contoh environment variables
└── go.mod / go.sum         # Go module files
```

## Menjalankan Server

```bash
# Copy env file
cp .env.example .env

# Install dependencies
go mod tidy

# Jalankan server
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## Endpoint API

### Auth (Public)
| Method | Endpoint           | Keterangan        |
|--------|--------------------|-------------------|
| POST   | /api/v1/auth/login | Login user        |
| POST   | /api/v1/auth/logout| Logout user       |

### Users (Protected)
| Method | Endpoint           | Keterangan        |
|--------|--------------------|-------------------|
| GET    | /api/v1/users      | List semua user   |
| GET    | /api/v1/users/:id  | Detail user       |
| POST   | /api/v1/users      | Buat user baru    |
| DELETE | /api/v1/users/:id  | Hapus user        |

### Products (Protected)
| Method | Endpoint              | Keterangan        |
|--------|-----------------------|-------------------|
| GET    | /api/v1/products      | List semua produk |
| GET    | /api/v1/products/:id  | Detail produk     |
| POST   | /api/v1/products      | Buat produk baru  |
| PUT    | /api/v1/products/:id  | Update produk     |
| DELETE | /api/v1/products/:id  | Hapus produk      |

### Stock & Opname (Protected)
| Method | Endpoint              | Keterangan          |
|--------|-----------------------|---------------------|
| GET    | /api/v1/stock         | List kartu stok     |
| GET    | /api/v1/stock/:id     | Detail kartu stok   |
| POST   | /api/v1/stock/opname  | Input data opname   |

## Contoh Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@inventory.com","password":"admin123"}'
```

## Tech Stack
- **Go 1.21+**
- **Gin** - HTTP framework
- **JWT** - Autentikasi (siap diintegrasikan)
- **PostgreSQL** - Database (siap diintegrasikan)
