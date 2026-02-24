# Fitur Lupa Password dengan Email - Inventory System

## ✅ Implementasi Selesai

### Backend (Go)

1. **Model PasswordReset**
   - `models/models.go` - Model untuk menyimpan token reset password
   - Token expire dalam 1 jam
   - Flag `used` untuk mencegah reuse token

2. **Email Service**
   - `config/email.go` - Service untuk mengirim email via SMTP
   - Support Gmail dan SMTP lainnya
   - Development mode: email di-log ke console (tidak perlu SMTP config)

3. **Controllers**
   - `POST /api/v1/auth/forgot-password` - Request reset password
   - `POST /api/v1/auth/reset-password` - Reset password dengan token
   - Email HTML template dengan tombol reset yang bagus

4. **Routes**
   - Endpoint forgot-password dan reset-password sudah ditambahkan

### Frontend (React)

1. **ForgotPassword.js** - Updated
   - Integrasi dengan API backend
   - Error handling
   - Success message dengan email yang dimasukkan

2. **ResetPassword.js** - New Page
   - Ambil token dari URL query parameter
   - Form untuk password baru dengan konfirmasi
   - Show/hide password toggle
   - Validasi password match
   - Success page dengan redirect ke login

3. **Routing**
   - `/forgot-password` - Halaman request reset
   - `/reset-password?token=xxx` - Halaman reset dengan token

## 🔧 Setup & Testing

### 1. Update Backend Dependencies

```bash
cd backend
go mod tidy
```

### 2. Setup Environment (Optional untuk Production)

Buat file `.env` di folder `backend/`:

```env
# Email Configuration (Leave empty for development)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=noreply@inventory.com
FROM_NAME=Inventory System
```

**Note:** Untuk development, email akan di-log ke console. Tidak perlu setup SMTP.

### 3. Jalankan Backend

```bash
cd backend
go run .
```

### 4. Jalankan Frontend

```bash
cd frontend
npm start
```

## 📧 Email Template

Email yang dikirim berisi:

- Header dengan logo
- Greeting dengan nama user
- Tombol "Reset Password" yang mengarah ke halaman reset
- Link text sebagai alternatif
- Warning bahwa link expire dalam 1 jam
- Footer dengan copyright

## 🔄 Flow Lengkap

### 1. Forgot Password Flow:

```
User → Klik "Forgot password?" →
Input email → Submit →
Backend cek user → Generate token →
Simpan ke database → Kirim email →
User dapat email → Klik tombol/link
```

### 2. Reset Password Flow:

```
User klik link → Buka /reset-password?token=xxx →
Input password baru → Confirm password →
Submit → Backend validasi token →
Update password → Redirect ke login
```

## 🧪 Testing

### Test Forgot Password

```bash
curl -X POST http://localhost:8888/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@inventory.com"}'
```

Response:

```json
{
  "message": "If the email exists, a password reset link has been sent"
}
```

### Test Reset Password

```bash
curl -X POST http://localhost:8888/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token":"your-token-here",
    "new_password":"newpassword123"
  }'
```

Response:

```json
{
  "message": "Password has been reset successfully"
}
```

## 📝 Catatan Penting

### Development Mode

- Tanpa SMTP config, email akan di-print ke console backend
- Lihat console backend untuk melihat email yang "terkirim"
- Token dan link reset akan terlihat di console

### Production Mode

- Setup SMTP credentials di `.env`
- Untuk Gmail, gunakan "App Password" bukan password biasa
- Enable "Less secure app access" atau gunakan OAuth2

### Security Features

✅ Token expire dalam 1 jam
✅ Token hanya bisa dipakai sekali
✅ Password di-hash dengan bcrypt
✅ Email validation
✅ Tidak reveal apakah email terdaftar atau tidak

## 🔐 Gmail Setup (Production)

1. Go to Google Account Settings
2. Security → 2-Step Verification (enable)
3. App Passwords → Generate new
4. Use generated password di `SMTP_PASSWORD`

## 🎯 URL Examples

- Forgot Password: `http://localhost:3000/forgot-password`
- Reset Password: `http://localhost:3000/reset-password?token=abc123...`
- Login: `http://localhost:3000/login`

## 📊 Database Table

Table `password_resets`:

- `id` - Primary key
- `email` - Email user
- `token` - Reset token (unique)
- `expires_at` - Waktu expire
- `used` - Boolean flag
- `created_at` - Timestamp

## ✨ Features

✅ Email HTML template yang menarik
✅ Token security dengan expiration
✅ Password validation (min 6 characters)
✅ Password confirmation match
✅ Show/hide password toggle
✅ Error handling di semua endpoint
✅ Success messages yang jelas
✅ Development & production mode
✅ Mobile responsive design
✅ Loading states
✅ Security best practices

## 🚀 Next Steps (Optional)

1. Implement rate limiting untuk forgot password
2. Add CAPTCHA untuk mencegah spam
3. Email queue system untuk async sending
4. Password strength meter
5. Send email notification saat password berhasil diubah
6. Add email verification untuk new users
