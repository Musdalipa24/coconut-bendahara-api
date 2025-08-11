# API Bendahara Desa

API ini digunakan untuk pengelolaan keuangan di sebuah desa. API ini dibangun menggunakan Golang dan menggunakan database MYSQL untuk menyimpan data transaksi pemasukan, pengeluaran, serta laporan keuangan.

## Cara Menjalankan API

### 1. Persiapan
- Pastikan Anda telah menginstal **Golang**
- Pastikan Anda memiliki **database MYSQL** yang telah dikonfigurasi
- Clone repository ini ke lokal:
  ```sh
  git clone -b backend https://github.com/ameliaendino/bendahara-inodes.git
  cd bendahara-inodes
  ```
- Instal dependensi dengan menggunakan:
  ```sh
  go mod tidy
  ```

### 2. Konfigurasi Database
Buat file `.env` dan isi dengan konfigurasi database yang sesuai. Contoh:
```
DB_NAME=database_name
APP_PORT=8087
JWT_SECRET=your_secret_key
```

### 3. Menjalankan API
Jalankan perintah berikut untuk menjalankan server:
```sh
go run main.go
```
Server akan berjalan di port yang telah dikonfigurasi dalam file `.env` (misalnya, `http://localhost:8080`).

## Endpoint API dan Fitur

### 1. **Admin**
- **Daftar Admin**
  - Endpoint: `POST /api/admin/daftar`
  - Deskripsi: Mendaftarkan admin baru dengan autentikasi JWT
- **Login Admin**
  - Endpoint: `POST /api/admin/login`
  - Deskripsi: Admin dapat login untuk mendapatkan token JWT
- **Cari Admin Berdasarkan NIK**
  - Endpoint: `GET /api/admin/:nik`
  - Deskripsi: Mendapatkan data admin berdasarkan NIK

### 2. **Pemasukan**
- **Tambah Pemasukan**
  - Endpoint: `POST /api/pemasukan/add`
  - Deskripsi: Menambahkan data pemasukan
- **Perbarui Pemasukan**
  - Endpoint: `PUT /api/pemasukan/update/:id`
  - Deskripsi: Memperbarui data pemasukan berdasarkan ID
- **Lihat Semua Pemasukan**
  - Endpoint: `GET /api/pemasukan/getall`
  - Deskripsi: Mengambil semua data pemasukan
- **Lihat Detail Pemasukan**
  - Endpoint: `GET /api/pemasukan/get/:id`
  - Deskripsi: Mengambil data pemasukan berdasarkan ID
- **Hapus Pemasukan**
  - Endpoint: `DELETE /api/pemasukan/delete/:id`
  - Deskripsi: Menghapus data pemasukan berdasarkan ID

### 3. **Pengeluaran**
- **Tambah Pengeluaran**
  - Endpoint: `POST /api/pengeluaran/add`
  - Deskripsi: Menambahkan data pengeluaran
- **Perbarui Pengeluaran**
  - Endpoint: `PUT /api/pengeluaran/update/:id`
  - Deskripsi: Memperbarui data pengeluaran berdasarkan ID
- **Lihat Semua Pengeluaran**
  - Endpoint: `GET /api/pengeluaran/getall`
  - Deskripsi: Mengambil semua data pengeluaran
- **Lihat Detail Pengeluaran**
  - Endpoint: `GET /api/pengeluaran/get/:id`
  - Deskripsi: Mengambil data pengeluaran berdasarkan ID
- **Hapus Pengeluaran**
  - Endpoint: `DELETE /api/pengeluaran/delete/:id`
  - Deskripsi: Menghapus data pengeluaran berdasarkan ID

### 4. **Transaksi**
- **Lihat Semua Transaksi**
  - Endpoint: `GET /api/transaksi/getall`
  - Deskripsi: Mengambil semua data transaksi
- **Lihat Transaksi Terakhir**
  - Endpoint: `GET /api/transaksi/getlast`
  - Deskripsi: Mengambil transaksi terakhir

### 5. **Laporan Keuangan**
- **Lihat Semua Laporan Keuangan**
  - Endpoint: `GET /api/laporan/getall`
  - Deskripsi: Mengambil semua laporan keuangan
- **Lihat Saldo Terakhir**
  - Endpoint: `GET /api/laporan/saldo`
  - Deskripsi: Mengambil saldo terakhir
- **Total Pengeluaran**
  - Endpoint: `GET /api/laporan/pengeluaran`
  - Deskripsi: Menghitung total pengeluaran
- **Total Pemasukan**
  - Endpoint: `GET /api/laporan/pemasukan`
  - Deskripsi: Menghitung total pemasukan

### 6. **Upload Berkas**
- **Akses Berkas yang Diupload**
  - Endpoint: `GET /api/uploads/*filepath`
  - Deskripsi: Melayani berkas yang diunggah di direktori `./uploads/`

## Middleware
API ini menggunakan middleware JWT untuk beberapa endpoint yang memerlukan autentikasi. Pastikan setiap permintaan yang memerlukan autentikasi menyertakan token dalam header:
```sh
Authorization: Bearer <your_token>
```

## Kesimpulan
API Bendahara Desa ini bertujuan untuk mempermudah pengelolaan keuangan desa dengan fitur pemasukan, pengeluaran, transaksi, dan laporan keuangan. Pastikan mengikuti langkah-langkah di atas untuk menjalankan dan mengakses API dengan benar.