/**
 * DisplayManager - LCD I2C 16x2 management
 * 
 * Handles all LCD display operations with helper functions
 * for common display patterns.
 */

#ifndef DISPLAY_MANAGER_H
#define DISPLAY_MANAGER_H

#include <Arduino.h>
#include <Wire.h>
#include <LiquidCrystal_I2C.h>
#include "pins.h"

class DisplayManager {
public:
  DisplayManager();

  /**
   * Initialize LCD with auto-detect I2C address
   * Tries 0x27 first, then 0x3F
   */
  void begin();

  /**
   * Clear the display
   */
  void clear();

  /**
   * Show custom message on both lines
   * @param line1 Text for first line (max 16 chars)
   * @param line2 Text for second line (max 16 chars, optional)
   */
  void showMessage(const char* line1, const char* line2 = "");

  /**
   * Show idle screen with time
   * Line 1: "Tap Kartu..."
   * Line 2: Current time (HH:MM:SS)
   * @param time Formatted time string
   */
  void showIdle(const char* time);

  /**
   * Show pairing mode screen
   * Line 1: "Mode Pairing"
   * Line 2: Student name (truncated if needed)
   * @param studentName Name of student to pair
   */
  void showPairingMode(const char* studentName);

  /**
   * Show success screen
   * Line 1: Student name (truncated if needed)
   * Line 2: Status message
   * @param name Student name
   * @param status Status text (e.g., "Tepat Waktu")
   */
  void showSuccess(const char* name, const char* status);

  /**
   * Show error screen
   * Line 1: "Error"
   * Line 2: Error message
   * @param message Error description
   */
  void showError(const char* message);

  /**
   * Show WiFi status
   * @param connected Connection status
   * @param ip IP address (optional, shown if connected)
   */
  void showWiFiStatus(bool connected, const char* ip = "");

  /**
   * Show connecting animation
   * @param ssid WiFi SSID being connected to
   */
  void showConnecting(const char* ssid);

private:
  LiquidCrystal_I2C* lcd;
  uint8_t lcdAddress;
  bool initialized;

  /**
   * Truncate string to fit LCD width
   * @param src Source string
   * @param dest Destination buffer
   * @param maxLen Maximum length (default 16)
   */
  void truncateString(const char* src, char* dest, int maxLen = LCD_COLS);

  /**
   * Center text on LCD line
   * @param text Text to center
   * @param dest Destination buffer
   * @param width Line width (default 16)
   */
  void centerText(const char* text, char* dest, int width = LCD_COLS);

  /**
   * Detect I2C address of LCD
   * @return Detected address or 0 if not found
   */
  uint8_t detectI2CAddress();
};

#endif // DISPLAY_MANAGER_H
