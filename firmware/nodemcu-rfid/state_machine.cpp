/**
 * StateMachine Implementation
 */

#include "state_machine.h"

StateMachine::StateMachine() 
  : currentState(DeviceState::INITIALIZING), 
    previousState(DeviceState::INITIALIZING),
    pairingStudentName(""),
    stateEnteredAt(0) {
}

void StateMachine::setState(DeviceState newState) {
  if (newState != currentState) {
    previousState = currentState;
    currentState = newState;
    stateEnteredAt = millis();
    
    Serial.print(F("[STATE] "));
    Serial.print(getStateName());
    Serial.println();
  }
}

DeviceState StateMachine::getState() {
  return currentState;
}

String StateMachine::getStateName() {
  switch (currentState) {
    case DeviceState::INITIALIZING:
      return F("INITIALIZING");
    case DeviceState::CONNECTING_WIFI:
      return F("CONNECTING_WIFI");
    case DeviceState::VALIDATING_API:
      return F("VALIDATING_API");
    case DeviceState::IDLE:
      return F("IDLE");
    case DeviceState::PROCESSING_CARD:
      return F("PROCESSING_CARD");
    case DeviceState::PAIRING_MODE:
      return F("PAIRING_MODE");
    case DeviceState::SHOWING_RESULT:
      return F("SHOWING_RESULT");
    case DeviceState::ERROR_WIFI:
      return F("ERROR_WIFI");
    case DeviceState::ERROR_API:
      return F("ERROR_API");
    default:
      return F("UNKNOWN");
  }
}

bool StateMachine::isInPairingMode() {
  return currentState == DeviceState::PAIRING_MODE;
}

void StateMachine::enterPairingMode(const String& studentName) {
  pairingStudentName = studentName;
  setState(DeviceState::PAIRING_MODE);
  
  Serial.print(F("[PAIRING] Entered pairing mode for: "));
  Serial.println(studentName);
}

void StateMachine::exitPairingMode() {
  pairingStudentName = "";
  setState(DeviceState::IDLE);
  
  Serial.println(F("[PAIRING] Exited pairing mode"));
}

String StateMachine::getCurrentStudentName() {
  return pairingStudentName;
}

unsigned long StateMachine::getTimeInState() {
  return millis() - stateEnteredAt;
}

bool StateMachine::hasTimedOut(unsigned long timeout) {
  return getTimeInState() >= timeout;
}

bool StateMachine::isInErrorState() {
  return currentState == DeviceState::ERROR_WIFI || 
         currentState == DeviceState::ERROR_API;
}

bool StateMachine::isReadyForCard() {
  return currentState == DeviceState::IDLE || 
         currentState == DeviceState::PAIRING_MODE;
}
