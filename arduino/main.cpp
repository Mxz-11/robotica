#include <Adafruit_Sensor.h>
#include <Arduino.h>
#include <ArduinoJson.h>
#include <DHT.h>
#include <DHT_U.h>

#define DHTTYPE DHT22
#define PIR_0_PIN 2
#define DHT_0_PIN 7
#define FIRE_0_PIN A0

DHT_Unified dht(DHT_0_PIN, DHTTYPE);
StaticJsonDocument<256> doc;

void setup()
{
  Serial.begin(9600);
  pinMode(PIR_0_PIN, INPUT);
  pinMode(FIRE_0_PIN, INPUT);
  dht.begin();
  Serial.println("Sistema iniciado.");
}

void loop()
{
  // --------- FUEGO ---------
  int fire_reading = analogRead(FIRE_0_PIN);
  bool fire = fire_reading < 300;

  // --------- PRESENCIA ---------
  bool presence = digitalRead(PIR_0_PIN) == HIGH;

  // --------- TEMPERATURA Y HUMEDAD ---------
  sensors_event_t temp_event, hum_event;
  dht.temperature().getEvent(&temp_event);
  dht.humidity().getEvent(&hum_event);

  float temp = !isnan(temp_event.temperature) ? temp_event.temperature : 0.0;
  float hum = !isnan(hum_event.relative_humidity) ? hum_event.relative_humidity : 0.0;

  // --------- JSON OUTPUT ---------
  doc.clear();
  doc["temperature"] = temp;
  doc["humidity"] = hum;
  doc["presence"] = presence;
  doc["fire"] = fire;

  serializeJson(doc, Serial);
  Serial.println();

  delay(35000);
}
