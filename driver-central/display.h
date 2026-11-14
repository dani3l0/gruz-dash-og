#include <U8g2lib.h>
#include <Wire.h>

// Wrapper to use instead of unsafe sprintf
#define mkbuf(buf, ...) snprintf(buf, sizeof(buf), __VA_ARGS__)

// Display initializer class
U8G2_SH1106_128X64_NONAME_F_HW_I2C u8g2(U8G2_R0, U8X8_PIN_NONE);

// Initializes connection to OLED display
void initDisplay() {
  u8g2.begin();
  Wire.setClock(400000);
}

// Some runtime variables
uint32_t updateDisplayTime = 0;
uint8_t displayId = 0;
const uint8_t displaysCount = 6;
uint32_t tripOffset = 0;

// Updates display with all data received?
void updateDisplay(bool buttonClick, bool buttonHold) {
  // Race-condition helper
  OutGaugePacket data;
  mutex_enter_blocking(&packetMutex);
  memcpy(&data, &packet, sizeof(OutGaugePacket));
  mutex_exit(&packetMutex);

  // Button events
  if (data.header[0] == specialChar) {
    if (buttonClick && ++displayId >= displaysCount) displayId = 0;
    if (displayId == 2 && buttonHold) tripOffset = data.odometer;
  }

  // Throttle refresh speed
  uint32_t now = millis();
  if (now - updateDisplayTime <= 50) return;
  updateDisplayTime = now;
  u8g2.clearBuffer();

  // Welcome message
  if (data.header[0] != specialChar) {
    u8g2.setFont(u8g2_font_mystery_quest_24_tf);
    u8g2.drawStr(15, 32, "GruzDash");
    u8g2.setFont(u8g2_font_luRS08_tf);
    u8g2.drawStr(1, 56, "Connect PC to continue");
    u8g2.sendBuffer();
    return;
  }

  // Buffer
  char buf[32];

  // Odometer / Mileage display
  u8g2.setFont(u8g2_font_spleen16x32_mu);
  mkbuf(buf, "%05" PRIu32, data.odometer / 1000);
  u8g2.drawStr(28, 32, buf);
  u8g2.setFont(u8g2_font_missingplanet_tf);
  mkbuf(buf, "%02" PRIu32, (data.odometer % 1000) / 10);
  u8g2.drawStr(112, 21, buf);
  u8g2.drawStr(112, 32, "km");

  // Bottom bars
  u8g2.drawHLine(0, 45, 128);
  u8g2.drawVLine(108, 45, 18);

  // Gear
  const char* gear;
  char gearBuf[4];
  if (data.gear == 0) gear = "R";
  else if (data.gear == 1) gear = "N";
  else {
    mkbuf(gearBuf, "%d", data.gear - 1);
    gear = gearBuf;
  }
  u8g2.setFont(u8g2_font_koleeko_tf);
  u8g2.drawStr(116, 60, gear);

  // Trip calculations
  if (tripOffset == 0) tripOffset = data.odometer;
  uint32_t trip_diff = data.odometer - tripOffset;
  float trip = trip_diff / 1000.0;
  if (trip >= 1000) {
    tripOffset += 1000000;
    trip -= 1000;
  }

  // Switchable display
  u8g2.setFont(u8g2_font_6x13_tf);
  if (displayId == 5 && !(data.flags & OG_TURBO)) displayId = 0;
  if (displayId == 0) {
    time_t t = data.time;
    struct tm *tm_info = gmtime(&t);
    if (tm_info != NULL) strftime(buf, sizeof(buf), "%d %b  %H:%M", tm_info);
    else mkbuf(buf, "-- ---  --:--");
  }
  else if (displayId == 1) {
    mkbuf(buf, "Speed: %.0f km/h", data.speed * 3.6);
  }
  else if (displayId == 2) {
    mkbuf(buf, "Trip: %.2f km", trip);
  }
  else if (displayId == 3) {
    mkbuf(buf, "Fuel: %.1f%%", data.fuel * 100);
  }
  else if (displayId == 4) {
    mkbuf(buf, "Engine: %.1f C", data.engTemp);
  }
  else if (displayId == 5) {
    mkbuf(buf, "Boost: %.2f bar", data.turbo);
  }
  u8g2.drawStr(3, 60, buf);

  // Update display
  u8g2.sendBuffer();
}
