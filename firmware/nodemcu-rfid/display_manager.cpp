/**
 * DisplayManager Implementation
 */

#include "display_manager.h"

DisplayManager::DisplayManager() : lcd(nullptr), lcdAddress(0), initialized(false) {
}

void DisplayManager::begin() {
  // Initialize I2C
  Wire.begin(LCD_SDA_PIN, LCD_SCL_PIN);
  
  // Auto-detect I2C address
  lcdAddress = detectI2CAddress();
  
  if (lcdAddress == 0) {
    Serial.println(F("[ERROR] LCD not found on I2C bus"));
    return;
  }
  
  Serial.print(F("[OK] LCD found at address 0x"));
  Serial.println(lcdAddress, HEX);
  
  // Initialize LCD
  lcd = new LiquidCrystal_I2C(lcdAddress, LCD_COLS, LCD_ROWS);
  lcd->init();
  lcd->backlight();
  lcd->clear();
  
  initialized = true;
  Serial.println(F("[OK] DisplayManager initialized"));
  
  // Show startup message
  showMessage("RFID Attendance", "Starting...");
}

uint8_t DisplayManager::detectI2CAddress() {
  // Common LCD I2C addresses
  uint8_t addresses[] = {0x27, 0x3F, 0x20, 0x38};
  
  for (int i = 0; i < 4; i++) {
    Wire.beginTransmission(addresses[i]);
    if (Wire.endTransmission() == 0) {
      return addresses[i];
    }
  }
  
  return 0; // Not found
}

void DisplayManager::clear() {
  if (!initialized || lcd == nullptr) return;
  lcd->clear();
}

void DisplayManager::showMessage(const char* line1, const char* line2) {
  if (!initialized || lcd == nullptr) return;
  
  char buf1[LCD_COLS + 1];
  char buf2[LCD_COLS + 1];
  
  truncateString(line1, buf1);
  truncateString(line2, buf2);
  
  lcd->clear();
  lcd->setCursor(0, 0);
  lcd->print(buf1);
  lcd->setCursor(0, 1);
  lcd->print(buf2);
}

void DisplayManager::showIdle(const char* time) {
  if (!initialized || lcd == nullptr) return;
  
  // Line 1: "Tap Kartu..." - AC2.1
  lcd->setCursor(0, 0);
  lcd->print(F("Tap Kartu...    "));
  
  // Line 2: Time display - AC4.3
  char timeBuf[LCD_COLS + 1];
  truncateString(time, timeBuf);
  
  lcd->setCursor(0, 1);
  lcd->print(timeBuf);
  // Pad with spaces to clear old content
  for (int i = strlen(timeBuf); i < LCD_COLS; i++) {
    lcd->print(' ');
  }
}

void DisplayManager::showPairingMode(const char* studentName) {
  if (!initialized || lcd == nullptr) return;
  
  char nameBuf[LCD_COLS + 1];
  truncateString(studentName, nameBuf);
  
  lcd->clear();
  lcd->setCursor(0, 0);
  lcd->print(F("Mode Pairing"));
  lcd->setCursor(0, 1);
  lcd->print(nameBuf);
}

void DisplayManager::showSuccess(const char* name, const char* status) {
  if (!initialized || lcd == nullptr) return;
  
  char nameBuf[LCD_COLS + 1];
  char statusBuf[LCD_COLS + 1];
  
  truncateString(name, nameBuf);
  truncateString(status, statusBuf);
  
  lcd->clear();
  lcd->setCursor(0, 0);
  lcd->print(nameBuf);
  lcd->setCursor(0, 1);
  lcd->print(statusBuf);
}

void DisplayManager::showError(const char* message) {
  if (!initialized || lcd == nullptr) return;
  
  char msgBuf[LCD_COLS + 1];
  truncateString(message, msgBuf);
  
  lcd->clear();
  lcd->setCursor(0, 0);
  lcd->print(F("Error"));
  lcd->setCursor(0, 1);
  lcd->print(msgBuf);
}

void DisplayManager::showWiFiStatus(bool connected, const char* ip) {
  if (!initialized || lcd == nullptr) return;
  
  lcd->clear();
  lcd->setCursor(0, 0);
  
  if (connected) {
    lcd->print(F("WiFi: Connected"));
    if (ip != nullptr && strlen(ip) > 0) {
      char ipBuf[LCD_COLS + 1];
      truncateString(ip, ipBuf);
      lcd->setCursor(0, 1);
      lcd->print(ipBuf);
    }
  } else {
    lcd->print(F("WiFi: Terputus"));
    lcd->setCursor(0, 1);
    lcd->print(F("Reconnecting..."));
  }
}

void DisplayManager::showConnecting(const char* ssid) {
  if (!initialized || lcd == nullptr) return;
  
  char ssidBuf[LCD_COLS + 1];
  truncateString(ssid, ssidBuf);
  
  lcd->clear();
  lcd->setCursor(0, 0);
  lcd->print(F("Connecting WiFi"));
  lcd->setCursor(0, 1);
  lcd->print(ssidBuf);
}

void DisplayManager::truncateString(const char* src, char* dest, int maxLen) {
  if (src == nullptr) {
    dest[0] = '\0';
    return;
  }
  
  int srcLen = strlen(src);
  
  if (srcLen <= maxLen) {
    strcpy(dest, src);
  } else {
    // Truncate and add ellipsis if too long
    strncpy(dest, src, maxLen - 3);
    dest[maxLen - 3] = '.';
    dest[maxLen - 2] = '.';
    dest[maxLen - 1] = '.';
    dest[maxLen] = '\0';
  }
}

void DisplayManager::centerText(const char* text, char* dest, int width) {
  int textLen = strlen(text);
  int padding = (width - textLen) / 2;
  
  if (padding < 0) padding = 0;
  
  // Fill with spaces
  memset(dest, ' ', width);
  dest[width] = '\0';
  
  // Copy text to center
  int copyLen = (textLen > width) ? width : textLen;
  strncpy(dest + padding, text, copyLen);
}
