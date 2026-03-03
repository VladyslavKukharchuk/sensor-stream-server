#include "DHT.h"
#include <WiFi.h>
#include <HTTPClient.h>
#include "secrets.h"
#include <time.h>
#include <SPIFFS.h>
#include <ArduinoJson.h>

#define DHTPIN 4
#define DHTTYPE DHT22
DHT dht(DHTPIN, DHTTYPE);

const char* ntpServer = "pool.ntp.org";
const long gmtOffset_sec = 7200; // UTC+2
const int daylightOffset_sec = 0;

String deviceId;

void setupTime() {
  configTime(gmtOffset_sec, daylightOffset_sec, ntpServer);
}

String getTimestamp() {
  struct tm timeinfo;
  if(!getLocalTime(&timeinfo)){
    Serial.println("Failed to obtain time");
    return "";
  }
  char buf[25];
  strftime(buf, sizeof(buf), "%Y-%m-%dT%H:%M:%SZ", &timeinfo);
  return String(buf);
}

void connectToWiFi() {
  Serial.print("Connecting to Wi-Fi...");
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println();
  Serial.print("Connected! IP: ");
  Serial.println(WiFi.localIP());
}

bool readDeviceIdFromFS() {
  if(!SPIFFS.begin(true)) {
    Serial.println("Failed to mount SPIFFS");
    return false;
  }
  if(SPIFFS.exists("/device_id.txt")) {
    File f = SPIFFS.open("/device_id.txt", "r");
    if(f){
      deviceId = f.readString();
      f.close();
      Serial.print("Loaded device_id from SPIFFS: ");
      Serial.println(deviceId);
      return true;
    }
  }
  return false;
}

void saveDeviceIdToFS(const String& id) {
  File f = SPIFFS.open("/device_id.txt", "w");
  if(f){
    f.print(id);
    f.close();
    Serial.println("Saved device_id to SPIFFS");
  }
}

bool registerDevice(String mac_address) {
  HTTPClient http;
  http.begin(String(SERVER_URL) + "/api/v1/devices");
  http.addHeader("Content-Type", "application/json");

  String payload = "{\"mac\":\"" + mac_address + "\"}";

  int httpResponseCode = http.POST(payload);

  if (httpResponseCode > 0) {
    String response = http.getString();
    Serial.print("Register response: ");
    Serial.println(response);

    StaticJsonDocument<200> doc;
    DeserializationError error = deserializeJson(doc, response);

    if (error) {
      Serial.print("JSON parse error: ");
      Serial.println(error.c_str());
      http.end();
      return false;
    }

    if (doc.containsKey("id")) {
      deviceId = doc["id"].as<String>();
      saveDeviceIdToFS(deviceId);

      http.end();
      return true;
    } else {
      Serial.println("No 'id' field in response");
    }
  } else {
    Serial.print("Failed to register device. HTTP response code: ");
    Serial.println(httpResponseCode);
  }

  http.end();
  return false;
}

void sendData(float temperature, float humidity) {
  if (WiFi.status() == WL_CONNECTED && deviceId.length() > 0) {
    HTTPClient http;
    http.begin(String(SERVER_URL) + "/api/v1/measurements");
    http.addHeader("Content-Type", "application/json");

    String timestamp = getTimestamp();

    StaticJsonDocument<200> doc;

    doc["device_id"] = deviceId;
    doc["temperature"] = temperature;
    doc["humidity"] = humidity;
    doc["timestamp"] = timestamp;

    String payload;
    serializeJson(doc, payload);

    int httpResponseCode = http.POST(payload);

    if (httpResponseCode > 0) {
      Serial.print("HTTP Response code: ");
      Serial.println(httpResponseCode);
    } else {
      Serial.print("Error on sending POST: ");
      Serial.println(http.errorToString(httpResponseCode));
    }

    http.end();
  } else {
    Serial.println("Wi-Fi not connected or device_id empty");
  }
}

void setup() {
  Serial.begin(115200);
  Serial.println("ESP32 start!");
  dht.begin();

  connectToWiFi();
  setupTime();

  if(!readDeviceIdFromFS()){
    Serial.println("Device ID not found, registering...");
    if(!registerDevice(WiFi.macAddress())){
      Serial.println("Failed to register device. Retry after restart.");
    }
  }
}

unsigned long lastSendTime = 0;

void loop() {
  unsigned long now = millis();

  if (now - lastSendTime >= SEND_INTERVAL_SEC * 1000UL) {
    lastSendTime = now;

    float h = dht.readHumidity();
    float t = dht.readTemperature();

    if (isnan(h) || isnan(t)) {
      Serial.println("Error read from DHT22!");
    } else {
      Serial.print("Humidity: ");
      Serial.print(h);
      Serial.print("%  |  Temperature: ");
      Serial.print(t);
      Serial.println("°C");

      sendData(t, h);
    }
  }
}
