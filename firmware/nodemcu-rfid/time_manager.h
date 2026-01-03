/**
 * TimeManager - NTP time synchronization
 * 
 * Handles time synchronization with NTP server.
 * Timezone set to WITA (UTC+8).
 */

#ifndef TIME_MANAGER_H
#define TIME_MANAGER_H

#include <Arduino.h>
#include <WiFiUdp.h>
#include <NTPClient.h>

// NTP settings
#define NTP_SERVER          "pool.ntp.org"
#define NTP_SYNC_INTERVAL   3600000   // 1 hour in milliseconds
#define NTP_UTC_OFFSET_WITA 28800     // UTC+8 in seconds (8 * 3600) - WITA timezone

class TimeManager {
public:
  TimeManager();
  ~TimeManager();

  /**
   * Initialize NTP client
   */
  void begin();

  /**
   * Update time from NTP server
   * Call this periodically in loop()
   */
  void update();

  /**
   * Force sync with NTP server
   * @return true if sync successful
   */
  bool forceSync();

  /**
   * Get formatted time string
   * @return Time in HH:MM:SS format
   */
  String getFormattedTime();

  /**
   * Get formatted date string
   * @return Date in DD/MM/YYYY format
   */
  String getFormattedDate();

  /**
   * Get formatted datetime string
   * @return DateTime in DD/MM/YYYY HH:MM:SS format
   */
  String getFormattedDateTime();

  /**
   * Check if time is synced
   * @return true if synced with NTP
   */
  bool isSynced();

  /**
   * Get current hour (0-23)
   */
  int getHours();

  /**
   * Get current minute (0-59)
   */
  int getMinutes();

  /**
   * Get current second (0-59)
   */
  int getSeconds();

  /**
   * Get epoch time
   * @return Unix timestamp
   */
  unsigned long getEpochTime();

private:
  WiFiUDP* ntpUDP;
  NTPClient* timeClient;
  unsigned long lastSync;
  bool synced;

  /**
   * Pad number with leading zero
   */
  String padZero(int num);
};

#endif // TIME_MANAGER_H
