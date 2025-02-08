#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <WiFiClient.h>
#include <WEMOS_SHT3X.h>  //https://github.com/wemos/WEMOS_SHT3x_Arduino_Library
#include <Wire.h>

#define LIQUID_SENSOR_PIN D5  // GPIO14

const char* ssid = "<WIFI_SSID>";
const char* password = "<WIFI_PASSWORD>";

//Your Domain name with URL path or IP address with path
const char* serverName = "<URL>";

const int EMPTY_TANK = 0 ;
const int FULL_TANK = 100;

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

void loop() {
  sht30.get();  // Provides temp = sht30.ctemp or sht30.ftemp and sht30.humidity


  char buffer[100];
  float humidity = sht30.humidity;
  float temperature = sht30.fTemp;
  float celciusTemperature = ((temperature - 32) * 5) / 9;
  int fluidLevel = EMPTY_TANK;
  
  if(digitalRead(LIQUID_SENSOR_PIN) == LOW) {
    fluidLevel = FULL_TANK;
  }

  sprintf(buffer, "Temperature: %.2f°C (%.2f°F), Humidity: %.1f%%, Fluid Level: %d%%", celciusTemperature, temperature, humidity, fluidLevel);
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
    int httpResponseCode = http.POST("{\"temperature\": " + String(celciusTemperature) + ",\"humidity\": " + String(humidity) + ", \"fluid_level\": " + String(fluidLevel) + ", \"bat_level\": 0}");

    Serial.print("HTTP Response code: ");
    Serial.println(httpResponseCode);

    // Free resources
    http.end();
  } else {
    Serial.println("WiFi Disconnected");
  }

  delay(timerDelay);
}
