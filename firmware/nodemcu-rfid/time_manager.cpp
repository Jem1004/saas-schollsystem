/**
 * TimeManager Implementation
 */

#include "time_manager.h"

TimeManager::TimeManager() 
  : ntpUDP(nullptr), timeClient(nullptr), lastSync(0), synced(false) {
}

TimeManager::~TimeManager() {
  if (timeClient != nullptr) {
    delete timeClient;
  }
  if (ntpUDP != nullptr) {
    delete ntpUDP;
  }
}

void TimeManager::begin() {
  // Create UDP instance
  ntpUDP = new WiFiUDP();
  
  // Create NTP client with WITA timezone (UTC+8)
  timeClient = new NTPClient(*ntpUDP, NTP_SERVER, NTP_UTC_OFFSET_WITA, NTP_SYNC_INTERVAL);
  
  // Start NTP client
  timeClient->begin();
  
  Serial.println(F("[OK] TimeManager initialized"));
  Serial.print(F("     NTP Server: "));
  Serial.println(NTP_SERVER);
  Serial.println(F("     Timezone: WITA (UTC+8)"));
  
  // Initial sync - AC6.1
  forceSync();
}

void TimeManager::update() {
  if (timeClient == nullptr) return;
  
  // Update NTP client (handles periodic sync internally)
  timeClient->update();
  
  // Check if we need to resync - AC6.3
  unsigned long now = millis();
  if (synced && (now - lastSync >= NTP_SYNC_INTERVAL)) {
    Serial.println(F("[INFO] TimeManager: Periodic resync..."));
    forceSync();
  }
}

bool TimeManager::forceSync() {
  if (timeClient == nullptr) return false;
  
  Serial.println(F("[INFO] TimeManager: Syncing with NTP..."));
  
  // Force update
  bool success = timeClient->forceUpdate();
  
  if (success) {
    synced = true;
    lastSync = millis();
    Serial.print(F("[OK] TimeManager: Time synced - "));
    Serial.println(getFormattedDateTime());
  } else {
    // AC6.4 - Handle NTP failure gracefully
    Serial.println(F("[WARN] TimeManager: NTP sync failed, using last known time"));
    // Keep synced status if we had a previous sync
  }
  
  return success;
}

String TimeManager::getFormattedTime() {
  if (timeClient == nullptr) return F("--:--:--");
  
  // Format: HH:MM:SS
  return padZero(getHours()) + ":" + padZero(getMinutes()) + ":" + padZero(getSeconds());
}

String TimeManager::getFormattedDate() {
  if (timeClient == nullptr || !synced) return F("--/--/----");
  
  // Get epoch time and calculate date
  unsigned long epochTime = timeClient->getEpochTime();
  
  // Convert epoch to date components
  // This is a simplified calculation
  time_t rawTime = epochTime;
  struct tm* timeInfo = gmtime(&rawTime);
  
  // Adjust for local timezone (already applied in NTP offset)
  int day = timeInfo->tm_mday;
  int month = timeInfo->tm_mon + 1;
  int year = timeInfo->tm_year + 1900;
  
  // Format: DD/MM/YYYY
  return padZero(day) + "/" + padZero(month) + "/" + String(year);
}

String TimeManager::getFormattedDateTime() {
  return getFormattedDate() + " " + getFormattedTime();
}

bool TimeManager::isSynced() {
  return synced;
}

int TimeManager::getHours() {
  if (timeClient == nullptr) return 0;
  return timeClient->getHours();
}

int TimeManager::getMinutes() {
  if (timeClient == nullptr) return 0;
  return timeClient->getMinutes();
}

int TimeManager::getSeconds() {
  if (timeClient == nullptr) return 0;
  return timeClient->getSeconds();
}

unsigned long TimeManager::getEpochTime() {
  if (timeClient == nullptr) return 0;
  return timeClient->getEpochTime();
}

String TimeManager::padZero(int num) {
  if (num < 10) {
    return "0" + String(num);
  }
  return String(num);
}
