package outgauge

// OutGauge raw UDP packet structure
type OutGaugePacket struct {
	Time        uint32   // time in milliseconds
	Car         [4]byte  // Car name
	Flags       uint16   // Info (OG_x flags)
	Gear        uint8    // Reverse:0, Neutral:1, First:2...
	Plid        uint8    // Unique ID of viewed player
	Speed       float32  // M/S
	RPM         float32  // RPM
	Turbo       float32  // BAR
	EngTemp     float32  // C
	Fuel        float32  // 0 to 1
	OilPressure float32  // BAR
	OilTemp     float32  // C
	DashLights  uint32   // Dash lights available
	ShowLights  uint32   // Dash lights currently switched on
	Throttle    float32  // 0 to 1
	Brake       float32  // 0 to 1
	Clutch      float32  // 0 to 1
	Display1    [16]byte // Usually Fuel
	Display2    [16]byte // Usually Settings
	ID          int32    // OutGauge ID (optional, but standard in full struct)
}

// OutGauge stripped-down UDP packet for Serial communication
type OutGaugeStrippedDownPacket struct {
	Header      [4]byte
	Time        uint32
	Flags       uint16
	Gear        uint8
	Speed       float32
	RPM         float32
	Turbo       float32
	EngTemp     float32
	Fuel        float32
	OilPressure float32
	OilTemp     float32
	DashLights  uint32
	ShowLights  uint32
	Throttle    float32
	Brake       float32
	Clutch      float32
	Odometer    uint32
	Checksum    uint8
}

// Flags (OG_x)
const (
	OG_SHIFT = 1 << 0
	OG_CTRL  = 1 << 1
	OG_TURBO = 1 << 13
	OG_KM    = 1 << 14
	OG_BAR   = 1 << 15
)

// Dash Lights (DL_x)
const (
	DL_SHIFT     = 1 << 0
	DL_FULLBEAM  = 1 << 1
	DL_HANDBRAKE = 1 << 2
	DL_PITSPEED  = 1 << 3
	DL_TC        = 1 << 4
	DL_SIGNAL_L  = 1 << 5
	DL_SIGNAL_R  = 1 << 6
	DL_OILWARN   = 1 << 8
	DL_BATTERY   = 1 << 9
	DL_ABS       = 1 << 10
)
