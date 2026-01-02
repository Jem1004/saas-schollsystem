/**
 * Pin Definitions for NodeMCU V3 (ESP8266)
 * 
 * Hardware Configuration:
 * - RFID RC522 (SPI)
 * - LCD I2C 16x2
 * - Buzzer
 */

#ifndef PINS_H
#define PINS_H

// ============================================
// RFID RC522 - SPI Interface
// ============================================
#define RFID_SS_PIN   D4    // GPIO2  - SDA/SS
#define RFID_RST_PIN  D3    // GPIO0  - RST
// SCK  -> D5 (GPIO14) - Hardware SPI
// MOSI -> D7 (GPIO13) - Hardware SPI
// MISO -> D6 (GPIO12) - Hardware SPI
// 3.3V -> 3V3
// GND  -> GND

// ============================================
// LCD I2C 16x2
// ============================================
#define LCD_SDA_PIN   D2    // GPIO4  - I2C SDA
#define LCD_SCL_PIN   D1    // GPIO5  - I2C SCL
#define LCD_ADDRESS   0x27  // Default I2C address (bisa juga 0x3F)
#define LCD_COLS      16
#define LCD_ROWS      2

// ============================================
// Buzzer
// ============================================
#define BUZZER_PIN    D8    // GPIO15

// ============================================
// Built-in LED (untuk debugging)
// ============================================
#define LED_BUILTIN_PIN D0  // GPIO16 - Built-in LED (active LOW)

#endif // PINS_H
