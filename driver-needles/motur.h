class Motur {
private:
  uint8_t pA1, pA2, pB1, pB2;
  float totalSteps;
  
  float currentPos = 5.0;
  float targetPos = 0.0;
  
  unsigned long lastStepTime = 0;
  uint16_t stepIntervalMicros = 500; 
  float lastAppliedStep = 0.0;

  float sinTable[100];
  float getSin(int index400) {
    int phase = index400 / 100;
    int idx = index400 % 100;
    switch (phase) {
      case 0: return sinTable[idx];
      case 1: return sinTable[99 - idx];
      case 2: return -sinTable[idx];
      case 3: return -sinTable[99 - idx];
      default: return 0.0;
    }
  }

  void drive(float stepNumber) {
    float phase = fmod(stepNumber, 4.0);
    if (phase < 0) phase += 4.0;

    int indexA = (int)(phase * 100.0);
    if (indexA > 399) indexA = 399;
    int indexB = (indexA + 100) % 400;

    float valA = getSin(indexA);
    float valB = getSin(indexB);

    if (valA >= 0) {
      analogWrite(pA1, round(valA * 255.0));
      analogWrite(pA2, 0);
    } else {
      analogWrite(pA1, 0);
      analogWrite(pA2, round(-valA * 255.0));
    }

    if (valB >= 0) {
      analogWrite(pB1, round(valB * 255.0));
      analogWrite(pB2, 0);
    } else {
      analogWrite(pB1, 0);
      analogWrite(pB2, round(-valB * 255.0));
    }
  }

public:
  Motur(float steps, uint8_t pinA1, uint8_t pinA2, uint8_t pinB1, uint8_t pinB2) 
    : totalSteps(steps), pA1(pinA1), pA2(pinA2), pB1(pinB1), pB2(pinB2) {}

  void begin() {
    pinMode(pA1, OUTPUT);
    pinMode(pA2, OUTPUT);
    pinMode(pB1, OUTPUT);
    pinMode(pB2, OUTPUT);

    for (int i = 0; i < 100; i++) {
      float phase = (float)i / 100.0;
      float rad = phase * (PI / 2.0);
      sinTable[i] = sin(rad);
    }
  }

  void setPosition(float percentage) {
    if (percentage < 0.0) percentage = 0.0;
    if (percentage > 100.0) percentage = 100.0;
    targetPos = (percentage / 100.0) * totalSteps;
  }

  float getPosition() {
    return (currentPos / totalSteps) * 100.0;
  }

  void update() {
    unsigned long now = micros();
    if (now - lastStepTime < stepIntervalMicros) {
      return; 
    }
    lastStepTime = now;

    float diff = targetPos - currentPos;
    float absDiff = abs(diff);

    if (absDiff < 0.01) {
      currentPos = targetPos;
      lastAppliedStep = 0.0;
      drive(currentPos);
      return;
    }

    float maxStep = 0.05;
    float brakingZone = 8; 

    if (absDiff < brakingZone) {
      float factor = absDiff / brakingZone;
      maxStep *= (factor * factor); 
    }

    float maxAccelerationPerCycle = 0.002;
    if (maxStep > lastAppliedStep + maxAccelerationPerCycle) {
      maxStep = lastAppliedStep + maxAccelerationPerCycle;
    }

    if (maxStep < 0.004) {
      maxStep = 0.004;
    }
    
    lastAppliedStep = maxStep;
      
    if (diff > maxStep) {
      currentPos += maxStep;
    } else if (diff < -maxStep) {
      currentPos -= maxStep;
    } else {
      currentPos = targetPos;
      lastAppliedStep = 0.0;
    }

    drive(currentPos);
  }
};
