/**
 * credentials.h - Device Configuration
 * 
 * PENTING: File ini berisi credentials sensitif.
 * Jangan commit ke git! Sudah ada di .gitignore
 * 
 * Cara pakai:
 * 1. Copy credentials.example.h ke credentials.h
 * 2. Isi nilai sesuai environment kamu
 * 3. Upload ulang firmware
 */

#ifndef CREDENTIALS_H
#define CREDENTIALS_H

// WiFi Configuration
#define DEFAULT_WIFI_SSID     "NISA 2.4"
#define DEFAULT_WIFI_PASSWORD "makangabus"

// Backend Server Configuration  
#define DEFAULT_SERVER_URL    "http://192.168.100.43:3000"

// Device API Key (dari web admin saat register device)
#define DEFAULT_API_KEY       "61660fc7d308063dc701446ea74d804e55ddaaebaddbeedb1942283e08b2422c"

#endif // CREDENTIALS_H
