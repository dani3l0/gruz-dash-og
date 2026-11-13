#include "config.h"
#include "motur.h"
#include "motur-poor.h"
#include "outgauge.h"

Motur BigNeedle(BigNeedleSteps, BigNeedlePinA1, BigNeedlePinA2, BigNeedlePinB1, BigNeedlePinB2);
PoorMotur SmallNeedle(SmallNeedleSteps, SmallNeedlePinA1, SmallNeedlePinA2, SmallNeedlePinB1, SmallNeedlePinB2);

void setup() {
  // Arduino Nano annoying noise tweak
  TCCR1B = (TCCR1B & 0b11111000) | 0x01;
  TCCR2B = (TCCR2B & 0b11111000) | 0x01;

  // The rest of initialization stuff...
  Serial.begin(115200);
  BigNeedle.begin();
  SmallNeedle.begin();
}

void loop() {
  if (millis() > 2000) readOG();

  #if isRPM
    BigNeedle.setPosition(100.0 * (packet.rpm - BigNeedleMinValue) / (BigNeedleMaxValue - BigNeedleMinValue));
  #else
    BigNeedle.setPosition(100.0 * (packet.speed * 3.6 - BigNeedleMinValue) / (BigNeedleMaxValue - BigNeedleMinValue));
  #endif

  #if isCoolant
    SmallNeedle.setPosition(100.0 * (packet.engTemp - SmallNeedleMinValue) / (SmallNeedleMaxValue - SmallNeedleMinValue));
  #else
    SmallNeedle.setPosition(100.0 * (packet.fuel - SmallNeedleMinValue) / (SmallNeedleMaxValue - SmallNeedleMinValue));
  #endif
 
  BigNeedle.update();
  SmallNeedle.update();
}
