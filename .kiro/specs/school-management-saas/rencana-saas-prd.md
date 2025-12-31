📘 PRODUCT REQUIREMENTS DOCUMENT (PRD)
Sistem SaaS Manajemen Sekolah, Absensi RFID, Nilai, & Konseling (BK)

Versi Dokumen: 1.3 (FINAL – UPDATED)
Status: MVP+ (Core Feature Locked)
Fase: Development
Target: Sekolah Indonesia (SMP / SMA / sederajat)

1. RINGKASAN PRODUK

Produk ini adalah sistem SaaS manajemen sekolah berbasis web admin dan mobile app yang dirancang untuk mendukung pendataan sekolah secara rapi dan terstruktur, mengelola absensi siswa berbasis RFID, mencatat perkembangan akademik (nilai), serta menjalankan sistem konseling (BK) yang profesional dan terpisah.

Sistem ini menempatkan transparansi kepada orang tua sebagai nilai utama melalui notifikasi real-time dan monitoring terstruktur terhadap kehadiran, perilaku, akademik, dan izin siswa selama berada di sekolah.

2. LATAR BELAKANG MASALAH

Banyak sekolah masih mengandalkan pencatatan manual atau sistem terpisah untuk absensi, perilaku siswa, dan komunikasi dengan orang tua. Hal ini menyebabkan:

Data tidak konsisten dan sulit ditelusuri

Keterlambatan informasi kepada orang tua

Ketergantungan pada komunikasi informal (WA, buku penghubung manual)

Minimnya dokumentasi resmi untuk konseling dan izin siswa

Produk ini dibangun untuk menjawab kebutuhan tersebut dengan satu sistem terpadu yang rapi, aman, dan mudah digunakan oleh sekolah serta informatif bagi orang tua.

3. TUJUAN PRODUK

Tujuan utama sistem ini adalah:

Menyediakan sistem manajemen sekolah yang terstruktur

Memberikan transparansi real-time kepada orang tua

Memfasilitasi monitoring akademik, perilaku, dan kehadiran siswa

Menjadi pengganti buku penghubung manual

Menjadi fondasi sistem sekolah digital jangka panjang

4. VISI PRODUK

Dari perspektif sekolah, sistem ini menjadi alat kerja harian yang memudahkan pengelolaan data dan operasional tanpa kompleksitas teknis.
Dari perspektif orang tua, sistem ini menjadi sumber informasi yang tepercaya dan menenangkan tentang kondisi anak selama di sekolah.

5. RUANG LINGKUP PRODUK
Termasuk dalam MVP+ (Core):

Sistem multi-tenant (banyak sekolah)

Web Admin berbasis role

Mobile app untuk orang tua & siswa

Absensi RFID (ESP32)

Sistem BK terpisah & profesional

Input nilai oleh wali kelas

Catatan wali kelas untuk orang tua

Push notification real-time

Tidak termasuk (sementara):

Rapor resmi & UAS

Payment / billing

Offline-first penuh

AI & analytics lanjutan

6. ROLE & AKSES PENGGUNA (FINAL)
6.1 Super Admin (Web Only)

Mengelola platform SaaS secara global.

Akses:

Manajemen tenant sekolah

Manajemen admin sekolah

Manajemen device RFID global

Monitoring sistem

Batasan:

Tidak mengakses data siswa detail

6.2 Admin Sekolah (Web Only)

Pengelola operasional dan data sekolah.

Akses:

Master data (kelas, siswa, orang tua)

Manajemen akun guru, wali kelas, dan guru BK

Monitoring absensi

Monitoring data BK (read-only)

Pengaturan sekolah

6.3 Guru BK (Web Only – Halaman Khusus)

Pengelola konseling dan perilaku siswa.

Akses:

Input dan kelola pelanggaran siswa

Input dan kelola prestasi siswa (point)

Input izin keluar sekolah

Membuat catatan konseling (internal & ringkasan)

Melihat riwayat kasus siswa

6.4 Wali Kelas (Web Only – Fitur Utama)

Penghubung utama antara sekolah dan orang tua.

Akses:

Monitoring absensi kelas

Input nilai siswa (monitoring akademik)

Membuat catatan wali kelas

Monitoring data BK (read-only)

Monitoring izin keluar siswa

6.5 Guru (Opsional)

Akses terbatas.

Akses:

Melihat siswa

Menambahkan catatan umum (jika diizinkan)

6.6 Orang Tua (Mobile Only – Produk Utama)

Pengguna mobile application.

Karakteristik:

Login menggunakan NISN siswa

Read-only

Dapat terhubung ke lebih dari satu siswa

Akses:

Absensi anak

Nilai siswa

Catatan wali kelas

Informasi BK (ringkasan)

Riwayat izin keluar + dokumen

Notifikasi real-time

6.7 Siswa (Mobile Only – MVP+)

Pengguna monitoring mandiri.

Akses:

Profil diri

Riwayat absensi

Riwayat izin

Ringkasan prestasi & pelanggaran

Nilai monitoring

7. AUTENTIKASI & LOGIN

Web admin: email + password

Orang tua & siswa: NISN sebagai username

Password awal ditentukan sekolah

Wajib reset password saat login pertama

JWT + Refresh Token

8. FITUR UTAMA SISTEM (CORE FEATURES)
8.1 Absensi Siswa

RFID berbasis ESP32

Absensi manual (fallback)

Waktu masuk & keluar tercatat

Notifikasi ke orang tua

8.2 Sistem Konseling (BK)
8.2.1 Pelanggaran Siswa

Dicatat oleh Guru BK

Kategori & tingkat pelanggaran

Riwayat per siswa

Notifikasi ke orang tua

8.2.2 Prestasi Siswa

Dicatat oleh Guru BK

Point prestasi

Akumulasi point

Riwayat prestasi

8.2.3 Izin Keluar Sekolah

Dicatat oleh Guru BK / Wali Kelas

Alasan, waktu keluar, penanggung jawab

Dokumen izin (PDF/struk)

Notifikasi ke orang tua

8.2.4 Catatan Konseling

Catatan internal (privat BK)

Ringkasan untuk orang tua

8.3 Nilai Siswa (FITUR UTAMA)

Dicatat oleh Wali Kelas

Bukan rapor final

Digunakan sebagai monitoring & early warning

Orang tua & siswa dapat melihat

Notifikasi opsional

8.4 Catatan Wali Kelas (FITUR UTAMA)

Catatan resmi non-BK

Dibaca oleh orang tua & siswa

Menggantikan buku penghubung

9. PUSH NOTIFICATION

Notifikasi dikirim saat:

Absensi masuk / keluar

Pelanggaran BK

Prestasi

Izin keluar

Nilai baru

Catatan wali kelas

Notifikasi bersifat:

Event-driven

Asynchronous

Tidak menampilkan detail sensitif langsung

10. ARSITEKTUR TEKNOLOGI
Backend

Golang + Fiber

Modular Monolith

PostgreSQL + GORM

Redis (queue & cache)

JWT + Refresh Token

Firebase Cloud Messaging

Web Admin

Vue 3 + Vite

TypeScript

Ant Design Vue

Tema clean putih & orange

Mobile App

Flutter

Android & iOS

FCM Push Notification

Online-first

Hardware

ESP32 + RFID

REST API via HTTPS

11. FLOW SISTEM UTAMA
Flow Absensi

RFID → ESP32 → Backend → DB → Redis → Worker → FCM → Orang Tua

Flow BK / Nilai / Catatan

Web Admin → Backend → DB → Redis → Worker → FCM → Orang Tua

12. KEAMANAN & PRIVASI

RBAC ketat

Isolasi tenant

Hash password

API key untuk ESP32

Audit log untuk BK & nilai

Orang tua hanya melihat data anaknya

13. NON-FUNCTIONAL REQUIREMENTS

Stabil saat jam sibuk pagi

Response cepat

Kode maintainable

Mudah dikembangkan

Logging & error handling

14. KRITERIA KEBERHASILAN PRODUK

Absensi RFID stabil

Orang tua aktif menggunakan aplikasi

Wali kelas rutin input nilai & catatan

Buku penghubung manual tidak lagi dipakai

Sistem siap pilot & scale

15. ROADMAP PASCA MVP+

Rapor resmi

Analytics sekolah

Offline-first mobile

Payment & billing

Integrasi sistem eksternal

END OF PRODUCT REQUIREMENTS DOCUMENT v1.3