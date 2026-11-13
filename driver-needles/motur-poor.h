class PoorMotur {
private:
  uint8_t pA1, pA2, pB1, pB2;
  float totalSteps;
  
  float currentPos = 5.0;
  float targetPos = 0.0;
  
  unsigned long lastStepTime = 0;
  uint16_t stepIntervalMicros = 500;
  float lastAppliedStep = 0.0;

  const uint8_t halfStepTable[8][4] = {
    {1, 0, 0, 0},
    {1, 0, 1, 0},
    {0, 0, 1, 0},
    {0, 1, 1, 0},
    {0, 1, 0, 0},
    {0, 1, 0, 1},
    {0, 0, 0, 1},
    {1, 0, 0, 1}
  };

  void drive(float stepNumber) {
    float phase = fmod(stepNumber, 4.0);
    if (phase < 0) phase += 4.0;

    int stateIndex = (int)(phase * 2.0);
    if (stateIndex > 7) stateIndex = 7;

    digitalWrite(pA1, halfStepTable[stateIndex][0]);
    digitalWrite(pA2, halfStepTable[stateIndex][1]);
    digitalWrite(pB1, halfStepTable[stateIndex][2]);
    digitalWrite(pB2, halfStepTable[stateIndex][3]);
  }

public:
  PoorMotur(float steps, uint8_t pinA1, uint8_t pinA2, uint8_t pinB1, uint8_t pinB2) 
    : totalSteps(steps), pA1(pinA1), pA2(pinA2), pB1(pinB1), pB2(pinB2) {}

  void begin() {
    pinMode(pA1, OUTPUT);
    pinMode(pA2, OUTPUT);
    pinMode(pB1, OUTPUT);
    pinMode(pB2, OUTPUT);
    drive(currentPos);
  }

  void setPosition(float percentage) {
    if (percentage < 0.0) percentage = 0.0;
    if (percentage > 100.0) percentage = 100.0;
    targetPos = (percentage / 100.0) * totalSteps;
  }

  void update() {
    unsigned long now = micros();
    if (now - lastStepTime < stepIntervalMicros) {
      return; 
    }
    lastStepTime = now;

    float diff = targetPos - currentPos;
    float absDiff = abs(diff);

    if (absDiff < 0.05) {
      currentPos = targetPos;
      lastAppliedStep = 0.0;
      drive(currentPos);
      return;
    }

    float maxStep = 0.015;
    float brakingZone = 6.0;

    if (absDiff < brakingZone) {
      float factor = absDiff / brakingZone;
      maxStep *= (factor * factor); 
    }

    float maxAccelerationPerCycle = 0.001;
    if (maxStep > lastAppliedStep + maxAccelerationPerCycle) {
      maxStep = lastAppliedStep + maxAccelerationPerCycle;
    }

    if (maxStep < 0.003) {
      maxStep = 0.003;
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
