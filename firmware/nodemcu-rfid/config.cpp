/**
 * ConfigManager Implementation
 */

#include "config.h"

ConfigManager::ConfigManager() : initialized(false) {
  setDefaults();
}

bool ConfigManager::begin() {
  if (!LittleFS.begin()) {
    Serial.println(F("[ERROR] ConfigManager: Failed to mount LittleFS"));
    return false;
  }
  
  initialized = true;
  Serial.println(F("[OK] ConfigManager initialized"));
  return true;
}

bool ConfigManager::load() {
  if (!initialized) {
    Serial.println(F("[ERROR] ConfigManager: Not initialized"));
    return false;
  }

  // Check if config file exists
  if (!LittleFS.exists(CONFIG_FILE_PATH)) {
    Serial.println(F("[INFO] ConfigManager: No config file found, using defaults"));
    return false;
  }

  // Open config file
  File file = LittleFS.open(CONFIG_FILE_PATH, "r");
  if (!file) {
    Serial.println(F("[ERROR] ConfigManager: Failed to open config file"));
    return false;
  }

  // Parse JSON
  StaticJsonDocument<512> doc;
  DeserializationError error = deserializeJson(doc, file);
  file.close();

  if (error) {
    Serial.print(F("[ERROR] ConfigManager: JSON parse error: "));
    Serial.println(error.c_str());
    return false;
  }

  // Load values from JSON
  strlcpy(config.wifi_ssid, doc["wifi_ssid"] | "", CONFIG_WIFI_SSID_MAX);
  strlcpy(config.wifi_password, doc["wifi_password"] | "", CONFIG_WIFI_PASSWORD_MAX);
  strlcpy(config.api_key, doc["api_key"] | "", CONFIG_API_KEY_MAX);
  strlcpy(config.server_url, doc["server_url"] | "", CONFIG_SERVER_URL_MAX);
  config.configured = doc["configured"] | false;

  Serial.println(F("[OK] ConfigManager: Configuration loaded"));
  return true;
}

bool ConfigManager::save() {
  if (!initialized) {
    Serial.println(F("[ERROR] ConfigManager: Not initialized"));
    return false;
  }

  // Create JSON document
  StaticJsonDocument<512> doc;
  doc["wifi_ssid"] = config.wifi_ssid;
  doc["wifi_password"] = config.wifi_password;
  doc["api_key"] = config.api_key;
  doc["server_url"] = config.server_url;
  doc["configured"] = config.configured;

  // Open file for writing
  File file = LittleFS.open(CONFIG_FILE_PATH, "w");
  if (!file) {
    Serial.println(F("[ERROR] ConfigManager: Failed to create config file"));
    return false;
  }

  // Serialize JSON to file
  if (serializeJson(doc, file) == 0) {
    Serial.println(F("[ERROR] ConfigManager: Failed to write config file"));
    file.close();
    return false;
  }

  file.close();
  Serial.println(F("[OK] ConfigManager: Configuration saved"));
  return true;
}

void ConfigManager::reset() {
  setDefaults();
  
  // Delete config file if exists
  if (initialized && LittleFS.exists(CONFIG_FILE_PATH)) {
    LittleFS.remove(CONFIG_FILE_PATH);
  }
  
  Serial.println(F("[OK] ConfigManager: Configuration reset to defaults"));
}

Config& ConfigManager::get() {
  return config;
}

void ConfigManager::setWiFi(const char* ssid, const char* password) {
  strlcpy(config.wifi_ssid, ssid, CONFIG_WIFI_SSID_MAX);
  strlcpy(config.wifi_password, password, CONFIG_WIFI_PASSWORD_MAX);
  updateConfiguredStatus();
}

void ConfigManager::setAPIKey(const char* key) {
  strlcpy(config.api_key, key, CONFIG_API_KEY_MAX);
  updateConfiguredStatus();
}

void ConfigManager::setServerURL(const char* url) {
  strlcpy(config.server_url, url, CONFIG_SERVER_URL_MAX);
  updateConfiguredStatus();
}

bool ConfigManager::isConfigured() {
  return config.configured;
}

void ConfigManager::printConfig() {
  Serial.println(F("=== Current Configuration ==="));
  Serial.print(F("WiFi SSID: "));
  Serial.println(config.wifi_ssid[0] != '\0' ? config.wifi_ssid : "(not set)");
  Serial.print(F("WiFi Password: "));
  Serial.println(config.wifi_password[0] != '\0' ? "********" : "(not set)");
  Serial.print(F("API Key: "));
  if (config.api_key[0] != '\0') {
    // Show only first 8 chars for security
    Serial.print(config.api_key[0]);
    Serial.print(config.api_key[1]);
    Serial.print(config.api_key[2]);
    Serial.print(config.api_key[3]);
    Serial.println(F("...****"));
  } else {
    Serial.println(F("(not set)"));
  }
  Serial.print(F("Server URL: "));
  Serial.println(config.server_url[0] != '\0' ? config.server_url : "(not set)");
  Serial.print(F("Configured: "));
  Serial.println(config.configured ? "Yes" : "No");
  Serial.println(F("============================="));
}

void ConfigManager::setDefaults() {
  // Load from hardcoded config in config.h
  strlcpy(config.wifi_ssid, WIFI_SSID, CONFIG_WIFI_SSID_MAX);
  strlcpy(config.wifi_password, WIFI_PASSWORD, CONFIG_WIFI_PASSWORD_MAX);
  strlcpy(config.api_key, API_KEY, CONFIG_API_KEY_MAX);
  strlcpy(config.server_url, SERVER_URL, CONFIG_SERVER_URL_MAX);
  updateConfiguredStatus();
}

void ConfigManager::updateConfiguredStatus() {
  // Device is configured if all required fields are set
  config.configured = (
    config.wifi_ssid[0] != '\0' &&
    config.api_key[0] != '\0' &&
    config.server_url[0] != '\0'
  );
}
