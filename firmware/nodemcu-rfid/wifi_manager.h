/**
 * WiFiSetup - Simple WiFi connection
 */

#ifndef WIFI_MANAGER_H
#define WIFI_MANAGER_H

#include <Arduino.h>
#include <ESP8266WiFi.h>

#define WIFI_CONNECT_TIMEOUT 30000  // 30 seconds

class WiFiSetup {
public:
  WiFiSetup();
  bool begin(const char* ssid, const char* password);
  bool isConnected();
  String getIP();
  int getRSSI();
  String getSSID();

private:
  bool connected;
};

#endif
