📘 PRODUCT REQUIREMENTS DOCUMENT (DECOMPOSED – FINAL)
Sistem SaaS Manajemen Sekolah, Absensi RFID, Nilai, & Konseling (BK)
1. REQUIREMENTS / USER STORIES

Deskripsi fitur dari perspektif pengguna, apa yang ingin dicapai dan mengapa

1.1 Super Admin (Web)

Sebagai Super Admin,
saya ingin membuat dan mengelola sekolah (tenant),
sehingga satu platform dapat digunakan oleh banyak sekolah dengan data yang terisolasi dan aman.

Sebagai Super Admin,
saya ingin mengelola perangkat RFID (ESP32) secara terpusat,
sehingga setiap device memiliki API key yang valid dan dapat dikontrol keamanannya.

Sebagai Super Admin,
saya ingin memantau status sistem secara global,
sehingga saya dapat memastikan platform stabil untuk seluruh sekolah.

1.2 Admin Sekolah (Web)

Sebagai Admin Sekolah,
saya ingin memetakan data sekolah seperti kelas, siswa, dan orang tua,
sehingga seluruh sistem dapat berjalan dengan struktur data yang rapi dan konsisten.

Sebagai Admin Sekolah,
saya ingin mengelola akun guru, wali kelas, dan guru BK,
sehingga setiap peran memiliki akses yang sesuai tanggung jawabnya.

Sebagai Admin Sekolah,
saya ingin memantau absensi dan aktivitas siswa secara ringkas,
sehingga saya dapat mengetahui kondisi sekolah tanpa masuk ke detail sensitif konseling.

1.3 Guru BK (Web – Halaman Khusus)

Sebagai Guru BK,
saya ingin mencatat pelanggaran siswa,
sehingga perilaku siswa terdokumentasi secara resmi dan dapat ditindaklanjuti.

Sebagai Guru BK,
saya ingin mencatat prestasi siswa beserta poinnya,
sehingga siswa tidak hanya dinilai dari kesalahan tetapi juga dari pencapaian.

Sebagai Guru BK,
saya ingin mencatat izin siswa keluar sekolah saat jam pelajaran dan menghasilkan dokumen izin,
sehingga izin tersebut memiliki bukti administratif yang dapat dilihat oleh orang tua.

Sebagai Guru BK,
saya ingin membuat catatan konseling yang bersifat internal dan ringkasan untuk orang tua,
sehingga proses konseling tetap profesional tanpa membuka informasi sensitif.

1.4 Wali Kelas (Web – Fitur Utama)

Sebagai Wali Kelas,
saya ingin memantau absensi siswa di kelas saya,
sehingga saya mengetahui kedisiplinan harian siswa.

Sebagai Wali Kelas,
saya ingin menginput nilai siswa secara berkala,
sehingga orang tua dapat memantau perkembangan akademik anak sebelum rapor resmi.

Sebagai Wali Kelas,
saya ingin menulis catatan wali kelas yang bisa dibaca oleh orang tua,
sehingga komunikasi perkembangan siswa menjadi lebih terstruktur dibandingkan buku penghubung manual.

1.5 Orang Tua (Mobile – Produk Utama)

Sebagai Orang Tua,
saya ingin login menggunakan NISN anak saya,
sehingga saya tidak perlu proses pendaftaran akun yang rumit.

Sebagai Orang Tua,
saya ingin menerima notifikasi ketika anak saya hadir atau pulang sekolah,
sehingga saya merasa tenang mengetahui keberadaan anak saya.

Sebagai Orang Tua,
saya ingin mengetahui jika anak saya melakukan pelanggaran, mendapat prestasi, atau izin keluar sekolah,
sehingga saya dapat mengambil tindakan atau melakukan komunikasi lanjutan dengan sekolah.

Sebagai Orang Tua,
saya ingin membaca nilai dan catatan wali kelas,
sehingga saya mengetahui perkembangan akademik dan sikap anak saya di sekolah.

1.6 Siswa (Mobile – MVP+)

Sebagai Siswa,
saya ingin melihat riwayat absensi, nilai, prestasi, dan pelanggaran saya,
sehingga saya dapat memahami tanggung jawab dan perkembangan diri saya.

2. TECHNICAL CONTEXT

Stack teknologi, arsitektur, dan constraints

2.1 Stack Teknologi
Backend

Golang

Fiber Framework

Modular Monolith Architecture

PostgreSQL (GORM ORM)

Redis (Queue & Cache)

JWT + Refresh Token

Firebase Cloud Messaging (FCM)

Web Admin

Vue 3

Vite

TypeScript

Ant Design Vue

SPA (tanpa SSR)

Tema clean (putih + orange)

Mobile App

Flutter

Android & iOS

Online-first

Push Notification via FCM

Hardware

ESP32 + RFID

REST API via HTTPS

API Key per device

Development Environment

macOS

Docker (PostgreSQL & Redis)

2.2 Architecture Constraints

Sistem multi-tenant wajib

Isolasi data antar sekolah

Tidak ada business logic di ESP32

Semua notifikasi asynchronous

Mobile app hanya untuk monitoring, bukan input operasional

Tidak menggunakan microservices pada fase ini

Offline-first penuh belum diimplementasikan

3. API SPECIFICATIONS

Integrasi API & referensi OpenAPI

Backend menggunakan REST API + JWT Authentication.

Spesifikasi lengkap didefinisikan dalam OpenAPI:

#[[file:openapi.yaml]]

Domain Endpoint Utama

Auth

Attendance

Grades

Homeroom Notes

Counseling (BK)

Devices

Notifications

4. DATA MODELS

Struktur database utama

4.1 Core & Auth

users

id

school_id

role (super_admin, admin_school, guru_bk, wali_kelas, guru, parent, student)

username

password_hash

is_active

students

id

school_id

class_id

nis

nisn

name

4.2 Absensi

attendance

id

student_id

date

check_in_time

check_out_time

method (rfid/manual)

4.3 Akademik

grades

id

student_id

title

score

description

created_by (wali_kelas)

created_at

homeroom_notes

id

student_id

teacher_id

content

created_at

4.4 BK / Konseling

violations

id

student_id

category

level

description

created_by

achievements

id

student_id

title

point

description

counseling_notes

id

student_id

internal_note

parent_summary

permits

id

student_id

reason

exit_time

return_time

document_url

4.5 Device & Notification

devices

id

school_id

device_code

api_key

notifications

id

user_id

type

title

message

is_read

5. UI / UX MOCKUPS

Wireframe & design reference berbasis role

5.1 Web Admin (Desktop)
Global Layout

Sidebar kiri (role-based menu)

Header atas (profil & logout)

Content area putih

Super Admin

Dashboard tenant

List sekolah

Device RFID management

Admin Sekolah

Dashboard sekolah

Master data (kelas, siswa, orang tua)

User management

Guru BK

Dashboard BK

Tab per siswa:

Pelanggaran

Prestasi

Izin keluar

Catatan konseling

Modal & drawer untuk input data

Wali Kelas

Dashboard kelas

Halaman Nilai Siswa

Halaman Catatan Wali Kelas

Monitoring BK (read-only)

5.2 Mobile App (Flutter)
Orang Tua

Login (NISN)

Dashboard anak

Absensi

Nilai

Catatan wali kelas

Informasi BK

Izin keluar (detail + dokumen)

Notification center

Siswa

Profil

Riwayat absensi

Nilai

Prestasi & pelanggaran

5.3 Dokumen Izin Keluar

Format PDF / struk

Nama siswa

NIS / NISN

Kelas

Alasan

Waktu keluar

Guru penanggung jawab

Timestamp

✅ END OF DECOMPOSED PRD (FINAL)