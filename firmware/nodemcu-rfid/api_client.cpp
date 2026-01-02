/**
 * APIClient Implementation
 */

#include "api_client.h"

APIClient::APIClient() : serverUrl(""), apiKey(""), deviceId(0) {
}

void APIClient::begin(const char* serverUrl, const char* apiKey) {
  this->serverUrl = String(serverUrl);
  this->apiKey = String(apiKey);
  
  // Remove trailing slash if present
  if (this->serverUrl.endsWith("/")) {
    this->serverUrl = this->serverUrl.substring(0, this->serverUrl.length() - 1);
  }
  
  Serial.println(F("[OK] APIClient initialized"));
  Serial.print(F("     Server: "));
  Serial.println(this->serverUrl);
}

bool APIClient::isConfigured() {
  return serverUrl.length() > 0 && apiKey.length() > 0;
}

ValidationResponse APIClient::validateAPIKey() {
  ValidationResponse result;
  result.success = false;
  result.valid = false;
  result.deviceId = 0;
  result.schoolId = 0;
  result.message = "";

  if (!isConfigured()) {
    result.message = F("API not configured");
    return result;
  }

  // Build JSON payload
  StaticJsonDocument<128> doc;
  doc["api_key"] = apiKey;
  
  String payload;
  serializeJson(doc, payload);
  
  Serial.print(F("[API] Validate request: "));
  Serial.println(payload);

  // Send request - AC1.5 (using /public prefix to bypass auth)
  String response;
  int httpCode = sendPostRequest("/api/v1/public/devices/validate-key", payload, response);
  
  Serial.print(F("[API] Validate response: "));
  Serial.println(response);

  if (httpCode == 200) {
    // Parse response
    StaticJsonDocument<256> responseDoc;
    DeserializationError error = deserializeJson(responseDoc, response);
    
    if (!error) {
      result.success = responseDoc["success"] | false;
      if (result.success) {
        JsonObject data = responseDoc["data"];
        result.valid = data["valid"] | false;
        result.deviceId = data["device_id"] | 0;
        result.schoolId = data["school_id"] | 0;
        
        // Store device ID for later use (pairing status check)
        if (result.valid && result.deviceId > 0) {
          this->deviceId = result.deviceId;
          Serial.print(F("[API] Stored deviceId: "));
          Serial.println(this->deviceId);
        }
      }
      result.message = responseDoc["message"] | "";
    } else {
      result.message = F("JSON parse error");
    }
  } else {
    result.message = "HTTP error: " + String(httpCode);
  }

  return result;
}

AttendanceResponse APIClient::recordAttendance(const String& rfidCode) {
  AttendanceResponse result;
  result.success = false;
  result.studentName = "";
  result.status = "";
  result.message = "";
  result.errorCode = "";

  if (!isConfigured()) {
    result.message = F("API not configured");
    return result;
  }

  // Build JSON payload - AC2.3
  StaticJsonDocument<128> doc;
  doc["api_key"] = apiKey;
  doc["rfid_code"] = rfidCode;
  
  String payload;
  serializeJson(doc, payload);

  // Send request (using /public prefix to bypass auth)
  String response;
  int httpCode = sendPostRequest("/api/v1/public/attendance/rfid", payload, response);
  
  Serial.print(F("[API] Attendance response ("));
  Serial.print(httpCode);
  Serial.print(F("): "));
  Serial.println(response);

  if (httpCode == 200) {
    // Parse response
    StaticJsonDocument<512> responseDoc;
    DeserializationError error = deserializeJson(responseDoc, response);
    
    if (!error) {
      // Check top-level success first
      bool topSuccess = responseDoc["success"] | false;
      JsonObject data = responseDoc["data"];
      
      // Get student info from data
      result.studentName = data["student_name"] | "";
      result.status = data["type"] | data["status"] | "";  // Backend uses "type" field
      result.message = data["message"] | "";
      
      // Success if top-level success OR data.success is true
      result.success = topSuccess || (data["success"] | false);
      
      // If we got student name, consider it a success for display purposes
      if (result.studentName.length() > 0) {
        result.success = true;
      }
    } else {
      result.message = F("JSON parse error");
    }
  } else if (httpCode == 400 || httpCode == 404) {
    // Parse error response to get specific error code
    StaticJsonDocument<256> errorDoc;
    DeserializationError error = deserializeJson(errorDoc, response);
    
    if (!error) {
      JsonObject errorObj = errorDoc["error"];
      result.errorCode = errorObj["code"] | "";
      result.message = errorObj["message"] | "";
      
      Serial.print(F("[API] Error code: "));
      Serial.println(result.errorCode);
    } else {
      result.message = F("card_not_found");
    }
  } else {
    result.message = "HTTP error: " + String(httpCode);
  }

  return result;
}

PairingResponse APIClient::processPairing(const String& rfidCode) {
  PairingResponse result;
  result.success = false;
  result.studentName = "";
  result.message = "";
  result.sessionActive = false;

  if (!isConfigured()) {
    result.message = F("API not configured");
    return result;
  }

  // Build JSON payload - AC3.4
  StaticJsonDocument<128> doc;
  doc["api_key"] = apiKey;
  doc["rfid_code"] = rfidCode;
  
  String payload;
  serializeJson(doc, payload);
  
  Serial.print(F("[API] Pairing request: "));
  Serial.println(payload);

  // Send request (using /public prefix to bypass auth)
  String response;
  int httpCode = sendPostRequest("/api/v1/public/pairing/rfid", payload, response);
  
  Serial.print(F("[API] Pairing response ("));
  Serial.print(httpCode);
  Serial.print(F("): "));
  Serial.println(response);

  if (httpCode == 200) {
    // Parse response - backend returns flat structure, not nested in "data"
    StaticJsonDocument<256> responseDoc;
    DeserializationError error = deserializeJson(responseDoc, response);
    
    if (!error) {
      result.success = responseDoc["success"] | false;
      result.studentName = responseDoc["student_name"] | "";
      result.message = responseDoc["message"] | "";
    } else {
      result.message = F("JSON parse error");
    }
  } else if (httpCode == 409) {
    result.message = F("card_already_used");
  } else if (httpCode == 400) {
    // Parse error message from response
    StaticJsonDocument<256> responseDoc;
    if (deserializeJson(responseDoc, response) == DeserializationError::Ok) {
      result.message = responseDoc["message"] | "Bad request";
    } else {
      result.message = F("Bad request");
    }
  } else {
    result.message = "HTTP error: " + String(httpCode);
  }

  return result;
}

PairingResponse APIClient::checkPairingStatus() {
  PairingResponse result;
  result.success = false;
  result.studentName = "";
  result.message = "";
  result.sessionActive = false;

  if (!isConfigured() || deviceId == 0) {
    return result;
  }

  // Send GET request with device ID in path (using /public prefix)
  String endpoint = "/api/v1/public/pairing/status/" + String(deviceId);
  
  String response;
  int httpCode = sendGetRequest(endpoint.c_str(), response);

  if (httpCode == 200) {
    // Parse response
    StaticJsonDocument<512> responseDoc;
    DeserializationError error = deserializeJson(responseDoc, response);
    
    if (error) {
      return result;
    }
    
    result.success = responseDoc["success"] | false;
    
    if (result.success && responseDoc.containsKey("data")) {
      JsonObject data = responseDoc["data"];
      result.sessionActive = data["active"] | false;
      result.studentName = data["student_name"] | "";
      
      // Only log when session is active (reduces noise)
      if (result.sessionActive) {
        Serial.print(F("[PAIRING] Session active for: "));
        Serial.println(result.studentName);
      }
    }
  }

  return result;
}

String APIClient::buildURL(const char* endpoint) {
  return serverUrl + String(endpoint);
}

int APIClient::sendPostRequest(const char* endpoint, const String& payload, String& response) {
  int httpCode = -1;
  
  // Retry logic - AC5.3
  for (int retry = 0; retry < API_MAX_RETRIES; retry++) {
    if (retry > 0) {
      Serial.print(F("[INFO] APIClient: Retry "));
      Serial.print(retry);
      Serial.println(F("..."));
      delay(1000);  // Increased delay between retries
    }

    String url = buildURL(endpoint);
    
    // Close any previous connection
    http.end();
    delay(100);
    
    http.begin(wifiClient, url);
    http.setTimeout(API_TIMEOUT);  // AC2.9
    http.addHeader("Content-Type", "application/json");
    
    httpCode = http.POST(payload);
    
    if (httpCode > 0) {
      response = http.getString();
      http.end();
      
      Serial.print(F("[API] POST "));
      Serial.print(endpoint);
      Serial.print(F(" -> "));
      Serial.println(httpCode);
      
      return httpCode;
    }
    
    Serial.print(F("[ERROR] APIClient: Request failed, code: "));
    Serial.println(httpCode);
    http.end();
  }
  
  return httpCode;
}

int APIClient::sendGetRequest(const char* endpoint, String& response) {
  int httpCode = -1;
  
  // Retry logic - AC5.3
  for (int retry = 0; retry < API_MAX_RETRIES; retry++) {
    if (retry > 0) {
      delay(500);
    }

    String url = buildURL(endpoint);
    
    http.begin(wifiClient, url);
    http.setTimeout(API_TIMEOUT);
    
    httpCode = http.GET();
    
    if (httpCode > 0) {
      response = http.getString();
      http.end();
      return httpCode;
    }
    
    http.end();
  }
  
  return httpCode;
}
