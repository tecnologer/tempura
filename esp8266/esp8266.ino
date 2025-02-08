#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <WiFiClient.h>
#include <WEMOS_SHT3X.h>  //https://github.com/wemos/WEMOS_SHT3x_Arduino_Library
#include <Wire.h>

#define LIQUID_SENSOR_PIN D5  // GPIO14
#define R1 33000.0            // Resistor R1 = 33kΩ
#define R2 10000.0            // Resistor R2 = 10kΩ
#define MAX_ADC 1023.0        // 10-bit ADC resolution
#define REF_VOLTAGE 3.3       // ESP8266 ADC max voltage is 1V

const char* ssid = "<WIFI_SSID>";
const char* password = "<WIFI_PASSWORD>";

//Your Domain name with URL path or IP address with path
const char* serverName = "<URL>";

const int EMPTY_TANK = 0;
const int FULL_TANK = 100;
const float FULL_BAT_LEVEL = 4.4;
const float EMPTY_BAT_LEVEL = 3.4;
const float DANGER_BAT_LEVEL = 3.3;


unsigned long timerDelay = 30000;

SHT3X sht30(0x44);

void setup() {
  Serial.begin(115200);

  WiFi.begin(ssid, password);
  Serial.println("Connecting");
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("");
  Serial.print("Connected to WiFi network with IP Address: ");
  Serial.println(WiFi.localIP());

  Serial.println("Timer set to " + String(timerDelay / 1000) + " seconds (timerDelay variable), it will take " + String(timerDelay / 1000) + " seconds before publishing the first reading.");
  pinMode(LIQUID_SENSOR_PIN, INPUT_PULLUP);  // Use internal pull-up resistor
}

float calculateBatPercentage(float batLevel) {
  return ((batLevel - EMPTY_BAT_LEVEL) / (FULL_BAT_LEVEL - EMPTY_BAT_LEVEL)) * 100;
}

void loop() {
  sht30.get();  // Provides temp = sht30.ctemp or sht30.ftemp and sht30.humidity


  char buffer[100];
  float humidity = sht30.humidity;
  float temperature = sht30.fTemp;
  float celciusTemperature = ((temperature - 32) * 5) / 9;
  int fluidLevel = EMPTY_TANK;

  if (digitalRead(LIQUID_SENSOR_PIN) == LOW) {
    fluidLevel = FULL_TANK;
  }

  int adcValue = analogRead(A0);
  float voltageAtA0 = (adcValue / MAX_ADC) * REF_VOLTAGE;  // Voltage at A0

  // Calculate actual battery voltage based on resistor divider
  float batteryVoltage = voltageAtA0 * ((R1 + R2) / R2);
  float batPercentage = calculateBatPercentage(batteryVoltage);

  Serial.print("Battery Voltage: ");
  Serial.print(batteryVoltage);
  Serial.println(" V");
  Serial.print("Read value: ");
  Serial.println(adcValue);                                                                                     

  sprintf(buffer, "Temperature: %.2f°C (%.2f°F), Humidity: %.1f%%, Fluid Level: %d%%, Battery Level: %.2f%%", celciusTemperature, temperature, humidity, fluidLevel, batPercentage);
  Serial.println(buffer);

  //Check WiFi connection status
  if (WiFi.status() == WL_CONNECTED) {
    WiFiClient client;
    HTTPClient http;

    // Your Domain name with URL path or IP address with path
    http.begin(client, serverName);

    // If you need Node-RED/server authentication, insert user and password below
    //http.setAuthorization("REPLACE_WITH_SERVER_USERNAME", "REPLACE_WITH_SERVER_PASSWORD");

    // If you need an HTTP request with a content type: application/json, use the following:
    http.addHeader("Content-Type", "application/json");
    int httpResponseCode = http.POST("{\"temperature\": " + String(celciusTemperature) + ",\"humidity\": " + String(humidity) + ", \"fluid_level\": " + String(fluidLevel) + ", \"bat_level\": " + String(batPercentage) + "}");

    Serial.print("HTTP Response code: ");
    Serial.println(httpResponseCode);

    // Free resources
    http.end();
  } else {
    Serial.println("WiFi Disconnected");
  }

  delay(timerDelay);
}
