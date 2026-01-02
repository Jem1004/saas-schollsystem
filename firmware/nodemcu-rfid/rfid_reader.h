/**
 * RFIDReader - MFRC522 RFID Reader wrapper
 * 
 * Handles RFID card reading with consistent UID formatting.
 * Uses SPI interface on NodeMCU V3.
 */

#ifndef RFID_READER_H
#define RFID_READER_H

#include <Arduino.h>
#include <SPI.h>
#include <MFRC522.h>
#include "pins.h"

class RFIDReader {
public:
  RFIDReader();

  /**
   * Initialize RFID reader with SPI
   */
  void begin();

  /**
   * Check if a card is present on the reader
   * @return true if card detected
   */
  bool isCardPresent();

  /**
   * Read the UID of the card
   * Must call isCardPresent() first
   * @return UID as uppercase hex string (e.g., "A1B2C3D4")
   */
  String readCardUID();

  /**
   * Release the card and stop communication
   * Call after reading to allow next card
   */
  void halt();

  /**
   * Get reader firmware version
   * @return Version byte
   */
  byte getFirmwareVersion();

  /**
   * Check if reader is properly connected
   * @return true if reader responds
   */
  bool isConnected();

private:
  MFRC522 mfrc522;

  /**
   * Format UID bytes to hex string
   * @param uid UID byte array
   * @param size Size of UID (usually 4 or 7 bytes)
   * @return Uppercase hex string
   */
  String formatUID(byte* uid, byte size);
};

#endif // RFID_READER_H
