/**
 * BuzzerManager - Audio feedback management
 * 
 * Provides different beep patterns for various events:
 * - Success: 1x short beep
 * - Pairing OK: 2x short beep
 * - Error: 3x short beep
 * - Network Error: 2x long beep
 * - Startup: 1x long beep
 */

#ifndef BUZZER_MANAGER_H
#define BUZZER_MANAGER_H

#include <Arduino.h>
#include "pins.h"

class BuzzerManager {
public:
  /**
   * Initialize buzzer pin
   */
  void begin();

  /**
   * Success beep - 1x short (100ms)
   * Used for successful attendance recording
   */
  void beepSuccess();

  /**
   * Pairing OK beep - 2x short (100ms, 100ms gap)
   * Used for successful card pairing
   */
  void beepPairingOK();

  /**
   * Error beep - 3x short (100ms, 100ms gap)
   * Used for unknown card or pairing errors
   */
  void beepError();

  /**
   * Network error beep - 2x long (500ms, 200ms gap)
   * Used for network/API errors
   */
  void beepNetworkError();

  /**
   * Startup beep - 1x long (300ms)
   * Used when device successfully initializes
   */
  void beepStartup();

private:
  /**
   * Single beep with specified duration
   * @param duration Duration in milliseconds
   */
  void beep(int duration);

  /**
   * Pattern beep with multiple beeps
   * @param count Number of beeps
   * @param onTime Duration of each beep in ms
   * @param offTime Gap between beeps in ms
   */
  void beepPattern(int count, int onTime, int offTime);
};

#endif // BUZZER_MANAGER_H
