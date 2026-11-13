class Button {
private:
  uint8_t pin;
  bool lastRaw = false;
  bool stable = false;
  uint32_t debounceTime = 0;
  uint32_t pressTime = 0;

public:
  bool click = false;
  bool hold = false;

  Button(uint8_t p) : pin(p) {
    pinMode(p, INPUT_PULLUP);
  }

  // Read and update button state
  void update() {
    click = false;
    hold = false;
    bool raw = !digitalRead(pin);
    uint32_t now = millis();

    // Debounce filter, for poor buttons firing multiple times at once
    if (raw != lastRaw) debounceTime = now;
    lastRaw = raw;

    // If deboune filtered, assume this is the 'stable' button state
    if ((now - debounceTime) > 50 && raw != stable) {
      stable = raw;
      if (stable) pressTime = now;
      else if (pressTime) click = true;
    }

    // When button is down for 2s or more, fire hold event
    if (stable && pressTime && (now - pressTime) >= 2000) {
      hold = true;
      pressTime = 0;
    }
  }
};
