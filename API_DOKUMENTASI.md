# Dokumentasi API User Management (GoFiber MVC)

Seluruh endpoint menggunakan format JSON dan berjalan di server default: `http://localhost:3001`

---

## 1. Register
**POST** `/register`

**Body:**
```json
{
  "email": "user@email.com",
  "password": "password123"
}
```
**Response:**
- 201: Data user baru
- 400: Email sudah terdaftar atau format salah

---

## 2. Login (JWT)
**POST** `/login`

**Body:**
```json
{
  "email": "user@email.com",
  "password": "password123"
}
```
**Response:**
- 200: `{ "token": "<jwt_token>" }`
- 401: Email/password salah

---

## 3. Request Verifikasi Email
**POST** `/request-verification`

**Body:**
```json
{
  "email": "user@email.com"
}
```
**Response:**
- 200: Kode verifikasi dikirim ke email (simulasi/log)
- 400: Email tidak ditemukan/format salah

---

## 4. Verifikasi Email
**POST** `/verify-email`

**Body:**
```json
{
  "email": "user@email.com",
  "code": "123456"
}
```
**Response:**
- 200: Email berhasil diverifikasi
- 400: Kode salah/kadaluarsa

---

## 5. Profile (Protected, butuh JWT)
**GET** `/profile`

**Header:**
```
Authorization: <jwt_token>
```
**Response:**
- 200: Data user
- 401: Token tidak valid/tidak ditemukan
- 404: User tidak ditemukan

---

## 6. Update Profile
**PUT** `/profile`

**Header:**
```
Authorization: <jwt_token>
```
**Body:**
```json
{
  "email": "user@email.com",
  "status": "active"
}
```
**Response:**
- 200: Data user yang telah diperbarui
- 400: Data tidak valid
- 401: Token tidak valid/tidak ditemukan
- 404: User tidak ditemukan

---

## 7. Logout
**POST** `/logout`

**Header:**
```
Authorization: <jwt_token>
```
**Response:**
- 200: Logout berhasil
- 401: Token tidak valid/tidak ditemukan

---

## 8. Request Password Recovery
**POST** `/request-password-recovery`

**Body:**
```json
{
  "email": "user@email.com"
}
```
**Response:**
- 200: Kode recovery dikirim ke email (simulasi/log)
- 400: Email tidak ditemukan/format salah

---

## 9. Reset Password
**POST** `/reset-password`

**Body:**
```json
{
  "email": "user@email.com",
  "code": "123456",
  "new_password": "passwordBaru"
}
```
**Response:**
- 200: Password berhasil direset
- 400: Kode salah/kadaluarsa atau data tidak valid

---

## 10. Custom Attributes (Protected)
**GET** `/profile/custom`
**PUT** `/profile/custom`

**Header:**
```
Authorization: <jwt_token>
```
**GET Response:**
- 200: `{ "custom_attributes": { ... } }`
- 401: Token tidak valid/tidak ditemukan
- 404: User tidak ditemukan

**PUT Body:**
```json
{
  "custom_attributes": {
    "key": "value"
  }
}
```
**PUT Response:**
- 200: Custom attributes berhasil diupdate
- 400: Data tidak valid
- 401: Token tidak valid/tidak ditemukan
- 404: User tidak ditemukan

---

## 11. Login TOTP (2FA Step 2)
**POST** `/login/totp`

**Body:**
```json
{
  "email": "user@email.com",
  "totp_code": "123456"
}
```
**Response:**
- 200: Login 2FA berhasil, return token
- 400: Kode TOTP salah/kadaluarsa
- 401: Email/password salah

---

## 12. TOTP (2FA, Protected)
Semua endpoint di bawah ini membutuhkan JWT di header Authorization.

**POST** `/totp/setup`
- Setup TOTP untuk user, mengembalikan secret/key

**POST** `/totp/enable`
- Enable TOTP setelah setup, memerlukan kode verifikasi

---

## Struktur User (contoh)
```json
{
  "id": "uuid",
  "email": "user@email.com",
  "status": "active|pending|banned",
  "is_email_verified": true,
  "created_at": 1234567890,
  "updated_at": 1234567890,
  "custom_attributes": { "key": "value" }
}
```

---

## Catatan
- Semua response error akan dalam format `{ "error": "pesan error" }`
- Untuk testing verifikasi email dan recovery, kode akan muncul di log terminal (simulasi)
- JWT token berlaku 24 jam
- Semua endpoint protected harus menyertakan JWT di header Authorization
- Untuk pengujian TOTP, ikuti flow setup dan enable
