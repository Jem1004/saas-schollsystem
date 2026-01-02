# Requirements: NodeMCU RFID Attendance Firmware

## Overview
Firmware untuk perangkat absensi berbasis NodeMCU V3 (ESP8266) dengan RFID RC522, LCD I2C 16x2, dan Buzzer. Perangkat ini berfungsi sebagai terminal absensi siswa yang terintegrasi dengan backend sistem manajemen sekolah.

## Hardware Configuration
```
NodeMCU V3 (ESP8266)
├── RFID RC522 (SPI)
│   ├── SDA  → D4 (GPIO2)
│   ├── SCK  → D5 (GPIO14)
│   ├── MOSI → D7 (GPIO13)
│   ├── MISO → D6 (GPIO12)
│   ├── RST  → D3 (GPIO0)
│   ├── 3.3V → 3V3
│   └── GND  → GND
├── LCD I2C 16x2 (I2C)
│   ├── SDA → D2 (GPIO4)
│   ├── SCL → D1 (GPIO5)
│   ├── VCC → VIN (5V)
│   └── GND → GND
└── Buzzer
    ├── (+) → D8 (GPIO15)
    └── (-) → GND
```

## Backend API Endpoints
Firmware akan berkomunikasi dengan backend melalui endpoint berikut:

### 1. Validate API Key
- **Endpoint**: `POST /api/v1/devices/validate-key`
- **Request**: `{ "api_key": "string" }`
- **Response**: `{ "success": bool, "data": { "valid": bool, "device_id": uint, "school_id": uint } }`

### 2. Record RFID Attendance
- **Endpoint**: `POST /api/v1/attendance/rfid`
- **Request**: `{ "api_key": "string", "rfid_code": "string" }`
- **Response**: `{ "success": bool, "data": { "student_name": "string", "status": "on_time|late|very_late", "message": "string" } }`

### 3. Process RFID Pairing
- **Endpoint**: `POST /api/v1/pairing/rfid`
- **Request**: `{ "api_key": "string", "rfid_code": "string" }`
- **Response**: `{ "success": bool, "data": { "student_name": "string", "message": "string" } }`

---

## User Stories

### Story 1: Konfigurasi WiFi dan API Key
**Sebagai** teknisi/admin sekolah,
**Saya ingin** mengkonfigurasi WiFi credentials dan API key pada perangkat,
**Sehingga** perangkat dapat terhubung ke jaringan dan backend sistem.

#### Acceptance Criteria
- [ ] AC1.1: Perangkat menyimpan WiFi SSID, password, dan API key di EEPROM/flash memory
- [ ] AC1.2: Konfigurasi dapat diubah melalui Serial Monitor (115200 baud) dengan format command
- [ ] AC1.3: Perangkat menampilkan status koneksi WiFi di LCD saat startup
- [ ] AC1.4: Jika WiFi gagal connect dalam 30 detik, tampilkan error dan retry otomatis
- [ ] AC1.5: API key divalidasi ke backend saat startup, tampilkan error jika invalid

### Story 2: Mode Absensi (Normal Operation)
**Sebagai** siswa,
**Saya ingin** tap kartu RFID pada perangkat untuk mencatat kehadiran,
**Sehingga** absensi saya tercatat otomatis di sistem.

#### Acceptance Criteria
- [ ] AC2.1: LCD menampilkan "Tap Kartu..." saat idle dalam mode absensi
- [ ] AC2.2: Saat kartu terdeteksi, baca UID kartu dalam format hex (contoh: "A1B2C3D4")
- [ ] AC2.3: Kirim data ke endpoint `/api/v1/attendance/rfid` dengan api_key dan rfid_code
- [ ] AC2.4: Jika sukses, tampilkan nama siswa dan status (Tepat Waktu/Terlambat) di LCD selama 3 detik
- [ ] AC2.5: Jika sukses, buzzer berbunyi 1x pendek (100ms)
- [ ] AC2.6: Jika gagal (kartu tidak terdaftar), tampilkan "Kartu Tidak Dikenal" dan buzzer 3x pendek
- [ ] AC2.7: Jika error jaringan, tampilkan "Error Jaringan" dan buzzer 2x panjang
- [ ] AC2.8: Cooldown 2 detik setelah tap untuk mencegah double-tap
- [ ] AC2.9: Request timeout maksimal 5 detik, jika timeout tampilkan error

### Story 3: Mode Pairing (Pendaftaran Kartu)
**Sebagai** admin sekolah,
**Saya ingin** mendaftarkan kartu RFID baru ke siswa melalui perangkat,
**Sehingga** siswa dapat menggunakan kartu tersebut untuk absensi.

#### Acceptance Criteria
- [ ] AC3.1: Mode pairing diaktifkan dari web admin (backend mengelola session)
- [ ] AC3.2: Perangkat secara periodik (setiap 2 detik) cek status pairing ke backend
- [ ] AC3.3: Jika ada session pairing aktif, LCD menampilkan "Mode Pairing" dan nama siswa yang akan dipasangkan
- [ ] AC3.4: Saat kartu di-tap dalam mode pairing, kirim ke endpoint `/api/v1/pairing/rfid`
- [ ] AC3.5: Jika sukses, tampilkan "Kartu Terdaftar" + nama siswa, buzzer 2x pendek
- [ ] AC3.6: Jika kartu sudah dipakai siswa lain, tampilkan "Kartu Sudah Dipakai", buzzer 3x pendek
- [ ] AC3.7: Setelah pairing selesai/gagal, kembali ke mode absensi otomatis

### Story 4: Feedback Visual dan Audio
**Sebagai** pengguna perangkat,
**Saya ingin** mendapat feedback yang jelas dari LCD dan buzzer,
**Sehingga** saya tahu status operasi perangkat.

#### Acceptance Criteria
- [ ] AC4.1: LCD baris 1 menampilkan status/mode, baris 2 menampilkan informasi detail
- [ ] AC4.2: Pola buzzer berbeda untuk setiap kondisi:
  - Sukses absensi: 1 beep pendek (100ms)
  - Sukses pairing: 2 beep pendek (100ms, jeda 100ms)
  - Kartu tidak dikenal: 3 beep pendek (100ms, jeda 100ms)
  - Error jaringan: 2 beep panjang (500ms, jeda 200ms)
  - Startup sukses: 1 beep panjang (300ms)
- [ ] AC4.3: LCD menampilkan jam real-time (dari NTP) saat idle
- [ ] AC4.4: Karakter khusus Indonesia (nama dengan spasi) ditampilkan dengan benar

### Story 5: Reliability dan Error Handling
**Sebagai** operator sistem,
**Saya ingin** perangkat dapat menangani error dengan baik,
**Sehingga** sistem tetap berjalan meskipun ada gangguan.

#### Acceptance Criteria
- [ ] AC5.1: Jika WiFi terputus, perangkat mencoba reconnect otomatis setiap 10 detik
- [ ] AC5.2: Tampilkan indikator WiFi status di LCD (connected/disconnected)
- [ ] AC5.3: Jika backend tidak merespons, retry hingga 3x sebelum menampilkan error
- [ ] AC5.4: Watchdog timer aktif untuk auto-restart jika sistem hang (timeout 8 detik)
- [ ] AC5.5: Log error ke Serial untuk debugging
- [ ] AC5.6: Konfigurasi tersimpan persisten, tidak hilang saat restart/power off

### Story 6: Sinkronisasi Waktu
**Sebagai** sistem absensi,
**Saya ingin** waktu perangkat tersinkronisasi dengan server NTP,
**Sehingga** timestamp absensi akurat.

#### Acceptance Criteria
- [ ] AC6.1: Sinkronisasi waktu dengan NTP server (pool.ntp.org) saat startup
- [ ] AC6.2: Timezone diset ke WIB (UTC+7)
- [ ] AC6.3: Resync NTP setiap 1 jam untuk menjaga akurasi
- [ ] AC6.4: Jika NTP gagal, gunakan waktu terakhir yang valid dan tampilkan warning

---

## Non-Functional Requirements

### Performance
- Response time tap kartu ke tampilan hasil: < 3 detik (termasuk network latency)
- Boot time dari power on ke ready: < 10 detik
- Memory usage: < 80% dari available heap

### Reliability
- Uptime target: 99% (restart otomatis jika crash)
- Data tidak hilang saat power failure (konfigurasi di flash)

### Security
- API key tidak ditampilkan di LCD atau Serial log
- Komunikasi menggunakan HTTPS (jika backend support)
- API key disimpan terenkripsi di flash (opsional, nice-to-have)

### Maintainability
- Kode modular dengan pemisahan concern (WiFi, RFID, LCD, API, Config)
- Logging level configurable (DEBUG, INFO, ERROR)
- OTA update support (nice-to-have untuk future)

---

## Technical Constraints

### Hardware Limitations
- ESP8266 single-core, perlu non-blocking code
- RAM terbatas (~50KB heap), hindari String concatenation berlebihan
- Flash storage 4MB, gunakan LittleFS untuk config

### Library Dependencies
- ESP8266WiFi - WiFi connectivity
- MFRC522 - RFID RC522 reader
- LiquidCrystal_I2C - LCD I2C display
- ArduinoJson - JSON parsing
- ESP8266HTTPClient - HTTP requests
- NTPClient - Time synchronization
- LittleFS - File system untuk config storage

### Network Requirements
- WiFi 2.4GHz (ESP8266 tidak support 5GHz)
- Backend API accessible dari jaringan lokal sekolah
- Firewall allow outbound HTTP/HTTPS ke backend server

---

## Out of Scope
- Offline mode dengan local storage (untuk versi mendatang)
- Multiple card format (hanya support MIFARE Classic)
- Biometric integration
- Battery backup / UPS monitoring
- Camera integration untuk foto absensi
