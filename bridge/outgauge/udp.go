package outgauge

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"outgouge/config"
	"time"
)

// Some runtime variables
var udpStream *net.UDPConn
var packetSize = binary.Size(OutGaugePacket{})
var Buffer = make([]byte, packetSize)
var RawData OutGaugePacket

// Connects to target UDP address serving raw OutGauge protocol
func Connect(address string, port int) error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		fmt.Printf("Failed to resolve address: %v\n", err)
		return err
	}
	udpStream, err = net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Failed to listen on port: %v\n", err)
	}
	return err
}

// Receive OutGauge UDP packet
func Receive() error {
	n, _, err := udpStream.ReadFromUDP(Buffer)
	if err != nil {
		fmt.Printf("Error reading packet: %v\n", err)
		return err
	}

	// Verify we got the full struct payload
	if n < packetSize {
		fmt.Printf("Received partial packet: %d/%d bytes\n", n, packetSize)
		return err
	}

	// Unpack binary data directly into the struct
	reader := bytes.NewReader(Buffer[:n])
	err = binary.Read(reader, binary.LittleEndian, &RawData)
	if err != nil {
		fmt.Printf("Failed to parse packet: %v\n", err)
	}

	return err
}

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
