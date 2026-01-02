# Design: NodeMCU RFID Attendance Firmware

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Main Loop                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐ │
│  │  WiFi    │  │  RFID    │  │  State   │  │  Pairing Check   │ │
│  │  Check   │→ │  Scan    │→ │  Machine │→ │  (periodic)      │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
         │              │              │              │
         ▼              ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ WiFiManager  │ │ RFIDReader   │ │ DisplayMgr   │ │ APIClient    │
│              │ │              │ │              │ │              │
│ - connect()  │ │ - scan()     │ │ - showIdle() │ │ - validate() │
│ - reconnect()│ │ - getUID()   │ │ - showMsg()  │ │ - attendance│
│ - isConnected│ │              │ │ - showTime() │ │ - pairing()  │
└──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘
         │              │              │              │
         ▼              ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ ESP8266WiFi  │ │ MFRC522      │ │ LiquidCrystal│ │ HTTPClient   │
│ (Hardware)   │ │ (SPI)        │ │ _I2C         │ │ + ArduinoJson│
└──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘
```

## Module Design

### 1. ConfigManager (`config.h`, `config.cpp`)
Mengelola konfigurasi yang tersimpan di LittleFS.

```cpp
struct Config {
  char wifi_ssid[32];
  char wifi_password[64];
  char api_key[65];
  char server_url[128];
  bool configured;
};

class ConfigManager {
public:
  bool begin();
  bool load();
  bool save();
  Config& get();
  void setWiFi(const char* ssid, const char* password);
  void setAPIKey(const char* key);
  void setServerURL(const char* url);
  bool isConfigured();
private:
  Config config;
  const char* CONFIG_FILE = "/config.json";
};
```

### 2. WiFiManager (`wifi_manager.h`, `wifi_manager.cpp`)
Mengelola koneksi WiFi dengan auto-reconnect.

```cpp
class WiFiManager {
public:
  void begin(const char* ssid, const char* password);
  bool connect(unsigned long timeout = 30000);
  void checkConnection();
  bool isConnected();
  String getIP();
  int getRSSI();
private:
  const char* ssid;
  const char* password;
  unsigned long lastReconnectAttempt;
  const unsigned long RECONNECT_INTERVAL = 10000;
};
```

### 3. RFIDReader (`rfid_reader.h`, `rfid_reader.cpp`)
Wrapper untuk MFRC522 dengan format UID yang konsisten.

```cpp
// Pin definitions for NodeMCU V3
#define RFID_SS_PIN   D4  // GPIO2
#define RFID_RST_PIN  D3  // GPIO0

class RFIDReader {
public:
  void begin();
  bool isCardPresent();
  String readCardUID();
  void halt();
private:
  MFRC522 mfrc522;
  String formatUID(byte* uid, byte size);
};
```

### 4. DisplayManager (`display_manager.h`, `display_manager.cpp`)
Mengelola LCD I2C 16x2 dengan helper functions.

```cpp
// I2C address biasanya 0x27 atau 0x3F
#define LCD_ADDRESS 0x27
#define LCD_COLS 16
#define LCD_ROWS 2

class DisplayManager {
public:
  void begin();
  void clear();
  void showMessage(const char* line1, const char* line2 = "");
  void showIdle(const char* time);
  void showPairingMode(const char* studentName);
  void showSuccess(const char* name, const char* status);
  void showError(const char* message);
  void showWiFiStatus(bool connected, const char* ip = "");
private:
  LiquidCrystal_I2C lcd;
  void truncateString(const char* src, char* dest, int maxLen);
};
```

### 5. BuzzerManager (`buzzer_manager.h`, `buzzer_manager.cpp`)
Mengelola feedback audio dengan pola berbeda.

```cpp
#define BUZZER_PIN D8  // GPIO15

class BuzzerManager {
public:
  void begin();
  void beepSuccess();      // 1x pendek
  void beepPairingOK();    // 2x pendek
  void beepError();        // 3x pendek
  void beepNetworkError(); // 2x panjang
  void beepStartup();      // 1x panjang
private:
  void beep(int duration);
  void beepPattern(int count, int onTime, int offTime);
};
```

### 6. APIClient (`api_client.h`, `api_client.cpp`)
HTTP client untuk komunikasi dengan backend.

```cpp
struct AttendanceResponse {
  bool success;
  String studentName;
  String status;      // on_time, late, very_late
  String message;
};

struct PairingResponse {
  bool success;
  String studentName;
  String message;
  bool sessionActive;
};

struct ValidationResponse {
  bool valid;
  uint32_t deviceId;
  uint32_t schoolId;
};

class APIClient {
public:
  void begin(const char* serverUrl, const char* apiKey);
  ValidationResponse validateAPIKey();
  AttendanceResponse recordAttendance(const String& rfidCode);
  PairingResponse processPairing(const String& rfidCode);
  PairingResponse checkPairingStatus();
private:
  String serverUrl;
  String apiKey;
  HTTPClient http;
  const int TIMEOUT = 5000;
  const int MAX_RETRIES = 3;
  
  String buildURL(const char* endpoint);
  bool sendRequest(const char* endpoint, const String& payload, String& response);
};
```

### 7. TimeManager (`time_manager.h`, `time_manager.cpp`)
Sinkronisasi waktu dengan NTP.

```cpp
class TimeManager {
public:
  void begin();
  void update();
  String getFormattedTime();  // HH:MM:SS
  String getFormattedDate();  // DD/MM/YYYY
  bool isSynced();
private:
  WiFiUDP ntpUDP;
  NTPClient timeClient;
  unsigned long lastSync;
  const unsigned long SYNC_INTERVAL = 3600000; // 1 jam
  const long UTC_OFFSET = 7 * 3600; // WIB UTC+7
};
```

### 8. StateMachine (`state_machine.h`, `state_machine.cpp`)
Mengelola state aplikasi.

```cpp
enum class DeviceState {
  INITIALIZING,
  CONNECTING_WIFI,
  VALIDATING_API,
  IDLE,
  PROCESSING_CARD,
  PAIRING_MODE,
  ERROR_WIFI,
  ERROR_API
};

class StateMachine {
public:
  void setState(DeviceState newState);
  DeviceState getState();
  bool isInPairingMode();
  void enterPairingMode(const String& studentName);
  void exitPairingMode();
  String getCurrentStudentName();
private:
  DeviceState currentState;
  String pairingStudentName;
  unsigned long stateEnteredAt;
};
```

## File Structure

```
nodemcu-rfid-firmware/
├── nodemcu-rfid-firmware.ino    # Main sketch
├── config.h                      # Configuration structures
├── config.cpp
├── wifi_manager.h               # WiFi handling
├── wifi_manager.cpp
├── rfid_reader.h                # RFID RC522 wrapper
├── rfid_reader.cpp
├── display_manager.h            # LCD I2C wrapper
├── display_manager.cpp
├── buzzer_manager.h             # Buzzer patterns
├── buzzer_manager.cpp
├── api_client.h                 # HTTP API client
├── api_client.cpp
├── time_manager.h               # NTP time sync
├── time_manager.cpp
├── state_machine.h              # Application state
├── state_machine.cpp
├── pins.h                       # Pin definitions
└── README.md                    # Setup instructions
```

## Main Loop Flow

```cpp
void loop() {
  // 1. Check WiFi connection
  wifiManager.checkConnection();
  if (!wifiManager.isConnected()) {
    stateMachine.setState(DeviceState::ERROR_WIFI);
    displayManager.showError("WiFi Terputus");
    return;
  }

  // 2. Update time
  timeManager.update();

  // 3. Check pairing status (setiap 2 detik)
  static unsigned long lastPairingCheck = 0;
  if (millis() - lastPairingCheck > 2000) {
    checkPairingStatus();
    lastPairingCheck = millis();
  }

  // 4. Show idle display jika tidak ada aktivitas
  if (stateMachine.getState() == DeviceState::IDLE) {
    displayManager.showIdle(timeManager.getFormattedTime().c_str());
  }

  // 5. Scan RFID card
  if (rfidReader.isCardPresent()) {
    handleCardTap();
  }

  // 6. Feed watchdog
  yield();
}
```

## Sequence Diagrams

### Attendance Flow
```
User          NodeMCU         Backend
 │               │               │
 │  Tap Card     │               │
 │──────────────>│               │
 │               │  POST /attendance/rfid
 │               │──────────────>│
 │               │               │ Validate API Key
 │               │               │ Find Student by RFID
 │               │               │ Record Attendance
 │               │  Response     │
 │               │<──────────────│
 │  LCD: Name    │               │
 │  Buzzer: Beep │               │
 │<──────────────│               │
```

### Pairing Flow
```
Admin         Web Admin       Backend        NodeMCU
 │               │               │               │
 │ Start Pairing │               │               │
 │──────────────>│               │               │
 │               │ POST /pairing/start           │
 │               │──────────────>│               │
 │               │               │ Create Session│
 │               │  Session OK   │               │
 │               │<──────────────│               │
 │               │               │               │
 │               │               │  GET /pairing/status (polling)
 │               │               │<──────────────│
 │               │               │  Session Active
 │               │               │──────────────>│
 │               │               │               │ LCD: Pairing Mode
 │               │               │               │
User          │               │               │
 │ Tap Card     │               │               │
 │─────────────────────────────────────────────>│
 │               │               │  POST /pairing/rfid
 │               │               │<──────────────│
 │               │               │ Link RFID to Student
 │               │               │  Success      │
 │               │               │──────────────>│
 │               │               │               │ LCD: Terdaftar
 │               │               │               │ Buzzer: 2x beep
```

## Error Handling Strategy

| Error Type | Detection | Recovery | User Feedback |
|------------|-----------|----------|---------------|
| WiFi Disconnect | `WiFi.status() != WL_CONNECTED` | Auto reconnect setiap 10s | LCD: "WiFi Terputus" |
| API Timeout | HTTPClient timeout 5s | Retry 3x | LCD: "Error Jaringan" |
| Invalid API Key | Response `success: false` | Show error, require reconfig | LCD: "API Key Invalid" |
| RFID Read Error | `mfrc522.PICC_ReadCardSerial()` fail | Ignore, wait next card | - |
| Unknown Card | Response `success: false` | Show message | LCD: "Kartu Tidak Dikenal" |
| System Hang | Watchdog timeout 8s | Auto restart | - |

## Memory Optimization

1. **Avoid String concatenation** - Gunakan char arrays dan snprintf
2. **Use F() macro** - Untuk string literals di flash
3. **Minimize JSON buffer** - StaticJsonDocument dengan size yang tepat
4. **Release HTTP resources** - Selalu panggil `http.end()`
5. **Use yield()** - Untuk mencegah watchdog reset

## Configuration via Serial

Format command untuk konfigurasi:
```
WIFI:ssid:password
APIKEY:your_api_key_here
SERVER:http://192.168.1.100:8080
STATUS
RESET
```

Contoh:
```
WIFI:SekolahWiFi:password123
APIKEY:a1b2c3d4e5f6...
SERVER:http://api.sekolah.local:8080
```
