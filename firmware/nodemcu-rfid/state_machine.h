/**
 * StateMachine - Application state management
 * 
 * Manages device states and transitions between modes.
 */

#ifndef STATE_MACHINE_H
#define STATE_MACHINE_H

#include <Arduino.h>

/**
 * Device states
 */
enum class DeviceState {
  INITIALIZING,       // Device is starting up
  CONNECTING_WIFI,    // Connecting to WiFi
  VALIDATING_API,     // Validating API key
  IDLE,               // Ready for card tap (attendance mode)
  PROCESSING_CARD,    // Processing a card tap
  PAIRING_MODE,       // In pairing mode
  SHOWING_RESULT,     // Showing result on LCD
  ERROR_WIFI,         // WiFi connection error
  ERROR_API           // API error
};

class StateMachine {
public:
  StateMachine();

  /**
   * Set current state
   * @param newState New device state
   */
  void setState(DeviceState newState);

  /**
   * Get current state
   * @return Current device state
   */
  DeviceState getState();

  /**
   * Get state name as string (for debugging)
   * @return State name
   */
  String getStateName();

  /**
   * Check if device is in pairing mode
   * @return true if in pairing mode
   */
  bool isInPairingMode();

  /**
   * Enter pairing mode with student name
   * @param studentName Name of student to pair
   */
  void enterPairingMode(const String& studentName);

  /**
   * Exit pairing mode and return to idle
   */
  void exitPairingMode();

  /**
   * Get student name for current pairing session
   * @return Student name or empty string
   */
  String getCurrentStudentName();

  /**
   * Get time since state was entered
   * @return Milliseconds since state entry
   */
  unsigned long getTimeInState();

  /**
   * Check if state has timed out
   * @param timeout Timeout in milliseconds
   * @return true if timed out
   */
  bool hasTimedOut(unsigned long timeout);

  /**
   * Check if device is in error state
   * @return true if in error state
   */
  bool isInErrorState();

  /**
   * Check if device is ready for card tap
   * @return true if ready
   */
  bool isReadyForCard();

private:
  DeviceState currentState;
  DeviceState previousState;
  String pairingStudentName;
  unsigned long stateEnteredAt;
};

#endif // STATE_MACHINE_H
