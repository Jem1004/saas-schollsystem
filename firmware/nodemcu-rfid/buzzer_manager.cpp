/**
 * BuzzerManager Implementation
 */

#include "buzzer_manager.h"

void BuzzerManager::begin() {
  pinMode(BUZZER_PIN, OUTPUT);
  digitalWrite(BUZZER_PIN, LOW);
  Serial.println(F("[OK] BuzzerManager initialized"));
}

void BuzzerManager::beepSuccess() {
  // 1x short beep (100ms) - AC4.2
  beep(100);
}

void BuzzerManager::beepPairingOK() {
  // 2x short beep (100ms, 100ms gap) - AC4.2
  beepPattern(2, 100, 100);
}

void BuzzerManager::beepError() {
  // 3x short beep (100ms, 100ms gap) - AC4.2
  beepPattern(3, 100, 100);
}

void BuzzerManager::beepNetworkError() {
  // 2x long beep (500ms, 200ms gap) - AC4.2
  beepPattern(2, 500, 200);
}

void BuzzerManager::beepStartup() {
  // 1x long beep (300ms) - AC4.2
  beep(300);
}

void BuzzerManager::beep(int duration) {
  digitalWrite(BUZZER_PIN, HIGH);
  delay(duration);
  digitalWrite(BUZZER_PIN, LOW);
}

void BuzzerManager::beepPattern(int count, int onTime, int offTime) {
  for (int i = 0; i < count; i++) {
    beep(onTime);
    if (i < count - 1) {
      delay(offTime);
    }
  }
}
