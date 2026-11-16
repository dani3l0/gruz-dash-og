package outgauge

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// Some runtime variables
var udpStream *net.UDPConn
var packetSize = binary.Size(OutGaugePacket{})
var Buffer = make([]byte, packetSize)
var RawData = OutGaugePacket{
	Gear: 1,
}

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
