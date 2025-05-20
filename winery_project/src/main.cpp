#include <Adafruit_Sensor.h>
#include <Arduino.h>
#include <ArduinoJson.h>
#include <DHT.h>
#include <DHT_U.h>

typedef struct toggle_helper
{
  bool last_value;
  bool toggle_state;
} toggle_helper;

enum pins
{
  PIR_0_PIN = 2,
  DHT_0_PIN = 7,
  FIRE_0_PIN = A0
};

#define DHTTYPE DHT22
DHT_Unified dht (pins::DHT_0_PIN, DHTTYPE);
JsonDocument doc;
toggle_helper fire_toggle = { false, false };
toggle_helper presence_toggle = { false, false };

void send_fire ();
void send_presence ();
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
  send_fire ();
  send_presence ();
  float temp = get_temperature ();
  float hum = get_humidity ();
  doc["temperature"] = temp;
  doc["humidity"] = hum;
  serializeJson (doc, Serial);
  Serial.println ();
  doc.clear ();
  delay (5000);
}

void
send_fire ()
{
  bool fire = get_fire ();
  if ((fire && !fire_toggle.last_value) || !fire)
    {
      doc["fire"] = fire;
      serializeJson (doc, Serial);
      Serial.println ();
      doc.clear ();
    }
  fire_toggle.last_value = fire;
}

void
send_presence ()
{
  bool presence = get_presence ();
  if ((presence && !presence_toggle.last_value) || !presence)
    {
      doc["presence"] = presence;
      serializeJson (doc, Serial);
      Serial.println ();
      doc.clear ();
    }
  presence_toggle.last_value = presence;
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