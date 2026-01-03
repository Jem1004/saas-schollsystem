/**
 * Device Configuration
 * 
 * Edit nilai di bawah sesuai environment kamu, lalu upload firmware.
 */

#ifndef CONFIG_H
#define CONFIG_H

#include <Arduino.h>
#include <LittleFS.h>
#include <ArduinoJson.h>

// ============================================================
// KONFIGURASI DEVICE - EDIT DI SINI
// ============================================================

// WiFi Configuration
#define WIFI_SSID       "NISA 2.4"
#define WIFI_PASSWORD   "makangabus"

// Backend Server URL (tanpa trailing slash)
#define SERVER_URL      "http://192.168.100.43:3000"

// Device API Key (dari web admin saat register device)
#define API_KEY         "828c2b1b55befe80ad8a919ad7dfdcb8b6997a397b87b5f69425c49845cc29f4"

// ============================================================
// END KONFIGURASI - Jangan edit di bawah ini
// ============================================================

// Configuration limits
#define CONFIG_WIFI_SSID_MAX     32
#define CONFIG_WIFI_PASSWORD_MAX 64
#define CONFIG_API_KEY_MAX       65
#define CONFIG_SERVER_URL_MAX    128
#define CONFIG_FILE_PATH         "/config.json"

struct Config {
  char wifi_ssid[CONFIG_WIFI_SSID_MAX];
  char wifi_password[CONFIG_WIFI_PASSWORD_MAX];
  char api_key[CONFIG_API_KEY_MAX];
  char server_url[CONFIG_SERVER_URL_MAX];
  bool configured;
};

class ConfigManager {
public:
  ConfigManager();
  bool begin();
  bool load();
  bool save();
  void reset();
  Config& get();
  void setWiFi(const char* ssid, const char* password);
  void setAPIKey(const char* key);
  void setServerURL(const char* url);
  bool isConfigured();
  void printConfig();

private:
  Config config;
  bool initialized;
  void setDefaults();
  void updateConfiguredStatus();
};

#endif // CONFIG_H
