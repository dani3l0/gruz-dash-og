class Led {
private:
  uint8_t p;
  bool state = false;

public:
  Led(uint8_t pin) : p(pin) {
    pinMode(p, OUTPUT);
    digitalWrite(p, state);
  }

  void set(bool s) {
    state = s;
    digitalWrite(p, state);
  }
};
