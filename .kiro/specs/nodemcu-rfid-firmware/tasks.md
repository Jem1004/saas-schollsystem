# Tasks: NodeMCU RFID Attendance Firmware

## Phase 1: Project Setup & Basic Hardware

### Task 1.1: Create Project Structure
- [x] Buat folder `firmware/nodemcu-rfid/` di root project
- [x] Buat file `nodemcu-rfid-firmware.ino` sebagai main sketch
- [x] Buat file `pins.h` dengan definisi pin untuk NodeMCU V3
- [x] Setup Arduino IDE / PlatformIO dengan board ESP8266

**Files**: `firmware/nodemcu-rfid/nodemcu-rfid-firmware.ino`, `firmware/nodemcu-rfid/pins.h`

### Task 1.2: Implement BuzzerManager
- [x] Buat `buzzer_manager.h` dan `buzzer_manager.cpp`
- [x] Implement `begin()` untuk setup pin
- [x] Implement pola buzzer: `beepSuccess()`, `beepPairingOK()`, `beepError()`, `beepNetworkError()`, `beepStartup()`
- [x] Test dengan Serial command

**Files**: `firmware/nodemcu-rfid/buzzer_manager.h`, `firmware/nodemcu-rfid/buzzer_manager.cpp`
**Acceptance Criteria**: AC4.2

### Task 1.3: Implement DisplayManager
- [x] Buat `display_manager.h` dan `display_manager.cpp`
- [x] Implement `begin()` dengan auto-detect I2C address (0x27 atau 0x3F)
- [x] Implement `showMessage()`, `showIdle()`, `showSuccess()`, `showError()`
- [x] Implement `truncateString()` untuk handle nama panjang
- [x] Test tampilan dengan berbagai pesan

**Files**: `firmware/nodemcu-rfid/display_manager.h`, `firmware/nodemcu-rfid/display_manager.cpp`
**Acceptance Criteria**: AC4.1, AC4.4

### Task 1.4: Implement RFIDReader
- [x] Buat `rfid_reader.h` dan `rfid_reader.cpp`
- [x] Implement `begin()` dengan SPI setup untuk NodeMCU
- [x] Implement `isCardPresent()` dan `readCardUID()`
- [x] Format UID ke hex string uppercase (contoh: "A1B2C3D4")
- [x] Implement `halt()` untuk release card
- [x] Test baca kartu dan tampilkan UID di Serial

**Files**: `firmware/nodemcu-rfid/rfid_reader.h`, `firmware/nodemcu-rfid/rfid_reader.cpp`
**Acceptance Criteria**: AC2.2

---

## Phase 2: Connectivity & Configuration

### Task 2.1: Implement ConfigManager
- [x] Buat `config.h` dan `config.cpp`
- [x] Define struct `Config` dengan wifi_ssid, wifi_password, api_key, server_url
- [x] Implement `begin()` untuk init LittleFS
- [x] Implement `load()` dan `save()` dengan JSON format
- [x] Implement getter/setter untuk setiap config field
- [x] Test save dan load config

**Files**: `firmware/nodemcu-rfid/config.h`, `firmware/nodemcu-rfid/config.cpp`
**Acceptance Criteria**: AC1.1, AC5.6

### Task 2.2: Implement Serial Configuration
- [x] Tambahkan Serial command parser di main loop
- [x] Handle command: `WIFI:ssid:password`, `APIKEY:key`, `SERVER:url`, `STATUS`, `RESET`
- [x] Validasi input dan berikan feedback
- [x] Simpan config setelah setiap perubahan

**Files**: `firmware/nodemcu-rfid/nodemcu-rfid-firmware.ino`
**Acceptance Criteria**: AC1.2

### Task 2.3: Implement WiFiManager
- [x] Buat `wifi_manager.h` dan `wifi_manager.cpp`
- [x] Implement `begin()` dan `connect()` dengan timeout 30 detik
- [x] Implement `checkConnection()` dengan auto-reconnect setiap 10 detik
- [x] Implement `isConnected()`, `getIP()`, `getRSSI()`
- [x] Tampilkan status di LCD saat connecting
- [x] Test connect, disconnect, dan reconnect

**Files**: `firmware/nodemcu-rfid/wifi_manager.h`, `firmware/nodemcu-rfid/wifi_manager.cpp`
**Acceptance Criteria**: AC1.3, AC1.4, AC5.1, AC5.2

### Task 2.4: Implement TimeManager
- [x] Buat `time_manager.h` dan `time_manager.cpp`
- [x] Implement `begin()` dengan NTP client setup (pool.ntp.org)
- [x] Set timezone ke WIB (UTC+7)
- [x] Implement `update()` dengan resync setiap 1 jam
- [x] Implement `getFormattedTime()` dan `getFormattedDate()`
- [x] Handle NTP failure gracefully

**Files**: `firmware/nodemcu-rfid/time_manager.h`, `firmware/nodemcu-rfid/time_manager.cpp`
**Acceptance Criteria**: AC4.3, AC6.1, AC6.2, AC6.3, AC6.4

---

## Phase 3: API Integration

### Task 3.1: Implement APIClient Base
- [x] Buat `api_client.h` dan `api_client.cpp`
- [x] Define response structs: `ValidationResponse`, `AttendanceResponse`, `PairingResponse`
- [x] Implement `begin()` dengan server URL dan API key
- [x] Implement helper `buildURL()` dan `sendRequest()` dengan retry logic
- [x] Set timeout 5 detik dan max retry 3x

**Files**: `firmware/nodemcu-rfid/api_client.h`, `firmware/nodemcu-rfid/api_client.cpp`
**Acceptance Criteria**: AC2.9, AC5.3

### Task 3.2: Implement API Key Validation
- [x] Implement `validateAPIKey()` method
- [x] Call endpoint `POST /api/v1/devices/validate-key`
- [x] Parse response dan return `ValidationResponse`
- [x] Tampilkan error di LCD jika invalid

**Files**: `firmware/nodemcu-rfid/api_client.cpp`
**Acceptance Criteria**: AC1.5

### Task 3.3: Implement Attendance Recording
- [x] Implement `recordAttendance(rfidCode)` method
- [x] Call endpoint `POST /api/v1/attendance/rfid`
- [x] Parse response: student_name, status, message
- [x] Map status ke bahasa Indonesia (on_time → "Tepat Waktu", late → "Terlambat")

**Files**: `firmware/nodemcu-rfid/api_client.cpp`
**Acceptance Criteria**: AC2.3, AC2.4

### Task 3.4: Implement Pairing API
- [x] Implement `checkPairingStatus()` method
- [x] Implement `processPairing(rfidCode)` method
- [x] Call endpoint `POST /api/v1/pairing/rfid`
- [x] Parse response untuk success/failure dan student name

**Files**: `firmware/nodemcu-rfid/api_client.cpp`
**Acceptance Criteria**: AC3.2, AC3.4

---

## Phase 4: State Machine & Main Logic

### Task 4.1: Implement StateMachine
- [x] Buat `state_machine.h` dan `state_machine.cpp`
- [x] Define enum `DeviceState` dengan semua state
- [x] Implement `setState()`, `getState()`, `isInPairingMode()`
- [x] Implement `enterPairingMode()` dan `exitPairingMode()`
- [x] Track state entry time untuk timeout handling

**Files**: `firmware/nodemcu-rfid/state_machine.h`, `firmware/nodemcu-rfid/state_machine.cpp`

### Task 4.2: Implement Startup Sequence
- [x] Init semua manager di `setup()`
- [x] Load config dari LittleFS
- [x] Connect WiFi dengan progress di LCD
- [x] Sync NTP time
- [x] Validate API key
- [x] Beep startup jika sukses
- [x] Handle error di setiap step

**Files**: `firmware/nodemcu-rfid/nodemcu-rfid-firmware.ino`
**Acceptance Criteria**: AC1.3, AC1.4, AC1.5

### Task 4.3: Implement Main Loop - Idle State
- [x] Check WiFi connection dan reconnect jika perlu
- [x] Update NTP time
- [x] Tampilkan jam di LCD saat idle
- [x] Scan RFID card
- [x] Feed watchdog dengan `yield()`

**Files**: `firmware/nodemcu-rfid/nodemcu-rfid-firmware.ino`
**Acceptance Criteria**: AC2.1, AC4.3, AC5.4

### Task 4.4: Implement Card Tap Handler - Attendance Mode
- [x] Detect card tap dengan `isCardPresent()`
- [x] Read UID dan format ke hex
- [x] Implement cooldown 2 detik untuk prevent double-tap
- [x] Call `recordAttendance()` API
- [x] Handle success: tampilkan nama + status, beep success
- [x] Handle kartu tidak dikenal: tampilkan error, beep error
- [x] Handle network error: tampilkan error, beep network error
- [x] Return ke idle setelah 3 detik

**Files**: `firmware/nodemcu-rfid/nodemcu-rfid-firmware.ino`
**Acceptance Criteria**: AC2.2, AC2.3, AC2.4, AC2.5, AC2.6, AC2.7, AC2.8

### Task 4.5: Implement Pairing Mode
- [x] Periodic check pairing status (setiap 2 detik)
- [x] Jika session aktif, masuk pairing mode
- [x] Tampilkan "Mode Pairing" + nama siswa di LCD
- [x] Handle card tap: call `processPairing()` API
- [x] Handle success: tampilkan "Terdaftar", beep pairing OK
- [x] Handle kartu sudah dipakai: tampilkan error, beep error
- [x] Exit pairing mode setelah selesai/gagal

**Files**: `firmware/nodemcu-rfid/nodemcu-rfid-firmware.ino`
**Acceptance Criteria**: AC3.1, AC3.2, AC3.3, AC3.4, AC3.5, AC3.6, AC3.7

---

## Phase 5: Testing & Documentation

### Task 5.1: Integration Testing
- [x] Test full flow: startup → idle → tap card → attendance recorded
- [x] Test pairing flow: start pairing dari web → tap card → card registered
- [x] Test error scenarios: WiFi disconnect, invalid API key, unknown card
- [x] Test edge cases: rapid tapping, long names, network timeout
- [x] Verify memory usage < 80% heap

### Task 5.2: Create README Documentation
- [x] Hardware wiring diagram
- [x] Library dependencies dan versi
- [x] Arduino IDE / PlatformIO setup instructions
- [x] Configuration guide via Serial
- [x] Troubleshooting common issues

**Files**: `firmware/nodemcu-rfid/README.md`

### Task 5.3: Final Optimization
- [x] Review memory usage dan optimize jika perlu
- [x] Ensure semua String literals menggunakan F() macro
- [x] Verify watchdog tidak trigger dalam normal operation
- [x] Clean up debug logs untuk production

---

## Dependencies Checklist

Libraries yang perlu diinstall di Arduino IDE:
- [x] ESP8266 Board Package (via Board Manager)
- [x] MFRC522 by GithubCommunity
- [x] LiquidCrystal I2C by Frank de Brabander
- [x] ArduinoJson by Benoit Blanchon (v6.x)
- [x] NTPClient by Fabrice Weinberg
- [x] ESP8266HTTPClient (included with ESP8266 core)

---

## Estimated Timeline

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| Phase 1 | Project Setup & Basic Hardware | 2-3 jam |
| Phase 2 | Connectivity & Configuration | 3-4 jam |
| Phase 3 | API Integration | 3-4 jam |
| Phase 4 | State Machine & Main Logic | 4-5 jam |
| Phase 5 | Testing & Documentation | 2-3 jam |
| **Total** | | **14-19 jam** |
