# Login Integration - Inventory System

## Perubahan yang dilakukan:

### Frontend (React)

1. **Buat Service Layer**
   - `src/services/api.js` - Base HTTP client untuk semua API calls
   - `src/services/authService.js` - Service khusus untuk authentication

2. **Buat Auth Context**
   - `src/context/AuthContext.js` - React Context untuk state management user authentication
   - Menyediakan hooks `useAuth()` untuk akses ke login, logout, register functions

3. **Update Login Component**
   - `src/views/auth/Login.js` - Integrasi dengan backend API
   - Menambahkan state management untuk form data
   - Menambahkan error handling dan loading state
   - Menyimpan token dan user data ke localStorage

4. **Protected Route**
   - `src/components/ProtectedRoute.js` - Component untuk protect routes yang memerlukan authentication
   - Auto redirect ke `/login` jika user belum authenticated

5. **Update App.js**
   - Wrap aplikasi dengan `AuthProvider`
   - Protect semua routes kecuali login dan forgot password

6. **Environment Configuration**
   - `.env` - Konfigurasi API URL (http://localhost:8888/api/v1)

### Backend (Go)

1. **Update Login Controller**
   - `controllers/auth.go` - Update fungsi `Login` untuk validate dari database
   - Menggunakan bcrypt untuk compare password hash
   - Return user data dan token setelah berhasil login

2. **Update Routes**
   - `routes/routes.go` - Inject config ke login handler
   - Endpoint `/api/v1/auth/login` sekarang POST dengan body JSON

3. **Fix Endpoints**
   - Sudah menambahkan endpoint `/api/v1/auth/register`
   - CORS middleware sudah configured

## Cara Menggunakan:

### 1. Jalankan Backend

```bash
cd E:\inventory\backend
go run .
```

Backend akan berjalan di `http://localhost:8888`

### 2. Register User Baru (via API)

```bash
curl -X POST http://localhost:8888/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin User",
    "email": "admin@inventory.com",
    "password": "admin123",
    "role": "admin"
  }'
```

### 3. Jalankan Frontend

```bash
cd E:\inventory\frontend
npm start
```

Frontend akan berjalan di `http://localhost:3000`

### 4. Login

1. Buka browser ke `http://localhost:3000/login`
2. Masukkan credentials:
   - Email: admin@inventory.com
   - Password: admin123
3. Klik "Sign In"

## Flow Authentication:

1. **Login Process:**
   - User submit form di Login.js
   - Call `authService.login(email, password)`
   - API request ke `POST /api/v1/auth/login`
   - Backend validate email dan password dari database
   - Backend return token dan user data
   - Frontend simpan token dan user data ke localStorage
   - Context update state `user`
   - Redirect ke dashboard (`/`)

2. **Protected Routes:**
   - Setiap akses route yang protected, `ProtectedRoute` component check authentication
   - Jika `isAuthenticated` false, redirect ke `/login`
   - Jika true, render component yang diminta

3. **Logout Process:**
   - Call `authService.logout()`
   - Clear token dan user data dari localStorage
   - Update context state `user` menjadi null
   - Redirect ke `/login`

## API Endpoints:

- `POST /api/v1/auth/login` - Login user

  ```json
  Request: { "email": "user@example.com", "password": "password123" }
  Response: { "message": "Login successful", "token": "...", "user": {...} }
  ```

- `POST /api/v1/auth/logout` - Logout user

  ```json
  Response: { "message": "Logout successful" }
  ```

- `POST /api/v1/auth/register` - Register new user
  ```json
  Request: {
    "name": "User Name",
    "email": "user@example.com",
    "password": "password123",
    "role": "admin" // or "staff"
  }
  Response: { "message": "User registered successfully", "data": {...} }
  ```

## Testing:

### Test Register

```bash
curl -X POST http://localhost:8888/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"test123","role":"staff"}'
```

### Test Login

```bash
curl -X POST http://localhost:8888/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'
```

## Troubleshooting:

1. **CORS Error:**
   - Pastikan backend sudah running
   - Check `middleware/cors.go` sudah allow origin frontend

2. **401 Unauthorized:**
   - Pastikan email dan password benar
   - Check user sudah terdaftar di database

3. **Network Error:**
   - Pastikan `.env` file ada di frontend root folder
   - Pastikan `REACT_APP_API_URL` pointing ke backend yang benar
   - Restart frontend development server setelah ubah `.env`

4. **Can't login:**
   - Register user dulu sebelum login
   - Password di-hash dengan bcrypt, tidak bisa raw password

## Next Steps:

1. Implement real JWT token generation (sekarang masih mock token)
2. Add token refresh mechanism
3. Add remember me functionality
4. Implement forgot password flow
5. Add email verification
