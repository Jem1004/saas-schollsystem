/**
 * WiFiSetup Implementation - Simple approach
 */

#include "wifi_manager.h"

WiFiSetup::WiFiSetup() : connected(false) {
}

bool WiFiSetup::begin(const char* ssid, const char* password) {
  Serial.print(F("[WiFi] Connecting to "));
  Serial.println(ssid);
  
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);
  
  unsigned long start = millis();
  while (WiFi.status() != WL_CONNECTED) {
    if (millis() - start > WIFI_CONNECT_TIMEOUT) {
      Serial.println(F("\n[WiFi] Connection timeout!"));
      return false;
    }
    delay(500);
    Serial.print(".");
    yield();
  }
  
  connected = true;
  Serial.println();
  Serial.println(F("[WiFi] Connected!"));
  Serial.print(F("[WiFi] IP: "));
  Serial.println(WiFi.localIP());
  Serial.print(F("[WiFi] RSSI: "));
  Serial.print(WiFi.RSSI());
  Serial.println(F(" dBm"));
  
  return true;
}

bool WiFiSetup::isConnected() {
  return WiFi.status() == WL_CONNECTED;
}

String WiFiSetup::getIP() {
  return WiFi.localIP().toString();
}

int WiFiSetup::getRSSI() {
  return WiFi.RSSI();
}

String WiFiSetup::getSSID() {
  return WiFi.SSID();
}
