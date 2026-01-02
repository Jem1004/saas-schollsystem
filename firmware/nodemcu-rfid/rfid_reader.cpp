/**
 * RFIDReader Implementation
 */

#include "rfid_reader.h"

RFIDReader::RFIDReader() : mfrc522(RFID_SS_PIN, RFID_RST_PIN) {
}

void RFIDReader::begin() {
  // Initialize SPI bus
  SPI.begin();
  
  // Initialize MFRC522
  mfrc522.PCD_Init();
  
  // Small delay for initialization
  delay(50);
  
  // Check if reader is connected
  if (isConnected()) {
    Serial.println(F("[OK] RFIDReader initialized"));
    Serial.print(F("     Firmware version: 0x"));
    Serial.println(getFirmwareVersion(), HEX);
  } else {
    Serial.println(F("[ERROR] RFID reader not detected"));
  }
}

bool RFIDReader::isCardPresent() {
  // Check for new card
  if (!mfrc522.PICC_IsNewCardPresent()) {
    return false;
  }
  
  // Try to read card serial
  if (!mfrc522.PICC_ReadCardSerial()) {
    return false;
  }
  
  return true;
}

String RFIDReader::readCardUID() {
  // Format UID to hex string - AC2.2
  return formatUID(mfrc522.uid.uidByte, mfrc522.uid.size);
}

void RFIDReader::halt() {
  // Halt PICC
  mfrc522.PICC_HaltA();
  
  // Stop encryption on PCD
  mfrc522.PCD_StopCrypto1();
}

byte RFIDReader::getFirmwareVersion() {
  return mfrc522.PCD_ReadRegister(mfrc522.VersionReg);
}

bool RFIDReader::isConnected() {
  byte version = getFirmwareVersion();
  // Valid versions: 0x91 (v1.0), 0x92 (v2.0), 0x88 (clone)
  return (version == 0x91 || version == 0x92 || version == 0x88 || version == 0x12);
}

String RFIDReader::formatUID(byte* uid, byte size) {
  String result = "";
  
  for (byte i = 0; i < size; i++) {
    // Add leading zero if needed
    if (uid[i] < 0x10) {
      result += "0";
    }
    result += String(uid[i], HEX);
  }
  
  // Convert to uppercase - AC2.2
  result.toUpperCase();
  
  return result;
}
