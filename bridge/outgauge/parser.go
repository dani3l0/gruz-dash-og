package outgauge

import "outgouge/config"

// Pretty-printed
type OutGaugeData struct {
	Time       uint32
	Speed      float32
	RPM        int
	Gear       int
	Boost      float32
	EngineTemp float32
	OilTemp    float32
	Fuel       float32
	Handbrake  bool
	FullBeam   bool
	SignalL    bool
	SignalR    bool
	TC         bool
	ABS        bool
	Battery    bool
	Odometer   uint32
}

// Parses raw packet data and returns a nicer struct
func Parse(raw OutGaugePacket) OutGaugeData {
	var output OutGaugeData
	output.Time = raw.Time
	output.Speed = raw.Speed * 3.6
	output.RPM = int(raw.RPM)
	output.Gear = int(raw.Gear) - 1
	output.Boost = raw.Turbo
	output.EngineTemp = raw.EngTemp
	output.OilTemp = raw.OilTemp
	output.Fuel = raw.Fuel
	output.Handbrake = raw.ShowLights&DL_HANDBRAKE != 0
	output.FullBeam = raw.ShowLights&DL_FULLBEAM != 0
	output.SignalL = raw.ShowLights&DL_SIGNAL_L != 0
	output.SignalR = raw.ShowLights&DL_SIGNAL_R != 0
	output.TC = raw.ShowLights&DL_TC != 0
	output.ABS = raw.ShowLights&DL_ABS != 0
	output.Battery = raw.ShowLights&DL_BATTERY != 0
	output.Odometer = config.Main.Odometer
	return output
}
