# NodeMCU RFID Attendance Firmware

Firmware untuk perangkat absensi berbasis NodeMCU V3 (ESP8266) dengan RFID RC522, LCD I2C 16x2, dan Buzzer.

## Features

- ✅ Absensi dengan tap kartu RFID
- ✅ Mode pairing untuk mendaftarkan kartu baru
- ✅ LCD display dengan waktu real-time (NTP)
- ✅ Buzzer feedback untuk berbagai event
- ✅ Auto-reconnect WiFi
- ✅ Konfigurasi via Serial Monitor
- ✅ Persistent config di flash memory

## Hardware Requirements

### Components
- NodeMCU V3 (ESP8266)
- RFID RC522 Module
- LCD I2C 16x2 (address 0x27 atau 0x3F)
- Active Buzzer 5V
- Jumper wires

### Wiring Diagram

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
│
├── LCD I2C 16x2 (I2C)
│   ├── SDA → D2 (GPIO4)
│   ├── SCL → D1 (GPIO5)
│   ├── VCC → VIN (5V)
│   └── GND → GND
│
└── Buzzer
    ├── (+) → D8 (GPIO15)
    └── (-) → GND
```

### Pin Summary

| Component | Pin | NodeMCU | GPIO |
|-----------|-----|---------|------|
| RFID SDA  | SS  | D4      | 2    |
| RFID SCK  | SCK | D5      | 14   |
| RFID MOSI | MOSI| D7      | 13   |
| RFID MISO | MISO| D6      | 12   |
| RFID RST  | RST | D3      | 0    |
| LCD SDA   | SDA | D2      | 4    |
| LCD SCL   | SCL | D1      | 5    |
| Buzzer    | +   | D8      | 15   |

## Software Requirements

### Arduino IDE Setup

1. **Install Arduino IDE** (v1.8.x atau v2.x)

2. **Add ESP8266 Board Package**
   - File → Preferences → Additional Board Manager URLs
   - Add: `http://arduino.esp8266.com/stable/package_esp8266com_index.json`
   - Tools → Board → Board Manager → Search "ESP8266" → Install

3. **Select Board**
   - Tools → Board → ESP8266 Boards → NodeMCU 1.0 (ESP-12E Module)

4. **Board Settings**
   - Upload Speed: 115200
   - CPU Frequency: 80 MHz
   - Flash Size: 4MB (FS:2MB OTA:~1019KB)

### Required Libraries

Install via Arduino Library Manager (Sketch → Include Library → Manage Libraries):

| Library | Author | Version |
|---------|--------|---------|
| MFRC522 | GithubCommunity | 1.4.x |
| LiquidCrystal I2C | Frank de Brabander | 1.1.x |
| ArduinoJson | Benoit Blanchon | 6.x |
| NTPClient | Fabrice Weinberg | 3.2.x |

> **Note:** ESP8266WiFi, ESP8266HTTPClient, dan LittleFS sudah termasuk dalam ESP8266 core.

### PlatformIO Setup (Alternative)

```ini
; platformio.ini
[env:nodemcuv2]
platform = espressif8266
board = nodemcuv2
framework = arduino
monitor_speed = 115200
lib_deps =
    miguelbalboa/MFRC522@^1.4.10
    marcoschwartz/LiquidCrystal_I2C@^1.1.4
    bblanchon/ArduinoJson@^6.21.3
    arduino-libraries/NTPClient@^3.2.1
```


## Installation

1. **Clone/Download** repository ini

2. **Open Project**
   - Arduino IDE: File → Open → `nodemcu-rfid-firmware.ino`
   - PlatformIO: Open folder di VS Code

3. **Install Libraries** (lihat bagian Required Libraries)

4. **Connect NodeMCU** via USB

5. **Upload Firmware**
   - Arduino IDE: Sketch → Upload (Ctrl+U)
   - PlatformIO: Click Upload button

6. **Open Serial Monitor**
   - Baud rate: 115200
   - Line ending: Newline

## Configuration

### Via Serial Monitor

Buka Serial Monitor (115200 baud) dan gunakan command berikut:

```
WIFI:ssid:password    - Set WiFi credentials
APIKEY:your_api_key   - Set API key dari backend
SERVER:http://ip:port - Set server URL
VALIDATE              - Validate API key
STATUS                - Show device status
TIME                  - Show current time
RESET                 - Reset configuration
BEEP                  - Test buzzer patterns
HELP                  - Show available commands
```

### Configuration Example

```
WIFI:SekolahWiFi:password123
SERVER:http://192.168.1.100:8080
APIKEY:a1b2c3d4e5f6g7h8i9j0...
VALIDATE
```

### Configuration Storage

Konfigurasi disimpan di LittleFS (`/config.json`) dan akan tetap tersimpan meskipun device di-restart atau power off.

## Usage

### Normal Operation (Attendance Mode)

1. Device menampilkan "Tap Kartu..." dan waktu saat ini
2. Siswa tap kartu RFID
3. Device mengirim data ke backend
4. LCD menampilkan nama siswa dan status (Tepat Waktu/Terlambat)
5. Buzzer berbunyi 1x pendek
6. Kembali ke idle setelah 3 detik

### Pairing Mode (Card Registration)

1. Admin memulai pairing session dari web admin
2. Device otomatis mendeteksi session (polling setiap 2 detik)
3. LCD menampilkan "Mode Pairing" dan nama siswa
4. Tap kartu baru untuk mendaftarkan
5. LCD menampilkan "Kartu Terdaftar"
6. Buzzer berbunyi 2x pendek
7. Kembali ke mode absensi

### Buzzer Patterns

| Event | Pattern |
|-------|---------|
| Startup sukses | 1x panjang (300ms) |
| Absensi sukses | 1x pendek (100ms) |
| Pairing sukses | 2x pendek (100ms, 100ms gap) |
| Kartu tidak dikenal | 3x pendek (100ms, 100ms gap) |
| Error jaringan | 2x panjang (500ms, 200ms gap) |

## API Endpoints

Device berkomunikasi dengan backend melalui endpoint berikut:

### Validate API Key
```
POST /api/v1/devices/validate-key
Body: { "api_key": "string" }
Response: { "success": bool, "data": { "valid": bool, "device_id": uint, "school_id": uint } }
```

### Record Attendance
```
POST /api/v1/attendance/rfid
Body: { "api_key": "string", "rfid_code": "string" }
Response: { "success": bool, "data": { "student_name": "string", "status": "on_time|late|very_late" } }
```

### Check Pairing Status
```
GET /api/v1/pairing/status?api_key=string
Response: { "success": bool, "data": { "session_active": bool, "student_name": "string" } }
```

### Process Pairing
```
POST /api/v1/pairing/rfid
Body: { "api_key": "string", "rfid_code": "string" }
Response: { "success": bool, "data": { "student_name": "string" } }
```


## Troubleshooting

### LCD tidak menyala / tidak tampil

1. **Cek wiring** - Pastikan SDA ke D2, SCL ke D1, VCC ke VIN (5V)
2. **Cek I2C address** - Firmware auto-detect 0x27 dan 0x3F
3. **Adjust contrast** - Putar potensiometer di belakang LCD
4. **Cek power** - LCD butuh 5V, gunakan VIN bukan 3V3

### RFID tidak membaca kartu

1. **Cek wiring** - Terutama SDA (D4) dan RST (D3)
2. **Cek power** - RFID butuh 3.3V, JANGAN gunakan 5V
3. **Cek jarak** - Kartu harus dekat dengan reader (< 3cm)
4. **Cek Serial** - Lihat firmware version di startup log
5. **Kartu kompatibel** - Hanya support MIFARE Classic

### WiFi tidak connect

1. **Cek SSID/password** - Gunakan command `STATUS` untuk lihat config
2. **Cek signal** - ESP8266 hanya support 2.4GHz
3. **Cek router** - Pastikan tidak ada MAC filtering
4. **Reset config** - Gunakan command `RESET` lalu konfigurasi ulang

### API error / tidak bisa absensi

1. **Cek WiFi** - Pastikan connected (lihat `STATUS`)
2. **Cek server URL** - Format: `http://ip:port` (tanpa trailing slash)
3. **Cek API key** - Gunakan `VALIDATE` untuk test
4. **Cek network** - Pastikan device bisa reach server
5. **Cek backend** - Pastikan backend running dan accessible

### Device restart terus-menerus

1. **Cek power supply** - Gunakan USB yang stabil (min 500mA)
2. **Cek wiring** - Short circuit bisa menyebabkan restart
3. **Cek heap** - Gunakan `STATUS` untuk lihat free heap
4. **Reset config** - Config corrupt bisa menyebabkan crash

### Waktu tidak sinkron

1. **Cek WiFi** - NTP butuh koneksi internet
2. **Cek firewall** - Port 123 (UDP) harus open untuk NTP
3. **Tunggu sync** - NTP sync bisa butuh beberapa detik setelah WiFi connect

## File Structure

```
nodemcu-rfid/
├── nodemcu-rfid-firmware.ino  # Main sketch
├── pins.h                      # Pin definitions
├── config.h / config.cpp       # Configuration manager
├── wifi_manager.h / .cpp       # WiFi connection handler
├── rfid_reader.h / .cpp        # RFID RC522 wrapper
├── display_manager.h / .cpp    # LCD I2C handler
├── buzzer_manager.h / .cpp     # Buzzer patterns
├── api_client.h / .cpp         # HTTP API client
├── time_manager.h / .cpp       # NTP time sync
├── state_machine.h / .cpp      # Application state
└── README.md                   # This file
```

## Technical Specifications

| Spec | Value |
|------|-------|
| MCU | ESP8266 (80MHz) |
| Flash | 4MB |
| RAM | ~50KB heap |
| WiFi | 2.4GHz 802.11 b/g/n |
| RFID | 13.56MHz MIFARE |
| Display | 16x2 characters |
| Power | 5V via USB |
| Current | ~200mA typical |

## Memory Usage

Target: < 80% heap usage

Typical usage:
- Free heap at startup: ~35KB
- Free heap during operation: ~30KB

Tips untuk optimasi memory:
- Gunakan `F()` macro untuk string literals
- Hindari String concatenation berlebihan
- Release HTTP resources setelah request

## Version History

- **v1.0.0** - Initial release
  - Basic attendance recording
  - Pairing mode support
  - NTP time sync
  - Serial configuration

## License

This firmware is part of the School Management System project.

## Support

Jika ada masalah atau pertanyaan, silakan buat issue di repository atau hubungi tim development.
