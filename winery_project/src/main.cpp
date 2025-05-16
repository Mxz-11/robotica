#include <Adafruit_Sensor.h>
#include <Arduino.h>
#include <ArduinoJson.h>
#include <DHT.h>
#include <DHT_U.h>

enum pins
{
  PIR_0_PIN = 2,
  DHT_0_PIN = 7,
  FIRE_0_PIN = A0
};

#define DHTTYPE DHT22
DHT_Unified dht (pins::DHT_0_PIN, DHTTYPE);
JsonDocument doc;

bool get_fire ();
bool get_presence ();
float get_temperature ();
float get_humidity ();

void
setup ()
{
  Serial.begin (9600);
  pinMode (pins::PIR_0_PIN, INPUT);
  pinMode (pins::FIRE_0_PIN, INPUT);
  dht.begin ();
}

void
loop ()
{
  bool fire = get_fire ();
  bool presence = get_presence ();
  float temp = get_temperature ();
  float hum = get_humidity ();
  doc.clear ();
  doc["temperature"] = temp;
  doc["humidity"] = hum;
  doc["presence"] = presence;
  doc["fire"] = fire;
  serializeJson (doc, Serial);
  Serial.println ();
  delay (5000);
}

bool
get_fire ()
{
  return analogRead (pins::FIRE_0_PIN) < 300;
}

bool
get_presence ()
{
  return digitalRead (pins::PIR_0_PIN) == HIGH;
}

float
get_temperature ()
{
  sensors_event_t temp_event;
  dht.temperature ().getEvent (&temp_event);
  return !isnan (temp_event.temperature) ? temp_event.temperature : 0.0f;
}

float
get_humidity ()
{
  sensors_event_t hum_event;
  dht.humidity ().getEvent (&hum_event);
  return !isnan (hum_event.relative_humidity) ? hum_event.relative_humidity
                                              : 0.0f;
}