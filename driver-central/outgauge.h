#include <Arduino.h>

// Stripped-down OutGauge packet structure sent by outgouge via Serial
struct __attribute__((packed)) OutGaugePacket {
  char           header[4];
  uint32_t       time;
  uint16_t       flags;
  uint8_t        gear;
  float          speed;
  float          rpm;
  float          turbo;
  float          engTemp;
  float          fuel;
  float          oilPressure;
  float          oilTemp;
  uint32_t       dashLights;
  uint32_t       showLights;
  float          throttle;
  float          brake;
  float          clutch;
  uint32_t       odometer;
  uint8_t        checksum;
};

// Bits for flags detection
enum DashSetupBits : uint16_t {
  OG_SHIFT =     1,
  OG_CTRL  =     2,
  OG_TURBO =  8192,
  OG_KM    = 16384,
  OG_BAR   = 32768
};

// Bits for dashboard lights detection
enum DashLightsBits : uint32_t {
  DL_SHIFT      = (1UL << 0),
  DL_FULLBEAM   = (1UL << 1),
  DL_HANDBRAKE  = (1UL << 2),
  DL_PITSPEED   = (1UL << 3),
  DL_TC         = (1UL << 4),
  DL_SIGNAL_L   = (1UL << 5),
  DL_SIGNAL_R   = (1UL << 6),
  DL_SIGNAL_ANY = (1UL << 7),
  DL_OILWARN    = (1UL << 8),
  DL_BATTERY    = (1UL << 9),
  DL_ABS        = (1UL << 10)
};

// Main data storage
OutGaugePacket packet;
auto_init_mutex(packetMutex);

// Some runtime variables
#define packetSize sizeof(OutGaugePacket)
byte buff[packetSize];
byte lastChar1;
byte lastChar2;
byte lastChar3;
byte specialChar = 0xFE;
uint8_t buffIndex = packetSize + 1;
uint32_t lastPacketTime = 0;

// Reads stripped OutGauge packet via Serial
void readOG() {
  uint32_t now = millis();

  while (Serial.available()) {
    byte c = Serial.read();

    // If main device, echo received data as well
    #ifdef CENTRAL
      Serial1.write(c);
    #endif

    // Packet header detected
    if (c == specialChar && lastChar1 == c && lastChar2 == c && lastChar3 == c) {
      buffIndex = 0;
      while (buffIndex < 4) buff[buffIndex++] = c;
      c = 0x00;
    }

    // Incomplete packet
    else if (buffIndex < packetSize) {
      buff[buffIndex++] = c;
    }

    // Remember last chars
    lastChar3 = lastChar2;
    lastChar2 = lastChar1;
    lastChar1 = c;
  }

  // Packet is complete
  if (buffIndex == packetSize) {
    buffIndex++;
    // Calculate XOR checksum
    uint8_t calcChecksum = 0;
    for (unsigned int i = 0; i < packetSize - 1; i++) {
      calcChecksum ^= buff[i];
    }
    uint8_t receivedChecksum = buff[packetSize - 1];
    // If ok, load packet into memory
    if (calcChecksum == receivedChecksum) {
      mutex_enter_blocking(&packetMutex);
      memcpy(&packet, buff, packetSize);
      mutex_exit(&packetMutex);
      lastPacketTime = now;
    }
  }

  // If no data for more than 5s, zero some visible data
  if (now - lastPacketTime > 5000) {
    lastPacketTime = now;
    packet.speed = 0;
    packet.rpm = 0;
    packet.fuel = 0;
    packet.engTemp = 0;
    packet.dashLights = 0;
    packet.showLights = 0;
    packet.header[0] = 'X';
    packet.header[1] = 'D';
    packet.header[2] = 'D';
    packet.header[3] = 'D';
  }
}
