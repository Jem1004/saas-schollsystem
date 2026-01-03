/**
 * NodeMCU RFID Attendance Firmware
 * 
 * Edit konfigurasi di config.h sebelum upload.
 * 
 * @version 2.1.0
 */

#include <ESP8266WiFi.h>
#include <LittleFS.h>

#include "pins.h"
#include "config.h"
#include "buzzer_manager.h"
#include "display_manager.h"
#include "rfid_reader.h"
#include "wifi_manager.h"
#include "time_manager.h"
#include "api_client.h"
#include "state_machine.h"

// Global instances
BuzzerManager buzzerManager;
DisplayManager displayManager;
RFIDReader rfidReader;
ConfigManager configManager;
WiFiSetup wifiSetup;
TimeManager timeManager;
APIClient apiClient;
StateMachine stateMachine;

#define SERIAL_BAUD 115200
#define PAIRING_CHECK_INTERVAL 3000  // Check pairing every 3 seconds (reduced network load)
#define CARD_COOLDOWN 1000           // 1 second cooldown between taps (backend handles duplicates)
#define RESULT_DISPLAY_TIME 1500     // Show result for 1.5 seconds (enough to read)
#define DISPLAY_UPDATE_INTERVAL 1000

unsigned long lastPairingCheck = 0;
unsigned long lastCardTap = 0;
unsigned long lastDisplayUpdate = 0;
unsigned long resultDisplayStart = 0;
bool apiValidated = false;

// Forward declarations
void checkPairingStatus();
void handleCardTap();
void handleAttendanceTap(const String& uid);
void handlePairingTap(const String& uid);
String mapStatusToIndonesian(const String& status);

void setup() {
  Serial.begin(SERIAL_BAUD);
  Serial.println();
  Serial.println(F("================================"));
  Serial.println(F("NodeMCU RFID Attendance System"));
  Serial.println(F("Version 2.1.0"));
  Serial.println(F("================================"));

  // Show config
  Serial.println(F("\n[Config]"));
  Serial.print(F("  WiFi: ")); Serial.println(WIFI_SSID);
  Serial.print(F("  Server: ")); Serial.println(SERVER_URL);
  Serial.print(F("  API Key: ")); Serial.print(String(API_KEY).substring(0, 8)); Serial.println(F("..."));
  Serial.println();

  // Init hardware
  displayManager.begin();
  buzzerManager.begin();
  rfidReader.begin();

  // Connect WiFi
  displayManager.showMessage("Connecting...", WIFI_SSID);
  
  if (wifiSetup.begin(WIFI_SSID, WIFI_PASSWORD)) {
    displayManager.showWiFiStatus(true, wifiSetup.getIP().c_str());
    buzzerManager.beepSuccess();
    delay(1000);
    
    // Sync time
    timeManager.begin();
    
    // Init API client
    apiClient.begin(SERVER_URL, API_KEY);
    
    // Validate API key
    displayManager.showMessage("Validating...", "API Key");
    ValidationResponse validation = apiClient.validateAPIKey();
    
    if (validation.success && validation.valid) {
      Serial.println(F("[OK] API key valid"));
      Serial.print(F("  Device ID: ")); Serial.println(validation.deviceId);
      Serial.print(F("  School ID: ")); Serial.println(validation.schoolId);
      displayManager.showMessage("API Valid", "Device OK");
      buzzerManager.beepSuccess();
      apiValidated = true;
      delay(1000);
    } else {
      Serial.print(F("[ERROR] API validation: "));
      Serial.println(validation.message);
      displayManager.showError("API Invalid");
      buzzerManager.beepError();
      delay(2000);
    }
  } else {
    displayManager.showError("WiFi Failed");
    buzzerManager.beepNetworkError();
    delay(2000);
  }

  buzzerManager.beepStartup();
  stateMachine.setState(DeviceState::IDLE);
  displayManager.showMessage("System Ready", "Tap Kartu...");
  
  Serial.println(F("\n[OK] Setup complete\n"));
}

void loop() {
  // Check WiFi
  if (!wifiSetup.isConnected() && stateMachine.getState() == DeviceState::IDLE) {
    stateMachine.setState(DeviceState::ERROR_WIFI);
  } else if (wifiSetup.isConnected() && stateMachine.getState() == DeviceState::ERROR_WIFI) {
    stateMachine.setState(DeviceState::IDLE);
  }

  timeManager.update();

  // Check pairing status
  if (apiValidated && wifiSetup.isConnected()) {
    // Don't check pairing status while processing card to avoid WiFi conflicts
    if ((stateMachine.getState() == DeviceState::IDLE || stateMachine.getState() == DeviceState::PAIRING_MODE) &&
        stateMachine.getState() != DeviceState::PROCESSING_CARD) {
      if (millis() - lastPairingCheck > PAIRING_CHECK_INTERVAL) {
        checkPairingStatus();
        lastPairingCheck = millis();
      }
    }
  }

  // Handle result display timeout
  if (stateMachine.getState() == DeviceState::SHOWING_RESULT) {
    if (millis() - resultDisplayStart > RESULT_DISPLAY_TIME) {
      stateMachine.setState(DeviceState::IDLE);
    }
  }

  // Update display
  if (stateMachine.getState() == DeviceState::IDLE || 
      stateMachine.getState() == DeviceState::PAIRING_MODE) {
    if (millis() - lastDisplayUpdate > DISPLAY_UPDATE_INTERVAL) {
      if (stateMachine.isInPairingMode()) {
        displayManager.showPairingMode(stateMachine.getCurrentStudentName().c_str());
      } else {
        displayManager.showIdle(timeManager.isSynced() ? 
          timeManager.getFormattedTime().c_str() : "--:--:--");
      }
      lastDisplayUpdate = millis();
    }
  }

  // Scan RFID
  if (stateMachine.isReadyForCard() && rfidReader.isCardPresent()) {
    handleCardTap();
  }

  yield();
}

void checkPairingStatus() {
  if (!apiClient.isConfigured() || !wifiSetup.isConnected()) return;
  if (stateMachine.getState() != DeviceState::IDLE && stateMachine.getState() != DeviceState::PAIRING_MODE) return;
  
  PairingResponse response = apiClient.checkPairingStatus();
  
  if (response.success && response.sessionActive) {
    // Enter pairing mode if not already in it
    if (!stateMachine.isInPairingMode()) {
      stateMachine.enterPairingMode(response.studentName);
      displayManager.showPairingMode(response.studentName.c_str());
      Serial.print(F("[PAIRING] Active for: "));
      Serial.println(response.studentName);
    }
  } else if (stateMachine.isInPairingMode()) {
    // Exit pairing mode if session is no longer active
    Serial.println(F("[PAIRING] Session ended, exiting pairing mode"));
    stateMachine.exitPairingMode();
    displayManager.showIdle(timeManager.isSynced() ? 
      timeManager.getFormattedTime().c_str() : "--:--:--");
  }
}

void handleCardTap() {
  if (millis() - lastCardTap < CARD_COOLDOWN) {
    rfidReader.halt();
    return;
  }
  lastCardTap = millis();
  
  String uid = rfidReader.readCardUID();
  rfidReader.halt();
  
  Serial.print(F("[RFID] UID: "));
  Serial.println(uid);
  
  Serial.print(F("[DEBUG] Current state: "));
  Serial.print(stateMachine.getStateName());
  Serial.print(F(", isInPairingMode: "));
  Serial.println(stateMachine.isInPairingMode());
  
  // IMPORTANT: Save pairing mode status BEFORE changing state
  bool wasInPairingMode = stateMachine.isInPairingMode();
  String pairingStudentName = stateMachine.getCurrentStudentName();
  
  stateMachine.setState(DeviceState::PROCESSING_CARD);
  displayManager.showMessage("Memproses...", uid.c_str());
  
  if (!apiClient.isConfigured() || !apiValidated) {
    displayManager.showMessage("Card UID:", uid.c_str());
    buzzerManager.beepSuccess();
    resultDisplayStart = millis();
    stateMachine.setState(DeviceState::SHOWING_RESULT);
    return;
  }
  
  if (!wifiSetup.isConnected()) {
    displayManager.showError("WiFi Terputus");
    buzzerManager.beepNetworkError();
    resultDisplayStart = millis();
    stateMachine.setState(DeviceState::SHOWING_RESULT);
    return;
  }
  
  // Use saved pairing mode status
  if (wasInPairingMode) {
    Serial.println(F("[DEBUG] Calling handlePairingTap"));
    handlePairingTap(uid);
  } else {
    Serial.println(F("[DEBUG] Calling handleAttendanceTap"));
    handleAttendanceTap(uid);
  }
}

void handleAttendanceTap(const String& uid) {
  AttendanceResponse response = apiClient.recordAttendance(uid);
  
  if (response.success) {
    String status = mapStatusToIndonesian(response.status);
    displayManager.showSuccess(response.studentName.c_str(), status.c_str());
    buzzerManager.beepSuccess();
    Serial.print(F("[OK] ")); Serial.print(response.studentName);
    Serial.print(F(" - ")); Serial.println(status);
  } 
  else if (response.errorCode == "VAL_NO_SCHEDULE") {
    // No schedule configured for this time
    displayManager.showError("Diluar Jadwal");
    buzzerManager.beepError();
    Serial.println(F("[WARN] No schedule for current time"));
  }
  else if (response.errorCode == "VAL_ALREADY_CHECKED_IN") {
    // Already checked in for this schedule
    displayManager.showMessage("Sudah Absen", response.studentName.c_str());
    buzzerManager.beepError();
    Serial.println(F("[WARN] Already checked in for this schedule"));
  }
  else if (response.errorCode == "VAL_INVALID_RFID" || response.message == F("card_not_found")) {
    // Card not registered
    displayManager.showError("Kartu Tidak Dikenal");
    buzzerManager.beepError();
    Serial.println(F("[WARN] Card not registered"));
  }
  else if (response.errorCode == "VAL_ALREADY_CHECKED_OUT") {
    // Already checked out today
    displayManager.showError("Sudah Pulang");
    buzzerManager.beepError();
    Serial.println(F("[WARN] Already checked out"));
  }
  else if (response.errorCode.length() > 0) {
    // Other known error from backend
    displayManager.showError(response.message.c_str());
    buzzerManager.beepError();
    Serial.print(F("[ERROR] ")); Serial.println(response.errorCode);
  }
  else {
    // Network or unknown error
    displayManager.showError("Error Jaringan");
    buzzerManager.beepNetworkError();
    Serial.print(F("[ERROR] ")); Serial.println(response.message);
  }
  
  resultDisplayStart = millis();
  stateMachine.setState(DeviceState::SHOWING_RESULT);
}

void handlePairingTap(const String& uid) {
  PairingResponse response = apiClient.processPairing(uid);
  
  if (response.success) {
    displayManager.showSuccess("Kartu Terdaftar", response.studentName.c_str());
    buzzerManager.beepPairingOK();
    Serial.print(F("[OK] Registered to ")); Serial.println(response.studentName);
  }
  else if (response.message == F("card_already_used")) {
    displayManager.showError("Kartu Sudah Dipakai");
    buzzerManager.beepError();
  }
  else {
    displayManager.showError("Gagal Mendaftar");
    buzzerManager.beepError();
  }
  
  resultDisplayStart = millis();
  stateMachine.setState(DeviceState::SHOWING_RESULT);
}

String mapStatusToIndonesian(const String& status) {
  if (status == "on_time") return F("Tepat Waktu");
  if (status == "late") return F("Terlambat");
  if (status == "very_late") return F("Sangat Terlambat");
  if (status == "check_in") return F("Masuk");
  if (status == "check_out") return F("Pulang");
  if (status == "already_checked_in") return F("Sudah Masuk");
  if (status == "already_checked_out") return F("Sudah Pulang");
  return status;
}
