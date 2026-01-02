/**
 * APIClient - HTTP API client for backend communication
 * 
 * Handles all API calls to the attendance backend server.
 */

#ifndef API_CLIENT_H
#define API_CLIENT_H

#include <Arduino.h>
#include <ESP8266HTTPClient.h>
#include <WiFiClient.h>
#include <ArduinoJson.h>

// API settings
#define API_TIMEOUT     5000    // 5 seconds timeout - AC2.9
#define API_MAX_RETRIES 3       // Max retry attempts - AC5.3

/**
 * Response structure for API key validation
 */
struct ValidationResponse {
  bool success;
  bool valid;
  uint32_t deviceId;
  uint32_t schoolId;
  String message;
};

/**
 * Response structure for attendance recording
 */
struct AttendanceResponse {
  bool success;
  String studentName;
  String status;      // on_time, late, very_late
  String message;
  String errorCode;   // Error code from backend (e.g., VAL_NO_SCHEDULE, VAL_INVALID_RFID)
};

/**
 * Response structure for pairing operations
 */
struct PairingResponse {
  bool success;
  String studentName;
  String message;
  bool sessionActive;
};

class APIClient {
public:
  APIClient();

  /**
   * Initialize API client with server URL and API key
   * @param serverUrl Base URL of the backend server
   * @param apiKey Device API key
   */
  void begin(const char* serverUrl, const char* apiKey);

  /**
   * Validate API key with backend
   * @return ValidationResponse with result
   */
  ValidationResponse validateAPIKey();

  /**
   * Record attendance for RFID card
   * @param rfidCode RFID card UID
   * @return AttendanceResponse with result
   */
  AttendanceResponse recordAttendance(const String& rfidCode);

  /**
   * Process RFID pairing
   * @param rfidCode RFID card UID
   * @return PairingResponse with result
   */
  PairingResponse processPairing(const String& rfidCode);

  /**
   * Check if there's an active pairing session
   * @return PairingResponse with session status
   */
  PairingResponse checkPairingStatus();

  /**
   * Check if API client is configured
   * @return true if server URL and API key are set
   */
  bool isConfigured();

  /**
   * Get device ID (set after successful validation)
   */
  uint32_t getDeviceId() { return deviceId; }

  /**
   * Set device ID (from validation response)
   */
  void setDeviceId(uint32_t id) { deviceId = id; }

private:
  String serverUrl;
  String apiKey;
  uint32_t deviceId;
  WiFiClient wifiClient;
  HTTPClient http;

  /**
   * Build full URL from endpoint
   * @param endpoint API endpoint path
   * @return Full URL
   */
  String buildURL(const char* endpoint);

  /**
   * Send HTTP POST request with retry logic
   * @param endpoint API endpoint
   * @param payload JSON payload
   * @param response Output response string
   * @return HTTP status code or -1 on error
   */
  int sendPostRequest(const char* endpoint, const String& payload, String& response);

  /**
   * Send HTTP GET request with retry logic
   * @param endpoint API endpoint
   * @param response Output response string
   * @return HTTP status code or -1 on error
   */
  int sendGetRequest(const char* endpoint, String& response);
};

#endif // API_CLIENT_H
