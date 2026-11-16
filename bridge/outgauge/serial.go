package outgauge

import (
	"bytes"
	"encoding/binary"
	"outgouge/config"
	"time"
)

// Strips down 96-byte original UDP packet to smaller one
func StripDownPacket(p OutGaugePacket) OutGaugeStrippedDownPacket {
	now := time.Now()
	_, offset := now.Zone()
	var hd byte = 0xFE
	packt := OutGaugeStrippedDownPacket{
		Header:      [4]byte{hd, hd, hd, hd},
		Time:        uint32(now.Unix() + int64(offset)),
		Flags:       p.Flags,
		Gear:        p.Gear,
		Speed:       p.Speed,
		RPM:         p.RPM,
		Turbo:       p.Turbo,
		EngTemp:     p.EngTemp,
		Fuel:        p.Fuel,
		OilPressure: p.OilPressure,
		OilTemp:     p.OilTemp,
		DashLights:  p.DashLights,
		ShowLights:  p.ShowLights,
		Throttle:    p.Throttle,
		Brake:       p.Brake,
		Clutch:      p.Clutch,
		Odometer:    config.Main.Odometer,
		Checksum:    0,
	}

	// XOR checksum
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, &packt)
	raw := buf.Bytes()
	var calcChecksum uint8 = 0
	for i := 0; i < len(raw)-1; i++ {
		calcChecksum ^= raw[i]
	}
	packt.Checksum = calcChecksum

	return packt
}

// Prepares data for serial sending
func ToSerial(raw OutGaugePacket) []byte {
	p := StripDownPacket(raw)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, &p)
	return buf.Bytes()
}
