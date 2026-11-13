package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"outgouge/config"
	"outgouge/outgauge"
	"outgouge/serial"
	"time"
)

func main() {
	config.Init()
	err := outgauge.Connect(config.Main.Address, config.Main.Port)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	err = serial.Connect(config.Main.SerialPort, 115200)
	if err != nil {
		fmt.Printf("Failed to connect to port, but continuing\n")
		time.Sleep(time.Second * 3)
	}

	fmt.Printf("Waiting for UDP stream ...\n")

	lastLoopTime := time.Now().UnixMicro()
	odoFloat := float64(config.Main.Odometer)

	for {
		now := time.Now().UnixMicro()
		diff := now - lastLoopTime
		partOfTheSecond := float64(diff) / 1000.0 / 1000.0
		if err := outgauge.Receive(); err == nil {
			outgauge.LogData(outgauge.Parse(outgauge.RawData))
			p := outgauge.StripDownPacket(outgauge.RawData)
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, &p)
			byts := buf.Bytes()
			serial.Write(byts)
		}
		if outgauge.RawData.Speed > 0.05 {
			odoFloat += partOfTheSecond * float64(outgauge.RawData.Speed)
		}
		config.Main.Odometer = uint32(odoFloat)
		time.Sleep(time.Millisecond)
		lastLoopTime = now
	}
}
