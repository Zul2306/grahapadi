# Setup Mailtrap untuk Testing Email

## Langkah-langkah:

### 1. Buat Akun Mailtrap (Gratis)

- Buka https://mailtrap.io/
- Sign up dengan email atau GitHub
- Verifikasi email

### 2. Dapatkan SMTP Credentials

- Login ke Mailtrap
- Pilih **Email Testing** → **Inboxes**
- Klik inbox yang ingin digunakan (atau buat baru)
- Pilih tab **SMTP Settings**
- Pilih integration "Node.js - Nodemailer" atau "Other"
- Copy credentials:
  - **Host:** sandbox.smtp.mailtrap.io
  - **Port:** 2525 (atau 587, 465)
  - **Username:** (contoh: abc123def456)
  - **Password:** (contoh: xyz789uvw012)

### 3. Setup di Backend

Buat file `.env` di folder `backend/` (copy dari `.env.example`):

```env
# Server Configuration
PORT=8888

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=password
DB_NAME=gudang

# JWT Secret
JWT_SECRET=your-secret-key-change-this

# Email Configuration - MAILTRAP
SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USERNAME=paste-your-mailtrap-username-here
SMTP_PASSWORD=paste-your-mailtrap-password-here
FROM_EMAIL=noreply@inventory.com
FROM_NAME=Inventory System
```

### 4. Restart Backend

```bash
cd backend
go run .
```

### 5. Test Forgot Password

1. Buka frontend: http://localhost:3000/login
2. Klik "Forgot password?"
3. Input email user yang sudah terdaftar (misal: admin@inventory.com)
4. Submit

### 6. Cek Email di Mailtrap

- Buka Mailtrap inbox di browser
- Email akan muncul di inbox Mailtrap
- Klik email untuk melihat isi
- Klik tombol "Reset Password" di email
- Atau copy link reset password

### 7. Reset Password

- Link akan membuka: http://localhost:3000/reset-password?token=xxx
- Input password baru
- Submit
- Login dengan password baru

## Keuntungan Mailtrap:

✅ **Testing Real Email** - Lihat tampilan email yang sebenarnya
✅ **No Spam** - Email tidak dikirim ke inbox asli
✅ **HTML Preview** - Lihat rendering HTML email
✅ **Email Validation** - Check spam score, HTML issues
✅ **Multiple Inboxes** - Test berbagai scenario
✅ **Free Plan** - 500 emails/month gratis

## Production vs Development:

| Mode                        | Configuration       | Email goes to       |
| --------------------------- | ------------------- | ------------------- |
| **Development (No Config)** | Tanpa SMTP settings | Console log backend |
| **Testing (Mailtrap)**      | Mailtrap SMTP       | Mailtrap inbox      |
| **Production (Gmail/etc)**  | Real SMTP           | User's real email   |

## Troubleshooting:

### Email tidak muncul di Mailtrap?

1. Cek credentials di `.env` sudah benar
2. Restart backend setelah update `.env`
3. Cek console backend untuk error messages
4. Pastikan firewall tidak block port 2525

### Error "authentication failed"?

- Pastikan username dan password benar
- Try copy-paste ulang dari Mailtrap dashboard
- Pastikan tidak ada spasi di awal/akhir credentials

### Email masih di console?

- Pastikan file `.env` ada di folder `backend/`
- Pastikan SMTP_USERNAME dan SMTP_PASSWORD tidak kosong
- Restart backend

## Screenshot Mailtrap:

Tampilan email di Mailtrap akan menunjukkan:

- Subject: "Password Reset Request - Inventory System"
- From: Inventory System <noreply@inventory.com>
- Template HTML dengan logo, tombol, dan styling
- Link reset password yang clickable

## Note:

Mailtrap hanya untuk **testing/development**. Untuk production, gunakan:

- Gmail SMTP
- SendGrid
- Amazon SES
- Mailgun
- atau email service provider lainnya
