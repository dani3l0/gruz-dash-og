#define CENTRAL
#include "config.h"
#include "outgauge.h"
#include "display.h"
#include "led.h"
#include "button.h"

// Wow
Led LHeadlights(LED_HEADLIGHTS);
Led LCheckEngine(LED_CHECK_ENGINE);
Led LTickR(LED_TICK_R);
Led LTickL(LED_TICK_L);
Led LBacklight(LED_BACKLIGHT);
Led LHandbrake(LED_HANDBRAKE);
Led LABS(LED_ABS);
Led LTraction(LED_TRACTION);
Led LBattery(LED_BATTERY);
Button btn(LED_BUTTON);

// Core 0 //
void setup() {
  Serial.begin(SERIAL_BAUD_RATE);
  Serial1.setTX(TX_PIN);
  Serial1.setRX(RX_PIN);
  Serial1.begin(SERIAL_BAUD_RATE);
  LBacklight.set(true);
}

void loop() {
  readOG();
  LHeadlights.set(!(packet.showLights & DL_FULLBEAM));
  LCheckEngine.set(packet.fuel < 0.15 || packet.engTemp > 105 || packet.oilTemp > 120);
  LTickR.set(packet.showLights & DL_SIGNAL_R);
  LTickL.set(packet.showLights & DL_SIGNAL_L);
  LHandbrake.set(packet.showLights & DL_HANDBRAKE);
  LABS.set(packet.showLights & DL_ABS);
  LTraction.set(packet.showLights & DL_TC);
  LBattery.set(packet.showLights & DL_BATTERY);
}



// Offload display to Core 1 //
void setup1() {
  initDisplay();
}

void loop1() {
  btn.update();
  updateDisplay(btn.click, btn.hold);
}
